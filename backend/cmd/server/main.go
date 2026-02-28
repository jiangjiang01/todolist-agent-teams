package main

import (
	"log"
	"os"

	"todolist/internal/database"
	"todolist/internal/handlers"
	"todolist/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/todos.db"
	}

	if err := database.InitDB(dbPath); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	// 创建 Gin 路由
	r := gin.Default()

	// 使用 CORS 中间件
	r.Use(middleware.CORS())

	// API 路由组
	api := r.Group("/api")
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

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
