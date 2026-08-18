import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ArrowLeft } from 'lucide-react'
import PageHeader from '../components/PageHeader'
import StatusBadge from '../components/StatusBadge'
import DataTable, { type Column } from '../components/DataTable'
import Loading from '../components/Loading'
import { deviceApi } from '../api/services'
import { formatDate, formatDateOnly, money } from '../utils/format'
import type { Device, InspectionRecord, Repair } from '../types'

export default function DeviceDetailPage() {
  const { id } = useParams()
  const [device, setDevice] = useState<Device>()
  const [records, setRecords] = useState<InspectionRecord[]>([])
  const [repairs, setRepairs] = useState<Repair[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!id) return
    Promise.all([deviceApi.get(Number(id)), deviceApi.history(Number(id))])
      .then(([deviceData, history]) => {
        setDevice(deviceData)
        setRecords(history.records)
        setRepairs(history.repairs)
      })
      .finally(() => setLoading(false))
  }, [id])

  const recordColumns: Column<InspectionRecord>[] = [
    { key: 'executedAt', title: '执行时间', render: (row) => formatDate(row.executedAt) },
    { key: 'executor', title: '执行人', render: (row) => `用户 #${row.executorId}` },
    { key: 'status', title: '结果', render: (row) => <StatusBadge status={row.status} /> },
    { key: 'summary', title: '摘要', render: (row) => row.resultSummary },
  ]

  const repairColumns: Column<Repair>[] = [
    { key: 'type', title: '类型', render: (row) => row.repairType },
    { key: 'content', title: '维修内容', render: (row) => row.content },
    { key: 'vendor', title: '维修单位', render: (row) => row.vendor },
    { key: 'period', title: '维修周期', render: (row) => `${formatDateOnly(row.startedAt)} 至 ${formatDateOnly(row.endedAt)}` },
    { key: 'cost', title: '费用', render: (row) => money(row.cost) },
    { key: 'result', title: '结果', render: (row) => row.result },
  ]

  if (loading) return <Loading />
  if (!device) return <div className="empty-state">设备不存在</div>

  return (
    <>
      <PageHeader title={device.name} description={`设备编号 ${device.code}`} actions={<Link className="secondary-button" to="/devices"><ArrowLeft size={17} />返回列表</Link>} />
      <div className="detail-grid">
        <div className="detail-panel">
          <h3>基础信息</h3>
          <dl className="detail-list">
            <dt>类别</dt><dd>{device.category}</dd>
            <dt>型号</dt><dd>{device.model || '-'}</dd>
            <dt>序列号</dt><dd>{device.serialNumber || '-'}</dd>
            <dt>生产厂商</dt><dd>{device.manufacturer || '-'}</dd>
            <dt>安装位置</dt><dd>{device.location}</dd>
            <dt>启用日期</dt><dd>{formatDateOnly(device.installDate)}</dd>
            <dt>保修截止</dt><dd>{formatDateOnly(device.warrantyEnd)}</dd>
            <dt>负责人</dt><dd>{device.owner || '-'}</dd>
            <dt>状态</dt><dd><StatusBadge status={device.status} /></dd>
            <dt>备注</dt><dd>{device.remark || '-'}</dd>
          </dl>
        </div>
        <div className="detail-main">
          <section>
            <h3>巡检历史</h3>
            <DataTable columns={recordColumns} rows={records} emptyText="暂无巡检记录" />
          </section>
          <section>
            <h3>维修历史</h3>
            <DataTable columns={repairColumns} rows={repairs} emptyText="暂无维修记录" />
          </section>
        </div>
      </div>
    </>
  )
}
