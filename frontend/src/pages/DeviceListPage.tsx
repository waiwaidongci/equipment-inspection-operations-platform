import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Plus, Search } from 'lucide-react'
import PageHeader from '../components/PageHeader'
import DataTable, { type Column } from '../components/DataTable'
import StatusBadge from '../components/StatusBadge'
import Pagination from '../components/Pagination'
import Loading from '../components/Loading'
import DeviceFormModal, { type DeviceFormData } from '../components/DeviceFormModal'
import { deviceApi } from '../api/services'
import type { Device } from '../types'

export default function DeviceListPage() {
  const [data, setData] = useState<{ items: Device[]; total: number; page: number; totalPages: number }>({ items: [], total: 0, page: 1, totalPages: 1 })
  const [filters, setFilters] = useState({ query: '', category: '', status: '', owner: '', location: '' })
  const [dictionaries, setDictionaries] = useState<{ categories: string[]; locations: string[] }>({ categories: [], locations: [] })
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Device | undefined>()

  async function load() {
    setLoading(true)
    try {
      const params = new URLSearchParams()
      Object.entries(filters).forEach(([key, value]) => value && params.set(key, value))
      params.set('page', String(data.page))
      params.set('pageSize', '15')
      const result = await deviceApi.list(`?${params.toString()}`)
      setData({ items: result.items, total: result.total, page: result.page, totalPages: result.totalPages })
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    deviceApi.dictionaries().then(setDictionaries).catch(() => undefined)
  }, [])

  useEffect(() => {
    load()
  }, [filters, data.page])

  const columns = useMemo<Column<Device>[]>(() => [
    { key: 'code', title: '设备编号', render: (row) => <Link className="table-link" to={`/devices/${row.id}`}>{row.code}</Link> },
    { key: 'name', title: '名称', render: (row) => row.name },
    { key: 'category', title: '类别', render: (row) => row.category },
    { key: 'location', title: '位置', render: (row) => row.location },
    { key: 'owner', title: '负责人', render: (row) => row.owner || '-' },
    { key: 'status', title: '状态', render: (row) => <StatusBadge status={row.status} /> },
  ], [])

  async function save(form: DeviceFormData) {
    if (editing) {
      await deviceApi.update(editing.id, form)
    } else {
      await deviceApi.create(form)
    }
    setData((current) => ({ ...current, page: 1 }))
  }

  return (
    <>
      <PageHeader title="设备台账" description="建立并维护设备完整档案，停用设备不再生成巡检任务。" actions={
        <button className="primary-button" onClick={() => { setEditing(undefined); setModalOpen(true) }}><Plus size={17} />新增设备</button>
      } />
      <div className="toolbar">
        <label className="search-box"><Search size={17} /><input value={filters.query} onChange={(e) => setFilters((f) => ({ ...f, query: e.target.value }))} placeholder="搜索编号、名称、序列号" /></label>
        <select className="control compact" value={filters.category} onChange={(e) => setFilters((f) => ({ ...f, category: e.target.value }))}>
          <option value="">全部类别</option>
          {dictionaries.categories.map((item) => <option key={item}>{item}</option>)}
        </select>
        <select className="control compact" value={filters.location} onChange={(e) => setFilters((f) => ({ ...f, location: e.target.value }))}>
          <option value="">全部位置</option>
          {dictionaries.locations.map((item) => <option key={item}>{item}</option>)}
        </select>
        <select className="control compact" value={filters.status} onChange={(e) => setFilters((f) => ({ ...f, status: e.target.value }))}>
          <option value="">全部状态</option><option value="active">启用</option><option value="inactive">停用</option>
        </select>
      </div>
      {loading ? <Loading /> : <>
        <DataTable columns={columns} rows={data.items} />
        <Pagination page={data.page} totalPages={data.totalPages} total={data.total} onChange={(page) => setData((current) => ({ ...current, page }))} />
      </>}
      <DeviceFormModal open={modalOpen} device={editing} categories={dictionaries.categories} locations={dictionaries.locations} onClose={() => setModalOpen(false)} onSave={save} />
    </>
  )
}
