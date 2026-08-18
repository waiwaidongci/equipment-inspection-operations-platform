import { FormEvent, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { HardHat, Lock, User } from 'lucide-react'
import { authApi } from '../api/services'
import { saveSession } from '../utils/auth'

export default function LoginPage() {
  const navigate = useNavigate()
  const [username, setUsername] = useState('admin')
  const [password, setPassword] = useState('admin123')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function submit(event: FormEvent) {
    event.preventDefault()
    setError('')
    setLoading(true)
    try {
      const result = await authApi.login(username, password)
      saveSession(result.token, result.user)
      navigate('/devices')
    } catch (err) {
      setError(err instanceof Error ? err.message : '登录失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-page">
      <form className="login-card" onSubmit={submit}>
        <div className="login-brand">
          <HardHat size={34} />
          <h1>设备巡检与维护平台</h1>
          <p>登录后管理设备台账、巡检任务和异常闭环</p>
        </div>
        <label className="login-field">
          <User size={18} />
          <input value={username} onChange={(event) => setUsername(event.target.value)} placeholder="用户名" autoComplete="username" />
        </label>
        <label className="login-field">
          <Lock size={18} />
          <input type="password" value={password} onChange={(event) => setPassword(event.target.value)} placeholder="密码" autoComplete="current-password" />
        </label>
        {error && <div className="error-banner">{error}</div>}
        <button className="primary-button login-submit" disabled={loading}>{loading ? '登录中...' : '登录'}</button>
        <p className="login-hint">默认账号：admin / admin123</p>
      </form>
    </div>
  )
}
