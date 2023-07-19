package service

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"uvm-backend/model"
	"uvm-backend/utils"
)

type Ranking struct {
	Name        string `json:"name"`
	EnglishName string `json:"english_name"`
	Sale        int    `json:"sale"`
}

/**
从redis的zSet从获取对应售货柜和月份的商品排行榜
*/
func GetMonthlyRanking(machineId uint, date string) (rankings []Ranking, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.GetMonthlyRanking: %w", err)
		}
	}()
	// 获得zSet的名称
	zSet := utils.GetRankingZSetKey(machineId, date)
	log.Println("zSet: ", zSet)
	productIds, sales := model.ZsetRevRange(zSet, 0, 10) // 获得降序前10的商品id和对应销量
	for index, idStr := range productIds {
		// 根据productIds取出Name和EnglishName
		productId, err := strconv.Atoi(idStr)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		product := &model.Product{
			Id: uint(productId),
		}
		p, err := product.GetProductByStructQuery()
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// 将该排行榜记录加入排行榜list
		rankings = append(rankings, Ranking{p.Name, p.EnglishName, sales[index]})
	}
	return
}
func GetMachinesByBusinessId(businessId uint) (machines []model.Machine, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.GetMachinesByBusinessId: %w", err)
		}
	}()
	machine := &model.Machine{
		BusinessId: businessId,
	}
	machines, err = machine.GetMachineListByStructQuery()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return
}
func AddMachine(businessId uint, location string, modelId uint) (id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.AddMachine: %w", err)
		}
	}()

	machine := &model.Machine{
		BusinessId:   businessId,
		Location:     location,
		ModelId:      modelId,
		DeployedTime: time.Now(),
	}
	id, err = machine.AddMachine()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return
}
