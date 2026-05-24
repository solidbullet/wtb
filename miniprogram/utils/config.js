// ========================================
// 环境配置
// ========================================

// 服务器域名（已配置 HTTPS）
const DOMAIN = 'wtb.anzhitek.com'

// 开发工具自动用 localhost；真机（preview / 体验版 / 上线）自动用域名
const isDevTools = wx.getSystemInfoSync().platform === 'devtools'
const API_BASE = isDevTools ? 'http://localhost:8080' : `https://${DOMAIN}`

// 开发环境跳过扫码，直接使用默认桌号，方便调试
const SKIP_SCAN = isDevTools

module.exports = {
  DOMAIN,
  API_BASE,
  isDevTools,
  SKIP_SCAN
}
