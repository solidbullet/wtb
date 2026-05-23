import React, { useState, useEffect } from 'react'
import { get, post } from '../api'

export default function ActivityManagement() {
  const [announcements, setAnnouncements] = useState([])
  const [activities, setActivities] = useState([])

  const [annTitle, setAnnTitle] = useState('')
  const [annContent, setAnnContent] = useState('')

  const [actTitle, setActTitle] = useState('')
  const [actLoc, setActLoc] = useState('')
  const [actQuota, setActQuota] = useState('')

  const loadActivity = async () => {
    try {
      const anns = await get('/api/admin/activity/announcements')
      setAnnouncements(anns || [])

      const acts = await get('/api/admin/activity/list')
      setActivities(acts || [])
    } catch (e) {
      console.error('Activity error:', e)
    }
  }

  useEffect(() => {
    loadActivity()
  }, [])

  const addAnnouncement = async () => {
    if (!annTitle || !annContent) return alert('请填写完整信息')
    try {
      await post('/api/admin/activity/admin/announcement', { title: annTitle, content: annContent, is_published: true })
      setAnnTitle('')
      setAnnContent('')
      loadActivity()
    } catch (e) { alert(e.message) }
  }

  const addActivity = async () => {
    if (!actTitle) return alert('请输入活动标题')
    const quota = parseInt(actQuota) || 20
    try {
      await post('/api/admin/activity/admin/activity', { title: actTitle, location: actLoc, quota, max_participants: quota })
      setActTitle('')
      setActLoc('')
      setActQuota('')
      loadActivity()
    } catch (e) { alert(e.message) }
  }

  const statusText = { published: '已发布', draft: '草稿', ended: '已结束', cancelled: '已取消' }

  return (
    <div id="page-activity">
      <div className="section">
        <h3>公告管理</h3>
        <div className="inline-form">
          <input type="text" placeholder="标题" style={{width: '200px'}} value={annTitle} onChange={e => setAnnTitle(e.target.value)} />
          <input type="text" placeholder="内容" style={{width: '300px'}} value={annContent} onChange={e => setAnnContent(e.target.value)} />
          <button className="btn-primary" onClick={addAnnouncement}>+ 发布公告</button>
        </div>
        <table className="data-table">
          <thead><tr><th>ID</th><th>标题</th><th>内容</th><th>状态</th></tr></thead>
          <tbody>
            {announcements.length === 0 && <tr><td colSpan="4" style={{textAlign:'center',color:'#999'}}>暂无公告</td></tr>}
            {announcements.map(a => (
              <tr key={a.id}>
                <td>{a.id}</td><td>{a.title}</td><td>{a.content}</td><td>{a.status === 1 ? '已发布' : '下架'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <div className="section">
        <h3>活动管理</h3>
        <div className="inline-form">
          <input type="text" placeholder="活动标题" style={{width: '200px'}} value={actTitle} onChange={e => setActTitle(e.target.value)} />
          <input type="text" placeholder="地点" style={{width: '150px'}} value={actLoc} onChange={e => setActLoc(e.target.value)} />
          <input type="number" placeholder="名额" style={{width: '80px'}} value={actQuota} onChange={e => setActQuota(e.target.value)} />
          <button className="btn-primary" onClick={addActivity}>+ 创建活动</button>
        </div>
        <table className="data-table">
          <thead><tr><th>ID</th><th>标题</th><th>地点</th><th>名额</th><th>已报名</th><th>状态</th></tr></thead>
          <tbody>
            {activities.length === 0 && <tr><td colSpan="6" style={{textAlign:'center',color:'#999'}}>暂无活动</td></tr>}
            {activities.map(a => (
              <tr key={a.id}>
                <td>{a.id}</td>
                <td>{a.title}</td>
                <td>{a.location || '-'}</td>
                <td>{a.max_participants === -1 ? '不限' : (a.max_participants || 0)}</td>
                <td>{a.current_participants || 0}</td>
                <td>{statusText[a.status] || a.status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
