package model

import (
	"fmt"
	"log"
	"time"
)

type Business struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"size:50"`
	Info         string    `json:"info"`
	RegisterTime time.Time `json:"register_time"`
	Products     []Product `gorm:"ForeignKey:BusinessId; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Machines     []Machine `gorm:"ForeignKey:BusinessId; constraint:OnUpdate:CASCADE, OnDelete:SET NULL;"`
	Orders       []Order   `gorm:"ForeignKey:BusinessId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Business和Products为一对多
// Business和Machine为一对多

func (Business) TableName() string {
	return "business"
}

/**
根据商家id查询商家，只预加载Product表
*/
func (b *Business) GetBusinessProductById() (business Business, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetBusinessById: %w", err)
		}
	}()
	// 预加载产品
	result := DB.Preload("Products").Preload("Products.Goods").First(&business, b.ID)
	err = result.Error
	if err != nil {
		//log.Println(err)
		return Business{}, err
	}
	return
}

// 添加商家
func (b *Business) AddBusiness() (id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.AddBusiness: %w", err)
		}
	}()
	result := DB.Create(b)
	if result.Error != nil {
		log.Println(result.Error)
		return 0, result.Error
	}
	return b.ID, err
}

// 删除商家
func (b *Business) DeleteBusiness() (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.DeleteBusiness: %w", err)
		}
	}()
	// 判断数据库中是否有这个ID
	result := DB.First(&Business{}, b.ID)
	err = result.Error
	if err != nil {
		return err
	}
	result = DB.Delete(b)
	err = result.Error
	if err != nil {
		return err
	}
	return
}

// 更新商家
func (b *Business) UpdateBusiness() (business Business, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.UpdateBusiness: %w", err)
		}
	}()
	// 判断数据库中是否有这个ID
	result := DB.First(&business, b.ID)
	err = result.Error
	if err != nil {
		return Business{}, err
	}
	// 更新数据
	result = DB.Model(b).Updates(*b)
	err = result.Error
	if err != nil {
		return Business{}, err
	}
	return *b, nil
}
