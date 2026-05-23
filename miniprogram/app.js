// 初始化云开发（必须在 App 定义之前，且必须在调用云函数之前）
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
    userInfo: null
  },
  onLaunch() {
    console.log('App Launch')
    this.doLogin()
  },
  onShow() {
    console.log('App Show')
    // 切回前台时检查登录状态，未登录则自动重新登录
    if (!auth.checkLogin()) {
      this.doLogin()
    }
  },
  doLogin() {
    auth.wxLogin().then(data => {
      console.log('auto login success', data)
    }).catch(err => {
      console.error('auto login failed', err)
      // 仅在开发工具中提示，避免真机上每次启动都弹窗
      if (wx.getSystemInfoSync().platform === 'devtools') {
        wx.showToast({ title: err.message || '登录失败', icon: 'none' })
      }
    })
  }
})
