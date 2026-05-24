const cloud = require('wx-server-sdk')
const axios = require('axios')

cloud.init()

// 线上环境指向 HTTPS 域名
const API_BASE = 'https://wtb.anzhitek.com'

exports.main = async (event, context) => {
  const { path, method = 'GET', data = {} } = event
  const { OPENID } = cloud.getWXContext()

  if (!OPENID) {
    return { code: 401, message: '无法获取微信用户信息，请重新登录' }
  }

  try {
    const res = await axios({
      url: `${API_BASE}${path}`,
      method: method.toLowerCase(),
      data: method === 'GET' ? undefined : data,
      params: method === 'GET' ? data : undefined,
      headers: {
        'Content-Type': 'application/json',
        'X-OpenID': OPENID
      },
      timeout: 15000
    })
    return res.data
  } catch (err) {
    console.error('proxy error:', err.message, err.response?.status)
    if (err.response && err.response.data) {
      return err.response.data
    }
    return { code: 500, message: err.message || '后端服务连接失败' }
  }
}
