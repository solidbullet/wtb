const auth = require('./auth')
const { API_BASE } = require('./config')

// 开发阶段使用 HTTP 直连单体服务，上线后可切换为云函数
const USE_CLOUD = false

function httpRequest(options) {
  return new Promise((resolve, reject) => {
    const openid = auth.getOpenID()
    wx.request({
      url: API_BASE + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: {
        'Content-Type': 'application/json',
        'X-OpenID': openid || ''
      },
      success: (res) => {
        // HTTP 状态码检查
        if (res.statusCode === 401) {
          return reject(new Error('unauthorized'))
        }
        if (res.statusCode !== 200) {
          console.error('[API HTTP ERROR]', options.url, res.statusCode, res.data)
          return reject(new Error('服务器错误: ' + res.statusCode))
        }
        // 后端返回的 JSON 中没有 code 字段（可能收到了非 JSON 响应）
        if (!res.data || typeof res.data !== 'object') {
          console.error('[API FORMAT ERROR]', options.url, res.data)
          return reject(new Error('服务器返回格式错误'))
        }
        if (res.data.code !== 200) {
          const errMsg = res.data.message || res.data.msg || '请求失败'
          console.error('[API ERROR]', options.url, res.data.code, errMsg)
          return reject(new Error(errMsg))
        }
        resolve(res.data)
      },
      fail: (err) => {
        const errMsg = (err && err.errMsg) || (err && err.message) || '网络请求失败'
        console.error('[API FAIL]', options.url, errMsg, err)
        reject(new Error(errMsg))
      }
    })
  })
}

function cloudRequest(options) {
  return new Promise((resolve, reject) => {
    wx.cloud.callFunction({
      name: 'proxy',
      data: {
        path: options.url,
        method: options.method || 'GET',
        data: options.data || {}
      },
      success: (res) => {
        const data = res.result
        if (data.code === 401) {
          return reject(new Error('unauthorized'))
        }
        if (data.code !== 200) {
          const errMsg = data.message || data.msg || '请求失败'
          console.error('[CLOUD API ERROR]', options.url, data.code, errMsg)
          return reject(new Error(errMsg))
        }
        resolve(data)
      },
      fail: (err) => reject(new Error(err.message || '云函数调用失败'))
    })
  })
}

function request(options) {
  return USE_CLOUD ? cloudRequest(options) : httpRequest(options)
}

module.exports = {
  API_BASE,
  get: (url, data) => request({ url, method: 'GET', data }),
  post: (url, data) => request({ url, method: 'POST', data }),
  put: (url, data) => request({ url, method: 'PUT', data }),
  del: (url, data) => request({ url, method: 'DELETE', data })
}
