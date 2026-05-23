const api = require('../../utils/api')
const { resolveImageUrl } = require('../../utils/image')
const { SKIP_SCAN } = require('../../utils/config')

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
    memberLevel: 0,
    defaultImage: resolveImageUrl('/images/food.png'),
    errorMsg: '',
    scanning: false
  },
  onLoad(options) {
    this.loadCategories()
    // 优先使用 URL 传入的 seat_id
    if (options.seat_id) {
      this.setData({ seatId: options.seat_id })
      this.loadCart()
    } else if (SKIP_SCAN) {
      // 开发环境：跳过扫码，直接使用默认桌号
      this.setData({ seatId: 'dev-seat' })
      this.loadCart()
    } else {
      this.scanSeat()
    }
    // 延迟加载会员信息，确保自动登录已完成
    setTimeout(() => this.loadMemberInfo(), 500)
  },
  onShow() {
    if (!this.data.seatId && !this.data.scanning) {
      if (SKIP_SCAN) {
        this.setData({ seatId: 'dev-seat' })
        this.loadCart()
      } else {
        this.scanSeat()
      }
    } else if (this.data.seatId) {
      this.loadCart()
      this.loadMemberInfo()
    }
  },
  async loadMemberInfo() {
    try {
      const res = await api.get('/api/user/member-info').catch(() => ({ data: {} }))
      this.setData({ memberLevel: (res.data && res.data.member_level) || 0 })
    } catch (e) { console.log('load member info error', e) }
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
      // 后端返回 images 字段，需要提取第一张图片并解析为完整URL
      list = list.map(d => {
        if (d.images) {
          d.image = resolveImageUrl(d.images.split(',')[0].replace(/[[\]"]/g, ''))
        }
        return d
      })
      this.setData({ dishes: list, errorMsg: '' })
      // 用当前菜品价格刷新购物车显示（防止后台改价后购物车显示旧价格）
      if (this.data.cartItems.length > 0) {
        const items = this.data.cartItems.map(it => {
          const dish = list.find(d => d.id === it.dish_id)
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
      }
    } catch (e) {
      console.error('load dishes error', e)
      this.setData({ errorMsg: e.message || '加载菜品失败' })
    }
  },
  async loadCart() {
    if (!this.data.seatId) return
    try {
      const res = await api.get('/api/order/cart/list?seat_id=' + this.data.seatId)
      let items = res.data || []
      const map = {}
      let total = 0, amount = 0, originAmount = 0
      items = items.map(it => {
        const dish = this.data.dishes.find(d => d.id === it.dish_id)
        if (dish) {
          it.origin_price = dish.price
          it.member_price = dish.member_price || dish.price
          // 重新计算当前会员等级下的实际单价
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
  // 扫码获取桌号
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
          // not a URL, use raw result
        }
        // 调用后端验证座位
        api.get('/api/seat/scan?code=' + encodeURIComponent(result)).then((verifyRes) => {
          const verifiedSeatId = verifyRes.data && verifyRes.data.seat_id ? verifyRes.data.seat_id : seatId
          this.setData({ seatId: verifiedSeatId, errorMsg: '' })
          this.loadCart()
          wx.showToast({ title: '桌号: ' + verifiedSeatId, icon: 'none' })
        }).catch((err) => {
          // 验证失败也允许进入，使用解析出的 seatId
          console.log('seat verify error:', err)
          this.setData({ seatId: seatId, errorMsg: '' })
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
  // 重新扫码换桌
  reScan() {
    if (this.data.scanning) return
    wx.showModal({
      title: '换桌',
      content: '重新扫码将清空当前购物车，确定换桌吗？',
      success: (res) => {
        if (res.confirm) {
          this.setData({ seatId: '', cartMap: {}, cartItems: [], cartTotal: 0, cartAmount: 0, cartOriginAmount: 0, cartSavedAmount: 0, showCartPanel: false, scanning: false })
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
  // 菜品列表中的 +
  plus(e) {
    if (!this.data.seatId) {
      wx.showToast({ title: '请先扫码选择桌号', icon: 'none' })
      return
    }
    const id = parseInt(e.currentTarget.dataset.id)
    const dish = this.data.dishes.find(d => d.id === id)
    if (!dish || typeof dish.price !== 'number') return
    const unitPrice = (this.data.memberLevel >= 1 && dish.member_price > 0) ? dish.member_price : dish.price
    const q = (this.data.cartMap[id] || 0) + 1
    const map = { ...this.data.cartMap, [id]: q }
    this.setData({ cartMap: map, cartTotal: this.data.cartTotal + 1, cartAmount: this.data.cartAmount + unitPrice })
    this.syncCart(id, q, { name: dish.name, price: unitPrice })
  },
  // 菜品列表中的 -
  minus(e) {
    const id = parseInt(e.currentTarget.dataset.id)
    const dish = this.data.dishes.find(d => d.id === id)
    if (!dish || typeof dish.price !== 'number') return
    const unitPrice = (this.data.memberLevel >= 1 && dish.member_price > 0) ? dish.member_price : dish.price
    const q = Math.max(0, (this.data.cartMap[id] || 0) - 1)
    const map = { ...this.data.cartMap }
    if (q === 0) delete map[id]
    else map[id] = q
    this.setData({ cartMap: map, cartTotal: Math.max(0, this.data.cartTotal - 1), cartAmount: Math.max(0, this.data.cartAmount - unitPrice) })
    this.syncCart(id, q, { name: dish.name, price: unitPrice })
  },
  syncCart(dishId, quantity, dish) {
    api.post('/api/order/cart/add', {
      seat_id: this.data.seatId,
      dish_id: dishId,
      quantity: quantity,
      dish_name: dish.name,
      unit_price: dish.price
    }).then(() => {
      this.loadCart()
    }).catch((err) => {
      console.log('sync cart error:', err)
      wx.showToast({ title: err.message || '添加失败', icon: 'none' })
    })
  },
  // 购物车弹层
  toggleCartPanel() {
    if (this.data.cartAmount <= 0) return
    this.setData({ showCartPanel: !this.data.showCartPanel })
  },
  closeCartPanel() {
    this.setData({ showCartPanel: false })
  },
  // 弹层中的 +
  panelPlus(e) {
    const dishId = parseInt(e.currentTarget.dataset.id)
    const item = this.data.cartItems.find(it => it.dish_id === dishId)
    if (!item) return
    const dish = this.data.dishes.find(d => d.id === dishId)
    const unitPrice = (this.data.memberLevel >= 1 && dish && dish.member_price > 0)
      ? dish.member_price
      : ((dish && typeof dish.price === 'number') ? dish.price : item.unit_price)
    const q = item.quantity + 1
    this.syncCart(dishId, q, { name: item.dish_name, price: unitPrice })
  },
  // 弹层中的 -
  panelMinus(e) {
    const dishId = parseInt(e.currentTarget.dataset.id)
    const item = this.data.cartItems.find(it => it.dish_id === dishId)
    if (!item) return
    const dish = this.data.dishes.find(d => d.id === dishId)
    const unitPrice = (this.data.memberLevel >= 1 && dish && dish.member_price > 0)
      ? dish.member_price
      : ((dish && typeof dish.price === 'number') ? dish.price : item.unit_price)
    const q = Math.max(0, item.quantity - 1)
    this.syncCart(dishId, q, { name: item.dish_name, price: unitPrice })
  },
  // 弹层中的删除
  panelRemove(e) {
    const dishId = parseInt(e.currentTarget.dataset.id)
    const item = this.data.cartItems.find(it => it.dish_id === dishId)
    if (!item) return
    this.syncCart(dishId, 0, { name: item.dish_name, price: item.unit_price })
  },
  // 清空购物车
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
  // 去结算
  goConfirm() {
    if (this.data.cartAmount <= 0) return
    if (!this.data.seatId) {
      wx.showToast({ title: '请先扫码选择桌号', icon: 'none' })
      return
    }
    this.setData({ showCartPanel: false })
    wx.navigateTo({
      url: '/pages/order/confirm/index?seat_id=' + this.data.seatId
    })
  }
})
