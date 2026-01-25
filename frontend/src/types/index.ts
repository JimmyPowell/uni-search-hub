// 用户相关类型
export interface User {
  id: number
  username: string
  display_name: string
  email: string
  role: number // 1: 普通用户, 10: 管理员, 100: Root
  status: number
  created_at: string
  updated_at: string
}

export interface LoginForm {
  username: string
  password: string
  turnstile_token?: string
}

export interface RegisterForm {
  username: string
  password: string
  email: string
  display_name?: string
  turnstile_token?: string
}

// Token相关类型
export interface Token {
  id: number
  user_id: number
  name: string
  key: string
  status: number
  expired_at: string
  unlimited_quota: boolean
  remain_quota: number
  used_quota: number
  created_at: string
  updated_at: string
}

export interface TokenForm {
  id?: number
  name: string
  expired_at?: string
  unlimited_quota: boolean
  remain_quota: number
}

// Option配置类型
export interface Option {
  key: string
  value: string
}

// API响应类型
export interface ApiResponse<T = any> {
  success: boolean
  message: string
  data?: T
}

// 与后端 PageInfo 结构对齐
export interface PageInfo<T> {
  page: number
  page_size: number
  total: number
  items: T[]
}

export interface PaginatedResponse<T> {
  success: boolean
  message: string
  data: T[]
  total: number
}

// 用户角色枚举
export const UserRole = {
  User: 1,
  Admin: 10,
  Root: 100,
} as const

export type UserRoleType = typeof UserRole[keyof typeof UserRole]
