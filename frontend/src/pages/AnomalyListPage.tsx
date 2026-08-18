import { useEffect, useMemo, useState } from 'react'
import { ClipboardCheck, UserCheck, Workflow } from 'lucide-react'
import PageHeader from '../components/PageHeader'
import DataTable, { type Column } from '../components/DataTable'
import StatusBadge from '../components/StatusBadge'
import Pagination from '../components/Pagination'
import Loading from '../components/Loading'
import AnomalyActionModal from '../components/AnomalyActionModal'
import { anomalyApi, userApi } from '../api/services'
import { formatDate } from '../utils/format'
import type { Anomaly, User } from '../types'

export default function AnomalyListPage() {
  const [data, setData] = useState<{ items: Anomaly[]; total: number; page: number; totalPages: number }>({ items: [], total: 0, page: 1, totalPages: 1 })
  const [users, setUsers] = useState<User[]>([])
  const [status, setStatus] = useState('')
  const [severity, setSeverity] = useState('')
  const [loading, setLoading] = useState(true)
  const [modal, setModal] = useState<{ anomaly: Anomaly; action: 'assign' | 'progress' | 'close' } | null>(null)

  async function load() {
    setLoading(true)
    try {
      const params = new URLSearchParams({ page: String(data.page), pageSize: '15' })
      if (status) params.set('status', status)
      if (severity) params.set('severity', severity)
      const result = await anomalyApi.list(`?${params.toString()}`)
      setData({ items: result.items, total: result.total, page: result.page, totalPages: result.totalPages })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    userApi.list().then(setUsers).catch(() => undefined)
  }, [])

  useEffect(() => { load() }, [status, severity, data.page])

  const columns = useMemo<Column<Anomaly>[]>(() => [
    { key: 'device', title: '设备', render: (row) => `设备 #${row.deviceId}` },
    { key: 'description', title: '现象描述', render: (row) => row.description },
    { key: 'severity', title: '严重程度', render: (row) => <StatusBadge status={row.severity} /> },
    { key: 'discoverer', title: '发现人', render: (row) => `用户 #${row.discovererId}` },
    { key: 'discoveredAt', title: '发现时间', render: (row) => formatDate(row.discoveredAt) },
    { key: 'status', title: '状态', render: (row) => <StatusBadge status={row.status} /> },
    { key: 'actions', title: '操作', render: (row) => <div className="row-actions">
      {row.status === 'open' && <button className="secondary-button small" onClick={() => setModal({ anomaly: row, action: 'assign' })}><UserCheck size={14} />指派</button>}
      {row.status !== 'closed' && <button className="secondary-button small" onClick={() => setModal({ anomaly: row, action: 'progress' })}><Workflow size={14} />进展</button>}
      {row.status !== 'closed' && <button className="danger-button small" onClick={() => setModal({ anomaly: row, action: 'close' })}><ClipboardCheck size={14} />关闭</button>}
      {row.status === 'closed' && <span className="muted">已闭环</span>}
    </div> },
  ], [])

  async function save(body: unknown) {
    if (!modal) return
    if (modal.action === 'assign') await anomalyApi.assign(modal.anomaly.id, (body as { assigneeId: number }).assigneeId)
    if (modal.action === 'progress') await anomalyApi.progress(modal.anomaly.id, (body as { progress: string }).progress)
    if (modal.action === 'close') await anomalyApi.close(modal.anomaly.id, body)
    load()
  }

  return (
    <>
      <PageHeader title="异常管理" description="巡检异常自动进入列表，支持指派、进展更新、原因分析和关闭验证。" />
      <div className="toolbar">
        <select className="control compact" value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="">全部状态</option><option value="open">待处理</option><option value="assigned">已指派</option><option value="in_progress">处理中</option><option value="closed">已关闭</option>
        </select>
        <select className="control compact" value={severity} onChange={(e) => setSeverity(e.target.value)}>
          <option value="">全部严重程度</option><option value="low">低</option><option value="medium">中</option><option value="high">高</option><option value="critical">紧急</option>
        </select>
      </div>
      {loading ? <Loading /> : <>
        <DataTable columns={columns} rows={data.items} />
        <Pagination page={data.page} totalPages={data.totalPages} total={data.total} onChange={(page) => setData((current) => ({ ...current, page }))} />
      </>}
      {modal && <AnomalyActionModal open={!!modal} anomaly={modal.anomaly} users={users} action={modal.action} onClose={() => setModal(null)} onSave={save} />}
    </>
  )
}
