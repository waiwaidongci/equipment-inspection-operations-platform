import { useState } from 'react'
import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { Boxes, CalendarCheck, ClipboardList, HardHat, LogOut, Menu, ShieldAlert, X } from 'lucide-react'
import { clearSession } from '../utils/auth'

const navItems = [
  { to: '/devices', label: '设备台账', icon: Boxes },
  { to: '/plans', label: '巡检计划', icon: CalendarCheck },
  { to: '/tasks', label: '巡检任务', icon: ClipboardList },
  { to: '/anomalies', label: '异常管理', icon: ShieldAlert },
  { to: '/repairs', label: '维修记录', icon: HardHat },
]

export default function AppLayout() {
  const [open, setOpen] = useState(false)
  const navigate = useNavigate()

  function logout() {
    clearSession()
    navigate('/login')
  }

  return (
    <div className="app-shell">
      <aside className={`sidebar ${open ? 'open' : ''}`}>
        <div className="brand">
          <HardHat size={24} />
          <span>巡检维护平台</span>
        </div>
        <nav>
          {navItems.map((item) => {
            const Icon = item.icon
            return (
              <NavLink key={item.to} to={item.to} className={({ isActive }) => `nav-link ${isActive ? 'active' : ''}`} onClick={() => setOpen(false)}>
                <Icon size={18} />
                <span>{item.label}</span>
              </NavLink>
            )
          })}
        </nav>
        <button className="logout-button" onClick={logout}>
          <LogOut size={18} />
          <span>退出登录</span>
        </button>
      </aside>
      <div className="main-area">
        <header className="topbar">
          <button className="icon-button mobile-menu" onClick={() => setOpen((value) => !value)}>
            {open ? <X size={20} /> : <Menu size={20} />}
          </button>
          <span className="topbar-title">设备巡检与维护平台</span>
        </header>
        <main className="content">
          <Outlet />
        </main>
      </div>
    </div>
  )
}
