package model

import (
	"fmt"
	"log"
)

/**
具体售货柜的商品
*/
type Goods struct {
	// 对MachineId和ProductId创建复合索引
	Id        uint `json:"id" gorm:"primaryKey"`
	MachineId uint `json:"machine_id" gorm:"index:idx_goods"`
	Number    int  `json:"number" gorm:"not null"`
	ProductId uint `json:"product_id" gorm:"index:idx_goods"`
}

/**
添加商品
*/
func (g *Goods) AddGoods() (id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.AddGoods: %w", err)
		}
	}()
	result := DB.Create(g)
	if result.Error != nil {
		log.Println(result.Error)
		return 0, result.Error
	}
	return g.Id, err
}

/**
结构体查询Goods
*/
func (g *Goods) GetGoodsByStructQuery() (goods Goods, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetGoodsByStructQuery: %w", err)
		}
	}()
	result := DB.Where(g).First(&goods)
	err = result.Error
	if err != nil {
		log.Println(err)
		return Goods{}, err
	}
	return
}

/**
结构体查询GoodsList
*/
func (g *Goods) GetGoodsListByStructQuery() (goods []Goods, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetGoodsListByStructQuery: %w", err)
		}
	}()
	result := DB.Where(g).Find(&goods)
	err = result.Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return
}

/**
更新Goods
*/
func (g *Goods) UpdateGoods() (goods Goods, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.UpdateGoods: %w", err)
		}
	}()
	// 查询是否有该记录
	result := DB.First(&goods, g.Id)
	if result.Error != nil {
		log.Println(result.Error)
		return Goods{}, err
	}
	// 更新
	result = DB.Model(g).Updates(*g)
	if result.Error != nil {
		log.Println(result.Error)
		return Goods{}, err
	}
	return *g, nil
}
