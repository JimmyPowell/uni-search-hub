import request from './request'
import type { ApiResponse, PageInfo, User, LoginForm, RegisterForm } from '../types'

// 用户登录
export function login(data: LoginForm) {
  return request.post<ApiResponse<User>>('/user/login', data)
}

// 用户注册
export function register(data: RegisterForm) {
  return request.post<ApiResponse<User>>('/user/register', data)
}

// 用户登出
export function logout() {
  return request.get<ApiResponse>('/user/logout')
}

// 获取当前用户信息
export function getSelf() {
  return request.get<ApiResponse<User>>('/user/self')
}

// 更新当前用户信息
export function updateSelf(data: Partial<User>) {
  return request.put<ApiResponse>('/user/self', data)
}

// 删除当前账户
export function deleteSelf() {
  return request.delete<ApiResponse>('/user/self')
}

// 生成访问令牌
export function generateAccessToken() {
  return request.get<ApiResponse<string>>('/user/token')
}

// ============ 管理员接口 ============

// 获取所有用户
export function getAllUsers() {
  return request.get<ApiResponse<PageInfo<User>>>('/user/')
}

// 搜索用户
export function searchUsers(keyword: string) {
  return request.get<ApiResponse<PageInfo<User>>>('/user/search', { params: { keyword } })
}

// 获取指定用户
export function getUser(id: number) {
  return request.get<ApiResponse<User>>(`/user/${id}`)
}

// 创建用户
export function createUser(data: Partial<User> & { password: string }) {
  return request.post<ApiResponse>('/user/', data)
}

// 更新用户
export function updateUser(data: Partial<User>) {
  return request.put<ApiResponse>('/user/', data)
}

// 删除用户
export function deleteUser(id: number) {
  return request.delete<ApiResponse>(`/user/${id}`)
}

// 初始化Root用户
export function setupRootUser() {
  return request.get<ApiResponse>('/user/root')
}
