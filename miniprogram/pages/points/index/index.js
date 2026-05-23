const api = require('../../../utils/api')
const { resolveImageUrl } = require('../../../utils/image')

Page({
  data: {
    points: 0,
    goods: [],
    defaultImage: resolveImageUrl('/images/gift.png')
  },
  onLoad() {
    this.loadPoints()
    this.loadGoods()
  },
  onShow() {
    this.loadPoints()
  },
  async loadPoints() {
    try {
      const res = await api.get('/api/points/account')
      const account = res.data || {}
      const available = (account.total_points || 0) - (account.used_points || 0) - (account.frozen_points || 0)
      this.setData({ points: available })
    } catch (e) {
      console.log('load points error', e)
    }
  },
  async loadGoods() {
    try {
      const res = await api.get('/api/points/goods?pageSize=50')
      this.setData({ goods: (res.data || []).map(g => ({ ...g, image: resolveImageUrl(g.image) })) })
    } catch (e) {
      console.log('load goods error', e)
    }
  },
  exchange(e) {
    const id = e.currentTarget.dataset.id
    const pts = e.currentTarget.dataset.points
    if (this.data.points < pts) {
      wx.showToast({ title: '积分不足', icon: 'none' })
      return
    }
    wx.showModal({
      title: '确认兑换',
      content: '消耗 ' + pts + ' 积分兑换此商品？',
      success: (res) => {
        if (res.confirm) {
          api.post('/api/points/exchange', { goods_id: id }).then(() => {
            wx.showToast({ title: '兑换成功' })
            this.loadPoints()
          }).catch(err => {
            wx.showToast({ title: err.message || '兑换失败', icon: 'none' })
          })
        }
      }
    })
  }
})
