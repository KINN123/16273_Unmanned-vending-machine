package service

import (
	"fmt"
	"log"
	"uvm-backend/model"
)

type GoodsInfo struct {
	//售货柜对应商品的商品详细信息
	ProductId   uint    `json:"product_id"` // 产品编号
	Name        string  `json:"name"`
	EnglishName string  `json:"english_name"`
	Info        string  `json:"info"`
	Number      int     `json:"number"`
	Price       float64 `json:"price"`
	ImageUrl    string  `json:"image_url"`
}

/**
根据售货柜Id，返回售货柜商品完整信息（包括goods和product的信息）
*/
func GetGoodsByMachineId(id uint) (goodsInfoList []GoodsInfo, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.GetMachineById: %w", err)
		}
	}()
	g := &model.Goods{
		MachineId: id,
	}
	goodsList, err := g.GetGoodsListByStructQuery()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// goods中只保存了售货柜中商品数量这一信息，故查询对应的product，将完整信息返回
	for _, goods := range goodsList {
		product := &model.Product{
			Id: goods.ProductId,
		}
		p, err := product.GetProductByStructQuery()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		goodsInfo := GoodsInfo{goods.ProductId, p.Name, p.EnglishName, p.Info, goods.Number, p.Price, p.ImageUrl}
		goodsInfoList = append(goodsInfoList, goodsInfo)
	}
	// 根据goods
	return
}

/**
新增Goods，需要更新对应Product的number
*/
func AddGoods(machineId uint, productId uint, number int) (id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.AddGoods: %w", err)
		}
	}()
	// 添加Goods对象
	goods := &model.Goods{
		MachineId: machineId,
		ProductId: productId,
		Number:    number,
	}
	id, err = goods.AddGoods()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	// 增加Product的Number
	_, err = UpdateProductNumber(uint(productId), number, false)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return
}

/**
根据ProductId和MachineId查询Goods，更新Goods
change为库存改变量
isConsume为1时表示消费该商品，change为消耗量；为0时表示上货，change为上货量
*/
func UpdateGoods(productId uint, machineId uint, change int, isConsume bool) (id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.GetGoodsByProductNMachineId: %w", err)
		}
	}()
	// 根据productId和machineId查找goods
	goods := &model.Goods{
		ProductId: productId,
		MachineId: machineId,
	}
	g, err := goods.GetGoodsByStructQuery()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	// 更新goods的number
	var number int
	if isConsume {
		number = g.Number - change
	} else {
		number = g.Number + change
	}
	updatedGoods := &model.Goods{
		Id:     g.Id,
		Number: number,
	}
	updatedG, err := updatedGoods.UpdateGoods()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return updatedG.Id, nil
}
