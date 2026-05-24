if (!wx.cloud) {
  console.error('请使用 2.2.3 或以上的基础库以使用云能力')
} else {
  wx.cloud.init({
    env: wx.cloud.DYNAMIC_CURRENT_ENV,
    traceUser: true
  })
}

const auth = require('./utils/auth')

App({
  globalData: {
    apiBase: 'http://localhost:8080',
    userInfo: null,
    loginReady: false,
    loginCallbacks: []
  },

  onLaunch() {
    console.log('App Launch')
    this.doLogin()
  },

  onShow() {
    console.log('App Show')
    if (!auth.checkLogin()) {
      this.doLogin()
    }
  },

  doLogin() {
    auth.wxLogin().then(data => {
      console.log('auto login success', data)
      this.globalData.loginReady = true
      this.globalData.memberLevel = (data.user && data.user.member_level) || 0
      const cbs = this.globalData.loginCallbacks
      this.globalData.loginCallbacks = []
      cbs.forEach(cb => cb(data))
    }).catch(err => {
      console.error('auto login failed', err)
      if (wx.getSystemInfoSync().platform === 'devtools') {
        wx.showToast({ title: err.message || '登录失败', icon: 'none' })
      }
    })
  },

  onLoginReady(callback) {
    if (this.globalData.loginReady) {
      callback({ memberLevel: this.globalData.memberLevel })
    } else {
      this.globalData.loginCallbacks.push(callback)
    }
  }
})
