import request from './request'
import type { ApiResponse } from '../types'

export interface SetupStatus {
  rootExists: boolean
  setupRequired: boolean
  setupCompleted: boolean
}

export interface SetupRootRequest {
  username: string
  password: string
  confirmPassword: string
}

export function getSetupStatus() {
  return request.get<ApiResponse<SetupStatus>>('/setup/status')
}

export function setupRootUser(data: SetupRootRequest) {
  return request.post<ApiResponse>('/setup/root', data)
}

