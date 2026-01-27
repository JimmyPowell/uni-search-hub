// 用户相关类型
export interface User {
  id: number
  username: string
  display_name: string
  email: string
  role: number // 1: 普通用户, 10: 管理员, 100: Root
  status: number
  quota: number // 剩余配额
  used_quota: number // 已使用配额
  request_count: number // 请求次数
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
  group?: string
  created_at: string
  updated_at: string
}

export interface TokenForm {
  id?: number
  name: string
  expired_at?: string
  unlimited_quota: boolean
  remain_quota: number
  group?: string
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

export interface Channel {
  id: number
  name: string
  provider: string
  api_key?: string
  enabled: boolean
  weight: number
  fail_count?: number
  last_used_at?: number
  created_at?: string
  updated_at?: string
}

export interface RequestLog {
  id: number
  request_id: string
  user_id: number
  token_id: number
  channel_id: number
  provider: string
  action: string
  meta?: string
  endpoint: string
  status_code: number
  latency_ms: number
  cost: number
  error?: string
  created_at: string
}

export interface TokenUsage {
  object: 'token_usage'
  name: string
  total_granted: number
  total_used: number
  total_available: number
  unlimited_quota: boolean
  model_limits: Record<string, boolean>
  model_limits_enabled: boolean
  expires_at: number
}

export interface ServiceEndpoint {
  method: string
  path: string
  auth: 'token' | 'session'
  notes?: string
  example_headers?: Record<string, string>
  example_body?: any
}

export interface ServiceItem {
  id: string
  name: string
  provider: string
  enabled_channels_count: number
  endpoints: ServiceEndpoint[]
}

export interface ServiceCatalog {
  services: ServiceItem[]
}

// 兑换码相关类型
export interface Redemption {
  id: number
  user_id: number
  key: string
  status: number // 1: 未使用, 2: 已使用, 3: 已禁用
  name: string
  quota: number
  created_time: number
  redeemed_time: number
  used_user_id: number
  expired_time: number
}

export interface RedemptionForm {
  id?: number
  name: string
  quota: number
  count: number
  expired_time?: number
  status?: number
}
