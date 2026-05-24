export const API_BASE = ''

export function getToken() {
  return sessionStorage.getItem('admin_token') || ''
}

export function setToken(token) {
  sessionStorage.setItem('admin_token', token)
}

export function clearToken() {
  sessionStorage.removeItem('admin_token')
}

export async function api(method, path, body) {
  const token = getToken()
  const opts = { method, headers: { 'Content-Type': 'application/json', 'Authorization': token ? 'Bearer ' + token : '' } }
  if (body) opts.body = JSON.stringify(body)
  const res = await fetch(API_BASE + path, opts)
  const data = await res.json().catch(() => ({}))
  if (data.code !== 200) {
    if (data.code === 40101 || data.code === 40102) {
      clearToken()
      window.location.href = '/'
    }
    throw new Error(data.message || '请求失败')
  }
  return data.data
}

export const get = (path) => api('GET', path)
export const post = (path, body) => api('POST', path, body)
export const put = (path, body) => api('PUT', path, body)
export const del = (path) => api('DELETE', path)
