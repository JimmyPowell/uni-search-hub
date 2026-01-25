import request from './request'
import type { ApiResponse, Option } from '../types'

// 获取系统配置
export function getOptions() {
  return request.get<ApiResponse<Option[]>>('/option/')
}

// 更新配置
export function updateOption(data: Option) {
  return request.put<ApiResponse>('/option/', data)
}
