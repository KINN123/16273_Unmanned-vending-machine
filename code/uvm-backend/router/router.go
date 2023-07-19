package router

import (
	. "uvm-backend/controller"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	//// 用户
	//登录
	r.POST("/user/wxLogin", WXLogin)
	r.GET("/user/getUserInfo", GetUserInfo)
	// 订单
	r.POST("/user/order/add", AddOrder)
	r.GET("/user/order/orderList", GetOrdersByUserId)
	r.GET("/user/order/errorOrderList", GetErrorOrdersByUserId)
	// 商品
	r.GET("/user/machine/goodsList", GetGoodsByMachineId)
	// 生成购物车
	r.GET("/user/getCart", GetCart)
	//扫码开门
	r.POST("/user/openDoor", OpenDoor)

	//// 商家
	r.POST("/business/add", AddBusiness)
	r.POST("/business/delete", DeleteBusiness)
	r.POST("/business/update", UpdateBusiness)
	// 商品
	r.GET("/business/product/productList", GetProductsByBusinessId)
	r.POST("/business/product/add", AddProduct)
	r.POST("/business/product/update", UpdateProduct)
	r.GET("/business/product/getInfoByEN", GetProductInfoByEN)
	// 订单
	r.GET("/business/order/orderList", GetOrdersByDateNStatus)
	r.GET("/business/order/detail", GetOrderByOrderNumber)
	// 售货柜
	r.POST("/business/machine/add", AddMachine)
	r.GET("/business/machine/machineList", GetMachinesByBusinessId)
	r.GET("/business/machine/goodsList", GetGoodsByMachineId)
	r.GET("/business/machine/orderStatistics", GetOrderStatisticsByMachineId)
	r.GET("/business/machine/monthlyRanking", GetMonthlyRanking)
	//// 商品图片下载
	r.GET("/product/image/download", Download)
	r.Run(":8004")

	return r
}
