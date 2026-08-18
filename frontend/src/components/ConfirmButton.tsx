import { useState } from 'react'
import type { ReactNode } from 'react'

export default function ConfirmButton({ children, message, onConfirm, className = '' }: {
  children: ReactNode
  message: string
  onConfirm: () => void
  className?: string
}) {
  const [asking, setAsking] = useState(false)
  if (asking) {
    return (
      <span className="confirm-group">
        <span>{message}</span>
        <button className="primary-button small" onClick={() => { onConfirm(); setAsking(false) }}>确认</button>
        <button className="secondary-button small" onClick={() => setAsking(false)}>取消</button>
      </span>
    )
  }
  return <button className={className} onClick={() => setAsking(true)}>{children}</button>
}
