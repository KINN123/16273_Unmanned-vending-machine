
package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"path"
	"strconv"
	"uvm-backend/service"
	"uvm-backend/utils"
)

/**
根据商品的英文名查找商品，并返回信息
*/
func GetProductInfoByEN(ctx *gin.Context) {
	englishName := ctx.Query("EnglishName")
	product, err := service.GetProductInfoByEN(englishName)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "根据英文名查找商品成功",
		"product": product,
	})
}

/**
增加商品
*/
func AddProduct(ctx *gin.Context) {
	// 读取form data
	businessId, _ := strconv.Atoi(ctx.PostForm("BusinessId"))
	name := ctx.PostForm("Name")
	englishName := ctx.PostForm("EnglishName")
	info := ctx.PostForm("Info")
	price, _ := strconv.ParseFloat(ctx.PostForm("Price"), 64)
	// 读取文件
	_, header, err := ctx.Request.FormFile("upload")
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	filename := header.Filename
	// 把文件存入/upload/img/EnglishName.扩展名 路径下
	_, suffix := utils.GetFileNameAndSuffix(filename)
	filePath := "./upload/img/" + englishName + suffix
	log.Println("controller.AddProduct:\t", "uploadFileName:\t", filename, "saveFilePath:\t", filePath)
	err = ctx.SaveUploadedFile(header, filePath)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	// 将商品信息存入数据库
	log.Println("controller.AddProduct:\t", "BusinessId:\t", businessId, " Name:\t", name, "EnglishName:\t", englishName, " Info:\t", info, " Price:\t", price, "imgUrl:\t", filePath)
	productId, err := service.AddProduct(uint(businessId), name, englishName, info, price, filePath)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "商品增加成功",
		"id":      productId,
	})
}

/**
更新商品基本信息：name，englishName，info, price
*/
func UpdateProduct(ctx *gin.Context) {
	// 读取form data
	id, _ := strconv.Atoi(ctx.PostForm("id"))
	businessId, _ := strconv.Atoi(ctx.PostForm("BusinessId"))
	name := ctx.PostForm("Name")
	englishName := ctx.PostForm("EnglishName")
	info := ctx.PostForm("Info")
	price, _ := strconv.ParseFloat(ctx.PostForm("Price"), 64)
	// 读取文件
	_, header, err := ctx.Request.FormFile("upload")
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	filename := header.Filename
	// 把文件存入/upload/img/EnglishName.扩展名 路径下
	_, suffix := utils.GetFileNameAndSuffix(filename)
	filePath := "./upload/img/" + englishName + suffix
	log.Println("controller.UpdateProduct:\t", "uploadFileName:\t", filename, "saveFilePath:\t", filePath)
	err = ctx.SaveUploadedFile(header, filePath)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	// 更新数据库中的商品信息
	log.Println("controller.UpdateProduct:\t", "id:\t", id, " BusinessId:\t", businessId, " Name:\t", name, "EnglishName:\t", englishName, "Info:\t", info, "Price:\t", price, "imgUrl:\t", filePath)
	productId, err := service.UpdateProduct(uint(id), uint(businessId), name, englishName, info, price, filePath)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "商品更新成功",
		"id":      productId,
	})
}

/**
服务器端商品图片下载接口
*/
func Download(ctx *gin.Context) {
	// *匹配不到./……手动加上（
	filePath := ctx.Query("url")
	fileName := path.Base(filePath)
	log.Println("controller.Download:\t", "filePath:\t", filePath, "fileName:\t", fileName)
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.File(filePath)
	//return
}
