package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"uvm-backend/service"
)

/**
查询Machine某个月份的商品TOP10排行榜
*/
func GetMonthlyRanking(ctx *gin.Context) {
	machineId, err := strconv.Atoi(ctx.Query("MachineId"))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	date := ctx.Query("Date") // 2021-07
	rankings, err := service.GetMonthlyRanking(uint(machineId), date)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message":  "获取月份排行榜成功",
		"rankings": rankings,
	})
}

/**
根据商家Id查询Machine列表
*/
func GetMachinesByBusinessId(ctx *gin.Context) {
	businessId, err := strconv.Atoi(ctx.Query("BusinessId"))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	machines, err := service.GetMachinesByBusinessId(uint(businessId))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message":  "售货柜列表查询成功",
		"machines": machines,
	})
}

/**
增加售货柜
*/
func AddMachine(ctx *gin.Context) {
	var data struct {
		BusinessId uint
		Location   string
		ModelId    uint
	}
	ctx.ShouldBindJSON(&data)
	id, err := service.AddMachine(data.BusinessId, data.Location, data.ModelId)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "售货柜增加成功",
		"id":      id,
	})
}
