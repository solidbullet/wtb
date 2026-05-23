const api = require('../../../utils/api')

Page({
  data: {
    records: [],
    loading: false,
    hasMore: true,
    page: 1,
    pageSize: 20
  },
  onLoad() {
    this.loadRecords()
  },
  onShow() {
    this.setData({ page: 1, records: [], hasMore: true })
    this.loadRecords()
  },
  async loadRecords() {
    if (this.data.loading || !this.data.hasMore) return
    this.setData({ loading: true })
    try {
      const res = await api.get('/api/user/recharge-records', {
        page: this.data.page,
        pageSize: this.data.pageSize
      })
      const data = res.data || {}
      const rawList = data.list || []
      // 在 JS 里格式化数据，WXML 不支持调用自定义方法
      const list = rawList.map(item => ({
        ...item,
        _date: this.fmtDate(item.created_at),
        _amount: (item.amount / 100).toFixed(2),
        _gifted: (item.gifted_amount / 100).toFixed(2),
        _total: ((item.amount + item.gifted_amount) / 100).toFixed(2)
      }))
      this.setData({
        records: this.data.page === 1 ? list : this.data.records.concat(list),
        page: this.data.page + 1,
        hasMore: list.length === this.data.pageSize,
        loading: false
      })
    } catch (err) {
      this.setData({ loading: false })
      wx.showToast({ title: err.message || '加载失败', icon: 'none' })
    }
  },
  onReachBottom() {
    this.loadRecords()
  },
  fmtDate(datetime) {
    if (!datetime) return ''
    const d = new Date(datetime)
    return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
  }
})
