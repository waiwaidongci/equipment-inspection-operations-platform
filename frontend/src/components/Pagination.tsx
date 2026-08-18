export default function Pagination({ page, totalPages, total, onChange }: {
  page: number
  totalPages: number
  total: number
  onChange: (page: number) => void
}) {
  if (totalPages <= 1) return null
  return (
    <div className="pagination">
      <span>共 {total} 条</span>
      <button className="secondary-button" disabled={page <= 1} onClick={() => onChange(page - 1)}>上一页</button>
      <span>{page} / {totalPages}</span>
      <button className="secondary-button" disabled={page >= totalPages} onClick={() => onChange(page + 1)}>下一页</button>
    </div>
  )
}
