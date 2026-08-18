import { useEffect, useState } from 'react'
import Modal from './Modal'
import { Field, Select, TextArea, TextInput } from './FormField'
import type { Anomaly, Device, Repair } from '../types'

export interface RepairFormData {
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
}

const emptyForm: RepairFormData = {
  deviceId: 0, repairType: '', content: '', vendor: '',
  startedAt: '', endedAt: '', cost: 0, result: '', remark: '',
}

export default function RepairFormModal({ open, repair, devices, closedAnomalies, onClose, onSave }: {
  open: boolean
  repair?: Repair
  devices: Device[]
  closedAnomalies: Anomaly[]
  onClose: () => void
  onSave: (form: RepairFormData) => Promise<void>
}) {
  const [form, setForm] = useState<RepairFormData>(emptyForm)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (repair) {
      setForm({
        deviceId: repair.deviceId, anomalyId: repair.anomalyId, repairType: repair.repairType,
        content: repair.content, vendor: repair.vendor, startedAt: repair.startedAt,
        endedAt: repair.endedAt, cost: repair.cost, result: repair.result, remark: repair.remark,
      })
    } else {
      setForm(emptyForm)
    }
    setError('')
  }, [repair, open])

  function update<K extends keyof RepairFormData>(key: K, value: RepairFormData[K]) {
    setForm((current) => ({ ...current, [key]: value }))
  }

  async function submit() {
    if (!form.deviceId || !form.repairType || !form.content || !form.result || !form.startedAt || !form.endedAt) {
      setError('设备、类型、内容、结果和起止时间为必填项')
      return
    }
    setLoading(true)
    try {
      await onSave(form)
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : '保存失败')
    } finally {
      setLoading(false)
    }
  }

  const linkedAnomalies = closedAnomalies.filter((item) => !form.deviceId || item.deviceId === form.deviceId)

  return (
    <Modal title={repair ? '编辑维修记录' : '新增维修记录'} open={open} onClose={onClose} footer={
      <>
        <button className="secondary-button" onClick={onClose}>取消</button>
        <button className="primary-button" onClick={submit} disabled={loading}>{loading ? '保存中...' : '保存'}</button>
      </>
    }>
      <div className="form-grid">
        <Field label="设备" required><Select value={form.deviceId} onChange={(e) => update('deviceId', Number(e.target.value))}><option value={0}>请选择设备</option>{devices.map((device) => <option key={device.id} value={device.id}>{device.code} · {device.name}</option>)}</Select></Field>
        <Field label="关联异常"><Select value={form.anomalyId || 0} onChange={(e) => update('anomalyId', Number(e.target.value) || undefined)}><option value={0}>不关联</option>{linkedAnomalies.map((item) => <option key={item.id} value={item.id}>异常 #{item.id} · {item.description}</option>)}</Select></Field>
        <Field label="维修类型" required><TextInput value={form.repairType} onChange={(e) => update('repairType', e.target.value)} placeholder="预防性维护 / 故障维修" /></Field>
        <Field label="维修单位"><TextInput value={form.vendor} onChange={(e) => update('vendor', e.target.value)} /></Field>
        <Field label="开始日期" required><TextInput type="date" value={form.startedAt} onChange={(e) => update('startedAt', e.target.value)} /></Field>
        <Field label="结束日期" required><TextInput type="date" value={form.endedAt} onChange={(e) => update('endedAt', e.target.value)} /></Field>
        <Field label="费用"><TextInput type="number" min="0" step="0.01" value={form.cost} onChange={(e) => update('cost', Number(e.target.value))} /></Field>
        <Field label="结果" required><TextInput value={form.result} onChange={(e) => update('result', e.target.value)} /></Field>
        <Field label="维修内容" required><TextArea value={form.content} onChange={(e) => update('content', e.target.value)} /></Field>
        <Field label="备注"><TextArea value={form.remark} onChange={(e) => update('remark', e.target.value)} /></Field>
      </div>
      {error && <div className="error-banner">{error}</div>}
    </Modal>
  )
}
