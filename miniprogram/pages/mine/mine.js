const api = require('../../utils/api')
const auth = require('../../utils/auth')
const { resolveImageUrl } = require('../../utils/image')

// 默认头像
const DEFAULT_AVATAR = resolveImageUrl('/images/avatar.png')

Page({
  data: { user: {}, points: 0, memberInfo: {}, loading: true, errorMsg: '', avatarUrl: DEFAULT_AVATAR },
  onShow() {
    this.loadProfile()
    this.loadMemberInfo()
  },
  async loadProfile() {
    this.setData({ loading: true, errorMsg: '' })
    try {
      const [userRes, pointsRes] = await Promise.all([
        api.get('/api/user/profile'),
        api.get('/api/points/account').catch(() => ({ data: {} }))
      ])
      const user = userRes.data || {}
      const account = pointsRes.data || {}
      const available = (account.total_points || 0) - (account.used_points || 0) - (account.frozen_points || 0)
      this.setData({
        user: user,
        avatarUrl: resolveImageUrl(user.avatar_url) || DEFAULT_AVATAR,
        points: available,
        loading: false,
        errorMsg: ''
      })
    } catch (e) {
      console.log('load profile error', e)
      const msg = e.message === 'unauthorized' ? '登录已过期' : '加载失败'
      this.setData({ loading: false, errorMsg: msg })
    }
  },
  async loadMemberInfo() {
    try {
      const res = await api.get('/api/user/member-info')
      this.setData({ memberInfo: res.data || {} })
    } catch (e) {
      console.log('load member info error', e)
    }
  },
  // 手动重新登录
  async doRelogin() {
    wx.showLoading({ title: '登录中...' })
    try {
      wx.removeStorageSync('token')
      const data = await auth.wxLogin()
      wx.hideLoading()
      wx.showToast({ title: '登录成功', icon: 'success' })
      console.log('relogin success', data)
      this.setData({ errorMsg: '' })
      this.loadProfile()
      this.loadMemberInfo()
    } catch (e) {
      wx.hideLoading()
      console.error('relogin error:', e)
      const msg = e.message || '登录失败'
      if (msg.length > 20) {
        wx.showModal({ title: '登录失败', content: msg, showCancel: false })
      } else {
        wx.showToast({ title: msg, icon: 'none', duration: 3000 })
      }
    }
  },
  // 获取手机号并绑定
  async onGetPhoneNumber(e) {
    if (e.detail.errMsg && e.detail.errMsg.includes('fail')) {
      wx.showToast({ title: '您取消了授权', icon: 'none' })
      return
    }
    const code = e.detail.code
    if (!code) {
      wx.showToast({ title: '请升级微信版本后重试', icon: 'none' })
      return
    }

    wx.showLoading({ title: '绑定中...' })
    try {
      const res = await wx.cloud.callFunction({
        name: 'getPhoneNumber',
        data: { code }
      })
      const result = res.result
      if (result.code !== 200 || !result.phoneNumber) {
        throw new Error(result.message || '获取手机号失败')
      }
      const phone = result.phoneNumber

      await api.post('/api/user/bind-phone', {
        phone: phone,
        nickname: '用户' + phone.slice(-4)
      })

      wx.hideLoading()
      wx.showToast({ title: '绑定成功', icon: 'success' })
      this.loadProfile()
    } catch (err) {
      wx.hideLoading()
      console.error('bind phone error:', err)
      wx.showToast({ title: err.message || '绑定失败', icon: 'none', duration: 3000 })
    }
  },
  // 选择微信头像
  async onChooseAvatar(e) {
    const avatarUrl = e.detail.avatarUrl
    if (!avatarUrl) return

    wx.showLoading({ title: '保存中...' })
    try {
      await api.post('/api/user/avatar', { avatar_url: avatarUrl })
      this.setData({ avatarUrl: avatarUrl })
      wx.hideLoading()
      wx.showToast({ title: '头像已更新', icon: 'success' })
    } catch (err) {
      wx.hideLoading()
      wx.showToast({ title: '保存失败', icon: 'none' })
    }
  },
  // 右上角设置：弹出操作菜单
  showSettings() {
    wx.showActionSheet({
      itemList: ['退出登录'],
      itemColor: '#FF4D4D',
      success: (res) => {
        if (res.tapIndex === 0) {
          this.confirmLogout()
        }
      }
    })
  },
  // 确认退出登录
  confirmLogout() {
    wx.showModal({
      title: '确认退出',
      content: '退出后需要重新登录',
      confirmColor: '#FF4D4D',
      success: (res) => {
        if (res.confirm) {
          auth.logout()
          this.setData({ user: {}, avatarUrl: DEFAULT_AVATAR, points: 0, memberInfo: {}, errorMsg: '登录已过期' })
        }
      }
    })
  },
  goOrders() { wx.navigateTo({ url: '/pages/mine/orders/index' }) },
  goPoints() { wx.navigateTo({ url: '/pages/points/index/index' }) },
  goPets() { wx.navigateTo({ url: '/pages/mine/pets/index' }) },
  goActivities() { wx.navigateTo({ url: '/pages/activity/list/index' }) },
  goRegistrations() { wx.showToast({ title: '报名记录开发中', icon: 'none' }) },
  goRecharges() { wx.showToast({ title: '充值记录开发中', icon: 'none' }) }
})
