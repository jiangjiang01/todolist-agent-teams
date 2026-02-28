package tests

import (
	"database/sql"
	"net/http/httptest"
	"testing"

	"todolist/internal/database"
	"todolist/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// TestSuite 测试套件
type TestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *sql.DB
}

// SetupSuite 测试套件初始化
func (s *TestSuite) SetupSuite() {
	// 设置 Gin 为测试模式
	gin.SetMode(gin.TestMode)
}

// SetupTest 每个测试前执行
func (s *TestSuite) SetupTest() {
	// 初始化内存数据库，添加 _loc=auto 参数以正确处理 time.Time
	var err error
	s.db, err = sql.Open("sqlite3", ":memory:?_loc=auto")
	s.Require().NoError(err, "无法打开内存数据库")

	// 验证连接可用
	err = s.db.Ping()
	s.Require().NoError(err, "数据库连接失败")

	// 设置全局数据库连接
	database.DB = s.db

	// 创建表结构
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			completed BOOLEAN NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);

		CREATE INDEX IF NOT EXISTS idx_todos_completed ON todos(completed);
		CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);
	`)
	s.Require().NoError(err, "创建表失败")

	// 设置测试路由
	s.setupRouter()
}

// TearDownTest 每个测试后执行
func (s *TestSuite) TearDownTest() {
	if s.db != nil {
		s.db.Close()
	}
}

// setupRouter 设置测试路由
func (s *TestSuite) setupRouter() {
	s.router = gin.New()
	s.router.Use(gin.Recovery())

	api := s.router.Group("/api")
	{
		todos := api.Group("/todos")
		{
			todos.GET("", handlers.GetTodos)
			todos.GET("/:id", handlers.GetTodo)
			todos.POST("", handlers.CreateTodo)
			todos.PUT("/:id", handlers.UpdateTodo)
			todos.DELETE("/:id", handlers.DeleteTodo)
		}
	}
}

// makeRequest 辅助方法：发起 HTTP 请求
func (s *TestSuite) makeRequest(method, path string, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	if body != "" {
		req = httptest.NewRequest(method, path, nil)
		req.Header.Set("Content-Type", "application/json")
		// 在实际测试中会设置 body
	}
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)
	return w
}

// TestTestSuite 运行测试套件
func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
