import axios from 'axios'

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

export interface StatusResponse {
  ip: string
  verified: boolean
  has_record: boolean
}

export interface AuthRequest {
  password: string
  remoteIP: string
}

export interface CredentialStat {
  username: string
  password: string
  count: number
}

export interface IPStat {
  ip: string
  count: number
}

export interface VerifiedIP {
  ip: string
  method: string
  created_at: string
}

export interface Stats {
  total_ssh_attempts: number
  total_rejected: number
  total_honeypot: number
  total_forwarded: number
  total_auth_attempts: number
  total_auth_success: number
  unique_honeypot_creds: number
  unique_rejected_ips: number
}

export interface PageResponse<T> {
  data: T[]
  total: number
  page: number
}

export const getStatus = () =>
  http.get<StatusResponse>('/status')

export const postAuth = (data: AuthRequest) =>
  http.post('/auth', data)

export const postAuthWithHeader = (token: string) =>
  http.post('/auth', {}, { headers: { Authorization: token } })

export const getHoneypotCredentials = (page = 1, pageSize = 20) =>
  http.get<PageResponse<CredentialStat>>('/admin/honeypot', { params: { page, page_size: pageSize } })

export const getRejectedIPs = (page = 1, pageSize = 20) =>
  http.get<PageResponse<IPStat>>('/admin/rejected', { params: { page, page_size: pageSize } })

export const getVerifiedIPs = (page = 1, pageSize = 20) =>
  http.get<PageResponse<VerifiedIP>>('/admin/verified', { params: { page, page_size: pageSize } })

export const getStats = () =>
  http.get<Stats>('/admin/stats')
