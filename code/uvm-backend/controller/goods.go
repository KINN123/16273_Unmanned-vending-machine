package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"uvm-backend/service"
)

/**
根据售货柜Id，获取商品列表
*/
func GetGoodsByMachineId(ctx *gin.Context) {
	machineId, err := strconv.Atoi(ctx.Query("MachineId"))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	goodsInfoList, err := service.GetGoodsByMachineId(uint(machineId))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message":    "读取售货柜商品列表成功",
		"goods_list": goodsInfoList,
	})

}
