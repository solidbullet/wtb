const api = require('../../utils/api')

Page({
  data: {
    balance: 0,
    memberLevel: 0,
    totalRecharged: 0,
    plans: [],
    upgradeInfo: null,
    loading: false
  },

  onShow() {
    this.loadData()
  },

  async loadData() {
    this.setData({ loading: true })
    try {
      await Promise.all([
        this.loadProfile(),
        this.loadPlans(),
        this.loadUpgradeInfo()
      ])
    } catch (e) {
      console.error('load data error', e)
    }
    this.setData({ loading: false })
  },

  async loadProfile() {
    try {
      const res = await api.get('/api/user/profile')
      this.setData({
        balance: res.data.balance || 0,
        memberLevel: res.data.member_level || 0
      })
    } catch (e) {
      console.log('load profile error', e)
    }
  },

  async loadPlans() {
    try {
      const res = await api.get('/api/user/recharge-plans')
      this.setData({ plans: res.data || [] })
    } catch (e) {
      console.log('load plans error', e)
      // fallback 到本地定义
      this.setData({
        plans: [
          { id: 1, name: '开通会员', amount: 19900, final_amount: 0, gift_amount: 0, member_level: 1, desc: '享受会员价折扣，余额不变' },
          { id: 2, name: '预充值 ¥1000', amount: 100000, final_amount: 120000, gift_amount: 20000, member_level: 2, desc: '到账 ¥1200（送 ¥200）' },
          { id: 3, name: '预充值 ¥2000', amount: 200000, final_amount: 240000, gift_amount: 40000, member_level: 2, desc: '到账 ¥2400（送 ¥400）' },
          { id: 4, name: '预充值 ¥5000', amount: 500000, final_amount: 600000, gift_amount: 100000, member_level: 2, desc: '到账 ¥6000（送 ¥1000）' }
        ]
      })
    }
  },

  async loadUpgradeInfo() {
    try {
      const res = await api.get('/api/user/upgrade-info')
      const info = res.data || null
      this.setData({
        upgradeInfo: info,
        totalRecharged: info ? info.total_recharged : 0
      })
    } catch (e) {
      console.log('load upgrade info error', e)
    }
  },

  // 计算档位显示状态
  getPlanStatus(plan) {
    const { totalRecharged } = this.data
    if (totalRecharged >= plan.amount) {
      return { status: 'completed', label: '已完成', disabled: true }
    }
    const diff = plan.amount - totalRecharged
    if (totalRecharged > 0 && diff > 0) {
      return { status: 'upgrade', label: `补差价 ¥${(diff / 100).toFixed(0)} 升级`, disabled: false }
    }
    return { status: 'available', label: '立即充值', disabled: false }
  },

  // 点击充值/升级
  async doRecharge(e) {
    const planId = parseInt(e.currentTarget.dataset.planId)
    const plan = this.data.plans.find(p => p.id === planId)
    if (!plan) return

    const { totalRecharged } = this.data
    const payAmount = plan.amount - totalRecharged

    // 如果已经是该档位或更高，提示无需充值
    if (payAmount <= 0) {
      wx.showToast({ title: '您已达到该档位', icon: 'none' })
      return
    }

    // 构建确认弹窗内容
    const isUpgrade = totalRecharged > 0
    const title = isUpgrade ? '升级充值' : plan.name
    const contentLines = []
    if (isUpgrade) {
      contentLines.push(`您已累计充值 ¥${(totalRecharged / 100).toFixed(0)}`)
      contentLines.push(`本次支付 ¥${(payAmount / 100).toFixed(0)} 升级到 ${plan.name}`)
    } else {
      contentLines.push(`支付 ¥${(plan.amount / 100).toFixed(0)}`)
    }
    if (plan.final_amount > 0) {
      contentLines.push(`到账 ¥${(plan.final_amount / 100).toFixed(0)}`)
      if (plan.gift_amount > 0) {
        contentLines.push(`（赠送 ¥${(plan.gift_amount / 100).toFixed(0)}）`)
      }
    } else {
      contentLines.push('余额不变，享受会员价折扣')
    }

    wx.showModal({
      title: title,
      content: contentLines.join('\n'),
      success: async (res) => {
        if (res.confirm) {
          wx.showLoading({ title: '处理中...' })
          try {
            const result = await api.post('/api/user/recharge', {
              plan_id: planId,
              channel: 'wxpay'
            })
            wx.hideLoading()
            wx.showToast({
              title: isUpgrade ? '升级成功' : (plan.id === 1 ? '会员开通成功' : '充值成功'),
              icon: 'success'
            })
            // 刷新数据
            this.loadData()
          } catch (e) {
            wx.hideLoading()
            wx.showToast({ title: e.message || '操作失败', icon: 'none', duration: 3000 })
          }
        }
      }
    })
  }
})
