export default function Loading({ text = '加载中...' }: { text?: string }) {
  return <div className="loading">{text}</div>
}
