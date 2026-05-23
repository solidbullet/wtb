const api = require('../../../utils/api')
const { API_BASE } = require('../../../utils/config')

const BREED_OPTIONS = ['柯基', '边牧', '金毛', '拉布拉多', '柴犬', '哈士奇', '泰迪', '比熊', '博美', '萨摩耶', '阿拉斯加', '法斗', '雪纳瑞', '其他']

Page({
  data: {
    name: '',
    breed: '',
    breedIndex: -1,
    gender: '',
    weight: '',
    birthday: '',
    photoUrl: '',
    notes: '',
    breedOptions: BREED_OPTIONS
  },

  onNameInput(e) { this.setData({ name: e.detail.value }) },

  onBreedChange(e) {
    const idx = parseInt(e.detail.value)
    this.setData({ breedIndex: idx, breed: BREED_OPTIONS[idx] })
  },

  onGenderChange(e) { this.setData({ gender: e.currentTarget.dataset.gender }) },

  onWeightInput(e) { this.setData({ weight: e.detail.value }) },

  onBirthdayChange(e) { this.setData({ birthday: e.detail.value }) },

  onNotesInput(e) { this.setData({ notes: e.detail.value }) },

  // 选择并上传宠物照片
  async choosePhoto() {
    try {
      const res = await wx.chooseMedia({ count: 1, mediaType: ['image'] })
      const file = res.tempFiles[0]
      wx.showLoading({ title: '上传中...' })
      const token = wx.getStorageSync('token') || ''
      const uploadRes = await this.uploadFilePromise({
        url: `${API_BASE}/api/menu/admin/upload`,
        filePath: file.tempFilePath,
        name: 'file',
        header: { Authorization: token ? 'Bearer ' + token : '' }
      })
      wx.hideLoading()
      const data = JSON.parse(uploadRes.data)
      if (data.code === 200) {
        this.setData({ photoUrl: data.data.url })
      } else {
        wx.showToast({ title: data.message || '上传失败', icon: 'none' })
      }
    } catch (err) {
      wx.hideLoading()
      wx.showToast({ title: err.message || '上传失败', icon: 'none' })
    }
  },

  uploadFilePromise(options) {
    return new Promise((resolve, reject) => {
      const task = wx.uploadFile({
        ...options,
        success: resolve,
        fail: reject
      })
    })
  },

  // 提交添加宠物
  async submit() {
    const { name, breed, gender, weight, birthday, photoUrl, notes } = this.data
    if (!name.trim()) {
      wx.showToast({ title: '请输入宠物名字', icon: 'none' })
      return
    }
    if (!breed) {
      wx.showToast({ title: '请选择品种', icon: 'none' })
      return
    }
    wx.showLoading({ title: '保存中...' })
    try {
      await api.post('/api/user/pets', {
        name: name.trim(),
        breed,
        gender,
        weight: parseFloat(weight) || 0,
        birthday,
        photo_url: photoUrl,
        notes
      })
      wx.hideLoading()
      wx.showToast({ title: '添加成功', icon: 'success' })
      setTimeout(() => wx.navigateBack(), 1000)
    } catch (err) {
      wx.hideLoading()
      wx.showToast({ title: err.message || '添加失败', icon: 'none' })
    }
  }
})
