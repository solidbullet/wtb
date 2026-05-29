const api = require('../../../utils/api')

Page({
  data: {
    seatId: '',
    items: [],
    cartTotal: 0,
    cartAmount: 0,
    cartOriginAmount: 0,
    cartSavedAmount: 0,
    memberLevel: 0,
    remark: '',
    submitting: false
  },
  async onLoad(options) {
    let seatId = options.seat_id
    if (!seatId) {
      seatId = wx.getStorageSync('seat_id')
    }
    if (!seatId) {
      seatId = '2'
    }
    this.setData({ seatId })
    await this.loadMemberInfo()   // 先等待会员信息加载完成
    await this.loadCart(seatId)   // 再加载购物车（此时 memberLevel 已正确）
  },
  async onShow() {
    await this.loadMemberInfo()
    await this.loadCart(this.data.seatId)
  },
  async loadMemberInfo() {
    try {
      const res = await api.get('/api/user/member-info').catch(() => ({ data: {} }))
      this.setData({ memberLevel: (res.data && res.data.member_level) || 0 })
    } catch (e) { console.log('load member info error', e) }
  },
  async loadCart(seatId) {
    let sid = seatId || this.data.seatId
    if (!sid) {
      sid = wx.getStorageSync('seat_id') || '2'
    }
    try {
      // 自己获取会员等级，不依赖外部顺序
      let memberLevel = 0
      try {
        const memberRes = await api.get('/api/user/member-info').catch(() => ({ data: {} }))
        memberLevel = (memberRes.data && memberRes.data.member_level) || 0
        this.setData({ memberLevel })
        console.log('[confirm] memberLevel loaded:', memberLevel)
      } catch (e) {
        console.log('[confirm] load member info error', e)
      }

      const res = await api.get('/api/order/cart/list?seat_id=' + sid)
      let items = res.data || []
      console.log('[confirm] cart items from backend:', items.map(it => ({ id: it.dish_id, name: it.dish_name, qty: it.quantity, unit_price: it.unit_price })))

      // 查询菜品最新价格和会员价
      const dishIds = items.map(it => it.dish_id)
      if (dishIds.length > 0) {
        try {
          const dishRes = await api.post('/api/menu/dishes/batch', { dish_ids: dishIds })
          const dishMap = {}
          ;(dishRes.data || []).forEach(d => {
            dishMap[d.id] = d
          })
          console.log('[confirm] dishMap:', Object.keys(dishMap).map(k => ({ id: dishMap[k].id, price: dishMap[k].price, member_price: dishMap[k].member_price })))

          items = items.map(it => {
            const d = dishMap[it.dish_id]
            if (d && typeof d.price === 'number') {
              it.origin_price = d.price
              it.member_price = d.member_price || d.price
              const realPrice = (memberLevel >= 1 && d.member_price > 0) ? d.member_price : d.price
              it.unit_price = realPrice
              it.saved_per_unit = d.price - realPrice
              it.saved_total = it.saved_per_unit * it.quantity
              console.log('[confirm] dish calc:', it.dish_name, 'memberLevel=', memberLevel, 'origin=', d.price, 'member=', d.member_price, 'real=', realPrice)
            }
            return it
          })
        } catch (e) {
          console.log('[confirm] fetch dish prices error', e)
        }
      }
      let total = 0, amount = 0, originAmount = 0
      items.forEach(it => {
        total += it.quantity
        amount += it.quantity * it.unit_price
        originAmount += it.quantity * (it.origin_price || it.unit_price)
      })
      const savedAmount = originAmount - amount
      console.log('[confirm] total amount:', amount, 'origin:', originAmount, 'saved:', savedAmount)
      this.setData({ items, cartTotal: total, cartAmount: amount, cartOriginAmount: originAmount, cartSavedAmount: savedAmount })
      if (items.length === 0) {
        wx.showToast({ title: '购物车为空', icon: 'none' })
        setTimeout(() => wx.navigateBack(), 1500)
      }
    } catch (e) {
      console.log('[confirm] load cart error', e)
    }
  },
  onRemarkInput(e) {
    this.setData({ remark: e.detail.value })
  },
  async removeItem(e) {
    const dishId = e.currentTarget.dataset.id
    try {
      await api.post('/api/order/cart/remove', {
        seat_id: this.data.seatId,
        dish_id: dishId
      })
      this.loadCart()
    } catch (e) {
      wx.showToast({ title: e.message || '删除失败', icon: 'none' })
    }
  },
  async submitOrder() {
    if (this.data.submitting) return
    if (this.data.items.length === 0) {
      wx.showToast({ title: '购物车为空', icon: 'none' })
      return
    }
    this.setData({ submitting: true })
    try {
      // 创建订单
      const orderRes = await api.post('/api/order/create', {
        seat_id: this.data.seatId,
        remark: this.data.remark
      })
      const order = orderRes.data || {}
      wx.showToast({ title: '订单创建成功' })
      // 跳转到支付页或订单详情页
      wx.redirectTo({
        url: '/pages/mine/orders/index'
      })
    } catch (e) {
      wx.showToast({ title: e.message || '下单失败', icon: 'none' })
      this.setData({ submitting: false })
    }
  }
})
