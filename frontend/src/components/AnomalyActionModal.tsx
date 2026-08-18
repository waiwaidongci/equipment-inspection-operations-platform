import { useState } from 'react'
import Modal from './Modal'
import { Field, Select, TextArea } from './FormField'
import type { Anomaly, User } from '../types'

export default function AnomalyActionModal({ open, anomaly, users, action, onClose, onSave }: {
  open: boolean
  anomaly: Anomaly
  users: User[]
  action: 'assign' | 'progress' | 'close'
  onClose: () => void
  onSave: (body: unknown) => Promise<void>
}) {
  const [assigneeId, setAssigneeId] = useState(0)
  const [progress, setProgress] = useState('')
  const [causeAnalysis, setCauseAnalysis] = useState('')
  const [closeNote, setCloseNote] = useState('')
  const [verificationResult, setVerificationResult] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function submit() {
    if (action === 'assign' && assigneeId <= 0) {
      setError('请选择负责人')
      return
    }
    if (action === 'progress' && !progress) {
      setError('请填写处理进展')
      return
    }
    if (action === 'close' && (!causeAnalysis || !closeNote || !verificationResult)) {
      setError('关闭原因、关闭说明和验证结果均为必填')
      return
    }
    setLoading(true)
    try {
      if (action === 'assign') await onSave({ assigneeId })
      if (action === 'progress') await onSave({ progress })
      if (action === 'close') await onSave({ causeAnalysis, closeNote, verificationResult })
      onClose()
    } catch (err) {
      setError(err instanceof Error ? err.message : '操作失败')
    } finally {
      setLoading(false)
    }
  }

  const title = action === 'assign' ? '指派异常' : action === 'progress' ? '更新进展' : '关闭异常'
  return (
    <Modal title={title} open={open} onClose={onClose} footer={
      <>
        <button className="secondary-button" onClick={onClose}>取消</button>
        <button className="primary-button" onClick={submit} disabled={loading}>{loading ? '提交中...' : '确认'}</button>
      </>
    }>
      <div className="anomaly-context">
        <strong>{anomaly.description}</strong>
        <span>设备 #{anomaly.deviceId} · {anomaly.severity}</span>
      </div>
      {action === 'assign' && <Field label="负责人" required><Select value={assigneeId} onChange={(e) => setAssigneeId(Number(e.target.value))}><option value={0}>请选择</option>{users.map((user) => <option key={user.id} value={user.id}>{user.displayName}</option>)}</Select></Field>}
      {action === 'progress' && <Field label="处理进展" required><TextArea value={progress} onChange={(e) => setProgress(e.target.value)} /></Field>}
      {action === 'close' && <>
        <Field label="原因分析" required><TextArea value={causeAnalysis} onChange={(e) => setCauseAnalysis(e.target.value)} /></Field>
        <Field label="关闭说明" required><TextArea value={closeNote} onChange={(e) => setCloseNote(e.target.value)} /></Field>
        <Field label="验证结果" required><TextArea value={verificationResult} onChange={(e) => setVerificationResult(e.target.value)} /></Field>
      </>}
      {error && <div className="error-banner">{error}</div>}
    </Modal>
  )
}
