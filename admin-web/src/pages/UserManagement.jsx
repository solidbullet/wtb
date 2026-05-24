import React, { useState, useEffect } from 'react'
import { get } from '../api'

export default function UserManagement() {
  const [users, setUsers] = useState([])
  const [loading, setLoading] = useState(false)
  const [keyword, setKeyword] = useState('')
  const [detail, setDetail] = useState(null)
  const [showModal, setShowModal] = useState(false)

  const loadUsers = async () => {
    setLoading(true)
    try {
      const res = await get('/api/admin/users?pageSize=100')
      setUsers(res?.list || [])
    } catch (e) {
      alert('加载用户失败: ' + e.message)
    } finally {
      setLoading(false)
    }
  }

  const searchUsers = async () => {
    setLoading(true)
    try {
      const res = await get('/api/admin/users?keyword=' + encodeURIComponent(keyword) + '&pageSize=100')
      setUsers(res?.list || [])
    } catch (e) {
      alert('搜索失败: ' + e.message)
    } finally {
      setLoading(false)
    }
  }

  const showDetail = async (userId) => {
    try {
      const data = await get('/api/admin/users/' + userId)
      setDetail(data)
      setShowModal(true)
    } catch (e) {
      alert('加载详情失败: ' + e.message)
    }
  }

  useEffect(() => {
    loadUsers()
  }, [])

  const fmtPrice = (p) => p ? '¥' + (p / 100).toFixed(2) : '¥0.00'
  const fmtLevel = (l) => {
    if (l === 2) return '充值客户'
    if (l === 1) return '会员客户'
    return '普通客户'
  }
  const fmtTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : '-'

  return (
    <div id="page-users">
      <div className="section">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
          <h3>用户管理</h3>
          <div className="inline-form">
            <input type="text" placeholder="搜索昵称/手机号" value={keyword} onChange={e => setKeyword(e.target.value)} onKeyDown={e => e.key === 'Enter' && searchUsers()} />
            <button className="btn-primary" onClick={searchUsers}>搜索</button>
            <button className="btn-primary" style={{ background: '#999' }} onClick={loadUsers}>重置</button>
          </div>
        </div>

        {loading && <div style={{ textAlign: 'center', color: '#999', padding: 20 }}>加载中...</div>}
        {!loading && users.length === 0 && <div style={{ textAlign: 'center', color: '#999', padding: '40px 0' }}>暂无用户</div>}

        {!loading && users.length > 0 && (
          <table className="data-table">
            <thead>
              <tr><th>ID</th><th>昵称</th><th>手机号</th><th>会员等级</th><th>余额</th><th>累计消费</th><th>订单数</th><th>注册时间</th><th>操作</th></tr>
            </thead>
            <tbody>
              {users.map(u => (
                <tr key={u.id}>
                  <td>{u.id}</td>
                  <td>{u.nickname || '-'}</td>
                  <td>{u.phone || '-'}</td>
                  <td>{fmtLevel(u.member_level)}</td>
                  <td>{fmtPrice(u.balance)}</td>
                  <td>{fmtPrice(u.total_consumption)}</td>
                  <td>{u.total_orders || 0}</td>
                  <td>{fmtTime(u.created_at)}</td>
                  <td><button className="btn-primary" onClick={() => showDetail(u.id)}>查看详情</button></td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {showModal && detail && (
        <div style={{ position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', zIndex: 100, display: 'flex', alignItems: 'center', justifyContent: 'center' }} onClick={() => setShowModal(false)}>
          <div style={{ background: '#fff', borderRadius: 12, padding: 24, width: 700, maxHeight: '90vh', overflowY: 'auto', boxShadow: '0 10px 40px rgba(0,0,0,0.2)' }} onClick={e => e.stopPropagation()}>
            <h3 style={{ marginBottom: 20 }}>用户详情</h3>
            <div style={{ marginBottom: 20, padding: 16, background: '#f5f5f5', borderRadius: 8 }}>
              <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12, fontSize: 14 }}>
                <div><b>ID:</b> {detail.user?.id}</div>
                <div><b>昵称:</b> {detail.user?.nickname || '-'}</div>
                <div><b>手机号:</b> {detail.user?.phone || '-'}</div>
                <div><b>会员等级:</b> {fmtLevel(detail.user?.member_level)}</div>
                <div><b>余额:</b> {fmtPrice(detail.user?.balance)}</div>
                <div><b>累计消费:</b> {fmtPrice(detail.user?.total_consumption)}</div>
                <div><b>订单数:</b> {detail.user?.total_orders || 0}</div>
                <div><b>注册时间:</b> {fmtTime(detail.user?.created_at)}</div>
              </div>
            </div>

            <div style={{ marginBottom: 20 }}>
              <h4 style={{ marginBottom: 10 }}>宠物 ({detail.pets?.length || 0})</h4>
              {(!detail.pets || detail.pets.length === 0) && <div style={{ color: '#999', fontSize: 14 }}>暂无宠物</div>}
              {detail.pets?.length > 0 && (
                <table className="data-table" style={{ fontSize: 13 }}>
                  <thead><tr><th>名字</th><th>品种</th><th>性别</th><th>体重</th></tr></thead>
                  <tbody>
                    {detail.pets.map(p => (
                      <tr key={p.id}><td>{p.name}</td><td>{p.breed || '-'}</td><td>{p.gender === 'male' ? '公' : p.gender === 'female' ? '母' : '-'}</td><td>{p.weight > 0 ? p.weight + 'kg' : '-'}</td></tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>

            <div style={{ marginBottom: 20 }}>
              <h4 style={{ marginBottom: 10 }}>订单 ({detail.order_total || 0})</h4>
              {(!detail.orders || detail.orders.length === 0) && <div style={{ color: '#999', fontSize: 14 }}>暂无订单</div>}
              {detail.orders?.length > 0 && (
                <table className="data-table" style={{ fontSize: 13 }}>
                  <thead><tr><th>订单号</th><th>金额</th><th>状态</th><th>时间</th></tr></thead>
                  <tbody>
                    {detail.orders.map(o => (
                      <tr key={o.id}>
                        <td>{o.order_no}</td>
                        <td>{fmtPrice(o.pay_amount)}</td>
                        <td>{o.status === 'pending' ? '待支付' : o.status === 'paid' ? '已支付' : o.status === 'completed' ? '已完成' : o.status}</td>
                        <td>{fmtTime(o.created_at)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>

            <div>
              <h4 style={{ marginBottom: 10 }}>充值记录 ({detail.recharge_records?.length || 0})</h4>
              {(!detail.recharge_records || detail.recharge_records.length === 0) && <div style={{ color: '#999', fontSize: 14 }}>暂无充值记录</div>}
              {detail.recharge_records?.length > 0 && (
                <table className="data-table" style={{ fontSize: 13 }}>
                  <thead><tr><th>金额</th><th>赠送</th><th>渠道</th><th>状态</th><th>时间</th></tr></thead>
                  <tbody>
                    {detail.recharge_records.map(r => (
                      <tr key={r.id}>
                        <td>{fmtPrice(r.amount)}</td>
                        <td>{fmtPrice(r.gifted_amount)}</td>
                        <td>{r.channel}</td>
                        <td>{r.status === 'success' ? '成功' : r.status}</td>
                        <td>{fmtTime(r.created_at)}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              )}
            </div>

            <div style={{ display: 'flex', justifyContent: 'flex-end', marginTop: 20 }}>
              <button className="btn-primary" style={{ background: '#999' }} onClick={() => setShowModal(false)}>关闭</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
