import request from './request'
import type { ApiResponse, PageInfo, Redemption, RedemptionForm } from '../types'

// 获取所有兑换码
export function getAllRedemptions() {
  return request.get<ApiResponse<PageInfo<Redemption>>>('/redemption/')
}

// 搜索兑换码
export function searchRedemptions(keyword: string) {
  return request.get<ApiResponse<PageInfo<Redemption>>>('/redemption/search', { params: { keyword } })
}

// 获取单个兑换码
export function getRedemption(id: number) {
  return request.get<ApiResponse<Redemption>>(`/redemption/${id}`)
}

// 添加兑换码
export function addRedemption(data: RedemptionForm) {
  return request.post<ApiResponse<string[]>>('/redemption/', data)
}

// 更新兑换码
export function updateRedemption(data: RedemptionForm, statusOnly = false) {
  return request.put<ApiResponse<Redemption>>('/redemption/', data, { params: { status_only: statusOnly ? '1' : '' } })
}

// 删除兑换码
export function deleteRedemption(id: number) {
  return request.delete<ApiResponse>(`/redemption/${id}`)
}

// 删除无效兑换码
export function deleteInvalidRedemptions() {
  return request.delete<ApiResponse<number>>('/redemption/invalid')
}