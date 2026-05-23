const api = require('../../../utils/api')
const { resolveImageUrl } = require('../../../utils/image')

Page({
  data: {
    pets: [],
    defaultAvatar: resolveImageUrl('/images/pet.png')
  },
  onLoad() {
    this.loadPets()
  },
  onShow() {
    this.loadPets()
  },
  async loadPets() {
    try {
      const res = await api.get('/api/user/pets')
      this.setData({ pets: (res.data || []).map(p => ({ ...p, avatar: resolveImageUrl(p.avatar) })) })
    } catch (e) {
      // 静默失败，显示空状态
    }
  },
  addPet() {
    wx.navigateTo({ url: '/pages/mine/pets/add' })
  }
})
