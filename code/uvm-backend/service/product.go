package service

import (
	"fmt"
	"log"
	"time"
	"uvm-backend/model"
)

/**
根据英文名查找商品，并返回列表（不同售货柜下不同项）
*/
func GetProductInfoByEN(englishName string) (p model.Product, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.GetProductInfoByEN: %w", err)
		}
	}()
	product := &model.Product{
		EnglishName: englishName,
	}
	p, err = product.GetProductByStructQuery()
	if err != nil {
		log.Println(err)
		return model.Product{}, err
	}
	return
}

/**
增加商品
*/
func AddProduct(businessId uint, name string, englishName string, info string, price float64, imageUrl string) (id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.AddProduct: %w", err)
		}
	}()
	product := &model.Product{
		BusinessId:  businessId,
		Name:        name,
		EnglishName: englishName,
		Info:        info,
		Price:       price,
		ImageUrl:    imageUrl,
		UpdateTime:  time.Now(),
	}
	id, err = product.AddProduct()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return
}

/**
更新商品基本信息
*/
func UpdateProduct(id uint, businessId uint, name string, englishName string, info string, price float64, imageUrl string) (productId uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.UpdateProduct: %w", err)
		}
	}()
	product := &model.Product{
		Id:          id,
		BusinessId:  businessId,
		Name:        name,
		EnglishName: englishName,
		Info:        info,
		Price:       price,
		ImageUrl:    imageUrl,
		UpdateTime:  time.Now(),
	}
	productId, err = product.UpdateProduct()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return
}

/**
更新商品數量
change为库存改变量
isConsume为1时表示消费该商品，change为消耗量；为0时表示上货，change为上货量

*/
func UpdateProductNumber(id uint, change int, isConsume bool) (productId uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.UpdateProduct: %w", err)
		}
	}()
	// 查找商品，得到原number
	product := &model.Product{
		Id: id,
	}
	p, err := product.GetProductByStructQuery()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	// 更新product的number
	var number int
	if isConsume {
		number = p.Number - change
	} else {
		number = p.Number + change
	}
	// 更新商品
	updatedProduct := &model.Product{
		Id:     id,
		Number: number,
	}
	_, err = updatedProduct.UpdateProduct()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return
}
