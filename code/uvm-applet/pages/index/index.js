// pages/index/index.js
const app = getApp()
const config = require("../../config.js");
const cart = require("../cart/cart.js");
const utils = require("../../utils/utils.js");
Page({

  /**
   * 页面的初始数据
   */
  data: {
    imagesPath: [
      "/images/pictures/test_001.jpg",
      "/images/pictures/test_002.jpeg",
      "/images/pictures/test_003.jpg",
      "/images/pictures/test_004.jpg"
    ],
    scanImage: "/images/icons/scan_open.png",
    carts:[
      {
        productId: 1, 
        name: "贝奇野菜复合蔬果汁饮品", 
        price: 4,
        number: 5,
        imageUrl: "./upload/img/BQYCFS.jpg",
      },
      {
        productId: 2, 
        name: "百事可乐Pepsi", 
        price: 4.5,
        number: 3,
        imageUrl: "./upload/img/BSKL.jpg",
      },
      // {
      //   productId: 4, 
      //   name: "百怡百香果饮料", 
      //   price: 5.5,
      //   number: 3,
      //   imageUrl: "./upload/img/BYBXG.jpg",
      // }
    ]
  },

  onTap: function(e) {
    console.log("开门")
    let item = e.currentTarget.dataset.item
    let obj = {
      [`${item}`]: 1
    }
    let payload = JSON.stringify(obj)
    JSON.parse
    console.log(payload)
    wx.showLoading()
    wx.cloud.callFunction({
      name: 'iothub-publish',
      data: {
        SecretId: app.globalData.secretId,
        SecretKey: app.globalData.secretKey,
        ProductId: app.globalData.productId,
        DeviceName: app.globalData.deviceName,
        Topic: app.globalData.productId + "/" + app.globalData.deviceName + "/data",
        Payload: payload,
      },

     
      success: res => {
        wx.showToast({
          icon: 'none',
          title: 'publish完成',
        })
        console.log("res:", res)
         //3s关门  
        //  setTimeout(function() {
        //     console.log('我是xx')
        //   }, 3000); 
        setTimeout(delayedFunction(e), 300000);
        // setTimeout({}, 30);
        // delayedFunction(e);
      },
      fail: err => {
        wx.showToast({
          icon: 'none',
          title: 'publish失败，请连接设备',
        })
        console.error('[云函数] [iotexplorer] 调用失败：', err)
      }
    })  
    
    // wx.request({
    //   url: 'http://192.168.1.109:8000/openDoor',
    //   method:'POST',
    //   data:{
    //     "username": "rocketeerli",
    //     "password": "rocketeerli"
    //   },
    //   header: {
    //     'content-type': 'application/json'
    //   },
    //   success: function(res) {
    //     console.log(res.data)
    //   }
    // })
    function delayedFunction(e) {
        setTimeout(console.log("关门"), 3000);
        console.log("关门")
    let item = e.currentTarget.dataset.item
    let obj = {
      [`${item}`]: 0
    }
    let payload = JSON.stringify(obj)
    JSON.parse
    console.log(payload)
    wx.showLoading()
    wx.cloud.callFunction({
      name: 'iothub-publish',
      data: {
        SecretId: app.globalData.secretId,
        SecretKey: app.globalData.secretKey,
        ProductId: app.globalData.productId,
        DeviceName: app.globalData.deviceName,
        Topic: app.globalData.productId + "/" + app.globalData.deviceName + "/data",
        Payload: payload,
      },

     
      success: res => {
        wx.showToast({
          icon: 'none',
          title: 'publish完成',
        })
        console.log("res:", res)
         
      },
      fail: err => {
        wx.showToast({
          icon: 'none',
          title: 'publish失败，请连接设备',
        })
        console.error('[云函数] [iotexplorer] 调用失败：', err)
      }
    })  

  }
  },
  

  /**
   * 生命周期函数--监听页面加载
   */
  onLoad: function (options) {

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
  /**
   * 扫码
   */
  scanCodeEvent: function(){
    //若用户已登录，授权扫码
    //不能登录则提醒先登录
    if(wx.getStorageSync('login') == 1){
      if(wx.getStorageSync('errorOrder') == 1){
        wx.scanCode({
          // onlyFromCamera: true,
          success(res){
            console.log("扫码成功：" + JSON.stringify(res))
            console.log("开门")
            wx.request({
              // url: 'http://192.168.1.109:8000/openDoor',
              url: config.baseUrl + '/openDoor',
              method:'POST',
              data:{
                "username": "rocketeerli",
                "password": "rocketeerli",
                "openDoorUrl": res.result,
              },
              header: {
                // 'content-type': 'application/json'
                'content-type':'application/x-www-form-urlencoded'
              },
              success: function(res) {
                console.log(res.data)
              },
              fail: function(res) {
                console.log(res)
              }
            })
          },
          fail(res){
              wx.showToast({
                title: '扫码失败！',
                icon: 'error'
              })
              console.log("扫码失败：" + JSON.stringify(res))
          }
        })
      }
      else{
        wx.showToast({
          title: '您还有未支付的订单，请先去历史订单中完成付款',
          icon: 'none'
        })
        //未支付完成的订单强制跳出
        var body = "测试订单"; // 订单支付内容
        var that = this;
        var orderId = wx.getStorageSync('errorOrderId');
        wx.cloud.callFunction({
          name: "pay",
          data: {
            body: body,
            orderId: "" + orderId, // 不大于32位的商户内部订单号
            // money: that.data.totalMoney * 100, // 支付金额，接口单位是分，故*100
            money: 0.00 * 100, // 支付金额，接口单位是分，故*100
            nonceStr: utils.randomStr(32), // 不长于32位的随机字符串
          },
          success:(res) => {
            wx.hideLoading({
              complete: (res) => {},
            })
            console.log("订单提交成功！", res.result)
            cart.pay(res.result, orderId)
          },
          fail:(res) => {
            wx.hideLoading({
              complete: (res) => {},
            })
            console.log("订单提交失败！", res)
          }
        })
      }
    }
    else{
      wx.showToast({
        title: '请先登录',
        icon: 'error',
      })
    }
  },

  //cart.js中的支付API
  // pay: function(payData, orderId){
  //   const payment = payData.payment;
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
  //   wx.requestPayment({
  //     ...payment,
  //     success:(res) => {
  //       console.log("支付成功！", res)
  //       // 保存订单
  //       wx.request({
  //         url: config.baseUrl + '/order/add',
  //         method: 'POST',
  //         data: {
  //           // 订单号、内容、总价、订单状态、用户&商家&售货柜ID
  //           OrderNumber: "" + orderId,
  //           OrderContent: JSON.stringify(carts),
  //           Price: totalMoney,
  //           Number: totalNumber,
  //           UserId: userId, 
  //           MachineId: machineId, 
  //           BusinessId: businessId,
  //           Status: true,// 已完成
  //         },
  //         header: {
  //           "Content-Type": "application/json",
  //         },
  //         success: () => {
  //           console.log("保存订单成功！")
  //         },
  //         fail: () => {
  //           console.log("保存订单失败！")
  //         }
  //       })    
  //       // 跳转到支付成功页面
  //     },
  //     fail:(res) => {
  //       console.log(payment)
  //       console.log("支付失败！", res)
  //       wx.setStorageSync('errorOrderId', orderId)
  //       // 跳转到支付失败页面
  //       wx.request({
  //         url: config.baseUrl + '/order/add',
  //         method: 'POST',
  //         data: {
  //           // 订单号、内容、总价、订单状态、用户&商家&售货柜ID
  //           OrderNumber: "" + orderId,
  //           OrderContent: JSON.stringify(carts),
  //           Price: totalMoney,
  //           Number: totalNumber,
  //           UserId: userId, 
  //           MachineId: machineId, 
  //           BusinessId: businessId,
  //           Status: false,// 未完成
  //         },
  //         header: {
  //           "Content-Type": "application/json",
  //         },
  //         success: () => {
  //           console.log("保存订单成功！")
  //         },
  //         fail: () => {
  //           console.log("保存订单失败！")
  //         }
  //       })    
  //       // 跳转到支付成功页面     
  //     }
  //   })
  // }
})