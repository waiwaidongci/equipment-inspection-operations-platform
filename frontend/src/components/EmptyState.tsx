export default function EmptyState({ text = '暂无数据' }: { text?: string }) {
  return <div className="empty-state">{text}</div>
}
