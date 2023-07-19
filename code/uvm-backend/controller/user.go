package controller

import (
	"log"
	"strconv"
	"uvm-backend/service"

	"github.com/gin-gonic/gin"
)

/**
用户保持登录态；根据用户ID获取avatarUrl和nickName
*/
func GetUserInfo(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Query("userId"))
	// 数据类型转换
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	log.Println(userId)
	avatarUrl, nickName, err := service.GetUserInfo(uint(userId))
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message":   "读取用户基本信息成功",
		"avatarUrl": avatarUrl,
		"nickName":  nickName,
	})
}

/**
用户授权登录：保存用户openID、avatarUrl和nickName
*/
func WXLogin(ctx *gin.Context) {
	// 报文格式：x-www-form-urlencoded
	code := ctx.PostForm("code")
	avatarUrl := ctx.PostForm("avatarUrl")
	nickName := ctx.PostForm("nickName")
	log.Println(code)
	// 根据code获取敏感信息：openID和session_key
	wxLoginResp, err := service.GetWXSession(code)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	// 根据openID创建/更新用户
	id, err := service.UserLogin(wxLoginResp.OpenID, avatarUrl, nickName)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "登录成功",
		"id":      id,
	})
}

/**
获取购物车信息
*/
func GetCart(ctx *gin.Context) {
	cartResp, err := service.GetCartResp()
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "自动生成购物车成功",
		"change":  cartResp,
	})

}

/**
请求开门
*/
func OpenDoor(ctx *gin.Context) {
	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	openDoorUrl := ctx.PostForm("openDoorUrl")
	err := service.OpenDoorreq(username, password, openDoorUrl)
	if err != nil {
		log.Println(err)
		ErrorResponse(ctx, err)
		return
	}
	SuccessResponse(ctx, gin.H{
		"message": "开门成功",
	})
}
