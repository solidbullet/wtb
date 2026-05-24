import React, { useState, useEffect, useRef, useCallback } from 'react'
import { get, API_BASE } from '../api'

const POLL_INTERVAL = 15000

export default function OrderAlert() {
  const [orders, setOrders] = useState([])
  const [audioEnabled, setAudioEnabled] = useState(false)
  const seenOrderIds = useRef(new Set())
  const audioRef = useRef(null)
  const loadingRef = useRef(false)

  useEffect(() => {
    audioRef.current = new Audio(API_BASE + '/assets/notification.wav')
  }, [])

  const playSound = useCallback(() => {
    if (audioRef.current && audioEnabled) {
      audioRef.current.currentTime = 0
      audioRef.current.play().catch(() => {})
    }
  }, [audioEnabled])

  const loadOrders = useCallback(async () => {
    if (loadingRef.current) return
    loadingRef.current = true
    try {
      const list = await get('/api/order/admin/today-paid')
      if (!Array.isArray(list)) return
      const newOrders = list.filter(o => !seenOrderIds.current.has(o.id))
      if (newOrders.length > 0) {
        newOrders.forEach(o => seenOrderIds.current.add(o.id))
        playSound()
      }
      setOrders(list)
    } catch (e) {
      console.error('load orders failed', e)
    } finally {
      loadingRef.current = false
    }
  }, [playSound])

  useEffect(() => {
    loadOrders()
    const timer = setInterval(loadOrders, POLL_INTERVAL)
    return () => clearInterval(timer)
  }, [loadOrders])

  const fmtPrice = (p) => p ? '¥' + (p / 100).toFixed(2) : '-'
  const fmtTime = (t) => t ? new Date(t).toLocaleTimeString('zh-CN') : '-'

  return (
    <div id="page-order-alert">
      <div className="section">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
          <h3>订单提醒（当天已支付）</h3>
          <div>
            {!audioEnabled ? (
              <button className="btn-primary" onClick={() => setAudioEnabled(true)}>开启声音提醒</button>
            ) : (
              <span style={{ color: '#52c41a', fontSize: 14 }}>声音已开启</span>
            )}
          </div>
        </div>
        {orders.length === 0 && (
          <div style={{ textAlign: 'center', color: '#999', padding: '40px 0' }}>暂无当天已支付的订单</div>
        )}
        {orders.length > 0 && (
          <table className="data-table">
            <thead>
              <tr><th>订单号</th><th>时间</th><th>金额</th><th>座位</th><th>状态</th><th>菜品</th></tr>
            </thead>
            <tbody>
              {orders.map(o => (
                <tr key={o.id}>
                  <td>{o.order_no}</td>
                  <td>{fmtTime(o.created_at)}</td>
                  <td>{fmtPrice(o.pay_amount)}</td>
                  <td>{o.seat_id}</td>
                  <td>{o.status === 'paid' ? <span style={{ color: '#faad14' }}>已支付</span> : o.status === 'completed' ? <span style={{ color: '#52c41a' }}>已完成</span> : <span>{o.status}</span>}</td>
                  <td>{o.items?.map((it, idx) => (<div key={idx} style={{ fontSize: 12 }}>{it.dish_name} x{it.quantity}</div>)) || '-'}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
