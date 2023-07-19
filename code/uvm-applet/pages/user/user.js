const config = require("../../config.js");

// pages/user/user.js
Page({

  /**
   * 页面的初始数据
   */
  data: {
    canIUseGetUserProfile: false, // 是否可以使用getUserProfile
    nickName: null,
    avatarUrl: null,
    options: [
      // 注意，主页面之间的跳转不能使用 navigator
      {
        option: "历史订单",
        url: "/pages/history/history", 
      }, {
        option: "联系客服",
        url: "/pages/user/user"
      }, {
        option: "关于我们",
        url: "/pages/user/user"
      }
    ]
  },
  /**
   * 用户保持登录态时根据userID从后台获取用户信息
   */
login: function() {
  // 已经授权但还未获得用户信息
  if(wx.getStorageSync('hasAuth')){
    var userId = wx.getStorageSync('uvmUserId');
    wx.request({
      url: config.baseUrl + '/getUserInfo',
      method: 'GET',
      header: {                                               
        "Content-Type": "application/x-www-form-urlencoded"
      },
      data: {
        userId: userId
      },
      success: res => {
        if(res.data.code == 1){
          wx.setStorageSync('login', 1)
          // 成功读取用户数据
          this.setData({
            avatarUrl: res.data.data.avatarUrl,
            nickName: res.data.data.nickName
          });
        }
        // 读取用户数据失败，仍然需要授权登录
      },
      fail: () => {
        wx.setStorageSync('login', 0)
        console.log("读取后台用户信息失败！")
      }
    })
  }
},
/**
 * 用户登录授权，具体内容包括：
 * 获取用户基本信息：小程序接口更新推荐使用getUserProfile替代getUserInfo来获取用户基本信息；
 * 后台得到用户信息，返回UserID
 */
bindGetUserProfile: function(e){
  wx.getUserProfile({
    // 授权显示用户基本信息
    desc: '授权读取用户信息',
    success: (res) => {
      // console.log(res),
      wx.setStorageSync('hasAuth', 1);
      var avatarUrl = res.userInfo.avatarUrl;
      var nickName = res.userInfo.nickName;
      wx.login({
        // 用户登录：将code和用户信息发送到后台，接收用户ID
        success: res =>  {
        console.log(res.code)
        if(res.code){
          wx.request({
            url: config.baseUrl + '/wxLogin',
            method: 'POST',
            header: {
              "Content-Type": "application/x-www-form-urlencoded"
            },
            data: {
              code: res.code,
              nickName: nickName,
              avatarUrl: avatarUrl
            },
            success: res => {
              console.log(res)
              // fail处理不了404错误，用code来判断
              if(res.data.code == 1){
                // 把用户ID存入本地
                wx.setStorageSync('uvmUserId', res.data.data.id)
                this.setData({
                  avatarUrl: avatarUrl,
                  nickName: nickName
                });          
              }else{
                // 后端登录失败
                wx.showToast({
                  title: '登录失败',
                  icon: 'none'
                })
                console.log("后端登录失败")
              }
            }
          });
        }
      }
      })
    },
    fail: () => {
      console.log("用户拒绝授权");
    }
  })
},
bindGetUserInfo: function(e) {
  wx.setStorageSync('hasAuth', 1);
  var avatarUrl = e.detail.userInfo.avatarUrl;
  var nickName = e.detail.userInfo.nickName;
  wx.login({
    // 用户登录：将code和用户信息发送到后台，接收用户ID
    success: res => {
    console.log(res.code)
    if(res.code){
      wx.request({
        url: config.baseUrl + '/wxLogin',
        method: 'POST',
        header: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        data: {
          code: res.code,
          nickName: nickName,
          avatarUrl: avatarUrl
        },
        success: res => {
          console.log(res)
          // fail处理不了404错误，用code来判断
          if(res.data.code == 1){
            // 把用户ID存入本地
            wx.setStorageSync('uvmUserId', res.data.data.id)
            this.setData({
              avatarUrl: avatarUrl,
              nickName: nickName
            });    
          }else{
            // 后端登录失败
            wx.showToast({
              title: '登录失败',
              icon: 'none'
            })
            console.log("后端登录失败")
          }
        },
      });
    }
  }
  })

},

checkOrder: function() {
  var userId = wx.getStorageSync('uvmUserId');
  wx.request({
    url: config.baseUrl + '/order/errorOrderList',
    method: 'GET', 
    data: {
      UserId: userId
    },
    success: res => {
      console.log(res.data.data);
    //   if(res.data.data.iferror == 1){
    //     wx.setStorageSync('errorOrder', 1)
    //   }
    //   else{
    //     wx.setStorageSync('errorOrder', 0)
    //   }
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
    if(wx.getUserProfile){
      // 版本允许则使用getUserProfile
      this.setData({
        canIUseGetUserProfile: true
      })
      // wx.setStorageSync('canIUseGetUserProfile', 1)
    }
    this.login();

    this.checkOrder();
    console.log(wx.getStorageSync('errorOrder'))
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