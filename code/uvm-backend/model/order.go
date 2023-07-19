package model

import (
	"fmt"
	"log"
	"time"
)

type Order struct {
	Id           uint      `json:"id" gorm:"primaryKey"`
	UserId       uint      `json:"user_id"`
	BusinessId   uint      `json:"business_id"`
	MachineId    uint      `json:"machine_id"`
	OrderNumber  string    `json:"order_number" gorm:"size:50; not null"`
	Status       bool      `json:"status" gorm:"not null"`
	OrderContent string    `json:"order_content" gorm:"size:1000; not null"` // ProductId|Number|Price|Name
	Price        float64   `json:"price" gorm:"default:0; not null"`
	Number       int       `json:"number" gorm:"default:0; not null"`
	CreatedTime  time.Time `json:"created_time"`
}

func (Order) TableName() string {
	return "order"
}

/**
查询商家在某一时间区间内指定交易状态的OrderList
*/
func (o *Order) GetBusinessOrdersWithinSpan(businessId uint, status bool, startTime time.Time, endTime time.Time) (orders []Order, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetBusinessOrdersWithinSpan: %w", err)
		}
	}()
	// select *from orders where business_id = businessId and status = status and created_time between startTime and endTime
	result := DB.Where("business_id = ? AND status = ? AND created_time BETWEEN ? AND ?", businessId, status, startTime, endTime).Find(&orders)
	err = result.Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return
}

/**
查询售货柜在某一时间区间内指定交易状态的OrderList
*/
func (o *Order) GetMachineOrdersWithinSpan(machineId uint, status bool, startTime time.Time, endTime time.Time) (orders []Order, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetMachineOrdersWithinSpan: %w", err)
		}
	}()
	// select *from orders where business_id = businessId and status = status and created_time between startTime and endTime
	result := DB.Where("machine_id = ? AND status = ? AND created_time BETWEEN ? AND ?", machineId, status, startTime, endTime).Find(&orders)
	err = result.Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return
}

/**
结构体查询得到Order
*/
func (o *Order) GetOrderByStructQuery() (orders Order, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetOrderByStructQuery: %w", err)
		}
	}()
	result := DB.Where(o).First(&orders)
	err = result.Error
	if err != nil {
		log.Println(err)
		return Order{}, err
	}
	return
}

/**
结构体查询得到OrderList
*/
func (o *Order) GetOrderListByStructQuery() (orders []Order, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetOrderListByStructQuery: %w", err)
		}
	}()
	result := DB.Where(o).Find(&orders)
	err = result.Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return
}
func (o *Order) AddOrder() (Id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.AddOrder: %w", err)
		}
	}()
	result := DB.Create(o)
	err = result.Error
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return o.Id, nil
}
