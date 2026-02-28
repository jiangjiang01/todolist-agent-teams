import axios from 'axios'
import type {
  Todo,
  CreateTodoRequest,
  UpdateTodoRequest,
  ApiResponse
} from '../types'

// 创建 axios 实例
const apiClient = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json'
  },
  timeout: 10000
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const message = error.response?.data?.message || error.message || '请求失败'
    return Promise.reject(new Error(message))
  }
)

// API 方法
export const todoApi = {
  // 获取所有 todos
  getAll: async (): Promise<Todo[]> => {
    const response = await apiClient.get<ApiResponse<Todo[]>>('/todos')
    if (response.data.code === 0 && response.data.data) {
      return response.data.data
    }
    throw new Error(response.data.message || '获取列表失败')
  },

  // 获取单个 todo
  getById: async (id: number): Promise<Todo> => {
    const response = await apiClient.get<ApiResponse<Todo>>(`/todos/${id}`)
    if (response.data.code === 0 && response.data.data) {
      return response.data.data
    }
    throw new Error(response.data.message || '获取详情失败')
  },

  // 创建 todo
  create: async (data: CreateTodoRequest): Promise<Todo> => {
    const response = await apiClient.post<ApiResponse<Todo>>('/todos', data)
    if (response.data.code === 0 && response.data.data) {
      return response.data.data
    }
    throw new Error(response.data.message || '创建失败')
  },

  // 更新 todo
  update: async (id: number, data: UpdateTodoRequest): Promise<Todo> => {
    const response = await apiClient.put<ApiResponse<Todo>>(`/todos/${id}`, data)
    if (response.data.code === 0 && response.data.data) {
      return response.data.data
    }
    throw new Error(response.data.message || '更新失败')
  },

  // 删除 todo
  delete: async (id: number): Promise<void> => {
    const response = await apiClient.delete<ApiResponse<null>>(`/todos/${id}`)
    if (response.data.code !== 0) {
      throw new Error(response.data.message || '删除失败')
    }
  }
}

export default apiClient
