const api = require('../../../utils/api')
const { resolveImageUrl } = require('../../../utils/image')

Page({
  data: {
    activities: [],
    defaultImage: resolveImageUrl('/images/activity.png')
  },
  onLoad() {
    this.loadActivities()
  },
  onShow() {
    this.loadActivities()
  },
  async loadActivities() {
    try {
      const res = await api.get('/api/activity/list?pageSize=50')
      const list = (res.data || []).map(a => ({
        ...a,
        image: resolveImageUrl(a.image),
        date: a.created_at ? a.created_at.split('T')[0] : '',
        quota: a.max_participants || 0,
        registered: a.current_participants || 0
      }))
      this.setData({ activities: list })
    } catch (e) {
      console.log('load activities error', e)
    }
  },
  register(e) {
    const id = e.currentTarget.dataset.id
    wx.showModal({
      title: '确认报名',
      content: '确定报名参加该活动吗？',
      success: (res) => {
        if (res.confirm) {
          api.post('/api/activity/' + id + '/register', {}).then(() => {
            wx.showToast({ title: '报名成功' })
            this.loadActivities()
          }).catch(err => {
            wx.showToast({ title: err.message || '报名失败', icon: 'none' })
          })
        }
      }
    })
  }
})
