import axios, { type AxiosInstance, type AxiosResponse, type InternalAxiosRequestConfig } from 'axios'
import type { ApiResponse } from '../types'

const request: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 30000,
  withCredentials: true,
})

// 请求拦截器
request.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>) => {
    return response
  },
  (error) => {
    if (error.response?.status === 401) {
      // 未授权，清除状态并跳转登录（排除登录相关接口）
      const url = error.config?.url || ''
      if (!url.includes('/user/login') && !url.includes('/user/self')) {
        window.location.href = '/login'
      }
    }
    return Promise.reject(error)
  }
)

export default request
