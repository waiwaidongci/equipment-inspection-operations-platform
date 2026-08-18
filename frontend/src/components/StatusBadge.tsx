const labels: Record<string, string> = {
  active: '启用',
  inactive: '停用',
  enabled: '启用',
  paused: '暂停',
  terminated: '终止',
  pending: '待执行',
  draft: '草稿',
  completed: '已完成',
  abnormal: '异常',
  open: '待处理',
  assigned: '已指派',
  in_progress: '处理中',
  closed: '已关闭',
  low: '低',
  medium: '中',
  high: '高',
  critical: '紧急',
}

export default function StatusBadge({ status }: { status: string }) {
  return <span className={`status-badge status-${status}`}>{labels[status] || status}</span>
}
