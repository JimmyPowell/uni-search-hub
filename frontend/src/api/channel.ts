import request from './request'
import type { ApiResponse, Channel, PageInfo } from '../types'

export function getChannels(params?: { p?: number; page_size?: number; size?: number }) {
  return request.get<ApiResponse<PageInfo<Channel>>>('/channel/', { params })
}

export function createChannel(data: Pick<Channel, 'name' | 'provider' | 'api_key'> & { enabled?: boolean; weight?: number }) {
  return request.post<ApiResponse<Channel>>('/channel/', data)
}

export function updateChannel(data: Partial<Channel> & { id: number }) {
  return request.put<ApiResponse<Channel>>('/channel/', data)
}

export function deleteChannel(id: number) {
  return request.delete<ApiResponse>(`/channel/${id}`)
}

export function testChannel(id: number, params?: { q?: string }) {
  return request.post<ApiResponse>(`/channel/${id}/test`, null, { params })
}
