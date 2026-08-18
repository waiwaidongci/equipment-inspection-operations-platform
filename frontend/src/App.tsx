import { Navigate, Route, Routes } from 'react-router-dom'
import AppLayout from './components/AppLayout'
import LoginPage from './pages/LoginPage'
import DeviceListPage from './pages/DeviceListPage'
import DeviceDetailPage from './pages/DeviceDetailPage'
import PlanListPage from './pages/PlanListPage'
import TaskListPage from './pages/TaskListPage'
import TaskExecutePage from './pages/TaskExecutePage'
import AnomalyListPage from './pages/AnomalyListPage'
import RepairListPage from './pages/RepairListPage'
import { getToken } from './utils/auth'

function RequireAuth({ children }: { children: JSX.Element }) {
  return getToken() ? children : <Navigate to="/login" replace />
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/" element={<RequireAuth><AppLayout /></RequireAuth>}>
        <Route index element={<Navigate to="/devices" replace />} />
        <Route path="devices" element={<DeviceListPage />} />
        <Route path="devices/:id" element={<DeviceDetailPage />} />
        <Route path="plans" element={<PlanListPage />} />
        <Route path="tasks" element={<TaskListPage />} />
        <Route path="tasks/:id/execute" element={<TaskExecutePage />} />
        <Route path="anomalies" element={<AnomalyListPage />} />
        <Route path="repairs" element={<RepairListPage />} />
      </Route>
    </Routes>
  )
}
