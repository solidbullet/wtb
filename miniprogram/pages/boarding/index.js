Page({
  data: {
    scenes: [
      { icon: '🌳', title: '杨梅树下', desc: '夏日乘凉，树荫斑驳，狗狗在天然氧吧里打盹' },
      { icon: '🏕️', title: '千平草坪', desc: '白天尽情奔跑撒欢，释放天性，不做笼中犬' },
      { icon: '🛏️', title: '温馨寄养间', desc: '夜晚独立空间安睡，空调恒温，专人看护' },
      { icon: '📹', title: '24h 直播', desc: '随时随地手机查看，主人安心，狗狗开心' }
    ],
    schedule: [
      { time: '08:00', event: ' Morning 放风', detail: '打开围栏，狗狗们涌入大草坪' },
      { time: '12:00', event: '午餐 & 午休', detail: '回到寄养间用餐，午后树下小憩' },
      { time: '15:00', event: '下午茶时光', detail: '草坪互动游戏，社交玩耍' },
      { time: '18:00', event: '晚餐 & 回家', detail: '营养晚餐后回寄养间休息' }
    ]
  },
  goBack() { wx.navigateBack() },
  makeCall() {
    wx.makePhoneCall({ phoneNumber: '13800138000' })
  }
})
