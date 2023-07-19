package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
	"uvm-backend/model"
	"uvm-backend/utils"
)

type Item struct {
	// 订单内容
	ProductId uint
	Number    int
	Price     float64
	Name      string
	ImageUrl  string
}
type OrderStatistics struct {
	// 每个月的订单统计数据
	StartTime  string  `json:"start_time"` // 2006-01
	OrderCount int     `json:"order_count"`
	GoodsCount int     `json:"goods_count"`
	Income     float64 `json:"income"`
}

/**
根据订单号查询订单
*/
func GetOrderByOrderNumber(orderNumber string) (order model.Order, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("srevice.GetOrderByOrderNumber: %w", err)
		}
	}()
	o := model.Order{
		OrderNumber: orderNumber,
	}
	order, err = o.GetOrderByStructQuery()
	if err != nil {
		log.Println(err)
		return model.Order{}, err
	}
	return
}

/**
查询每个月的订单数量、商品数量和总收入
*/
func GetOrderStatisticsByMachineId(machineId uint) (orderStatistics []OrderStatistics, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("srevice.GetOrderByOrderNumber: %w", err)
		}
	}()
	// 查询该售货柜投入运营的时间，即统计开始的时间
	machine := &model.Machine{
		Id: machineId,
	}
	m, err := machine.GetMachineByStructQuery()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// 从售货柜部署当月开始，按月份统计
	deployedTime := m.DeployedTime // 售货柜部署准确时间
	startTime := deployedTime.AddDate(0, 0, -deployedTime.Day()+1)
	for startTime.Before(time.Now()) {
		// 到查询当月为止
		endTime := startTime.AddDate(0, 1, 0)
		log.Println("startTime: ", startTime, " endTime: ", endTime)
		order := &model.Order{}
		orders, err := order.GetMachineOrdersWithinSpan(machineId, true, startTime, endTime) // 该月所有交易成功的订单
		if err != nil {
			log.Println(err)
			return nil, err
		}
		// 该月订单统计信息
		var orderSta OrderStatistics
		orderSta.StartTime = startTime.Format(utils.GetTimeFormat(false)) // 月份
		orderSta.OrderCount = len(orders)                                 // 订单数量
		goodsCount := 0
		income := 0.0
		for _, o := range orders {
			goodsCount += o.Number
			income += o.Price
		}
		orderSta.GoodsCount = goodsCount
		orderSta.Income = income
		orderStatistics = append(orderStatistics, orderSta)
		startTime = endTime
	}
	return
}

/**
查询某商家某日期及某状态的订单列表
*/
func GetOrderByDateNStatus(businessId uint, date string, status bool) (orders []model.Order, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("srevice.GetOrderByDateNStatus: %w", err)
		}
	}()
	// 求时间区间
	startTime, err := time.ParseInLocation(utils.GetTimeFormat(true), date, time.Local) // 当天零点
	if err != nil {
		log.Println(err)
		return nil, err
	}
	endTime := startTime.AddDate(0, 0, 1) // 第二天零点
	log.Println("startTime: ", startTime, " endTime: ", endTime)
	order := &model.Order{}
	orders, err = order.GetBusinessOrdersWithinSpan(businessId, status, startTime, endTime)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return
}
func GetOrdersByUserId(userId uint) (orders []model.Order, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("srevice.GetOrdersByUserId: %w", err)
		}
	}()
	o := &model.Order{
		UserId: userId,
	}
	orders, err = o.GetOrderListByStructQuery()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return
}

func GetErrorOrdersByUserId(userId uint) (orders []model.Order, err error, iferror bool) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("srevice.GetErrorOrdersByUserId: %w", err)
		}
	}()
	o := &model.Order{
		UserId: userId,
	}
	iferror = false
	orders, err = o.GetOrderListByStructQuery()
	for _, order := range orders {
		if !order.Status {
			iferror = true
		}
	}

	if err != nil {
		log.Println(err)
		return nil, err, true
	}
	return
}

/**
添加新的订单：
1. 保存订单数据；
2. 根据订单调整Goods和Product对应记录的Number；
3. 调整user的deal_times和last_deal_time；
4. 更新redis排行榜zset"machineId-year-month"
*/
func AddOrder(orderNumber string, userId uint, machineId uint, businessId uint, orderContent string, number int, price float64, status bool) (Id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("service.AddOrder: %w", err)
		}
	}()
	createdTime := time.Now()
	// 1. 保存订单
	order := &model.Order{
		OrderNumber:  orderNumber,
		UserId:       userId,
		BusinessId:   businessId,
		MachineId:    machineId,
		OrderContent: orderContent,
		Number:       number,
		Price:        price,
		Status:       status,
		CreatedTime:  createdTime,
	}
	Id, err = order.AddOrder()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	// 2. 调整user的DealTimes和LastDealTime
	user := &model.User{
		Id: order.UserId,
	}
	_, err = user.UpdateUserDealInfo(createdTime)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	items := []Item{} // orderContent反序列化
	err = json.Unmarshal([]byte(order.OrderContent), &items)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	// 3. 根据订单调整Goods、Product的Number
	// 4. 更新redis排行榜zset
	for _, item := range items {
		// 根据ProductId和MachineId找到对应Goods，调整Number
		_, err = UpdateGoods(item.ProductId, order.MachineId, item.Number, true)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		// 调整Product的Number
		_, err = UpdateProductNumber(item.ProductId, item.Number, true)
		if err != nil {
			log.Println(err)
			return 0, err
		}
		// 更新redis排行榜
		zSet := utils.GetRankingZSetKey(order.MachineId, createdTime.Format(utils.GetTimeFormat(false)))
		err = model.ZsetAdd(zSet, strconv.Itoa(int(item.ProductId)), item.Number)
		if err != nil {
			log.Println(err)
			return 0, err
		}
	}
	return
}
