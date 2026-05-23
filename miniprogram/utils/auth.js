const { API_BASE, isDevTools } = require('./config')

// 开发环境固定 openid，避免每次编译创建新用户
const DEV_OPENID = 'dev_openid_dev_test_199'

function wxLogin() {
  return new Promise((resolve, reject) => {
    if (isDevTools) {
      const openid = DEV_OPENID
      wx.setStorageSync('openid', openid)
      resolve({ openid, token: '', user: null })
      return
    }

    // 真机环境：先尝试云函数获取 openid
    // 如果云开发未配置或失败，自动 fallback 到 wx.login
    wx.cloud.callFunction({
      name: 'getOpenID',
      success: (res) => {
        const openid = res.result && res.result.openid
        if (openid) {
          backendLogin({ openid }, resolve, reject)
        } else {
          console.warn('云函数未返回 openid，fallback 到 wx.login')
          wxLoginFallback(resolve, reject)
        }
      },
      fail: (err) => {
        console.error('云函数调用失败，fallback 到 wx.login:', err.message || err)
        wxLoginFallback(resolve, reject)
      }
    })
  })
}

// fallback：使用 wx.login 获取 code，再调后端登录
function wxLoginFallback(resolve, reject) {
  wx.login({
    success: (res) => {
      const code = res.code || 'fb_' + Date.now()
      backendLogin({ code }, resolve, reject)
    },
    fail: reject
  })
}

// 调用后端登录/注册接口（自动创建用户）
function backendLogin(payload, resolve, reject) {
  console.log('[LOGIN] calling backend', API_BASE + '/api/user/wx-login', payload)
  wx.request({
    url: API_BASE + '/api/user/wx-login',
    method: 'POST',
    data: payload,
    header: { 'Content-Type': 'application/json' },
    success: (r) => {
      console.log('[LOGIN] response', r.statusCode, r.data)
      if (r.statusCode !== 200) {
        console.error('[LOGIN] HTTP error', r.statusCode, r.data)
        return reject(new Error('服务器错误: ' + r.statusCode))
      }
      if (!r.data || typeof r.data !== 'object') {
        console.error('[LOGIN] format error', r.data)
        return reject(new Error('服务器返回格式错误'))
      }
      if (r.data.code === 200) {
        const user = r.data.data && r.data.data.user
        const token = r.data.data && r.data.data.token
        const openid = user ? user.openid : payload.openid
        wx.setStorageSync('openid', openid)
        resolve({ openid, token: token || '', user: user || null })
      } else {
        const errMsg = (r.data && (r.data.message || r.data.msg)) || '登录失败'
        console.error('[LOGIN] failed', r.statusCode, errMsg, r.data)
        reject(new Error(errMsg))
      }
    },
    fail: (err) => {
      console.error('[LOGIN] request failed', err)
      const errMsg = (err && err.errMsg) || (err && err.message) || '网络请求失败，请检查服务器域名配置'
      reject(new Error(errMsg))
    }
  })
}

function checkLogin() {
  return !!wx.getStorageSync('openid')
}

function getOpenID() {
  return wx.getStorageSync('openid') || (isDevTools ? DEV_OPENID : '')
}

function logout() {
  wx.removeStorageSync('openid')
  wx.showToast({ title: '已退出登录', icon: 'none' })
}

module.exports = {
  wxLogin,
  checkLogin,
  getOpenID,
  logout
}
