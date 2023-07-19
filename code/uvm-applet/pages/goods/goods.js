const config = require("../../config.js");
// pages/goods/goods.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    goodsList: null,
    imgDownloadPath:  config.ImageDownloadUrl,
    imgTsValue: "", // 请求url一致时图片会被小程序缓存，无法及时更新；故添加一个时间戳参数
  },
  getGoodsList: function() {
    var machineId = 1;
    this.setData({
      imgTsValue: "&ts=" + new Date().getTime() // 获取当前时间戳
    })
    wx.request({
      url: config.baseUrl + '/machine/goodsList',
      method: 'GET',
      data: {
        "MachineId": machineId
      },
      success: res => {
        console.log(res)
        this.setData({
          goodsList: res.data.data.goods_list
        })
      },
      fail:()=>{
        wx.showToast({
          title: '读取商品列表失败！',
          icon: 'error',
        })
      }
    })
  },
  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.getGoodsList();
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