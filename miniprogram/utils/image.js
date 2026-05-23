const { API_BASE } = require('./config')

/**
 * 处理图片URL，将相对路径转为完整网络URL
 * @param {string} url - 图片路径，可能是相对路径（如 /images/xxx.png）或网络URL
 * @returns {string} 完整的图片URL
 */
function resolveImageUrl(url) {
  if (!url) return ''
  // 已经是完整网络URL
  if (url.startsWith('http://') || url.startsWith('https://')) {
    return url
  }
  // 相对路径，拼接服务器域名
  // 确保路径以 / 开头
  const path = url.startsWith('/') ? url : '/' + url
  return API_BASE + path
}

module.exports = {
  resolveImageUrl
}
