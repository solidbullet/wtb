import React from 'react'
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Layout from './components/Layout'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import MenuManagement from './pages/MenuManagement'
import OrderManagement from './pages/OrderManagement'
import OrderAlert from './pages/OrderAlert'
import ActivityManagement from './pages/ActivityManagement'
import PointsManagement from './pages/PointsManagement'
import RechargePlanManagement from './pages/RechargePlanManagement'
import PetManagement from './pages/PetManagement'
import SeatManagement from './pages/SeatManagement'
import UserManagement from './pages/UserManagement'

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route element={<Layout />}>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/menu" element={<MenuManagement />} />
          <Route path="/orders" element={<OrderManagement />} />
          <Route path="/order-alert" element={<OrderAlert />} />
          <Route path="/users" element={<UserManagement />} />
          <Route path="/activity" element={<ActivityManagement />} />
          <Route path="/points" element={<PointsManagement />} />
          <Route path="/recharge" element={<RechargePlanManagement />} />
          <Route path="/pets" element={<PetManagement />} />
          <Route path="/seats" element={<SeatManagement />} />
        </Route>
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </BrowserRouter>
  )
}
