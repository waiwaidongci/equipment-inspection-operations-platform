import { useEffect, useState } from 'react'
import Modal from './Modal'
import { Field, Select, TextArea, TextInput } from './FormField'
import type { Device } from '../types'

export interface DeviceFormData {
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
}

const emptyForm: DeviceFormData = {
  code: '', name: '', category: '电气设备', model: '', serialNumber: '',
  manufacturer: '', location: '生产车间', installDate: '', warrantyEnd: '',
  owner: '', status: 'active', remark: '',
}

export default function DeviceFormModal({ open, device, categories, locations, onClose, onSave }: {
  open: boolean
  device?: Device
  categories: string[]
  locations: string[]
  onClose: () => void
  onSave: (form: DeviceFormData) => Promise<void>
}) {
  const [form, setForm] = useState<DeviceFormData>(emptyForm)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (device) {
      setForm({
        code: device.code, name: device.name, category: device.category,
        model: device.model, serialNumber: device.serialNumber,
        manufacturer: device.manufacturer, location: device.location,
        installDate: device.installDate, warrantyEnd: device.warrantyEnd,
        owner: device.owner, status: device.status, remark: device.remark,
      })
    } else {
      setForm(emptyForm)
    }
    setError('')
  }, [device, open])

  function update<K extends keyof DeviceFormData>(key: K, value: DeviceFormData[K]) {
    setForm((current) => ({ ...current, [key]: value }))
  }

  async function submit() {
    if (!form.code || !form.name || !form.location) {
      setError('设备编号、名称和安装位置为必填项')
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

  return (
    <Modal title={device ? '编辑设备' : '新增设备'} open={open} onClose={onClose} footer={
      <>
        <button className="secondary-button" onClick={onClose}>取消</button>
        <button className="primary-button" onClick={submit} disabled={loading}>{loading ? '保存中...' : '保存'}</button>
      </>
    }>
      <div className="form-grid">
        <Field label="设备编号" required><TextInput value={form.code} onChange={(e) => update('code', e.target.value)} /></Field>
        <Field label="名称" required><TextInput value={form.name} onChange={(e) => update('name', e.target.value)} /></Field>
        <Field label="类别"><Select value={form.category} onChange={(e) => update('category', e.target.value)}>{categories.map((item) => <option key={item}>{item}</option>)}</Select></Field>
        <Field label="型号"><TextInput value={form.model} onChange={(e) => update('model', e.target.value)} /></Field>
        <Field label="序列号"><TextInput value={form.serialNumber} onChange={(e) => update('serialNumber', e.target.value)} /></Field>
        <Field label="生产厂商"><TextInput value={form.manufacturer} onChange={(e) => update('manufacturer', e.target.value)} /></Field>
        <Field label="安装位置" required><Select value={form.location} onChange={(e) => update('location', e.target.value)}>{locations.map((item) => <option key={item}>{item}</option>)}</Select></Field>
        <Field label="启用日期"><TextInput type="date" value={form.installDate} onChange={(e) => update('installDate', e.target.value)} /></Field>
        <Field label="保修截止日期"><TextInput type="date" value={form.warrantyEnd} onChange={(e) => update('warrantyEnd', e.target.value)} /></Field>
        <Field label="负责人"><TextInput value={form.owner} onChange={(e) => update('owner', e.target.value)} /></Field>
        <Field label="状态"><Select value={form.status} onChange={(e) => update('status', e.target.value as DeviceFormData['status'])}><option value="active">启用</option><option value="inactive">停用</option></Select></Field>
        <Field label="备注" error={error}><TextArea value={form.remark} onChange={(e) => update('remark', e.target.value)} /></Field>
      </div>
    </Modal>
  )
}
