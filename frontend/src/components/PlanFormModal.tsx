import { useEffect, useState } from 'react'
import { Plus, Trash2 } from 'lucide-react'
import Modal from './Modal'
import { Field, Select, TextArea, TextInput } from './FormField'
import type { Device, InspectionPlan, PlanCheckItem } from '../types'

export interface PlanItemForm {
  name: string
  itemType: string
  standardRange: string
  required: boolean
  description: string
}

export interface PlanFormData {
  name: string
  scopeType: 'device' | 'category'
  deviceId?: number
  category: string
  period: 'daily' | 'weekly' | 'monthly'
  startDate: string
  endDate: string
  owner: string
  status: 'enabled' | 'paused' | 'terminated'
  items: PlanItemForm[]
}

const emptyItem: PlanItemForm = { name: '', itemType: 'text', standardRange: '', required: true, description: '' }
const emptyForm: PlanFormData = {
  name: '', scopeType: 'device', category: '', period: 'daily',
  startDate: '', endDate: '', owner: '', status: 'enabled', items: [{ ...emptyItem }],
}

export default function PlanFormModal({ open, plan, devices, categories, onClose, onSave }: {
  open: boolean
  plan?: InspectionPlan
  devices: Device[]
  categories: string[]
  onClose: () => void
  onSave: (form: PlanFormData) => Promise<void>
}) {
  const [form, setForm] = useState<PlanFormData>(emptyForm)
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (plan) {
      setForm({
        name: plan.name, scopeType: plan.scopeType, deviceId: plan.deviceId,
        category: plan.category, period: plan.period, startDate: plan.startDate,
        endDate: plan.endDate, owner: plan.owner, status: plan.status,
        items: (plan.items || []).map((item) => ({
          name: item.name, itemType: item.itemType, standardRange: item.standardRange,
          required: item.required, description: item.description,
        })),
      })
    } else {
      setForm(emptyForm)
    }
    setError('')
  }, [plan, open])

  function update<K extends keyof PlanFormData>(key: K, value: PlanFormData[K]) {
    setForm((current) => ({ ...current, [key]: value }))
  }

  function updateItem(index: number, patch: Partial<PlanItemForm>) {
    setForm((current) => ({
      ...current,
      items: current.items.map((item, i) => i === index ? { ...item, ...patch } : item),
    }))
  }

  async function submit() {
    if (!form.name || !form.owner || !form.startDate || !form.endDate) {
      setError('名称、负责人和起止日期为必填项')
      return
    }
    if (form.scopeType === 'device' && !form.deviceId) {
      setError('指定设备巡检必须选择设备')
      return
    }
    if (form.scopeType === 'category' && !form.category) {
      setError('按类别巡检必须选择类别')
      return
    }
    if (form.items.some((item) => !item.name)) {
      setError('检查项目名称不能为空')
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
    <Modal title={plan ? '编辑巡检计划' : '新增巡检计划'} open={open} onClose={onClose} footer={
      <>
        <button className="secondary-button" onClick={onClose}>取消</button>
        <button className="primary-button" onClick={submit} disabled={loading}>{loading ? '保存中...' : '保存'}</button>
      </>
    }>
      <div className="form-grid">
        <Field label="计划名称" required><TextInput value={form.name} onChange={(e) => update('name', e.target.value)} /></Field>
        <Field label="负责人" required><TextInput value={form.owner} onChange={(e) => update('owner', e.target.value)} /></Field>
        <Field label="适用方式"><Select value={form.scopeType} onChange={(e) => update('scopeType', e.target.value as PlanFormData['scopeType'])}><option value="device">指定设备</option><option value="category">按类别</option></Select></Field>
        {form.scopeType === 'device'
          ? <Field label="设备"><Select value={form.deviceId || ''} onChange={(e) => update('deviceId', Number(e.target.value))}><option value="">请选择设备</option>{devices.map((device) => <option key={device.id} value={device.id}>{device.code} · {device.name}</option>)}</Select></Field>
          : <Field label="设备类别"><Select value={form.category} onChange={(e) => update('category', e.target.value)}><option value="">请选择类别</option>{categories.map((item) => <option key={item}>{item}</option>)}</Select></Field>}
        <Field label="巡检周期"><Select value={form.period} onChange={(e) => update('period', e.target.value as PlanFormData['period'])}><option value="daily">每日</option><option value="weekly">每周</option><option value="monthly">每月</option></Select></Field>
        <Field label="开始日期" required><TextInput type="date" value={form.startDate} onChange={(e) => update('startDate', e.target.value)} /></Field>
        <Field label="结束日期" required><TextInput type="date" value={form.endDate} onChange={(e) => update('endDate', e.target.value)} /></Field>
        <Field label="状态"><Select value={form.status} onChange={(e) => update('status', e.target.value as PlanFormData['status'])}><option value="enabled">启用</option><option value="paused">暂停</option></Select></Field>
      </div>
      <div className="items-editor">
        <div className="items-editor-head"><strong>检查项目</strong><button className="secondary-button small" onClick={() => setForm((f) => ({ ...f, items: [...f.items, { ...emptyItem }] }))}><Plus size={15} />添加项目</button></div>
        {form.items.map((item, index) => (
          <div className="item-editor-row" key={index}>
            <TextInput placeholder="项目名称" value={item.name} onChange={(e) => updateItem(index, { name: e.target.value })} />
            <Select value={item.itemType} onChange={(e) => updateItem(index, { itemType: e.target.value })}><option value="text">文本</option><option value="number">数值</option><option value="boolean">是/否</option></Select>
            <TextInput placeholder="标准值范围" value={item.standardRange} onChange={(e) => updateItem(index, { standardRange: e.target.value })} />
            <label className="inline-check"><input type="checkbox" checked={item.required} onChange={(e) => updateItem(index, { required: e.target.checked })} />必填</label>
            <TextArea placeholder="说明" value={item.description} onChange={(e) => updateItem(index, { description: e.target.value })} />
            <button className="icon-button danger" onClick={() => setForm((f) => ({ ...f, items: f.items.filter((_, i) => i !== index) }))}><Trash2 size={16} /></button>
          </div>
        ))}
      </div>
      {error && <div className="error-banner">{error}</div>}
    </Modal>
  )
}
