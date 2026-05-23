import React, { useState, useEffect } from 'react'
import { get, post, del, API_BASE } from '../api'

export default function MenuManagement() {
  const [categories, setCategories] = useState([])
  const [dishes, setDishes] = useState([])

  // 分类表单
  const [catName, setCatName] = useState('')
  const [catSort, setCatSort] = useState('')

  // 菜品模态框
  const [showModal, setShowModal] = useState(false)
  const [editingId, setEditingId] = useState(null)
  const [dishName, setDishName] = useState('')
  const [dishCat, setDishCat] = useState('')
  const [dishSubtitle, setDishSubtitle] = useState('')
  const [dishDesc, setDishDesc] = useState('')
  const [dishPrice, setDishPrice] = useState('')
  const [dishMemberPrice, setDishMemberPrice] = useState('')
  const [dishStock, setDishStock] = useState('')
  const [dishTags, setDishTags] = useState('')
  const [dishStatus, setDishStatus] = useState(1)
  const [dishImage, setDishImage] = useState('')
  const [previewUrl, setPreviewUrl] = useState('')
  const [uploading, setUploading] = useState(false)

  const loadMenu = async () => {
    try {
      const cats = await get('/api/admin/menu/categories')
      setCategories(cats || [])
      const d = await get('/api/admin/menu/dishes')
      setDishes(d?.list || [])
    } catch (e) {
      alert('加载菜单失败: ' + e.message)
    }
  }

  useEffect(() => {
    loadMenu()
  }, [])

  const catMap = categories.reduce((acc, c) => ({ ...acc, [c.id]: c.name }), {})

  // 分类增删
  const addCategory = async () => {
    if (!catName) return alert('请输入分类名称')
    try {
      await post('/api/admin/menu/admin/category', { name: catName, sort_order: parseInt(catSort) || 0 })
      setCatName('')
      setCatSort('')
      alert('分类添加成功！')
      loadMenu()
    } catch (e) { alert(e.message) }
  }

  const deleteCategory = async (id) => {
    if (!window.confirm('确定删除该分类？')) return
    try {
      await del('/api/admin/menu/admin/category/' + id)
      loadMenu()
    } catch (e) { alert(e.message) }
  }

  // 图片上传
  const handleFileChange = async (e) => {
    const file = e.target.files[0]
    if (!file) return
    setUploading(true)
    const url = URL.createObjectURL(file)
    setPreviewUrl(url)
    try {
      const token = localStorage.getItem('admin_token') || ''
      const form = new FormData()
      form.append('file', file)
      const res = await fetch(API_BASE + '/api/menu/admin/upload', {
        method: 'POST',
        headers: { 'Authorization': token ? 'Bearer ' + token : '' },
        body: form
      })
      const data = await res.json().catch(() => ({}))
      if (data.code === 200) {
        setDishImage(data.data.url)
      } else {
        alert(data.message || '上传失败')
      }
    } catch (err) {
      alert('上传失败: ' + err.message)
    } finally {
      setUploading(false)
    }
  }

  // 打开新增模态框
  const openAdd = () => {
    setEditingId(null)
    setDishName('')
    setDishCat('')
    setDishSubtitle('')
    setDishDesc('')
    setDishPrice('')
    setDishMemberPrice('')
    setDishStock('')
    setDishTags('')
    setDishStatus(1)
    setDishImage('')
    setPreviewUrl('')
    setShowModal(true)
  }

  // 提交菜品
  const submitDish = async () => {
    if (!dishName || !dishCat) return alert('请填写菜品名称和分类')
    const prices = []
    if (dishPrice) prices.push({ price_type: 'standard', price: parseInt(dishPrice) })
    if (dishMemberPrice) prices.push({ price_type: 'member', price: parseInt(dishMemberPrice) })

    const payload = {
      name: dishName,
      category_id: parseInt(dishCat),
      subtitle: dishSubtitle,
      description: dishDesc,
      images: dishImage,
      tags: dishTags,
      status: parseInt(dishStatus),
      prices,
      stock: parseInt(dishStock || -1)
    }

    try {
      if (editingId) {
        await post('/api/admin/menu/admin/dish/' + editingId, payload)
      } else {
        await post('/api/admin/menu/admin/dish', payload)
      }
      setShowModal(false)
      loadMenu()
    } catch (e) { alert(e.message) }
  }

  const deleteDish = async (id) => {
    if (!window.confirm('确定删除该菜品？')) return
    try {
      await del('/api/admin/menu/admin/dish/' + id)
      loadMenu()
    } catch (e) { alert(e.message) }
  }

  const fmtPrice = (p) => p ? '¥' + (p / 100).toFixed(2) : '-'

  return (
    <div id="page-menu">
      {/* 分类管理 */}
      <div className="section">
        <h3>菜品分类</h3>
        <div className="inline-form">
          <input type="text" placeholder="分类名称" value={catName} onChange={e => setCatName(e.target.value)} onKeyDown={e => e.key === 'Enter' && addCategory()} />
          <input type="number" placeholder="排序" style={{ width: '80px' }} value={catSort} onChange={e => setCatSort(e.target.value)} onKeyDown={e => e.key === 'Enter' && addCategory()} />
          <button className="btn-primary" onClick={addCategory}>+ 新增分类</button>
        </div>
        <table className="data-table">
          <thead><tr><th>ID</th><th>名称</th><th>排序</th><th>操作</th></tr></thead>
          <tbody>
            {categories.length === 0 && <tr><td colSpan="4" style={{ textAlign: 'center', color: '#999' }}>暂无数据</td></tr>}
            {categories.map(c => (
              <tr key={c.id}>
                <td>{c.id}</td><td>{c.name}</td><td>{c.sort_order}</td>
                <td><button className="btn-danger" onClick={() => deleteCategory(c.id)}>删除</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* 菜品列表 */}
      <div className="section">
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
          <h3>菜品列表</h3>
          <button className="btn-primary" onClick={openAdd}>+ 新增菜品</button>
        </div>
        <table className="data-table">
          <thead>
            <tr>
              <th>图片</th><th>ID</th><th>名称</th><th>分类</th><th>标准价</th><th>会员价</th>
              <th>库存</th><th>状态</th><th>标签</th><th>操作</th>
            </tr>
          </thead>
          <tbody>
            {dishes.length === 0 && <tr><td colSpan="10" style={{ textAlign: 'center', color: '#999' }}>暂无数据</td></tr>}
            {dishes.map(d => (
              <tr key={d.id}>
                <td>
                  {d.images ? (
                    <img src={API_BASE + d.images} alt="" style={{ width: 60, height: 45, objectFit: 'cover', borderRadius: 6 }} />
                  ) : (
                    <div style={{ width: 60, height: 45, background: '#eee', borderRadius: 6, display: 'flex', alignItems: 'center', justifyContent: 'center', fontSize: 12, color: '#999' }}>无图</div>
                  )}
                </td>
                <td>{d.id}</td>
                <td>
                  <div>{d.name}</div>
                  {d.subtitle && <div style={{ fontSize: 12, color: '#999' }}>{d.subtitle}</div>}
                </td>
                <td>{catMap[d.category_id] || '-'}</td>
                <td>{fmtPrice(d.price)}</td>
                <td>{fmtPrice(d.member_price)}</td>
                <td>{d.stock === -1 ? '不限' : d.stock}</td>
                <td>{d.status === 1 ? <span style={{ color: '#52c41a' }}>上架</span> : <span style={{ color: '#999' }}>下架</span>}</td>
                <td>{d.tags || '-'}</td>
                <td><button className="btn-danger" onClick={() => deleteDish(d.id)}>删除</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* 新增/编辑模态框 */}
      {showModal && (
        <div style={{
          position: 'fixed', inset: 0, background: 'rgba(0,0,0,0.5)', zIndex: 100,
          display: 'flex', alignItems: 'center', justifyContent: 'center'
        }} onClick={() => setShowModal(false)}>
          <div style={{
            background: '#fff', borderRadius: 12, padding: 24, width: 520, maxHeight: '90vh',
            overflowY: 'auto', boxShadow: '0 10px 40px rgba(0,0,0,0.2)'
          }} onClick={e => e.stopPropagation()}>
            <h3 style={{ marginBottom: 20 }}>{editingId ? '编辑菜品' : '新增菜品'}</h3>

            <div style={{ marginBottom: 12 }}>
              <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>菜品名称 *</label>
              <input type="text" style={{ width: '100%', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6 }} value={dishName} onChange={e => setDishName(e.target.value)} />
            </div>

            <div style={{ marginBottom: 12 }}>
              <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>分类 *</label>
              <select style={{ width: '100%', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6 }} value={dishCat} onChange={e => setDishCat(e.target.value)}>
                <option value="">选择分类</option>
                {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
              </select>
            </div>

            <div style={{ marginBottom: 12 }}>
              <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>副标题</label>
              <input type="text" style={{ width: '100%', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6 }} value={dishSubtitle} onChange={e => setDishSubtitle(e.target.value)} placeholder="如：招牌慢炖红烧肉" />
            </div>

            <div style={{ marginBottom: 12 }}>
              <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>描述</label>
              <textarea style={{ width: '100%', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6, minHeight: 60 }} value={dishDesc} onChange={e => setDishDesc(e.target.value)} placeholder="菜品详细介绍..." />
            </div>

            <div style={{ marginBottom: 12 }}>
              <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>菜品图片</label>
              <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
                {(previewUrl || dishImage) && (
                  <img src={previewUrl || (API_BASE + dishImage)} alt="preview" style={{ width: 100, height: 75, objectFit: 'cover', borderRadius: 6, border: '1px solid #eee' }} />
                )}
                <div>
                  <input type="file" accept="image/*" onChange={handleFileChange} />
                  {uploading && <div style={{ fontSize: 12, color: '#667eea', marginTop: 4 }}>上传中...</div>}
                </div>
              </div>
            </div>

            <div style={{ display: 'flex', gap: 12, marginBottom: 12 }}>
              <div style={{ flex: 1 }}>
                <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>标准价（分）*</label>
                <input type="number" style={{ width: '100%', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6 }} value={dishPrice} onChange={e => setDishPrice(e.target.value)} placeholder="如 4500 = ¥45.00" />
              </div>
              <div style={{ flex: 1 }}>
                <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>会员价（分）</label>
                <input type="number" style={{ width: '100%', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6 }} value={dishMemberPrice} onChange={e => setDishMemberPrice(e.target.value)} placeholder="如 3800 = ¥38.00" />
              </div>
            </div>

            <div style={{ display: 'flex', gap: 12, marginBottom: 12 }}>
              <div style={{ flex: 1 }}>
                <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>每日库存</label>
                <input type="number" style={{ width: '100%', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6 }} value={dishStock} onChange={e => setDishStock(e.target.value)} placeholder="-1 表示不限" />
              </div>
              <div style={{ flex: 1 }}>
                <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>状态</label>
                <select style={{ width: '100%', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6 }} value={dishStatus} onChange={e => setDishStatus(e.target.value)}>
                  <option value={1}>上架</option>
                  <option value={0}>下架</option>
                </select>
              </div>
            </div>

            <div style={{ marginBottom: 20 }}>
              <label style={{ display: 'block', fontSize: 13, color: '#666', marginBottom: 6 }}>标签</label>
              <input type="text" style={{ width: '100%', padding: '8px 12px', border: '1px solid #ddd', borderRadius: 6 }} value={dishTags} onChange={e => setDishTags(e.target.value)} placeholder="推荐,热门,新品（逗号分隔）" />
            </div>

            <div style={{ display: 'flex', gap: 12, justifyContent: 'flex-end' }}>
              <button className="btn-primary" style={{ background: '#999' }} onClick={() => setShowModal(false)}>取消</button>
              <button className="btn-primary" onClick={submitDish}>保存</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
