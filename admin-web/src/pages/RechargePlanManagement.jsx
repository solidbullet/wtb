import React, { useState, useEffect } from 'react'
import { get, post, put, del } from '../api'

export default function RechargePlanManagement() {
  const [plans, setPlans] = useState([])
  const [name, setName] = useState('')
  const [amount, setAmount] = useState('')
  const [finalAmount, setFinalAmount] = useState('')
  const [giftAmount, setGiftAmount] = useState('')
  const [sortOrder, setSortOrder] = useState('')
  const [editingId, setEditingId] = useState(null)

  const loadPlans = async () => {
    try {
      const data = await get('/api/admin/pricing/recharge-plans')
      setPlans(data || [])
    } catch (e) {
      console.error('Recharge plans error:', e)
    }
  }

  useEffect(() => {
    loadPlans()
  }, [])

  const resetForm = () => {
    setName('')
    setAmount('')
    setFinalAmount('')
    setGiftAmount('')
    setSortOrder('')
    setEditingId(null)
  }

  const savePlan = async () => {
    if (!name || !amount || !finalAmount) return alert('请填写完整信息')
    const payload = {
      name,
      amount: parseInt(amount) || 0,
      final_amount: parseInt(finalAmount) || 0,
      gift_amount: parseInt(giftAmount) || 0,
      sort_order: parseInt(sortOrder) || 0,
      status: 1
    }
    try {
      if (editingId) {
        await put('/api/admin/pricing/admin/recharge-plan/' + editingId, payload)
      } else {
        await post('/api/admin/pricing/admin/recharge-plan', payload)
      }
      resetForm()
      loadPlans()
    } catch (e) { alert(e.message) }
  }

  const editPlan = (plan) => {
    setEditingId(plan.id)
    setName(plan.name)
    setAmount(plan.amount)
    setFinalAmount(plan.final_amount)
    setGiftAmount(plan.gift_amount || '')
    setSortOrder(plan.sort_order || '')
  }

  const removePlan = async (id) => {
    if (!confirm('确定删除该充值方案吗？')) return
    try {
      await del('/api/admin/pricing/admin/recharge-plan/' + id)
      loadPlans()
    } catch (e) { alert(e.message) }
  }

  return (
    <div id="page-recharge">
      <div className="section">
        <h3>充值方案</h3>
        <div className="inline-form">
          <input type="text" placeholder="方案名称" style={{width: '150px'}} value={name} onChange={e => setName(e.target.value)} />
          <input type="number" placeholder="充值金额(分)" style={{width: '120px'}} value={amount} onChange={e => setAmount(e.target.value)} />
          <input type="number" placeholder="到账金额(分)" style={{width: '120px'}} value={finalAmount} onChange={e => setFinalAmount(e.target.value)} />
          <input type="number" placeholder="赠送金额(分)" style={{width: '120px'}} value={giftAmount} onChange={e => setGiftAmount(e.target.value)} />
          <input type="number" placeholder="排序" style={{width: '80px'}} value={sortOrder} onChange={e => setSortOrder(e.target.value)} />
          <button className="btn-primary" onClick={savePlan}>{editingId ? '更新' : '+ 新增方案'}</button>
          {editingId && <button className="btn-secondary" onClick={resetForm}>取消</button>}
        </div>
        <table className="data-table">
          <thead><tr><th>ID</th><th>名称</th><th>充值金额</th><th>到账金额</th><th>赠送</th><th>排序</th><th>状态</th><th>操作</th></tr></thead>
          <tbody>
            {plans.length === 0 && <tr><td colSpan="8" style={{textAlign:'center',color:'#999'}}>暂无充值方案</td></tr>}
            {plans.map(p => (
              <tr key={p.id}>
                <td>{p.id}</td>
                <td>{p.name}</td>
                <td>¥{(p.amount/100).toFixed(2)}</td>
                <td>¥{(p.final_amount/100).toFixed(2)}</td>
                <td>¥{(p.gift_amount/100).toFixed(2)}</td>
                <td>{p.sort_order}</td>
                <td>{p.status === 1 ? '上架' : '下架'}</td>
                <td>
                  <button className="btn-small" onClick={() => editPlan(p)}>编辑</button>
                  <button className="btn-small btn-danger" onClick={() => removePlan(p.id)}>删除</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
