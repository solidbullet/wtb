const cloud = require('wx-server-sdk')
cloud.init({ env: cloud.DYNAMIC_CURRENT_ENV })

exports.main = async (event, context) => {
  const { scene, page = 'pages/order/menu', envVersion = 'release', width = 280 } = event

  if (!scene) {
    return { code: 40001, message: 'scene 不能为空' }
  }

  try {
    // 调用微信接口生成小程序码（getwxacodeunlimited）
    const result = await cloud.openapi.wxacode.getUnlimited({
      scene,
      page,
      checkPath: false,
      envVersion,
      width
    })

    // 返回 base64，方便 Go 后台或其他端接收后保存
    const base64 = result.buffer.toString('base64')

    return {
      code: 200,
      message: 'ok',
      data: {
        base64,
        contentType: 'image/png'
      }
    }
  } catch (err) {
    console.error('generateWxaCode error:', err)
    return {
      code: 50001,
      message: err.message || '生成小程序码失败'
    }
  }
}
