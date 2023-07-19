package model

import "time"

type Model struct {
	Id                uint      `json:"id" gorm:"primaryKey"`
	Name              string    `json:"name" gorm:"size:100; not null"`
	Url               string    `json:"url" gorm:"size:150; not null"`
	ProductCategories string    `json:"product_categories" gorm:"size:2000 not null"`
	UpdateTime        time.Time `json:"update_time"`
}

func (Model) TableName() string {
	return "model"
}
