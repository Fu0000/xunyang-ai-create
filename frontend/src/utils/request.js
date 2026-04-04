/**
 * TASK-17: 统一 HTTP 请求封装
 *
 * 功能：
 * 1. 自动附加 Authorization Bearer Token
 * 2. 响应拦截：401 自动清除登录态并跳至首页
 * 3. 统一处理错误结构（{code, message} 格式）
 * 4. 30s 超时保护
 */
import axios from 'axios'

// 创建 axios 实例，不在这里引入 store 和 router，避免循环依赖
const request = axios.create({
  baseURL: '/api',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器：自动附加 Token
request.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error),
)

// 响应拦截器：统一错误处理
request.interceptors.response.use(
  (response) => response.data,
  (error) => {
    if (error.response) {
      const { status, data } = error.response

      // 401: Token 过期或无效，清除登录态
      if (status === 401) {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        // 触发全局重新加载，让 App 重置状态
        // 避免循环依赖，不直接引入 router/store
        if (typeof window !== 'undefined') {
          window.dispatchEvent(new CustomEvent('auth:logout'))
        }
      }

      // 提取统一错误信息（兼容 {error: "..."} 和 {message: "..."} 两种格式）
      const message = data?.message || data?.error || `请求失败 (${status})`
      const enrichedError = new Error(message)
      enrichedError.status = status
      enrichedError.code = data?.code || status
      enrichedError.responseData = data
      return Promise.reject(enrichedError)
    }

    if (error.request) {
      // 请求已发出但无响应（网络问题）
      return Promise.reject(new Error('网络连接失败，请检查网络'))
    }

    return Promise.reject(error)
  },
)

export default request
