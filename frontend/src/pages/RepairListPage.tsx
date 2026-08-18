import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus } from 'lucide-react'
import PageHeader from '../components/PageHeader'
import DataTable, { type Column } from '../components/DataTable'
import Pagination from '../components/Pagination'
import Loading from '../components/Loading'
import RepairFormModal, { type RepairFormData } from '../components/RepairFormModal'
import { anomalyApi, deviceApi, repairApi } from '../api/services'
import { formatDateOnly, money } from '../utils/format'
import type { Anomaly, Device, Repair } from '../types'

export default function RepairListPage() {
  const [data, setData] = useState<{ items: Repair[]; total: number; page: number; totalPages: number }>({ items: [], total: 0, page: 1, totalPages: 1 })
  const [devices, setDevices] = useState<Device[]>([])
  const [closedAnomalies, setClosedAnomalies] = useState<Anomaly[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Repair | undefined>()

  async function load() {
    setLoading(true)
    try {
      const result = await repairApi.list(`?page=${data.page}&pageSize=15`)
      setData({ items: result.items, total: result.total, page: result.page, totalPages: result.totalPages })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    Promise.all([deviceApi.list('?page=1&pageSize=1000'), anomalyApi.list('?status=closed&page=1&pageSize=1000')])
      .then(([devicePage, anomalyPage]) => {
        setDevices(devicePage.items)
        setClosedAnomalies(anomalyPage.items)
      })
      .catch(() => undefined)
  }, [])

  useEffect(() => { load() }, [data.page])

  const columns = useMemo<Column<Repair>[]>(() => [
    { key: 'device', title: '设备', render: (row) => <Link className="table-link" to={`/devices/${row.deviceId}`}>设备 #{row.deviceId}</Link> },
    { key: 'anomaly', title: '关联异常', render: (row) => row.anomalyId ? `异常 #${row.anomalyId}` : '-' },
    { key: 'type', title: '类型', render: (row) => row.repairType },
    { key: 'content', title: '维修内容', render: (row) => row.content },
    { key: 'vendor', title: '维修单位', render: (row) => row.vendor },
    { key: 'period', title: '起止日期', render: (row) => `${formatDateOnly(row.startedAt)} 至 ${formatDateOnly(row.endedAt)}` },
    { key: 'cost', title: '费用', render: (row) => money(row.cost) },
    { key: 'result', title: '结果', render: (row) => row.result },
    { key: 'actions', title: '操作', render: (row) => <button className="secondary-button small" onClick={() => { setEditing(row); setModalOpen(true) }}>编辑</button> },
  ], [])

  async function save(form: RepairFormData) {
    const body = { ...form, anomalyId: form.anomalyId || undefined }
    if (editing) await repairApi.update(editing.id, body)
    else await repairApi.create(body)
    setData((current) => ({ ...current, page: 1 }))
  }

  return (
    <>
      <PageHeader title="维修记录" description="记录设备维修过程与费用，已关闭异常可关联维修，形成完整追溯链。" actions={
        <button className="primary-button" onClick={() => { setEditing(undefined); setModalOpen(true) }}><Plus size={17} />新增维修</button>
      } />
      {loading ? <Loading /> : <>
        <DataTable columns={columns} rows={data.items} />
        <Pagination page={data.page} totalPages={data.totalPages} total={data.total} onChange={(page) => setData((current) => ({ ...current, page }))} />
      </>}
      <RepairFormModal open={modalOpen} repair={editing} devices={devices} closedAnomalies={closedAnomalies} onClose={() => setModalOpen(false)} onSave={save} />
    </>
  )
}
