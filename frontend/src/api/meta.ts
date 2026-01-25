import request from './request'
import type { ApiResponse, ServiceCatalog } from '../types'

export function getServiceCatalog() {
  return request.get<ApiResponse<ServiceCatalog>>('/meta/services')
}

