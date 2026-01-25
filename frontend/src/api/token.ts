import request from './request'
import type { ApiResponse, PageInfo, Token, TokenForm } from '../types'

// 获取所有Token
export function getAllTokens() {
  return request.get<ApiResponse<PageInfo<Token>>>('/token/')
}

// 搜索Token
export function searchTokens(keyword: string) {
  return request.get<ApiResponse<Token[]>>('/token/search', { params: { keyword } })
}

// 获取单个Token
export function getToken(id: number) {
  return request.get<ApiResponse<Token>>(`/token/${id}`)
}

// 添加Token
export function addToken(data: TokenForm) {
  return request.post<ApiResponse<Token>>('/token/', data)
}

// 更新Token
export function updateToken(data: TokenForm) {
  return request.put<ApiResponse>('/token/', data)
}

// 删除Token
export function deleteToken(id: number) {
  return request.delete<ApiResponse>(`/token/${id}`)
}

// 批量删除Token
export function deleteTokenBatch(ids: number[]) {
  return request.post<ApiResponse>('/token/batch', { ids })
}

// 获取Token用量
export function getTokenUsage() {
  return request.get<ApiResponse>('/usage/token/')
}
