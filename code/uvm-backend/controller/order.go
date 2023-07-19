package controller

import (
	"log"
	"strconv"
	"uvm-backend/service"

	"github.com/gin-gonic/gin"
)

/**
根据售货柜Id查询每个月订单数、商品数、总收入
*/
func GetOrderStatisticsByMachineId(ctx *gin.Context) {
	machineId, err := strconv.Atoi(ctx.Query("MachineId"))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	orderStatistics, err := service.GetOrderStatisticsByMachineId(uint(machineId))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message":          "获取订单统计数据成功！",
		"order_statistics": orderStatistics,
	})
}

/**
根据订单号查询订单
*/
func GetOrderByOrderNumber(ctx *gin.Context) {
	orderNumber := ctx.Query("OrderNumber")
	order, err := service.GetOrderByOrderNumber(orderNumber)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "获取订单成功",
		"order":   order,
	})
}

/**
通过UserId查询用户订单列表
*/
func GetOrdersByUserId(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Query("UserId"))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	orders, err := service.GetOrdersByUserId(uint(userId)) // 按id排序
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "获取用户订单列表成功",
		"orders":  orders,
	})
}

/**
查询某商家具体日期和状态的交易订单列表
*/
func GetOrdersByDateNStatus(ctx *gin.Context) {
	businessId, err := strconv.Atoi(ctx.Query("BusinessId"))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	date := ctx.Query("Date") // 格式形如2021-06-30
	// 订单状态
	statusInt, err := strconv.Atoi(ctx.Query("Status"))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	// 转换为bool
	var status bool
	if statusInt != 0 {
		status = true
	} else {
		status = false
	}
	orders, err := service.GetOrderByDateNStatus(uint(businessId), date, status)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "获取指定日期及状态的商家订单列表成功",
		"orders":  orders,
	})

}

/**
添加订单
*/
func AddOrder(ctx *gin.Context) {
	var data struct {
		OrderNumber  string
		UserId       uint
		BusinessId   uint
		MachineId    uint
		Status       bool
		OrderContent string
		Number       int
		Price        float64
	}
	ctx.ShouldBindJSON(&data)
	id, err := service.AddOrder(data.OrderNumber, data.UserId, data.MachineId, data.BusinessId, data.OrderContent, data.Number, data.Price, data.Status)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "订单添加成功",
		"id":      id,
	})

}

/**
通过UserId判断用户订单列表是否存在未完成订单
*/
func GetErrorOrdersByUserId(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Query("UserId"))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	orders, err, iferror := service.GetErrorOrdersByUserId(uint(userId)) // 按id排序
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "获取用户订单列表成功",
		"orders":  orders,
		"iferror": iferror,
	})
}
