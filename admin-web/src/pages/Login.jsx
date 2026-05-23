import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { post, setToken, getToken } from '../api'
import './Login.css'

export default function Login() {
  const [user, setUser] = useState('admin')
  const [pass, setPass] = useState('admin123')
  const [error, setError] = useState('')
  const navigate = useNavigate()

  useEffect(() => {
    if (getToken()) {
      navigate('/dashboard')
    }
  }, [navigate])

  const doLogin = async () => {
    try {
      const data = await post('/admin/login', { username: user, password: pass })
      setToken(data.token)
      navigate('/dashboard')
    } catch (e) {
      setError(e.message)
    }
  }

  return (
    <div id="login-page">
      <div className="login-box">
        <h1>汪托帮</h1>
        <p className="subtitle">后台管理系统</p>
        <div className="input-group">
          <input 
            type="text" 
            placeholder="用户名" 
            value={user} 
            onChange={e => setUser(e.target.value)} 
          />
        </div>
        <div className="input-group">
          <input 
            type="password" 
            placeholder="密码" 
            value={pass} 
            onChange={e => setPass(e.target.value)} 
          />
        </div>
        <button onClick={doLogin}>进入系统</button>
        <p className="error">{error}</p>
      </div>
    </div>
  )
}
