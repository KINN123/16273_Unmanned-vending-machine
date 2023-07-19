package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"uvm-backend/config"
)

var once sync.Once
var db *gorm.DB

func connectDatabase() {
	var err error
	conf := config.GetConfig()
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", conf.DBUsername, conf.DBPassword, conf.DBAddress, conf.DBName)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	once.Do(connectDatabase)
	return db
}
