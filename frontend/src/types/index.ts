export interface User {
  id: number
  username: string
  displayName: string
  createdAt: string
  updatedAt: string
}

export interface Page<T> {
  items: T[]
  total: number
  page: number
  pageSize: number
  totalPages: number
}

export interface Device {
  id: number
  code: string
  name: string
  category: string
  model: string
  serialNumber: string
  manufacturer: string
  location: string
  installDate: string
  warrantyEnd: string
  owner: string
  status: 'active' | 'inactive'
  remark: string
  createdAt: string
  updatedAt: string
}

export interface PlanCheckItem {
  id: number
  planId: number
  name: string
  itemType: string
  standardRange: string
  required: boolean
  description: string
  sortOrder: number
}

export interface InspectionPlan {
  id: number
  name: string
  scopeType: 'device' | 'category'
  deviceId?: number
  category: string
  period: 'daily' | 'weekly' | 'monthly'
  startDate: string
  endDate: string
  owner: string
  status: 'enabled' | 'paused' | 'terminated'
  items?: PlanCheckItem[]
  createdAt: string
  updatedAt: string
}

export interface TaskCheckResult {
  itemId: number
  name: string
  itemType: string
  standardRange: string
  required: boolean
  description: string
  value: string
  passed: boolean
  abnormal: boolean
  remark: string
}

export interface InspectionTask {
  id: number
  planId: number
  deviceId: number
  plannedAt: string
  status: 'pending' | 'draft' | 'completed' | 'abnormal'
  actualExecutorId?: number
  executedAt?: string
  results: TaskCheckResult[]
  remark: string
  createdAt: string
  updatedAt: string
}

export interface InspectionRecord {
  id: number
  taskId: number
  deviceId: number
  planId: number
  executorId: number
  executedAt: string
  status: string
  resultSummary: string
  createdAt: string
}

export interface Anomaly {
  id: number
  deviceId: number
  taskId?: number
  discovererId: number
  discoveredAt: string
  severity: 'low' | 'medium' | 'high' | 'critical'
  description: string
  status: 'open' | 'assigned' | 'in_progress' | 'closed'
  assigneeId?: number
  progress: string
  causeAnalysis: string
  closeNote: string
  verificationResult: string
  closedAt?: string
  createdAt: string
  updatedAt: string
}

export interface Repair {
  id: number
  deviceId: number
  anomalyId?: number
  repairType: string
  content: string
  vendor: string
  startedAt: string
  endedAt: string
  cost: number
  result: string
  remark: string
  createdBy: number
  createdAt: string
  updatedAt: string
}

export interface DeviceHistory {
  records: InspectionRecord[]
  repairs: Repair[]
}

export interface LoginResult {
  token: string
  user: User
}
