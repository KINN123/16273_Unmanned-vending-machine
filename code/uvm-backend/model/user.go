package model

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id           uint       `json:"id" gorm:"primaryKey"`
	Name         string     `json:"name" gorm:"size:50"`
	DealTimes    int        `json:"deal_times" gorm:"default:0;not null"`
	BusinessId   uint       `json:"business_id" gorm:"not null" sql:"type:integer constraint fk_product_business REFERENCES business(id)"`
	Business     Business   `json:"business" gorm:"ForeignKey:BusinessId;AssociationForeignKey:ID"`
	LastDealTime *time.Time `json:"last_deal_time"` // *time.Time允许空值
	WXOpenId     string     `json:"open_id" gorm:"not null;index:openid_idx"`
	Nickname     string     `json:"nickName" gorm:"size:50"`
	AvatarUrl    string     `json:"avatarUrl"`
	Orders       []Order    `gorm:"ForeignKey:UserId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

//SessionId     string    `json:"session_id" gorm:"not null;index:openid_idx"`

func (User) TableName() string {
	return "user"
}

//// 第三方自定义Session
//type CustomSession struct {
//	OpenId 		string		`json:"openid"`
//	SessionKey  string 		`json:"session_key"`
//}
/**
下单后，要更新用户的订单信息
*/
func (u *User) UpdateUserDealInfo(dealTime time.Time) (user User, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.UpdateUserDealInfo: %w", err)
		}
	}()
	// 查询user
	result := DB.Where(u).Find(&user)
	err = result.Error
	if err != nil {
		log.Println(err)
		return User{}, err
	}
	// 交易数+1，更新上一次交易时间
	user.DealTimes = user.DealTimes + 1
	user.LastDealTime = &dealTime
	result = DB.Model(u).Updates(user)
	err = result.Error
	if err != nil {
		log.Println(err)
		return User{}, err
	}
	return user, nil
}

/**
更新User记录
*/
func (u *User) UpdateUser() (user User, err error) {
	defer func() {
		if err != nil && err != gorm.ErrRecordNotFound {
			// 没有记录时不修改，否则service不好判断
			err = fmt.Errorf("model.UpdateUser: %w", err)
		}
	}()
	//result := DB.Debug().Save(u)
	result := DB.Model(u).Updates(*u)
	err = result.Error
	if err != nil {
		return User{}, err
	}
	return *u, nil
}

/**
根据WXOpenId或者Id（候选键）索引User对象
*/
func (u *User) GetUserByID() (user User, err error) {
	defer func() {
		if err != nil && err != gorm.ErrRecordNotFound {
			// 没有记录时不修改，否则service不好判断
			err = fmt.Errorf("model.GetUserByOpenId: %w", err)
		}
	}()
	result := DB.Where(u).First(&user)
	err = result.Error
	if err != nil {
		//log.Println(err)
		return User{}, err
	}
	return
}

/**
新增用户
*/
func (u *User) AddUser() (id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.AddUser: %w", err)
		}
	}()
	result := DB.Create(u)
	if result.Error != nil {
		log.Println(result.Error)
		return 0, result.Error
	}
	return u.Id, err
}
