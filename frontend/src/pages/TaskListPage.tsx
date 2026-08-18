import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Play, Search } from 'lucide-react'
import PageHeader from '../components/PageHeader'
import DataTable, { type Column } from '../components/DataTable'
import StatusBadge from '../components/StatusBadge'
import Pagination from '../components/Pagination'
import Loading from '../components/Loading'
import { taskApi } from '../api/services'
import { formatDate } from '../utils/format'
import type { InspectionTask } from '../types'

export default function TaskListPage() {
  const [data, setData] = useState<{ items: InspectionTask[]; total: number; page: number; totalPages: number }>({ items: [], total: 0, page: 1, totalPages: 1 })
  const [status, setStatus] = useState('')
  const [deviceId, setDeviceId] = useState('')
  const [from, setFrom] = useState('')
  const [to, setTo] = useState('')
  const [loading, setLoading] = useState(true)

  async function load() {
    setLoading(true)
    try {
      const params = new URLSearchParams({ page: String(data.page), pageSize: '15' })
      if (status) params.set('status', status)
      if (deviceId) params.set('deviceId', deviceId)
      if (from) params.set('from', from)
      if (to) params.set('to', to)
      const result = await taskApi.list(`?${params.toString()}`)
      setData({ items: result.items, total: result.total, page: result.page, totalPages: result.totalPages })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { load() }, [status, deviceId, from, to, data.page])

  const columns = useMemo<Column<InspectionTask>[]>(() => [
    { key: 'plannedAt', title: '计划执行时间', render: (row) => <span className={isOverdue(row) ? 'overdue-text' : ''}>{formatDate(row.plannedAt)}{isOverdue(row) && <em> 已过期</em>}</span> },
    { key: 'deviceId', title: '设备', render: (row) => `设备 #${row.deviceId}` },
    { key: 'planId', title: '计划', render: (row) => `计划 #${row.planId}` },
    { key: 'executor', title: '实际执行人', render: (row) => row.actualExecutorId ? `用户 #${row.actualExecutorId}` : '-' },
    { key: 'status', title: '状态', render: (row) => <StatusBadge status={row.status} /> },
    { key: 'actions', title: '操作', render: (row) => (row.status === 'pending' || row.status === 'draft')
      ? <Link className="primary-button small" to={`/tasks/${row.id}/execute`}><Play size={14} />执行</Link>
      : <span className="muted">已完成</span> },
  ], [])

  function isOverdue(task: InspectionTask) {
    return (task.status === 'pending' || task.status === 'draft') && new Date(task.plannedAt).getTime() < Date.now()
  }

  return (
    <>
      <PageHeader title="巡检任务" description="查看计划生成的待办与已完成任务，过期未执行任务会高亮提示。" />
      <div className="toolbar">
        <label className="search-box"><Search size={17} /><input value={deviceId} onChange={(e) => setDeviceId(e.target.value)} placeholder="设备 ID" /></label>
        <select className="control compact" value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="">全部状态</option><option value="pending">待执行</option><option value="draft">草稿</option><option value="completed">已完成</option><option value="abnormal">异常</option>
        </select>
        <input className="control compact" type="date" value={from} onChange={(e) => setFrom(e.target.value)} />
        <span className="muted">至</span>
        <input className="control compact" type="date" value={to} onChange={(e) => setTo(e.target.value)} />
      </div>
      {loading ? <Loading /> : <>
        <DataTable columns={columns} rows={data.items} />
        <Pagination page={data.page} totalPages={data.totalPages} total={data.total} onChange={(page) => setData((current) => ({ ...current, page }))} />
      </>}
    </>
  )
}
