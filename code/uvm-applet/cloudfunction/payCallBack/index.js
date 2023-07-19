// 云函数入口文件
const cloud = require('wx-server-sdk')

cloud.init({
  // api调用与云函数当前所在环境一致
  env: cloud.DYNAMIC_CURRENT_ENV
})

// 云函数入口函数
exports.main = async (event, context) => {
  // 回调函数
  return {"errcode": 0} // 一定要返回，否则认为回调失败
}