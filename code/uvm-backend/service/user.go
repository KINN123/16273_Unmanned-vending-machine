package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"uvm-backend/model"
	"uvm-backend/utils"

	"gorm.io/gorm"
)

//// wx.getUserProfile得到的完整信息
//type WXUser struct {
//	UserInfo      NonSensitiveUser `json:"user_info"`
//	RawData       string     `json:"raw_data"`
//	Signature     string     `json:"signature"`
//	EncryptedData string     `json:"encrypted_data"`
//	IV            string     `json:"iv"`
//}
//// WX UserInfo返回的不包含敏感信息的Non-Sensitive UserInfo
//type NonSensitiveUser struct {
//	OpenID		string		  `json:"open_id, omitempty"`
//	Name        string		  `json:"nick_name"`
//	UnionID		string		  `json:"union_id, omitempty"`
//	AvatarUrl   string 		  `json:"avatar_url"`
//	Gender      int       	  `json:"gender"`
//	Country     string    	  `json:"country"`
//	Province    string    	  `json:"province"`
//	City        string    	  `json:"city"`
//	Language    string        `json:"language"`
//}

// 微信接口服务返回值
// 注意查看api定义的json值
type WXLoginResp struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type ProductInfo struct {
	ProductId int     `json:"productId"`
	ImageUrl  string  `json:"imageUrl"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Number    int     `json:"number"`
}

type Cart struct {
	Change []ProductInfo `json:"change"`
}

/**
根据用户ID返回用户昵称和头像
*/
func GetUserInfo(userId uint) (avatarUrl, nickName string, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.GetUserInfo: %w", err)
		}
	}()
	user := &model.User{
		Id: userId,
	}
	u, err := user.GetUserByID()
	if err != nil {
		log.Println(err)
		return "", "", err
	}
	// 已有该用户记录
	return u.AvatarUrl, u.Nickname, nil
}

/**
请求微信接口服务得到openID和sessionKey等敏感信息
*/
func GetWXSession(code string) (w *WXLoginResp, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.WXLogin: %w", err)
		}
	}()
	appId := "wxece7a11ea50654c3"
	appSecret := "7265dd5dc9728f24f9bbc3d072ec0a68"
	url := "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
	// 合成url, appId和appSecret可以直接得到
	url = fmt.Sprintf(url, appId, appSecret, code)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	// 解析body到结构体中
	w = &WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(w); err != nil {
		log.Println(err)
		return nil, err
	}

	// 判断微信接口服务返回的是否为异常情况
	if w.ErrCode != 0 {
		err = errors.New(fmt.Sprintf("ErrCode: %d ErrMsg: %s", w.ErrCode, w.ErrMsg))
		return nil, err
	}
	// 解析
	return w, err
}

/**
请求中间服务端的购物车数据
*/
func GetCartResp() (cart *Cart, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.GetCartResp: %w", err)
		}
	}()

	url := "http://192.168.1.120:8001/result"
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	cart = &Cart{}
	decoder := json.NewDecoder(resp.Body)
	if err = decoder.Decode(cart); err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(cart)
	return cart, err
}

/**
根据OpenID创建或者更新用户，返回用户ID
*/
func UserLogin(openID, avatarUrl, nickName string) (id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.UserLogin: %w", err)
		}
	}()
	user := &model.User{
		WXOpenId: openID,
	}
	u, err := user.GetUserByID()
	if err == gorm.ErrRecordNotFound {
		// 没有该用户记录，则需初始化该用户
		user.Name = utils.GetUUID()
		user.BusinessId = 1
		user.AvatarUrl = avatarUrl
		user.Nickname = nickName
		id, err = user.AddUser()
		if err != nil {
			// 添加用户出现错误
			log.Println(err)
			return 0, err
		}
		return id, err
	}
	// 查询出现其他错误
	if err != nil {
		log.Println(err)
		return 0, err
	}
	// 已有该用户记录，更新用户
	u.AvatarUrl = avatarUrl
	u.Nickname = nickName
	u, err = u.UpdateUser()
	// 更新出现错误
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return u.Id, nil
}

/**
根据username和password请求开门
*/
func OpenDoorreq(username string, password string, openDoorUrl string) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.openDoorReq: %w", err)
		}
	}()
	data := url.Values{
		"username": {username},
		"password": {password},
	}
	resp, err := http.PostForm(openDoorUrl, data)
	defer resp.Body.Close()

	if err != nil {
		log.Println(err)
		return err
	}
	return err
}
