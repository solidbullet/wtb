const cloud = require('wx-server-sdk')

cloud.init({ env: cloud.DYNAMIC_CURRENT_ENV })

exports.main = async (event, context) => {
  const { code } = event

  if (!code) {
    return { code: 400, message: '缺少 code 参数' }
  }

  try {
    const res = await cloud.openapi.phonenumber.getPhoneNumber({
      code: code
    })
    return {
      code: 200,
      phoneNumber: res.phoneInfo.phoneNumber,
      purePhoneNumber: res.phoneInfo.purePhoneNumber,
      countryCode: res.phoneInfo.countryCode
    }
  } catch (err) {
    console.error('获取手机号失败:', err)
    return { code: 500, message: err.message || '获取手机号失败' }
  }
}
