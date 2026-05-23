import React, { useState, useEffect } from 'react'
import { get } from '../api'

export default function Dashboard() {
  const [stats, setStats] = useState({ dishes: 0, activities: 0, orders: 0, users: '1+' })

  useEffect(() => {
    async function loadDashboard() {
      try {
        const [dishes, activities, orders] = await Promise.all([
          get('/api/admin/menu/dishes').catch(() => ({ total: 0 })),
          get('/api/admin/activity/list').catch(() => []),
          get('/api/admin/order/list').catch(() => ({ total: 0 }))
        ])
        setStats({
          dishes: dishes?.total || 0,
          activities: activities?.length || 0,
          orders: orders?.total || 0,
          users: '1+'
        })
      } catch (e) {
        console.error('Dashboard error:', e)
      }
    }
    loadDashboard()
  }, [])

  return (
    <div id="page-dashboard">
      <div className="stats-grid">
        <div className="stat-card"><div className="stat-label">今日订单</div><div className="stat-value">{stats.orders}</div></div>
        <div className="stat-card"><div className="stat-label">注册用户</div><div className="stat-value">{stats.users}</div></div>
        <div className="stat-card"><div className="stat-label">菜品数量</div><div className="stat-value">{stats.dishes}</div></div>
        <div className="stat-card"><div className="stat-label">活动数量</div><div className="stat-value">{stats.activities}</div></div>
      </div>
    </div>
  )
}
