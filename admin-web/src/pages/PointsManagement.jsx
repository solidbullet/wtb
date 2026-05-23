import React, { useState, useEffect } from 'react'
import { get, post } from '../api'

export default function PointsManagement() {
  const [goods, setGoods] = useState([])
  const [name, setName] = useState('')
  const [points, setPoints] = useState('')
  const [stock, setStock] = useState('')

  const loadPoints = async () => {
    try {
      const data = await get('/api/admin/points/goods')
      setGoods(data || [])
    } catch (e) {
      console.error('Points error:', e)
    }
  }

  useEffect(() => {
    loadPoints()
  }, [])

  const addGoods = async () => {
    if (!name) return alert('请输入商品名称')
    try {
      await post('/api/admin/points/admin/goods', { 
        name, 
        points_price: parseInt(points) || 0, 
        stock: parseInt(stock) || 0, 
        type: 'physical', 
        status: 1 
      })
      setName('')
      setPoints('')
      setStock('')
      loadPoints()
    } catch (e) { alert(e.message) }
  }

  return (
    <div id="page-points">
      <div className="section">
        <h3>积分商品</h3>
        <div className="inline-form">
          <input type="text" placeholder="商品名称" style={{width: '200px'}} value={name} onChange={e => setName(e.target.value)} />
          <input type="number" placeholder="所需积分" style={{width: '100px'}} value={points} onChange={e => setPoints(e.target.value)} />
          <input type="number" placeholder="库存" style={{width: '80px'}} value={stock} onChange={e => setStock(e.target.value)} />
          <button className="btn-primary" onClick={addGoods}>+ 新增商品</button>
        </div>
        <table className="data-table">
          <thead><tr><th>ID</th><th>名称</th><th>积分</th><th>库存</th><th>状态</th></tr></thead>
          <tbody>
            {goods.length === 0 && <tr><td colSpan="5" style={{textAlign:'center',color:'#999'}}>暂无商品</td></tr>}
            {goods.map(g => (
              <tr key={g.id}>
                <td>{g.id}</td>
                <td>{g.name}</td>
                <td>{g.points_price}</td>
                <td>{g.stock}</td>
                <td>{g.status === 1 ? '上架' : '下架'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
