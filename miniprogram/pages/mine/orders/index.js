const api = require('../../../utils/api')

Page({
  data: {
    currentTab: 'all',
    orders: [],
    statusText: {
      pending: '待支付',
      paid: '进行中',
      completed: '已完成',
      cancelled: '已取消'
    }
  },
  onLoad() {
    this.loadOrders()
  },
  onShow() {
    this.loadOrders()
  },
  switchTab(e) {
    const tab = e.currentTarget.dataset.tab
    this.setData({ currentTab: tab })
    this.loadOrders()
  },
  async loadOrders() {
    try {
      const params = { pageSize: 50 }
      if (this.data.currentTab !== 'all') {
        params.status = this.data.currentTab
      }
      const res = await api.get('/api/order/list', params)
      this.setData({ orders: (res.data && res.data.list) ? res.data.list : [] })
    } catch (err) {
      wx.showToast({ title: err.message || '加载失败', icon: 'none' })
    }
  }
})
