// 云函数入口文件
const cloud = require('wx-server-sdk')

cloud.init({
  env: cloud.DYNAMIC_CURRENT_ENV
})

// 云函数入口函数
exports.main = async (event, context) => {
  // 调用统一下单接口
  const res = await cloud.cloudPay.unifiedOrder({
    "body": event.body, // 商品简单描述
    "outTradeNo": event.orderId, // 系统内部订单号，32个字符内
    "spbillCreateIp": "127.0.0.1", // 调用微信支付API的机器IP
    "subMchId": "1607170016", // 微信支付分配的子商户号
    "totalFee": parseInt(event.money), // 必须是int
    "envId": "cloud1", // 回调函数所属云环境Id
    "functionName": "payCallBack", // 回调函数名
    // 以下两个参数必传
    "nonceStr": event.nonceStr, // 不长于32位的随机字符串
    "tradeType": "JSAPI"
  })
  return res
}