import React, { useState, useEffect } from 'react'
import { get, post } from '../api'

export default function SeatManagement() {
  const [areas, setAreas] = useState([])
  const [seats, setSeats] = useState([])
  const [selectedArea, setSelectedArea] = useState('')

  // 新增区域
  const [areaName, setAreaName] = useState('')
  const [areaSort, setAreaSort] = useState('')

  // 新增座位
  const [seatName, setSeatName] = useState('')
  const [seatType, setSeatType] = useState('normal')
  const [seatCapacity, setSeatCapacity] = useState('4')

  // 生成结果
  const [genResult, setGenResult] = useState(null)

  const loadAreas = async () => {
    try {
      const data = await get('/api/seat/areas')
      setAreas(data || [])
      if (data && data.length > 0 && !selectedArea) {
        setSelectedArea(String(data[0].id))
      }
    } catch (e) {
      alert('加载区域失败: ' + e.message)
    }
  }

  const loadSeats = async () => {
    if (!selectedArea) return
    try {
      const data = await get('/api/seat/list?area_id=' + selectedArea)
      setSeats(data || [])
    } catch (e) {
      alert('加载座位失败: ' + e.message)
    }
  }

  useEffect(() => {
    loadAreas()
  }, [])

  useEffect(() => {
    loadSeats()
  }, [selectedArea])

  const addArea = async () => {
    if (!areaName) return alert('请输入区域名称')
    try {
      await post('/api/seat/areas', { name: areaName, sort_order: parseInt(areaSort) || 0 })
      setAreaName('')
      setAreaSort('')
      alert('区域添加成功')
      loadAreas()
    } catch (e) {
      alert(e.message)
    }
  }

  const addSeat = async () => {
    if (!seatName) return alert('请输入座位名称')
    if (!selectedArea) return alert('请先选择区域')
    try {
      await post('/api/seat/create', {
        area_id: parseInt(selectedArea),
        name: seatName,
        type: seatType,
        capacity: parseInt(seatCapacity) || 4
      })
      setSeatName('')
      setSeatType('normal')
      setSeatCapacity('4')
      alert('座位添加成功')
      loadSeats()
    } catch (e) {
      alert(e.message)
    }
  }

  const generateQrcodes = async () => {
    if (seats.length === 0) return alert('当前区域没有座位')
    try {
      const result = await post('/api/seat/qrcode/batch', { seat_ids: seats.map(s => s.id) })
      setGenResult(result || [])
      alert('二维码路径已生成，可复制下方路径到草料等工具生成微信小程序码')
      loadSeats()
    } catch (e) {
      alert(e.message)
    }
  }

  const copyPath = (path) => {
    navigator.clipboard.writeText(path).then(() => alert('已复制: ' + path))
  }

  const statusText = { available: '空闲', occupied: '占用', reserved: '预留' }

  return (
    <div id="page-seats">
      <div className="section">
        <h3>座位管理</h3>

        {/* 新增区域 */}
        <div className="card" style={{ marginBottom: 20, padding: 16, background: '#fffbf0', border: '1px solid #f0e0c0', borderRadius: 8 }}>
          <h4 style={{ margin: '0 0 12px 0' }}>➕ 新增区域</h4>
          <div className="inline-form">
            <input type="text" placeholder="区域名称，如：大厅区" style={{ width: '200px' }} value={areaName} onChange={e => setAreaName(e.target.value)} />
            <input type="number" placeholder="排序号" style={{ width: '100px' }} value={areaSort} onChange={e => setAreaSort(e.target.value)} />
            <button className="btn-primary" onClick={addArea}>添加区域</button>
          </div>
        </div>

        {/* 区域选择 & 新增座位 */}
        <div className="card" style={{ marginBottom: 20, padding: 16, background: '#fffbf0', border: '1px solid #f0e0c0', borderRadius: 8 }}>
          <h4 style={{ margin: '0 0 12px 0' }}>➕ 新增座位</h4>
          <div className="inline-form" style={{ flexWrap: 'wrap' }}>
            <select style={{ width: '160px', height: 36, borderRadius: 4, border: '1px solid #ddd' }} value={selectedArea} onChange={e => setSelectedArea(e.target.value)}>
              <option value="">选择区域</option>
              {areas.map(a => (
                <option key={a.id} value={a.id}>{a.name}</option>
              ))}
            </select>
            <input type="text" placeholder="座位名称，如：A01" style={{ width: '160px' }} value={seatName} onChange={e => setSeatName(e.target.value)} />
            <select style={{ width: '120px', height: 36, borderRadius: 4, border: '1px solid #ddd' }} value={seatType} onChange={e => setSeatType(e.target.value)}>
              <option value="normal">普通</option>
              <option value="vip">VIP</option>
              <option value="outdoor">户外</option>
            </select>
            <input type="number" placeholder="容纳人数" style={{ width: '100px' }} value={seatCapacity} onChange={e => setSeatCapacity(e.target.value)} />
            <button className="btn-primary" onClick={addSeat}>添加座位</button>
          </div>
        </div>

        {/* 座位列表 */}
        <div className="inline-form" style={{ marginBottom: 16 }}>
          <select style={{ width: '160px', height: 36, borderRadius: 4, border: '1px solid #ddd' }} value={selectedArea} onChange={e => setSelectedArea(e.target.value)}>
            <option value="">选择区域</option>
            {areas.map(a => (
              <option key={a.id} value={a.id}>{a.name}</option>
            ))}
          </select>
          <button className="btn-primary" onClick={loadSeats}>🔄 刷新</button>
          <button className="btn-secondary" onClick={generateQrcodes}>📷 批量生成小程序码路径</button>
        </div>

        {genResult && genResult.length > 0 && (
          <div className="card" style={{ marginBottom: 20, padding: 16, background: '#f0fff4', border: '1px solid #c0f0d0', borderRadius: 8 }}>
            <h4 style={{ margin: '0 0 12px 0' }}>✅ 生成结果</h4>
            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 20 }}>
              {genResult.map(r => (
                <div key={r.seat_id} style={{ textAlign: 'center', padding: 12, background: '#fff', borderRadius: 8, border: '1px solid #e0e0e0' }}>
                  <div style={{ fontWeight: 'bold', marginBottom: 8 }}>{r.seat_name}</div>
                  {r.has_image ? (
                    <img src={r.qrcode_url} alt={r.seat_name} style={{ width: 160, height: 160, objectFit: 'contain' }} />
                  ) : (
                    <div style={{ width: 160, height: 160, display: 'flex', alignItems: 'center', justifyContent: 'center', background: '#f5f5f5', color: '#999', fontSize: 12 }}>
                      未生成图片
                      <br />
                      {r.qrcode_url}
                    </div>
                  )}
                  <div style={{ marginTop: 8 }}>
                    <button className="btn-small" onClick={() => copyPath(r.wxa_path)}>复制路径</button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        <table className="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>区域</th>
              <th>座位名</th>
              <th>类型</th>
              <th>容纳人数</th>
              <th>状态</th>
              <th>小程序路径</th>
            </tr>
          </thead>
          <tbody>
            {seats.length === 0 && <tr><td colSpan="7" style={{ textAlign: 'center', color: '#999' }}>暂无座位</td></tr>}
            {seats.map(s => (
              <tr key={s.id}>
                <td>{s.id}</td>
                <td>{areas.find(a => a.id === s.area_id)?.name || '-'}</td>
                <td>{s.name}</td>
                <td>{s.type}</td>
                <td>{s.capacity}</td>
                <td>{statusText[s.status] || s.status}</td>
                <td style={{ fontFamily: 'monospace', fontSize: 12 }}>
                  {s.qrcode_url ? (
                    <span style={{ cursor: 'pointer', color: '#1890ff' }} onClick={() => copyPath(s.qrcode_url)} title="点击复制">
                      {s.qrcode_url}
                    </span>
                  ) : '-'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
