import { useEffect, useMemo, useState } from 'react'
import { Pause, Play, Plus, RefreshCw, Square } from 'lucide-react'
import PageHeader from '../components/PageHeader'
import DataTable, { type Column } from '../components/DataTable'
import StatusBadge from '../components/StatusBadge'
import Pagination from '../components/Pagination'
import Loading from '../components/Loading'
import ConfirmButton from '../components/ConfirmButton'
import PlanFormModal, { type PlanFormData } from '../components/PlanFormModal'
import { deviceApi, planApi } from '../api/services'
import { formatDateOnly } from '../utils/format'
import type { Device, InspectionPlan } from '../types'

export default function PlanListPage() {
  const [data, setData] = useState<{ items: InspectionPlan[]; total: number; page: number; totalPages: number }>({ items: [], total: 0, page: 1, totalPages: 1 })
  const [devices, setDevices] = useState<Device[]>([])
  const [categories, setCategories] = useState<string[]>([])
  const [status, setStatus] = useState('')
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<InspectionPlan | undefined>()
  const [message, setMessage] = useState('')

  async function load() {
    setLoading(true)
    try {
      const params = new URLSearchParams({ page: String(data.page), pageSize: '15' })
      if (status) params.set('status', status)
      const result = await planApi.list(`?${params.toString()}`)
      setData({ items: result.items, total: result.total, page: result.page, totalPages: result.totalPages })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    Promise.all([deviceApi.list('?page=1&pageSize=1000'), deviceApi.dictionaries()])
      .then(([devicePage, dict]) => {
        setDevices(devicePage.items.filter((item) => item.status === 'active'))
        setCategories(dict.categories)
      })
      .catch(() => undefined)
  }, [])

  useEffect(() => { load() }, [status, data.page])

  const columns = useMemo<Column<InspectionPlan>[]>(() => [
    { key: 'name', title: '计划名称', render: (row) => row.name },
    { key: 'scope', title: '适用范围', render: (row) => row.scopeType === 'device' ? `设备 #${row.deviceId}` : `类别：${row.category}` },
    { key: 'period', title: '周期', render: (row) => ({ daily: '每日', weekly: '每周', monthly: '每月' })[row.period] },
    { key: 'range', title: '起止日期', render: (row) => `${formatDateOnly(row.startDate)} 至 ${formatDateOnly(row.endDate)}` },
    { key: 'owner', title: '负责人', render: (row) => row.owner },
    { key: 'status', title: '状态', render: (row) => <StatusBadge status={row.status} /> },
    { key: 'actions', title: '操作', render: (row) => <div className="row-actions">
      <button className="secondary-button small" onClick={() => { setEditing(row); setModalOpen(true) }}>编辑</button>
      {row.status !== 'enabled' && <ConfirmButton className="secondary-button small" message="确认启用？" onConfirm={() => changeStatus(row.id, 'enabled')}><Play size={14} />启用</ConfirmButton>}
      {row.status === 'enabled' && <ConfirmButton className="secondary-button small" message="确认暂停？" onConfirm={() => changeStatus(row.id, 'paused')}><Pause size={14} />暂停</ConfirmButton>}
      {row.status !== 'terminated' && <ConfirmButton className="danger-button small" message="终止后不可恢复，确认？" onConfirm={() => changeStatus(row.id, 'terminated')}><Square size={14} />终止</ConfirmButton>}
    </div> },
  ], [])

  async function save(form: PlanFormData) {
    const body = { ...form, deviceId: form.scopeType === 'device' ? form.deviceId : undefined }
    if (editing) await planApi.update(editing.id, body)
    else await planApi.create(body)
    setData((current) => ({ ...current, page: 1 }))
  }

  async function changeStatus(id: number, nextStatus: string) {
    await planApi.setStatus(id, nextStatus)
    setMessage('计划状态已更新')
    load()
  }

  async function generate() {
    const result = await planApi.generate()
    setMessage(`已生成 ${result.created} 个新任务`)
    load()
  }

  return (
    <>
      <PageHeader title="巡检计划" description="按设备或类别配置周期计划，系统自动生成并去重未来任务。" actions={
        <div className="page-actions">
          <button className="secondary-button" onClick={generate}><RefreshCw size={17} />立即生成</button>
          <button className="primary-button" onClick={() => { setEditing(undefined); setModalOpen(true) }}><Plus size={17} />新增计划</button>
        </div>
      } />
      <div className="toolbar">
        <select className="control compact" value={status} onChange={(e) => setStatus(e.target.value)}>
          <option value="">全部状态</option><option value="enabled">启用</option><option value="paused">暂停</option><option value="terminated">终止</option>
        </select>
        {message && <span className="success-banner">{message}</span>}
      </div>
      {loading ? <Loading /> : <>
        <DataTable columns={columns} rows={data.items} />
        <Pagination page={data.page} totalPages={data.totalPages} total={data.total} onChange={(page) => setData((current) => ({ ...current, page }))} />
      </>}
      <PlanFormModal open={modalOpen} plan={editing} devices={devices} categories={categories} onClose={() => setModalOpen(false)} onSave={save} />
    </>
  )
}
