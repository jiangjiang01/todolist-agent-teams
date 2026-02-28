package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"todolist/internal/database"
	"todolist/internal/models"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess              = 0
	CodeValidationError      = 1001
	CodeNotFound            = 1002
	CodeDatabaseError       = 1003
	CodeInternalError       = 1004
)

// GetTodos 获取所有 Todo
func GetTodos(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    CodeDatabaseError,
			Message: "数据库查询失败",
			Data:    nil,
		})
		return
	}
	defer rows.Close()

	todos := []models.Todo{}
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Code:    CodeDatabaseError,
				Message: "数据解析失败",
				Data:    nil,
			})
			return
		}
		todos = append(todos, todo)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    CodeSuccess,
		Message: "success",
		Data:    todos,
	})
}

// GetTodo 获取单个 Todo
func GetTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    CodeValidationError,
			Message: "无效的 ID",
			Data:    nil,
		})
		return
	}

	var todo models.Todo
	err = database.DB.QueryRow(`
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = ?
	`, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    CodeNotFound,
			Message: "Todo 不存在",
			Data:    nil,
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    CodeDatabaseError,
			Message: "数据库查询失败",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    CodeSuccess,
		Message: "success",
		Data:    todo,
	})
}

// CreateTodo 创建新的 Todo
func CreateTodo(c *gin.Context) {
	var req models.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    CodeValidationError,
			Message: "参数验证失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	now := time.Now()
	result, err := database.DB.Exec(`
		INSERT INTO todos (title, description, completed, created_at, updated_at)
		VALUES (?, ?, 0, ?, ?)
	`, req.Title, req.Description, now, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    CodeDatabaseError,
			Message: "创建失败",
			Data:    nil,
		})
		return
	}

	id, _ := result.LastInsertId()

	var todo models.Todo
	err = database.DB.QueryRow(`
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = ?
	`, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    CodeDatabaseError,
			Message: "查询新建数据失败",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusCreated, models.APIResponse{
		Code:    CodeSuccess,
		Message: "created",
		Data:    todo,
	})
}

// UpdateTodo 更新 Todo
func UpdateTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    CodeValidationError,
			Message: "无效的 ID",
			Data:    nil,
		})
		return
	}

	var req models.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    CodeValidationError,
			Message: "参数验证失败: " + err.Error(),
			Data:    nil,
		})
		return
	}

	// 检查 Todo 是否存在
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM todos WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    CodeNotFound,
			Message: "Todo 不存在",
			Data:    nil,
		})
		return
	}

	// 构建更新语句
	updates := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.Title != nil {
		updates = append(updates, "title = ?")
		args = append(args, *req.Title)
		argIndex++
	}
	if req.Description != nil {
		updates = append(updates, "description = ?")
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Completed != nil {
		updates = append(updates, "completed = ?")
		completedInt := 0
		if *req.Completed {
			completedInt = 1
		}
		args = append(args, completedInt)
		argIndex++
	}

	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now())
	argIndex++

	args = append(args, id)

	query := "UPDATE todos SET " + string([]rune(joinStrings(updates, ", "))) + " WHERE id = ?"
	_, err = database.DB.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    CodeDatabaseError,
			Message: "更新失败",
			Data:    nil,
		})
		return
	}

	// 获取更新后的数据
	var todo models.Todo
	err = database.DB.QueryRow(`
		SELECT id, title, description, completed, created_at, updated_at
		FROM todos
		WHERE id = ?
	`, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    CodeDatabaseError,
			Message: "查询更新数据失败",
			Data:    nil,
		})
		return
	}
	c.JSON(http.StatusOK, models.APIResponse{
		Code:    CodeSuccess,
		Message: "updated",
		Data:    todo,
	})
}

// DeleteTodo 删除 Todo
func DeleteTodo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Code:    CodeValidationError,
			Message: "无效的 ID",
			Data:    nil,
		})
		return
	}

	// 检查 Todo 是否存在
	var exists bool
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM todos WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Code:    CodeNotFound,
			Message: "Todo 不存在",
			Data:    nil,
		})
		return
	}

	_, err = database.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Code:    CodeDatabaseError,
			Message: "删除失败",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Code:    CodeSuccess,
		Message: "deleted",
		Data:    nil,
	})
}

// joinStrings 连接字符串数组
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
