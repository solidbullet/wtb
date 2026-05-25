const api = require('../../../utils/api')
const { API_BASE } = require('../../../utils/config')

const BREED_OPTIONS = ['柯基', '边牧', '金毛', '拉布拉多', '柴犬', '哈士奇', '泰迪', '比熊', '博美', '萨摩耶', '阿拉斯加', '法斗', '雪纳瑞', '其他']

Page({
  data: {
    isEdit: false,
    editId: null,
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

  onLoad(options) {
    const tempPhoto = wx.getStorageSync('temp_pet_photo_url') || ''
    if (options.id) {
      this.setData({ isEdit: true, editId: parseInt(options.id) })
      this.loadPetDetail(parseInt(options.id))
    }
    // 恢复可能因开发者工具自动编译而丢失的上传图片（添加/编辑都生效）
    if (tempPhoto) {
      this.setData({ photoUrl: tempPhoto })
    }
  },

  async loadPetDetail(id) {
    try {
      const res = await api.get('/api/user/pets')
      const pet = (res.data || []).find(p => p.id === id)
      if (!pet) {
        wx.showToast({ title: '宠物信息不存在', icon: 'none' })
        return
      }
      const breedIndex = BREED_OPTIONS.indexOf(pet.breed)
      const tempPhoto = wx.getStorageSync('temp_pet_photo_url') || ''
      this.setData({
        name: pet.name || '',
        breed: pet.breed || '',
        breedIndex: breedIndex >= 0 ? breedIndex : -1,
        gender: pet.gender || '',
        weight: pet.weight ? String(pet.weight) : '',
        birthday: pet.birthday || '',
        // 优先显示未保存的临时上传图片，其次显示数据库里的
        photoUrl: tempPhoto || pet.photo_url || '',
        notes: pet.notes || ''
      })
    } catch (e) {
      wx.showToast({ title: '加载宠物信息失败', icon: 'none' })
    }
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
      const openid = wx.getStorageSync('openid') || ''
      console.log('[choosePhoto] start upload, file:', file.tempFilePath)
      const uploadRes = await this.uploadFilePromise({
        url: `${API_BASE}/api/menu/admin/upload`,
        filePath: file.tempFilePath,
        name: 'file',
        header: {
          Authorization: token ? 'Bearer ' + token : '',
          'X-OpenID': openid
        }
      })
      wx.hideLoading()
      console.log('[choosePhoto] upload response raw:', uploadRes.data)
      const data = JSON.parse(uploadRes.data)
      console.log('[choosePhoto] upload response parsed:', data)
      if (data.code === 200) {
        console.log('[choosePhoto] upload success, url:', data.data.url)
        this.setData({ photoUrl: data.data.url })
        wx.setStorageSync('temp_pet_photo_url', data.data.url)
      } else {
        console.error('[choosePhoto] upload failed, code:', data.code, 'msg:', data.message)
        wx.showToast({ title: data.message || '上传失败', icon: 'none' })
      }
    } catch (err) {
      wx.hideLoading()
      console.error('[choosePhoto] upload error:', err)
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

  // 提交（添加或编辑）
  async submit() {
    const { isEdit, editId, name, breed, gender, weight, birthday, photoUrl, notes } = this.data
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
      const payload = {
        name: name.trim(),
        breed,
        gender,
        weight: parseFloat(weight) || 0,
        birthday,
        photo_url: photoUrl,
        notes
      }
      console.log('[submit] payload:', JSON.stringify(payload))
      console.log('[submit] photoUrl from data:', this.data.photoUrl)
      if (isEdit) {
        await api.put(`/api/user/pets/${editId}`, payload)
        wx.hideLoading()
        wx.showToast({ title: '修改成功', icon: 'success' })
      } else {
        await api.post('/api/user/pets', payload)
        wx.hideLoading()
        wx.showToast({ title: '添加成功', icon: 'success' })
      }
      wx.removeStorageSync('temp_pet_photo_url')
      setTimeout(() => wx.navigateBack(), 1000)
    } catch (err) {
      wx.hideLoading()
      wx.showToast({ title: err.message || '保存失败', icon: 'none' })
    }
  }
})
