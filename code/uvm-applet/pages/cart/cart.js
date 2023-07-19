const utils = require("../../utils/utils.js");
const config = require("../../config.js");
const { cartUrl } = require("../../config.js");
// pages/cart/cart.js
require("../../utils/utils.js")
Page({
  /**
   * 页面的初始数据
   */
  data: {
    imgDownloadPath:  config.ImageDownloadUrl,

    totalMoney: 0,
    realTime: null,//实时数据对象(用于关闭实时刷新方法)
    change: {},
    // 订单内容
    carts:[
    //   {
    //     productId: 1, 
    //     name: "贝奇野菜复合蔬果汁饮品", 
    //     price: 4,
    //     number: 5,
    //     imageUrl: "./upload/img/BQYCFS.jpg",
    //   },
      {
        productId: 2, 
        name: "百岁山", 
        price: 3,
        number: 1,
        imageUrl: "./upload/img/BSS.jpg",
      }
    //   {
    //     productId: 4, 
    //     name: "百怡百香果饮料", 
    //     price: 5.5,
    //     number: 3,
    //     imageUrl: "./upload/img/BYBXG.jpg",
    //   }
    ]
  },

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {
    this.updateTotalPrice();
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
    var self = this
    /**
     * 防止用户拿不到最新数据(因为页面切换会重新计时)
     * 无条件请求一次最新数据
     */
    wx.request({
      url: config.cartUrl,
      success: function (res) {
        self.setData({
          carts: res.data.change,
        })
        self.updateTotalPrice();
      }
    })
    /**
     * 每隔一段时间请求服务器刷新数据(客户状态)
     * 当页面显示时开启定时器(开启实时刷新)
     * 每隔1分钟请求刷新一次
     * @注意：用户切换后页面会重新计时
     */
    this.data.realTime = setInterval(function () {
      // 请求服务器数据
      wx.request({
        // url: config.baseUrl + '/getCart',
        url: config.cartUrl,
        success: function (res) {
          console.log(res)
          self.setData({
            carts: res.data.change,
          })
        }
      })
      // 反馈提示
      wx.showToast({
        title: '数据已更新！'
      })
      self.updateTotalPrice();
    }, 6000)//间隔时间
    // 更新数据
    this.setData({
      realTime: this.data.realTime,//实时数据对象(用于关闭实时刷新方法)
    })
  },

  /**
   * 生命周期函数--监听页面隐藏
   */
  onHide: function () {
    /**
     * 当页面隐藏时关闭定时器(关闭实时刷新)
     * 切换到其他页面了
     */
    clearInterval(this.data.realTime)
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
/**
 * 基于云开发的微信支付
 */
// 调用云开发统一下单函数，得到支付API所需参数payment
  submitOrder: function() {
    // 测试订单添加
    // this.addOrderTest(this.getRandomOrderId());
    // 测试订单添加

    
    wx.showLoading({
      title: '加载中……',
    })
    let that = this;
    var body = "测试订单"; // 订单支付内容
    var orderId = that.getRandomOrderId();
    wx.cloud.callFunction({
      name: "pay",
      data: {
        body: body,
        orderId: "" + orderId, // 不大于32位的商户内部订单号
        // money: that.data.totalMoney * 100, // 支付金额，接口单位是分，故*100
        money: 0.01 * 100, // 支付金额，接口单位是分，故*100
        nonceStr: utils.randomStr(32), // 不长于32位的随机字符串
      },
      success:(res) => {
        wx.hideLoading({
          complete: (res) => {},
        })
        console.log("订单提交成功！", res.result)
        that.pay(res.result, orderId)
      },
      fail:(res) => {
        wx.hideLoading({
          complete: (res) => {},
        })
        console.log("订单提交失败！", res)
      }
    })
  },
  // // 添加订单测试
  // addOrderTest: function(orderId){
  //   var totalMoney = this.data.totalMoney;
  //   var userId = 1;
  //   var machineId = 1;
  //   var businessId = 1;
  //   var carts = this.data.carts;
  //   var totalNumber = 0;
  //   carts.forEach(function(item){
  //     totalNumber += item.number
  //   });
  //   console.log("totalNumber: ", totalNumber, " totalPrice: ", totalMoney)
  //   wx.request({
  //     url: config.baseUrl + '/order/add',
  //     method: 'POST',
  //     data: {
  //       // 订单号、内容、总价、订单状态、用户&商家&售货柜ID
  //       OrderNumber: "" + orderId,
  //       OrderContent: JSON.stringify(carts),
  //       Price: totalMoney,
  //       Number: totalNumber,
  //       UserId: userId, 
  //       MachineId: machineId, 
  //       BusinessId: businessId,
  //       Status: true,// 已完成
  //     },
  //     header: {
  //       "Content-Type": "application/json",
  //     },
  //     success: () => {
  //       console.log("保存订单成功！")
  //     },
  //     fail: () => {
  //       console.log("保存订单失败！")
  //     }
  //   })
  // },
// 调用支付API
  pay: function(payData, orderId){
    const payment = payData.payment;
    var totalMoney = this.data.totalMoney;
    var userId = 1;
    var machineId = 1;
    var businessId = 1;
    var carts = this.data.carts;
    var totalNumber = 0;
    carts.forEach(function(item){
      totalNumber += item.number
    });
    console.log("totalNumber: ", totalNumber, " totalPrice: ", totalMoney)
    wx.requestPayment({
      ...payment,
      success:(res) => {
        console.log("支付成功！", res)
        // 保存订单
        wx.request({
          url: config.baseUrl + '/order/add',
          method: 'POST',
          data: {
            // 订单号、内容、总价、订单状态、用户&商家&售货柜ID
            OrderNumber: "" + orderId,
            OrderContent: JSON.stringify(carts),
            Price: totalMoney,
            Number: totalNumber,
            UserId: userId, 
            MachineId: machineId, 
            BusinessId: businessId,
            Status: true,// 已完成
          },
          header: {
            "Content-Type": "application/json",
          },
          success: () => {
            console.log("保存订单成功！")
          },
          fail: () => {
            console.log("保存订单失败！")
          }
        })    
        // 跳转到支付成功页面
      },
      fail:(res) => {
        console.log(payment)
        console.log("支付失败！", res)
        wx.setStorageSync('errorOrderId', orderId)
        // 跳转到支付失败页面
        wx.request({
          url: config.baseUrl + '/order/add',
          method: 'POST',
          data: {
            // 订单号、内容、总价、订单状态、用户&商家&售货柜ID
            OrderNumber: "" + orderId,
            OrderContent: JSON.stringify(carts),
            Price: totalMoney,
            Number: totalNumber,
            UserId: userId, 
            MachineId: machineId, 
            BusinessId: businessId,
            Status: false,// 未完成
          },
          header: {
            "Content-Type": "application/json",
          },
          success: () => {
            console.log("保存订单成功！")
          },
          fail: () => {
            console.log("保存订单失败！")
          }
        })    
        // 跳转到支付成功页面     
      }
    })
  },
  /**
   * 随机数+系统时间戳生成不重复的商品订单号
   */
  getRandomOrderId: function() {
    const date = new Date();
    let year = `${date.getFullYear()}`;
    let month = `${date.getMonth() + 1 < 10 ? `0${date.getMonth() + 1}`: date.getMonth() + 1}`;
    let day = `${date.getDate() < 10 ? `0${date.getDate()}`: date.getDate()}`;
    let hour = `${date.getHours() < 10 ? `0${date.getHours()}`: date.getHours()}`;
    let minute = `${date.getMinutes() < 10 ? `0${date.getMinutes()}`: date.getMinutes()}`;
    let second = `${date.getSeconds() < 10 ? `0${date.getSeconds()}`: date.getSeconds()}`;
    let formatDate = `${year}${month}${day}${hour}${minute}${second}`;
    return `${Math.round(Math.random() * 1000)}${formatDate + Math.round(Math.random() * 89 + 100).toString()}`;
  },
  toPay: function() {
    wx.requestPayment(
      {
      'timeStamp': '',
      'nonceStr': '',
      'package': '',
      'signType': 'MD5',
      'paySign': '',
      'success':function(res){
        console.log(res)
      },
      fail: function(res) {
        wx.showModal({
          title:'支付提示',
          content:'<text>',
          showCancel: false
        })
      },
      'complete':function(res){}
      })
  },
  /**
   * 更新总价
   */
  updateTotalPrice: function() {
    var price = 0;
    if (this.data.carts) {
      this.data.carts.forEach(function(item) {
        price += item.price * item.number;
      });
    }
    this.setData({
      totalMoney: price
    })
  }
})

