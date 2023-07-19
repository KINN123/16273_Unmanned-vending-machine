package model

import (
	"fmt"
	"log"
	"time"
)

type Machine struct {
	Id           uint      `json:"id" gorm:"primaryKey"`
	BusinessId   uint      `json:"business_id"`
	Location     string    `json:"location" gorm:"size:100;not null"`
	Goods        []Goods   `gorm:"ForeignKey:MachineId; constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ModelId      uint      `json:"model_id" gorm:"not null"`
	Orders       []Order   `gorm:"ForeignKey:MachineId; constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	DeployedTime time.Time `json:"deployed_time" gorm:"not null"`
}

func (Machine) TableName() string {
	return "machine"
}
func (m *Machine) GetMachineListByStructQuery() (machines []Machine, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetMachineListByStructQuery: %w", err)
		}
	}()
	result := DB.Where(m).Find(&machines)
	err = result.Error
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return
}
func (m *Machine) GetMachineByStructQuery() (machine Machine, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetMachineByStructQuery: %w", err)
		}
	}()
	result := DB.Where(m).First(&machine)
	err = result.Error
	if err != nil {
		log.Println(err)
		return Machine{}, err
	}
	return
}

// 添加售货柜
func (m *Machine) AddMachine() (id uint, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.AddMachine: %w", err)
		}
	}()
	result := DB.Create(m)
	if result.Error != nil {
		log.Println(result.Error)
		return 0, result.Error
	}
	return m.Id, err
}

// 通过售货柜id查找售货柜信息，只预加载商品表
func (m *Machine) GetMachineGoodsById() (machine Machine, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.GetMachineById: %w", err)
		}
	}()
	// 预加载产品
	result := DB.Preload("Goods").First(&machine, m.Id)
	err = result.Error
	if err != nil {
		//log.Println(err)
		return Machine{}, err
	}
	return
}
