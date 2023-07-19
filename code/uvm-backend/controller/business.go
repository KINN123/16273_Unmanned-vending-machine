package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"uvm-backend/service"
)

/**
根据商家ID，获取其下商品信息
*/
func GetProductsByBusinessId(ctx *gin.Context) {
	businessId, err := strconv.Atoi(ctx.Query("BusinessId"))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	_, _, _, productList, err := service.GetBusinessProductById(uint(businessId))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message":      "读取商家商品列表成功",
		"product_list": productList,
	})
}
func AddBusiness(ctx *gin.Context) {
	var data struct {
		Name string
		Info string
	}
	ctx.ShouldBindJSON(&data)
	id, err := service.AddBusiness(data.Name, data.Info)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "增加成功",
		"id":      id,
		"name":    data.Name,
		"info":    data.Info,
	})
}

func DeleteBusiness(ctx *gin.Context) {
	var data struct {
		ID uint
	}
	ctx.ShouldBindJSON(&data)
	err := service.DeleteBusiness(data.ID)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "删除成功",
	})
}

func UpdateBusiness(ctx *gin.Context) {
	var data struct {
		ID   uint
		Name string
		Info string
	}
	ctx.ShouldBindJSON(&data)
	name, info, t, err := service.UpdateBusiness(data.ID, data.Name, data.Info)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message":       "更新成功",
		"id":            data.ID,
		"name":          name,
		"info":          info,
		"register time": t,
	})
}
