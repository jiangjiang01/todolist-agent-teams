package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"todolist/internal/models"

	"github.com/stretchr/testify/assert"
)

// 辅助函数：创建 JSON 请求体
func createJSONBody(data interface{}) (*bytes.Reader, error) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBytes), nil
}

// GetTodosSuite 测试获取所有 Todo
func (s *TestSuite) TestGetTodos() {
	s.Run("空列表返回", func() {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/todos", nil)
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code)
		assert.Contains(s.T(), w.Body.String(), `"code":0`)
		assert.Contains(s.T(), w.Body.String(), `"message":"success"`)

		var resp models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		// 空列表时返回空数组
		todos, ok := resp.Data.([]interface{})
		assert.True(s.T(), ok)
		assert.Empty(s.T(), todos)
	})

	s.Run("多个 todo 返回", func() {
		// 插入测试数据
		_, err := s.db.Exec(`
			INSERT INTO todos (title, description, completed, created_at, updated_at)
			VALUES ('测试标题1', '测试描述1', 0, datetime('now'), datetime('now')),
			       ('测试标题2', '测试描述2', 1, datetime('now'), datetime('now'))
		`)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/todos", nil)
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)

		todos, ok := resp.Data.([]interface{})
		assert.True(s.T(), ok)
		assert.Len(s.T(), todos, 2)
	})
}

// GetTodoSuite 测试获取单个 Todo
func (s *TestSuite) TestGetTodo() {
	s.Run("正常返回存在的 todo", func() {
		// 先创建一个 todo
		result, err := s.db.Exec(`
			INSERT INTO todos (title, description, completed, created_at, updated_at)
			VALUES ('测试标题', '测试描述', 0, datetime('now'), datetime('now'))
		`)
		assert.NoError(s.T(), err)
		id, _ := result.LastInsertId()

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/api/todos/%d", id), nil)
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 0, resp.Code)

		// 将 data 转换为 map
		data, ok := resp.Data.(map[string]interface{})
		assert.True(s.T(), ok)
		assert.Equal(s.T(), "测试标题", data["title"])
		assert.Equal(s.T(), "测试描述", data["description"])
	})

	s.Run("返回 404 当 todo 不存在", func() {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/todos/999999", nil)
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusNotFound, w.Code)

		var resp models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 1002, resp.Code)
		assert.Contains(s.T(), resp.Message, "不存在")
	})

	s.Run("无效 ID 返回 400", func() {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/todos/invalid", nil)
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	})
}

// CreateTodoSuite 测试创建 Todo
func (s *TestSuite) TestCreateTodo() {
	s.Run("正常创建（只有标题）", func() {
		// 先验证数据库连接正常
		var count int
		err := s.db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
		assert.NoError(s.T(), err, "数据库查询失败")

		bodyData := map[string]string{"title": "新待办事项"}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		// 如果失败，打印响应体
		if w.Code != http.StatusCreated {
			s.T().Logf("响应体: %s", w.Body.String())
		}

		assert.Equal(s.T(), http.StatusCreated, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 0, resp.Code)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(s.T(), ok)
		assert.Equal(s.T(), "新待办事项", data["title"])
		assert.Equal(s.T(), false, data["completed"])
	})

	s.Run("正常创建（标题+描述）", func() {
		bodyData := map[string]string{
			"title":       "新待办事项",
			"description": "这是详细描述",
		}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusCreated, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(s.T(), ok)
		assert.Equal(s.T(), "新待办事项", data["title"])
		assert.Equal(s.T(), "这是详细描述", data["description"])
	})

	s.Run("参数验证失败（空标题）", func() {
		bodyData := map[string]string{"title": ""}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusBadRequest, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 1001, resp.Code)
		assert.Contains(s.T(), resp.Message, "参数验证失败")
	})

	s.Run("参数验证失败（超长标题 > 200字符）", func() {
		longTitle := strings.Repeat("a", 201)
		bodyData := map[string]string{"title": longTitle}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusBadRequest, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 1001, resp.Code)
	})

	s.Run("参数验证失败（缺少标题）", func() {
		bodyData := map[string]string{}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusBadRequest, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
	})
}

// UpdateTodoSuite 测试更新 Todo
func (s *TestSuite) TestUpdateTodo() {
	s.Run("正常更新全部字段", func() {
		// 先创建一个 todo
		result, err := s.db.Exec(`
			INSERT INTO todos (title, description, completed, created_at, updated_at)
			VALUES ('原标题', '原描述', 0, datetime('now'), datetime('now'))
		`)
		assert.NoError(s.T(), err)
		id, _ := result.LastInsertId()

		bodyData := map[string]interface{}{
			"title":       "更新后标题",
			"description": "更新后描述",
			"completed":   true,
		}
		body, _ := createJSONBody(bodyData)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/todos/%d", id), body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(s.T(), ok)
		assert.Equal(s.T(), "更新后标题", data["title"])
		assert.Equal(s.T(), "更新后描述", data["description"])
		assert.Equal(s.T(), true, data["completed"])
	})

	s.Run("部分更新（仅 completed）", func() {
		// 先创建一个 todo
		result, err := s.db.Exec(`
			INSERT INTO todos (title, description, completed, created_at, updated_at)
			VALUES ('原标题', '原描述', 0, datetime('now'), datetime('now'))
		`)
		assert.NoError(s.T(), err)
		id, _ := result.LastInsertId()

		bodyData := map[string]interface{}{"completed": true}
		body, _ := createJSONBody(bodyData)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", fmt.Sprintf("/api/todos/%d", id), body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(s.T(), ok)
		// 标题和描述应该保持不变
		assert.Equal(s.T(), "原标题", data["title"])
		assert.Equal(s.T(), "原描述", data["description"])
		assert.Equal(s.T(), true, data["completed"])
	})

	s.Run("返回 404 当 todo 不存在", func() {
		bodyData := map[string]string{"title": "更新标题"}
		body, _ := createJSONBody(bodyData)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/todos/999999", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusNotFound, w.Code)

		var resp models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 1002, resp.Code)
	})
}

