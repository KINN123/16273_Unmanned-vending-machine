// pages/details/details.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    id: "",
    img: "",
    name: "",
    price: "",
    number: "",
    info: "",

    imgDownloadPath:  "http://127.0.0.1:8000/product/image/download?url=",
    imgTsValue: "", // 请求url一致时图片会被小程序缓存，无法及时更新；故添加一个时间戳参数
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.setData({
      imgTsValue: "&ts=" + new Date().getTime(), // 获取当前时间戳
      id: options.id,
      img: options.img,
      name: options.name,
      number: options.number,
      price: options.price,
      info: options.info
    })
  },

  /**
   * 生命周期函数--监听页面初次渲染完成
   */
  onReady: function () {

  },

  /**
   * 生命周期函数--监听页面显示
   */
  onShow: function () {

  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {

  },

  /**
   * 生命周期函数--监听页面卸载
   */
  onUnload: function () {

  },

  /**
   * 页面相关事件处理函数--监听用户下拉动作
   */
  onPullDownRefresh: function () {

  },

  /**
   * 页面上拉触底事件的处理函数
   */
  onReachBottom: function () {

  },

  /**
   * 用户点击右上角分享
   */
  onShareAppMessage: function () {

  }
})