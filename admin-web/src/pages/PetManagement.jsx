import React, { useState, useEffect } from 'react'
import { get, put, del } from '../api'

export default function PetManagement() {
  const [pets, setPets] = useState([])
  const [searchName, setSearchName] = useState('')
  const [searchPhone, setSearchPhone] = useState('')

  const [editingId, setEditingId] = useState(null)
  const [name, setName] = useState('')
  const [breed, setBreed] = useState('')
  const [gender, setGender] = useState('')
  const [weight, setWeight] = useState('')
  const [birthday, setBirthday] = useState('')
  const [photoUrl, setPhotoUrl] = useState('')
  const [notes, setNotes] = useState('')

  const loadPets = async () => {
    try {
      const params = new URLSearchParams()
      if (searchName) params.append('name', searchName)
      if (searchPhone) params.append('phone', searchPhone)
      const data = await get('/api/admin/pets?' + params.toString())
      setPets(data || [])
    } catch (e) {
      console.error('Pets error:', e)
    }
  }

  useEffect(() => {
    loadPets()
  }, [])

  const resetForm = () => {
    setEditingId(null)
    setName('')
    setBreed('')
    setGender('')
    setWeight('')
    setBirthday('')
    setPhotoUrl('')
    setNotes('')
  }

  const editPet = (pet) => {
    setEditingId(pet.id)
    setName(pet.name || '')
    setBreed(pet.breed || '')
    setGender(pet.gender || '')
    setWeight(pet.weight ? String(pet.weight) : '')
    setBirthday(pet.birthday || '')
    setPhotoUrl(pet.photo_url || '')
    setNotes(pet.notes || '')
    // scroll to form
    window.scrollTo({ top: 0, behavior: 'smooth' })
  }

  const savePet = async () => {
    if (!name) return alert('请输入宠物名字')
    const payload = {
      name: name.trim(),
      breed,
      gender,
      weight: parseFloat(weight) || 0,
      birthday,
      photo_url: photoUrl,
      notes
    }
    try {
      await put('/api/admin/pets/' + editingId, payload)
      resetForm()
      loadPets()
    } catch (e) { alert(e.message) }
  }

  const removePet = async (id, petName) => {
    if (!confirm(`确定删除宠物「${petName}」吗？`)) return
    try {
      await del('/api/admin/pets/' + id)
      loadPets()
    } catch (e) { alert(e.message) }
  }

  const genderText = { male: '♂ 公', female: '♀ 母', '': '-' }

  return (
    <div id="page-pets">
      <div className="section">
        <h3>宠物管理</h3>

        {/* 搜索 */}
        <div className="inline-form" style={{ marginBottom: 16 }}>
          <input type="text" placeholder="宠物名字" style={{ width: '160px' }} value={searchName} onChange={e => setSearchName(e.target.value)} />
          <input type="text" placeholder="主人手机号" style={{ width: '160px' }} value={searchPhone} onChange={e => setSearchPhone(e.target.value)} />
          <button className="btn-primary" onClick={loadPets}>🔍 查询</button>
          <button className="btn-secondary" onClick={() => { setSearchName(''); setSearchPhone(''); loadPets() }}>重置</button>
        </div>

        {/* 编辑表单 */}
        {editingId && (
          <div className="card" style={{ marginBottom: 20, padding: 16, background: '#fffbf0', border: '1px solid #f0e0c0', borderRadius: 8 }}>
            <h4 style={{ margin: '0 0 12px 0' }}>✏️ 编辑宠物信息</h4>
            <div className="inline-form" style={{ flexWrap: 'wrap' }}>
              <input type="text" placeholder="宠物名字" style={{ width: '140px' }} value={name} onChange={e => setName(e.target.value)} />
              <input type="text" placeholder="品种" style={{ width: '120px' }} value={breed} onChange={e => setBreed(e.target.value)} />
              <select style={{ width: '100px', height: 36, borderRadius: 4, border: '1px solid #ddd' }} value={gender} onChange={e => setGender(e.target.value)}>
                <option value="">性别</option>
                <option value="male">♂ 公</option>
                <option value="female">♀ 母</option>
              </select>
              <input type="number" placeholder="体重(kg)" style={{ width: '100px' }} value={weight} onChange={e => setWeight(e.target.value)} />
              <input type="date" placeholder="生日" style={{ width: '140px', height: 36, borderRadius: 4, border: '1px solid #ddd' }} value={birthday} onChange={e => setBirthday(e.target.value)} />
              <input type="text" placeholder="照片URL" style={{ width: '200px' }} value={photoUrl} onChange={e => setPhotoUrl(e.target.value)} />
              <button className="btn-primary" onClick={savePet}>保存</button>
              <button className="btn-secondary" onClick={resetForm}>取消</button>
            </div>
            <textarea placeholder="寄养注意事项" style={{ width: '100%', marginTop: 8, minHeight: 60, borderRadius: 4, border: '1px solid #ddd', padding: 8, boxSizing: 'border-box' }} value={notes} onChange={e => setNotes(e.target.value)} />
          </div>
        )}

        <table className="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>照片</th>
              <th>名字</th>
              <th>品种</th>
              <th>性别</th>
              <th>体重</th>
              <th>生日</th>
              <th>主人</th>
              <th>手机号</th>
              <th>注意事项</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {pets.length === 0 && <tr><td colSpan="11" style={{ textAlign: 'center', color: '#999' }}>暂无宠物</td></tr>}
            {pets.map(p => (
              <tr key={p.id}>
                <td>{p.id}</td>
                <td>
                  {p.photo_url ? (
                    <img src={p.photo_url} alt="" style={{ width: 48, height: 48, objectFit: 'cover', borderRadius: 4 }} />
                  ) : (
                    <span style={{ color: '#bbb' }}>-</span>
                  )}
                </td>
                <td>{p.name}</td>
                <td>{p.breed || '-'}</td>
                <td>{genderText[p.gender] || '-'}</td>
                <td>{p.weight ? p.weight + ' kg' : '-'}</td>
                <td>{p.birthday || '-'}</td>
                <td>{p.owner_name || '-'}</td>
                <td>{p.owner_phone || '-'}</td>
                <td style={{ maxWidth: 200, whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis' }} title={p.notes}>
                  {p.notes || '-'}
                </td>
                <td>
                  <button className="btn-small" onClick={() => editPet(p)}>编辑</button>
                  <button className="btn-small btn-danger" onClick={() => removePet(p.id, p.name)}>删除</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
