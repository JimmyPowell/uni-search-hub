import request from './request'
import type { ApiResponse, TokenUsage } from '../types'

export function getTokenUsageById(id: number) {
  return request.get<ApiResponse<TokenUsage>>(`/usage/token/${id}`)
}

