const config = require("../../config.js");
// pages/cart/cart.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    imgDownloadPath:  config.ImageDownloadUrl,
    orders: [],
  },
  /**
   * 根据用户Id查询历史订单
   */
  getOrderList: function() {
    var userId = wx.getStorageSync('uvmUserId');
    wx.request({
      url: config.baseUrl + '/order/orderList',
      method: 'GET', 
      data: {
        UserId: userId
      },
      success: res => {
        console.log(res);
        var orders = res.data.data.orders;
        // OrderContent反序列化
        orders.forEach(function(item){
          item.order_content = JSON.parse(item.order_content)
        })
        this.setData({
          orders: orders.reverse() // 时间倒序，新的订单在上面
        })
        console.log(orders)

      },
      fail: res => {
        console.log(res)
      }
    })
  },
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    // 根据userId请求历史订单
    var userId = wx.getStorageSync('uvmUserId');
    this.getOrderList();
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

  },
})