// DeleteTodoSuite 测试删除 Todo
func (s *TestSuite) TestDeleteTodo() {
	s.Run("正常删除", func() {
		// 先创建一个 todo
		result, err := s.db.Exec(`
			INSERT INTO todos (title, description, completed, created_at, updated_at)
			VALUES ('待删除项', '测试描述', 0, datetime('now'), datetime('now'))
		`)
		assert.NoError(s.T(), err)
		id, _ := result.LastInsertId()

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/todos/%d", id), nil)
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusOK, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 0, resp.Code)
		assert.Equal(s.T(), "deleted", resp.Message)

		// 验证确实被删除了
		var count int
		s.db.QueryRow("SELECT COUNT(*) FROM todos WHERE id = ?", id).Scan(&count)
		assert.Equal(s.T(), 0, count)
	})

	s.Run("返回 404 当 todo 不存在", func() {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/todos/999999", nil)
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusNotFound, w.Code)

		var resp models.APIResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 1002, resp.Code)
	})

	s.Run("重复删除返回 404", func() {
		// 先创建一个 todo
		result, err := s.db.Exec(`
			INSERT INTO todos (title, description, completed, created_at, updated_at)
			VALUES ('待删除项', '测试描述', 0, datetime('now'), datetime('now'))
		`)
		assert.NoError(s.T(), err)
		id, _ := result.LastInsertId()

		// 第一次删除
		w1 := httptest.NewRecorder()
		req1 := httptest.NewRequest("DELETE", fmt.Sprintf("/api/todos/%d", id), nil)
		s.router.ServeHTTP(w1, req1)
		assert.Equal(s.T(), http.StatusOK, w1.Code)

		// 第二次删除
		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest("DELETE", fmt.Sprintf("/api/todos/%d", id), nil)
		s.router.ServeHTTP(w2, req2)
		assert.Equal(s.T(), http.StatusNotFound, w2.Code)
	})
}

// BoundaryTestsSuite 边界测试
func (s *TestSuite) TestBoundaryTests() {
	s.Run("标题最大长度（200字符）", func() {
		maxTitle := strings.Repeat("测", 200)
		bodyData := map[string]string{"title": maxTitle}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusCreated, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 0, resp.Code)
	})

	s.Run("描述最大长度（1000字符）", func() {
		maxDesc := strings.Repeat("内", 1000)
		bodyData := map[string]string{
			"title":       "测试",
			"description": maxDesc,
		}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusCreated, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), 0, resp.Code)
	})

	s.Run("描述超长（>1000字符）返回验证失败", func() {
		longDesc := strings.Repeat("内", 1001)
		bodyData := map[string]string{
			"title":       "测试",
			"description": longDesc,
		}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	})
}

// SecurityTestsSuite 安全测试
func (s *TestSuite) TestSecurityTests() {
	s.Run("SQL 注入防护", func() {
		// 尝试 SQL 注入攻击
		maliciousTitle := "'; DROP TABLE todos; --"
		bodyData := map[string]string{"title": maliciousTitle}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		// 应该正常创建，而不是执行注入
		assert.Equal(s.T(), http.StatusCreated, w.Code)

		// 验证表仍然存在
		var count int
		err = s.db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)
		assert.NoError(s.T(), err)
		assert.Greater(s.T(), count, 0)

		// 验证插入的数据是原样存储
		var title string
		s.db.QueryRow("SELECT title FROM todos WHERE title = ?", maliciousTitle).Scan(&title)
		assert.Equal(s.T(), maliciousTitle, title)
	})

	s.Run("特殊字符处理", func() {
		specialTitle := "测试<script>alert('xss')</script>"
		bodyData := map[string]string{"title": specialTitle}
		body, err := createJSONBody(bodyData)
		assert.NoError(s.T(), err)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos", body)
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusCreated, w.Code)

		var resp models.APIResponse
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(s.T(), err)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(s.T(), ok)
		assert.Equal(s.T(), specialTitle, data["title"])
	})

	s.Run("JSON 格式验证", func() {
		// 发送无效的 JSON
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/todos",
			bytes.NewBufferString(`{invalid json}`))
		req.Header.Set("Content-Type", "application/json")
		s.router.ServeHTTP(w, req)

		assert.Equal(s.T(), http.StatusBadRequest, w.Code)
	})
}
