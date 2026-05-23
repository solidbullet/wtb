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
      this.setData({ pets: (res.data || []).map(p => ({
        ...p,
        photo_url: p.photo_url ? resolveImageUrl(p.photo_url) : ''
      })) })
    } catch (e) {
      // 静默失败，显示空状态
    }
  },
  addPet() {
    wx.navigateTo({ url: '/pages/mine/pets/add' })
  },
  // 点击宠物卡片查看详情
  showPetDetail(e) {
    const pet = e.currentTarget.dataset.pet
    const lines = [
      `🐾 名字：${pet.name}`,
      `品种：${pet.breed || '未知'}`,
      pet.gender ? `性别：${pet.gender === 'male' ? '♂ 公' : '♀ 母'}` : '',
      pet.weight ? `体重：${pet.weight} kg` : '',
      pet.birthday ? `生日：${pet.birthday}` : '',
      pet.notes ? `\n📋 寄养注意事项：\n${pet.notes}` : ''
    ].filter(Boolean)

    wx.showModal({
      title: pet.name,
      content: lines.join('\n'),
      showCancel: false,
      confirmText: '知道了'
    })
  }
})
