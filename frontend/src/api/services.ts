import { api } from './client'
import type {
  Anomaly, Device, DeviceHistory, InspectionPlan, InspectionRecord, InspectionTask,
  LoginResult, Page, Repair, User,
} from '../types'

export const authApi = {
  login: (username: string, password: string) => api.post<LoginResult>('/api/auth/login', { username, password }),
  me: () => api.get<User>('/api/auth/me'),
}

export const userApi = {
  list: () => api.get<User[]>('/api/users'),
}

export const deviceApi = {
  list: (params = '') => api.get<Page<Device>>(`/api/devices${params}`),
  get: (id: number) => api.get<Device>(`/api/devices/${id}`),
  create: (body: unknown) => api.post<Device>('/api/devices', body),
  update: (id: number, body: unknown) => api.put<Device>(`/api/devices/${id}`, body),
  history: (id: number) => api.get<DeviceHistory>(`/api/devices/${id}/history`),
  dictionaries: () => api.get<{ categories: string[]; locations: string[]; statuses: string[] }>('/api/dictionaries'),
}

export const planApi = {
  list: (params = '') => api.get<Page<InspectionPlan>>(`/api/plans${params}`),
  get: (id: number) => api.get<InspectionPlan>(`/api/plans/${id}`),
  create: (body: unknown) => api.post<InspectionPlan>('/api/plans', body),
  update: (id: number, body: unknown) => api.put<InspectionPlan>(`/api/plans/${id}`, body),
  setStatus: (id: number, status: string) => api.post<InspectionPlan>(`/api/plans/${id}/status`, { status }),
  generate: () => api.post<{ created: number }>('/api/plans/generate'),
}

export const taskApi = {
  list: (params = '') => api.get<Page<InspectionTask>>(`/api/tasks${params}`),
  get: (id: number) => api.get<InspectionTask>(`/api/tasks/${id}`),
  submit: (id: number, body: unknown) => api.post<InspectionTask>(`/api/tasks/${id}/submit`, body),
  draft: (id: number, body: unknown) => api.post<InspectionTask>(`/api/tasks/${id}/draft`, body),
}

export const recordApi = {
  list: (deviceId: number) => api.get<InspectionRecord[]>(`/api/records?deviceId=${deviceId}`),
}

export const anomalyApi = {
  list: (params = '') => api.get<Page<Anomaly>>(`/api/anomalies${params}`),
  get: (id: number) => api.get<Anomaly>(`/api/anomalies/${id}`),
  assign: (id: number, assigneeId: number) => api.post<Anomaly>(`/api/anomalies/${id}/assign`, { assigneeId }),
  progress: (id: number, progress: string) => api.post<Anomaly>(`/api/anomalies/${id}/progress`, { progress }),
  close: (id: number, body: unknown) => api.post<Anomaly>(`/api/anomalies/${id}/close`, body),
}

export const repairApi = {
  list: (params = '') => api.get<Page<Repair>>(`/api/repairs${params}`),
  get: (id: number) => api.get<Repair>(`/api/repairs/${id}`),
  create: (body: unknown) => api.post<Repair>('/api/repairs', body),
  update: (id: number, body: unknown) => api.put<Repair>(`/api/repairs/${id}`, body),
}
