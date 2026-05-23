import React, { useState, useEffect } from 'react'
import { get, put } from '../api'

export default function OrderManagement() {
  const [orders, setOrders] = useState([])

  const loadOrders = async () => {
    try {
      const data = await get('/api/admin/order/list')
      setOrders(data?.list || [])
    } catch (e) {
      console.error('Orders error:', e)
    }
  }

  useEffect(() => {
    loadOrders()
  }, [])

  const updateOrderStatus = async (id, status) => {
    try {
      await put('/api/admin/order/admin/status', { order_id: id, status })
      loadOrders()
    } catch (e) {
      alert(e.message)
    }
  }

  const statusText = { pending: '待支付', paid: '已支付', completed: '已完成', cancelled: '已取消' }

  const formatItems = (items) => {
    if (!items || items.length === 0) return '-'
    return items.map(it => `${it.dish_name} x${it.quantity}`).join('，')
  }

  return (
    <div id="page-orders">
      <div className="section">
        <h3>订单列表</h3>
        <table className="data-table">
          <thead>
            <tr>
              <th>订单号</th>
              <th>桌号</th>
              <th>用户</th>
              <th>菜品明细</th>
              <th>备注</th>
              <th>金额</th>
              <th>状态</th>
              <th>时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {orders.length === 0 && <tr><td colSpan="9" style={{textAlign:'center',color:'#999'}}>暂无订单</td></tr>}
            {orders.map(o => (
              <tr key={o.id}>
                <td>{o.order_no || o.id}</td>
                <td>{o.seat_id || '-'}</td>
                <td>{o.user_id}</td>
                <td style={{maxWidth:240,whiteSpace:'nowrap',overflow:'hidden',textOverflow:'ellipsis'}} title={formatItems(o.items)}>
                  {formatItems(o.items)}
                </td>
                <td>{o.remark || '-'}</td>
                <td>¥{(o.total_amount || 0) / 100}</td>
                <td>{statusText[o.status] || o.status}</td>
                <td>{o.created_at ? o.created_at.split('T')[0] : '-'}</td>
                <td>
                  <button className="btn-primary" onClick={() => updateOrderStatus(o.id, 'completed')}>完成</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
