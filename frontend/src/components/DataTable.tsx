import type { ReactNode } from 'react'

export interface Column<T> {
  key: string
  title: string
  render: (row: T) => ReactNode
  className?: string
}

export default function DataTable<T extends { id: number }>({ columns, rows, emptyText = '暂无数据' }: {
  columns: Column<T>[]
  rows: T[]
  emptyText?: string
}) {
  if (rows.length === 0) {
    return <div className="empty-state">{emptyText}</div>
  }
  return (
    <div className="table-wrap">
      <table className="data-table">
        <thead>
          <tr>
            {columns.map((column) => <th key={column.key} className={column.className}>{column.title}</th>)}
          </tr>
        </thead>
        <tbody>
          {rows.map((row) => (
            <tr key={row.id}>
              {columns.map((column) => <td key={column.key} className={column.className}>{column.render(row)}</td>)}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
