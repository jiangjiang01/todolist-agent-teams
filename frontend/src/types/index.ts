// Todo 项数据类型
export interface Todo {
  id: number
  title: string
  description: string | null
  completed: boolean
  created_at: string
  updated_at: string
}

// 创建 Todo 请求类型
export interface CreateTodoRequest {
  title: string
  description?: string
}

// 更新 Todo 请求类型
export interface UpdateTodoRequest {
  title?: string
  description?: string
  completed?: boolean
}

// API 响应类型
export interface ApiResponse<T> {
  code: number
  message: string
  data: T | null
}

// 加载状态
export type LoadingState = 'idle' | 'loading' | 'success' | 'error'
