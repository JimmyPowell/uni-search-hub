import request from './request'
import type { ApiResponse, PageInfo, RequestLog } from '../types'

export interface RequestLogQuery {
  p?: number
  page_size?: number

  provider?: string
  endpoint?: string
  status_code?: number
  user_id?: number
  token_id?: number
  channel_id?: number
  request_id?: string

  // Unix ms
  start_time?: number
  end_time?: number
}

export function getRequestLogs(params: RequestLogQuery) {
  return request.get<ApiResponse<PageInfo<RequestLog>>>('/request_log/', { params })
}

