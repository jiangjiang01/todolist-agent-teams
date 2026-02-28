import { useEffect, useState } from 'react'
import { TodoForm } from './TodoForm'
import { TodoItem } from './TodoItem'
import { todoApi } from '../services/api'
import type { Todo } from '../types'

export function TodoList() {
  const [todos, setTodos] = useState<Todo[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  // 获取所有 todos
  const fetchTodos = async () => {
    try {
      setLoading(true)
      setError(null)
      const data = await todoApi.getAll()
      setTodos(data)
    } catch (err) {
      setError(err instanceof Error ? err.message : '获取任务列表失败')
    } finally {
      setLoading(false)
    }
  }

  // 创建 todo
  const handleCreate = async (data: { title: string; description?: string }) => {
    try {
      setIsSubmitting(true)
      setError(null)
      const newTodo = await todoApi.create(data)
      setTodos((prev) => [newTodo, ...prev])
    } catch (err) {
      setError(err instanceof Error ? err.message : '创建任务失败')
    } finally {
      setIsSubmitting(false)
    }
  }

  // 更新 todo
  const handleUpdate = async (id: number, data: { title?: string; description?: string; completed?: boolean }) => {
    try {
      setError(null)
      const updatedTodo = await todoApi.update(id, data)
      setTodos((prev) =>
        prev.map((todo) => (todo.id === id ? updatedTodo : todo))
      )
    } catch (err) {
      setError(err instanceof Error ? err.message : '更新任务失败')
    }
  }

  // 删除 todo
  const handleDelete = async (id: number) => {
    try {
      setError(null)
      await todoApi.delete(id)
      setTodos((prev) => prev.filter((todo) => todo.id !== id))
    } catch (err) {
      setError(err instanceof Error ? err.message : '删除任务失败')
    }
  }

  // 初始加载
  useEffect(() => {
    fetchTodos()
  }, [])

  // 统计信息
  const stats = {
    total: todos.length,
    completed: todos.filter((t) => t.completed).length,
    active: todos.filter((t) => !t.completed).length
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100 py-8 px-4">
      <div className="max-w-2xl mx-auto">
        {/* 标题 */}
        <header className="text-center mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">待办事项</h1>
          <p className="text-gray-600">管理你的日常任务</p>
        </header>

        {/* 统计信息 */}
        {todos.length > 0 && (
          <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4 mb-6">
            <div className="flex justify-around text-center">
              <div>
                <div className="text-2xl font-bold text-gray-900">{stats.total}</div>
                <div className="text-sm text-gray-500">总计</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-blue-500">{stats.active}</div>
                <div className="text-sm text-gray-500">进行中</div>
              </div>
              <div>
                <div className="text-2xl font-bold text-green-500">{stats.completed}</div>
                <div className="text-sm text-gray-500">已完成</div>
              </div>
            </div>
          </div>
        )}

        {/* 添加表单 */}
        <TodoForm onSubmit={handleCreate} isLoading={isSubmitting} />

        {/* 错误提示 */}
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-4 flex items-start gap-2">
            <svg className="w-5 h-5 flex-shrink-0 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
            </svg>
            <span className="text-sm">{error}</span>
            <button
              onClick={() => setError(null)}
              className="ml-auto text-red-500 hover:text-red-700"
            >
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
              </svg>
            </button>
          </div>
        )}

        {/* 加载状态 */}
        {loading && (
          <div className="flex items-center justify-center py-12">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
            <span className="ml-3 text-gray-600">加载中...</span>
          </div>
        )}

        {/* Todo 列表 */}
        {!loading && todos.length === 0 ? (
          <div className="text-center py-12 bg-white rounded-lg border border-gray-200">
            <svg className="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <p className="text-gray-500">暂无任务，开始添加你的第一个待办事项吧！</p>
          </div>
        ) : (
          <div className="space-y-3">
            {todos.map((todo) => (
              <TodoItem
                key={todo.id}
                todo={todo}
                onDelete={handleDelete}
                onUpdate={handleUpdate}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
