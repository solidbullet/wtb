const api = require('../../utils/api')
const { resolveImageUrl } = require('../../utils/image')

Page({
  data: {
    announcements: [],
    activities: [],
    currentSwiper: 0
  },
  onLoad() {
    this.loadData()
  },
  async loadData() {
    try {
      const [annRes, actRes] = await Promise.all([
        api.get('/api/activity/announcements'),
        api.get('/api/activity/list?pageSize=3')
      ])
      let announcements = (annRes.data || []).map(a => ({ ...a, image: resolveImageUrl(a.image) }))
      if (!announcements || announcements.length === 0) {
        announcements = [
          { id: 1, title: '宠物欢乐派对', type: '活动', image: resolveImageUrl('/images/dog_party.png') },
          { id: 2, title: '生日特惠活动', type: '公告', image: resolveImageUrl('/images/birthday_party.png') }
        ]
      }
      this.setData({
        announcements: announcements,
        activities: (actRes.data || []).map(act => {
          if (act.image) act.image = resolveImageUrl(act.image)
          if (act.event_time) {
            const d = new Date(act.event_time)
            act.event_time = `${d.getMonth() + 1}月${d.getDate()}日 ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
          }
          return act
        })
      })
    } catch (e) {
      console.log('load data error', e)
    }
  },
  onSwiperChange(e) { this.setData({ currentSwiper: e.detail.current }) },

  /* ---- 四大服务跳转 ---- */
  goBoarding() { wx.navigateTo({ url: '/pages/boarding/index' }) },
  goCheckup() { wx.showToast({ title: '体检预约功能开发中', icon: 'none' }) },
  goDaycare() { wx.showToast({ title: '日托预约功能开发中', icon: 'none' }) },
  goFullcare() { wx.showToast({ title: '全托咨询：请致电 13800138000', icon: 'none' }) },

  /* ---- 快捷入口 ---- */
  goMenu() { wx.switchTab({ url: '/pages/order/menu' }) },
  goRecharge() { wx.switchTab({ url: '/pages/recharge/recharge' }) },
  goActivities() { wx.navigateTo({ url: '/pages/activity/list/index' }) },
  goActivityDetail(e) {
    wx.navigateTo({ url: '/pages/activity/list/index?id=' + e.currentTarget.dataset.id })
  }
})
