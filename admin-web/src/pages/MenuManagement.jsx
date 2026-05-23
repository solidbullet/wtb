import React, { useState, useEffect } from 'react'
import { get, post, del } from '../api'

export default function MenuManagement() {
  const [categories, setCategories] = useState([])
  const [dishes, setDishes] = useState([])
  const [catName, setCatName] = useState('')
  const [catSort, setCatSort] = useState('')
  const [dishName, setDishName] = useState('')
  const [dishCat, setDishCat] = useState('')
  const [dishPrice, setDishPrice] = useState('')
  const [dishStock, setDishStock] = useState('')
  const [dishTags, setDishTags] = useState('')

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

  const addDish = async () => {
    if (!dishName || !dishCat) return alert('请填写完整信息')
    try {
      await post('/api/admin/menu/admin/dish', {
        name: dishName,
        category_id: parseInt(dishCat),
        price: parseInt(dishPrice) || 0,
        stock: parseInt(dishStock) || 0,
        tags: dishTags,
        status: 1
      })
      setDishName('')
      setDishPrice('')
      setDishStock('')
      setDishTags('')
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

  const catMap = categories.reduce((acc, c) => ({ ...acc, [c.id]: c.name }), {})

  return (
    <div id="page-menu">
      <div className="section">
        <h3>菜品分类</h3>
        <div className="inline-form">
          <input 
            type="text" 
            placeholder="分类名称" 
            value={catName} 
            onChange={e => setCatName(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && addCategory()} 
          />
          <input 
            type="number" 
            placeholder="排序" 
            style={{width: '80px'}} 
            value={catSort} 
            onChange={e => setCatSort(e.target.value)}
            onKeyDown={e => e.key === 'Enter' && addCategory()} 
          />
          <button className="btn-primary" onClick={addCategory}>+ 新增分类</button>
        </div>
        <table className="data-table">
          <thead><tr><th>ID</th><th>名称</th><th>排序</th><th>操作</th></tr></thead>
          <tbody>
            {categories.length === 0 && <tr><td colSpan="4" style={{textAlign:'center',color:'#999'}}>暂无数据</td></tr>}
            {categories.map(c => (
              <tr key={c.id}>
                <td>{c.id}</td><td>{c.name}</td><td>{c.sort_order}</td>
                <td><button className="btn-danger" onClick={() => deleteCategory(c.id)}>删除</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="section">
        <h3>菜品列表</h3>
        <div className="inline-form">
          <input type="text" placeholder="菜品名称" value={dishName} onChange={e => setDishName(e.target.value)} />
          <select value={dishCat} onChange={e => setDishCat(e.target.value)}>
            <option value="">选择分类</option>
            {categories.map(c => <option key={c.id} value={c.id}>{c.name}</option>)}
          </select>
          <input type="number" placeholder="价格（分）" style={{width: '120px'}} value={dishPrice} onChange={e => setDishPrice(e.target.value)} />
          <input type="number" placeholder="库存" style={{width: '80px'}} value={dishStock} onChange={e => setDishStock(e.target.value)} />
          <input type="text" placeholder="标签" style={{width: '150px'}} value={dishTags} onChange={e => setDishTags(e.target.value)} />
          <button className="btn-primary" onClick={addDish}>+ 新增菜品</button>
        </div>
        <table className="data-table">
          <thead><tr><th>ID</th><th>名称</th><th>分类</th><th>价格</th><th>库存</th><th>标签</th><th>操作</th></tr></thead>
          <tbody>
            {dishes.length === 0 && <tr><td colSpan="7" style={{textAlign:'center',color:'#999'}}>暂无数据</td></tr>}
            {dishes.map(d => (
              <tr key={d.id}>
                <td>{d.id}</td>
                <td>{d.name}</td>
                <td>{catMap[d.category_id] || '-'}</td>
                <td>¥{(d.price || 0) / 100}</td>
                <td>{d.stock === -1 ? '不限' : (d.stock ?? '-')}</td>
                <td>{d.tags || '-'}</td>
                <td><button className="btn-danger" onClick={() => deleteDish(d.id)}>删除</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
