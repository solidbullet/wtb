const api = require('../../utils/api')
const { resolveImageUrl } = require('../../utils/image')
const { SKIP_SCAN } = require('../../utils/config')

const app = getApp()

Page({
  data: {
    categories: [],
    currentCat: 0,
    currentCatName: '全部',
    dishes: [],
    cartMap: {},
    cartItems: [],
    cartTotal: 0,
    cartAmount: 0,
    seatId: '',
    showCartPanel: false,
    showSeatPicker: false,
    seatInput: '',
    seatOptions: ['A01','A02','A03','A04','A05','A06','B01','B02','B03'],
    memberLevel: 0,
    defaultImage: resolveImageUrl('/images/food.png'),
    errorMsg: '',
    scanning: false
  },

  _resolveSeatId(options) {
    let seatId = ''
    if (options && options.seat_id) {
      seatId = options.seat_id
    }
    if (!seatId && options && options.scene) {
      try {
        const decoded = decodeURIComponent(options.scene)
        const match = decoded.match(/seat_id=([^&]+)/)
        if (match) {
          seatId = match[1]
        } else {
          seatId = decoded
        }
      } catch (e) {
        seatId = options.scene
      }
    }
    return seatId
  },

  onLoad(options) {
    this.loadCategories()
    let seatId = this._resolveSeatId(options)

    if (seatId) {
      this.setData({ seatId })
      wx.setStorageSync('seat_id', seatId)
      this.loadCart()
      app.onLoginReady(() => this.loadMemberInfo())
      return
    }

    seatId = wx.getStorageSync('seat_id')
    if (seatId) {
      this.setData({ seatId })
      this.loadCart()
      app.onLoginReady(() => this.loadMemberInfo())
      return
    }

    if (SKIP_SCAN) {
      this.setData({ seatId: 'dev-seat' })
      this.loadCart()
      app.onLoginReady(() => this.loadMemberInfo())
      return
    }

    // 没有桌号，显示选择弹窗
    this.setData({ showSeatPicker: true })
    app.onLoginReady(() => this.loadMemberInfo())
  },

  onShow() {
    if (this.data.seatId && !this.data.showSeatPicker) {
      this.loadCart()
      app.onLoginReady(() => this.loadMemberInfo())
    }
  },

  selectSeat(e) {
    const seatId = e.currentTarget.dataset.id
    this.setData({ seatId, showSeatPicker: false })
    wx.setStorageSync('seat_id', seatId)
    wx.showToast({ title: '桌号: ' + seatId, icon: 'none' })
    this.loadCart()
    app.onLoginReady(() => this.loadMemberInfo())
  },

  onSeatInput(e) {
    this.setData({ seatInput: e.detail.value.trim() })
  },

  confirmSeatInput() {
    const seatId = this.data.seatInput
    if (!seatId) {
      wx.showToast({ title: '请输入桌号', icon: 'none' })
      return
    }
    this.setData({ seatId, showSeatPicker: false, seatInput: '' })
    wx.setStorageSync('seat_id', seatId)
    wx.showToast({ title: '桌号: ' + seatId, icon: 'none' })
    this.loadCart()
    app.onLoginReady(() => this.loadMemberInfo())
  },

  closeSeatPicker() {
    // 不允许关闭，必须选桌号
    wx.showToast({ title: '请先选择桌号', icon: 'none' })
  },

  onPickerTap() {
    // 阻止事件冒泡到遮罩层
  },

  async loadMemberInfo() {
    try {
      const res = await api.get('/api/user/member-info').catch(() => ({ data: {} }))
      this.setData({ memberLevel: (res.data && res.data.member_level) || 0 })
    } catch (e) {
      console.log('load member info error', e)
    }
  },

  async loadCategories() {
    try {
      const res = await api.get('/api/menu/categories')
      const cats = res.data || []
      const allCats = [{ id: 0, name: '全部' }, ...cats]
      this.setData({
        categories: allCats,
        currentCat: 0,
        currentCatName: '全部'
      })
      this.loadDishes()
    } catch (e) {
      console.error('load categories error', e)
      this.setData({ errorMsg: e.message || '加载分类失败' })
    }
  },

  async loadDishes() {
    try {
      const cat = this.data.currentCat
      const url = cat === 0 ? '/api/menu/dishes?pageSize=50' : '/api/menu/dishes?pageSize=50&category_id=' + cat
      const res = await api.get(url)
      let list = (res.data && res.data.list) ? res.data.list : []
      list = list.map(d => {
        if (d.images) {
          d.image = resolveImageUrl(d.images.split(',')[0].replace(/[[\]"]/g, ''))
        }
        return d
      })
      this.setData({ dishes: list, errorMsg: '' })
      if (this.data.cartItems.length > 0) {
        this._refreshCartPrices(list)
      }
    } catch (e) {
      console.error('load dishes error', e)
      this.setData({ errorMsg: e.message || '加载菜品失败' })
    }
  },

  _refreshCartPrices(dishes) {
    const items = this.data.cartItems.map(it => {
      const dish = dishes.find(d => d.id === it.dish_id)
      if (dish && typeof dish.price === 'number') {
        it.origin_price = dish.price
        it.member_price = dish.member_price || dish.price
        const realPrice = (this.data.memberLevel >= 1 && dish.member_price > 0) ? dish.member_price : dish.price
        it.unit_price = realPrice
        it.saved_per_unit = dish.price - realPrice
        it.saved_total = it.saved_per_unit * it.quantity
      }
      return it
    })
    const map = {}
    let total = 0, amount = 0, originAmount = 0
    items.forEach(it => {
      map[it.dish_id] = it.quantity
      total += it.quantity
      amount += it.quantity * it.unit_price
      originAmount += it.quantity * (it.origin_price || it.unit_price)
    })
    const savedAmount = originAmount - amount
    this.setData({ cartMap: map, cartItems: items, cartTotal: total, cartAmount: amount, cartOriginAmount: originAmount, cartSavedAmount: savedAmount })
  },

  async loadCart() {
    if (!this.data.seatId) return
    try {
      const res = await api.get('/api/order/cart/list?seat_id=' + this.data.seatId)
      let items = res.data || []
      const map = {}
      let total = 0, amount = 0, originAmount = 0
      const dishes = this.data.dishes
      items = items.map(it => {
        const dish = dishes.find(d => d.id === it.dish_id)
        if (dish) {
          it.origin_price = dish.price
          it.member_price = dish.member_price || dish.price
          const realPrice = (this.data.memberLevel >= 1 && dish.member_price > 0) ? dish.member_price : dish.price
          it.unit_price = realPrice
          it.saved_per_unit = dish.price - realPrice
          it.saved_total = it.saved_per_unit * it.quantity
        }
        map[it.dish_id] = it.quantity
        total += it.quantity
        amount += it.quantity * it.unit_price
        originAmount += it.quantity * (it.origin_price || it.unit_price)
        return it
      })
      const savedAmount = originAmount - amount
      this.setData({ cartMap: map, cartItems: items, cartTotal: total, cartAmount: amount, cartOriginAmount: originAmount, cartSavedAmount: savedAmount })
    } catch (e) {
      console.log('load cart error', e)
    }
  },

  scanSeat() {
    if (this.data.scanning || this.data.seatId) return
    this.setData({ scanning: true })
    wx.scanCode({
      onlyFromCamera: true,
      success: (res) => {
        const result = res.result
        let seatId = result
        try {
          const url = new URL(result)
          seatId = url.searchParams.get('seat_id') || result
        } catch (e) {
          // not a URL
        }
        api.get('/api/seat/scan?code=' + encodeURIComponent(result)).then((verifyRes) => {
          const verifiedSeatId = verifyRes.data && verifyRes.data.seat_id ? verifyRes.data.seat_id : seatId
          this.setData({ seatId: verifiedSeatId, errorMsg: '' })
          wx.setStorageSync('seat_id', verifiedSeatId)
          this.loadCart()
          wx.showToast({ title: '桌号: ' + verifiedSeatId, icon: 'none' })
        }).catch(() => {
          this.setData({ seatId: seatId, errorMsg: '' })
          wx.setStorageSync('seat_id', seatId)
          this.loadCart()
        })
      },
      fail: (err) => {
        if (err.errMsg && err.errMsg.includes('cancel')) {
          this.setData({ errorMsg: '请点击扫码按钮选择桌号' })
        } else {
          wx.showToast({ title: '扫码失败', icon: 'none' })
        }
      },
      complete: () => {
        this.setData({ scanning: false })
      }
    })
  },

  reScan() {
    if (this.data.scanning) return
    wx.showModal({
      title: '换桌',
      content: '重新扫码将清空当前购物车，确定换桌吗？',
      success: (res) => {
        if (res.confirm) {
          wx.removeStorageSync('seat_id')
          this.setData({
            seatId: '', cartMap: {}, cartItems: [],
            cartTotal: 0, cartAmount: 0, cartOriginAmount: 0,
            cartSavedAmount: 0, showCartPanel: false, scanning: false
          })
          this.scanSeat()
        }
      }
    })
  },

  switchCat(e) {
    const id = e.currentTarget.dataset.id
    const cat = this.data.categories.find(c => c.id === id)
    this.setData({ currentCat: id, currentCatName: cat ? cat.name : '全部' })
    this.loadDishes()
  },

  _getUnitPrice(dish) {
    if (!dish || typeof dish.price !== 'number') return null
    return (this.data.memberLevel >= 1 && dish.member_price > 0) ? dish.member_price : dish.price
  },

  plus(e) {
    if (!this.data.seatId) {
      wx.showToast({ title: '请先选择桌号', icon: 'none' })
      return
    }
    const id = parseInt(e.currentTarget.dataset.id)
    const dish = this.data.dishes.find(d => d.id === id)
    const unitPrice = this._getUnitPrice(dish)
    if (unitPrice === null) return
    const q = (this.data.cartMap[id] || 0) + 1
    const map = { ...this.data.cartMap, [id]: q }
    this.setData({ cartMap: map, cartTotal: this.data.cartTotal + 1, cartAmount: this.data.cartAmount + unitPrice })
    this.syncCart(id, q, { name: dish.name, price: unitPrice })
  },

  minus(e) {
    const id = parseInt(e.currentTarget.dataset.id)
    const dish = this.data.dishes.find(d => d.id === id)
    const unitPrice = this._getUnitPrice(dish)
    if (unitPrice === null) return
    const q = Math.max(0, (this.data.cartMap[id] || 0) - 1)
    const map = { ...this.data.cartMap }
    if (q === 0) delete map[id]
    else map[id] = q
    this.setData({ cartMap: map, cartTotal: Math.max(0, this.data.cartTotal - 1), cartAmount: Math.max(0, this.data.cartAmount - unitPrice) })
    this.syncCart(id, q, { name: dish.name, price: unitPrice })
  },

  _cartTimers: {},
  _syncVersion: 0,
  _pendingSyncs: [],

  syncCart(dishId, quantity, dish) {
    var self = this
    var done = function () {
      var idx = self._pendingSyncs.indexOf(promise)
      if (idx >= 0) self._pendingSyncs.splice(idx, 1)
    }
    var promise = api.post('/api/order/cart/add', {
      seat_id: this.data.seatId,
      dish_id: dishId,
      quantity: quantity,
      dish_name: dish.name,
      unit_price: dish.price
    }).then(function () {
      self._scheduleCartRefresh()
      done()
    }).catch(function (err) {
      console.log('sync cart error:', err)
      wx.showToast({ title: err.message || '操作失败', icon: 'none' })
      self.loadCart()
      done()
    })
    this._pendingSyncs.push(promise)
  },

  _scheduleCartRefresh() {
    this._syncVersion++
    const version = this._syncVersion
    if (this._cartTimers['_refresh']) {
      clearTimeout(this._cartTimers['_refresh'])
    }
    this._cartTimers['_refresh'] = setTimeout(() => {
      delete this._cartTimers['_refresh']
      if (this._syncVersion === version) {
        this.loadCart()
      }
    }, 200)
  },

  toggleCartPanel() {
    if (this.data.cartAmount <= 0) return
    this.setData({ showCartPanel: !this.data.showCartPanel })
  },

  closeCartPanel() {
    this.setData({ showCartPanel: false })
  },

  panelPlus(e) {
    const dishId = parseInt(e.currentTarget.dataset.id)
    const item = this.data.cartItems.find(it => it.dish_id === dishId)
    if (!item) return
    const dish = this.data.dishes.find(d => d.id === dishId)
    const unitPrice = this._getUnitPrice(dish)
    if (unitPrice === null) return
    const q = item.quantity + 1
    this.syncCart(dishId, q, { name: item.dish_name, price: unitPrice })
  },

  panelMinus(e) {
    const dishId = parseInt(e.currentTarget.dataset.id)
    const item = this.data.cartItems.find(it => it.dish_id === dishId)
    if (!item) return
    const dish = this.data.dishes.find(d => d.id === dishId)
    const unitPrice = this._getUnitPrice(dish)
    if (unitPrice === null) return
    const q = Math.max(0, item.quantity - 1)
    this.syncCart(dishId, q, { name: item.dish_name, price: unitPrice })
  },

  panelRemove(e) {
    const dishId = parseInt(e.currentTarget.dataset.id)
    const item = this.data.cartItems.find(it => it.dish_id === dishId)
    if (!item) return
    this.syncCart(dishId, 0, { name: item.dish_name, price: item.unit_price })
  },

  async clearCart() {
    wx.showModal({
      title: '确认清空',
      content: '确定清空购物车吗？',
      success: (res) => {
        if (res.confirm) {
          const promises = this.data.cartItems.map(item =>
            api.post('/api/order/cart/add', {
              seat_id: this.data.seatId,
              dish_id: item.dish_id,
              quantity: 0,
              dish_name: item.dish_name,
              unit_price: item.unit_price
            })
          )
          Promise.allSettled(promises).then(() => {
            this.loadCart()
            this.setData({ showCartPanel: false })
          }).catch(() => {
            wx.showToast({ title: '清空失败', icon: 'none' })
          })
        }
      }
    })
  },

  async goConfirm() {
    if (this.data.cartAmount <= 0) return
    if (!this.data.seatId) {
      wx.showToast({ title: '请先选择桌号', icon: 'none' })
      return
    }
    this.setData({ showCartPanel: false })
    if (this._pendingSyncs.length > 0) {
      await Promise.all(this._pendingSyncs.map(function (p) {
        return p.catch(function () {})
      }))
    }
    if (this._cartTimers['_refresh']) {
      clearTimeout(this._cartTimers['_refresh'])
      delete this._cartTimers['_refresh']
    }
    wx.navigateTo({
      url: '/pages/order/confirm/index?seat_id=' + this.data.seatId
    })
  }
})
