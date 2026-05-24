import React from 'react'
import { NavLink, Outlet, useNavigate, useLocation } from 'react-router-dom'
import { clearToken } from '../api'

const pageTitles = {
  '/dashboard': '数据概览',
  '/menu': '菜单管理',
  '/orders': '订单管理',
  '/order-alert': '订单提醒',
  '/users': '用户管理',
  '/activity': '活动管理',
  '/points': '积分商品',
  '/recharge': '充值方案',
  '/pets': '宠物管理'
}

export default function Layout() {
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = () => {
    clearToken()
    navigate('/')
  }

  const title = pageTitles[location.pathname] || '管理后台'

  return (
    <div className="app">
      <aside className="sidebar">
        <div className="logo">汪托帮</div>
        <NavLink to="/dashboard" className={({isActive}) => `nav-item ${isActive ? 'active' : ''}`}>📊 数据概览</NavLink>
        <NavLink to="/menu" className={({isActive}) => `nav-item ${isActive ? 'active' : ''}`}>🍽️ 菜单管理</NavLink>
        <NavLink to="/orders" className={({isActive}) => `nav-item ${isActive ? 'active' : ''}`}>📋 订单管理</NavLink>
        <NavLink to="/order-alert" className={({isActive}) => `nav-item ${isActive ? 'active' : ''}`}>🔔 订单提醒</NavLink>
        <NavLink to="/users" className={({isActive}) => `nav-item ${isActive ? 'active' : ''}`}>👤 用户管理</NavLink>
        <NavLink to="/activity" className={({isActive}) => `nav-item ${isActive ? 'active' : ''}`}>🎉 活动管理</NavLink>
        <NavLink to="/points" className={({isActive}) => `nav-item ${isActive ? 'active' : ''}`}>⭐ 积分商品</NavLink>
        <NavLink to="/recharge" className={({isActive}) => `nav-item ${isActive ? 'active' : ''}`}>💰 充值方案</NavLink>
        <NavLink to="/pets" className={({isActive}) => `nav-item ${isActive ? 'active' : ''}`}>🐾 宠物管理</NavLink>
        <button className="nav-item" onClick={handleLogout}>🚪 退出登录</button>
      </aside>
      <div className="main">
        <header className="topbar"><span>{title}</span></header>
        <div className="content">
          <Outlet />
        </div>
      </div>
    </div>
  )
}
