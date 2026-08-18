import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import { ArrowLeft, Save, Send } from 'lucide-react'
import PageHeader from '../components/PageHeader'
import StatusBadge from '../components/StatusBadge'
import Loading from '../components/Loading'
import { taskApi } from '../api/services'
import { formatDate } from '../utils/format'
import type { InspectionTask } from '../types'

export default function TaskExecutePage() {
  const { id } = useParams()
  const navigate = useNavigate()
  const [task, setTask] = useState<InspectionTask>()
  const [values, setValues] = useState<Record<number, { value: string; abnormal: boolean; remark: string }>>({})
  const [overallRemark, setOverallRemark] = useState('')
  const [abnormal, setAbnormal] = useState(false)
  const [severity, setSeverity] = useState('medium')
  const [anomalyDescription, setAnomalyDescription] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    if (!id) return
    taskApi.get(Number(id)).then((data) => {
      setTask(data)
      const initial: Record<number, { value: string; abnormal: boolean; remark: string }> = {}
      data.results.forEach((item) => {
        initial[item.itemId] = { value: item.value, abnormal: item.abnormal, remark: item.remark }
      })
      setValues(initial)
      setOverallRemark(data.remark)
    }).finally(() => setLoading(false))
  }, [id])

  function update(itemId: number, patch: Partial<{ value: string; abnormal: boolean; remark: string }>) {
    setValues((current) => ({ ...current, [itemId]: { ...current[itemId], ...patch } }))
  }

  function buildBody() {
    return {
      results: (task?.results || []).map((item) => ({
        itemId: item.itemId,
        value: values[item.itemId]?.value || '',
        passed: !values[item.itemId]?.abnormal,
        abnormal: values[item.itemId]?.abnormal || false,
        remark: values[item.itemId]?.remark || '',
      })),
      remark: overallRemark,
      abnormal,
      severity,
      anomalyDescription,
    }
  }

  async function saveDraft() {
    if (!task) return
    setError('')
    setSaving(true)
    try {
      await taskApi.draft(task.id, buildBody())
      navigate('/tasks')
    } catch (err) {
      setError(err instanceof Error ? err.message : '保存失败')
    } finally {
      setSaving(false)
    }
  }

  async function submit() {
    if (!task) return
    setError('')
    for (const item of task.results) {
      if (item.required && !values[item.itemId]?.value) {
        setError(`必填项「${item.name}」未填写`)
        return
      }
    }
    setSaving(true)
    try {
      await taskApi.submit(task.id, buildBody())
      navigate('/tasks')
    } catch (err) {
      setError(err instanceof Error ? err.message : '提交失败')
    } finally {
      setSaving(false)
    }
  }

  if (loading) return <Loading />
  if (!task) return <div className="empty-state">任务不存在</div>
  const completed = task.status === 'completed' || task.status === 'abnormal'

  return (
    <>
      <PageHeader title="执行巡检任务" description={`计划执行时间：${formatDate(task.plannedAt)}`} actions={<Link className="secondary-button" to="/tasks"><ArrowLeft size={17} />返回任务</Link>} />
      <div className="task-summary">
        <StatusBadge status={task.status} />
        <span>设备 #{task.deviceId} · 计划 #{task.planId}</span>
      </div>
      {completed && <div className="info-banner">该任务已经提交，不能重复执行。</div>}
      <div className="check-form">
        {task.results.map((item) => (
          <div className="check-item" key={item.itemId}>
            <div className="check-item-head">
              <strong>{item.name}{item.required && <em>*</em>}</strong>
              {item.standardRange && <span>标准范围：{item.standardRange}</span>}
              {item.description && <small>{item.description}</small>}
            </div>
            <div className="check-item-body">
              <input className="control" disabled={completed} value={values[item.itemId]?.value || ''} onChange={(e) => update(item.itemId, { value: e.target.value })} placeholder="填写测量值或结果" />
              <input className="control" disabled={completed} value={values[item.itemId]?.remark || ''} onChange={(e) => update(item.itemId, { remark: e.target.value })} placeholder="现场备注" />
              <label className="inline-check"><input type="checkbox" disabled={completed} checked={values[item.itemId]?.abnormal || false} onChange={(e) => update(item.itemId, { abnormal: e.target.checked })} />标记异常</label>
            </div>
          </div>
        ))}
        {!completed && <>
          <div className="form-grid">
            <label className="field"><span className="field-label">现场备注</span><textarea className="control" value={overallRemark} onChange={(e) => setOverallRemark(e.target.value)} /></label>
            <label className="field"><span className="field-label">发现异常</span><input className="control" type="checkbox" checked={abnormal} onChange={(e) => setAbnormal(e.target.checked)} /></label>
            <label className="field"><span className="field-label">严重程度</span><select className="control" value={severity} onChange={(e) => setSeverity(e.target.value)}><option value="low">低</option><option value="medium">中</option><option value="high">高</option><option value="critical">紧急</option></select></label>
            <label className="field"><span className="field-label">异常说明</span><textarea className="control" value={anomalyDescription} onChange={(e) => setAnomalyDescription(e.target.value)} placeholder="异常现象和处理要求" /></label>
          </div>
          {error && <div className="error-banner">{error}</div>}
          <div className="form-actions">
            <button className="secondary-button" onClick={saveDraft} disabled={saving}><Save size={17} />保存草稿</button>
            <button className="primary-button" onClick={submit} disabled={saving}><Send size={17} />{saving ? '提交中...' : '提交结果'}</button>
          </div>
        </>}
      </div>
    </>
  )
}
