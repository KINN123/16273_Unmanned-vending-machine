package main

import (
	"fmt"
	"log"
	"strconv"
	"uvm-backend/config"
	"uvm-backend/database"
	"uvm-backend/model"
	"uvm-backend/router"
	"uvm-backend/service"
	"uvm-backend/utils"

	"github.com/360EntSecGroup-Skylar/excelize"
)

/*
*
手动添加数据：1个商家，1个售货柜，100类商品图像及信息
*/
func initDB() {
	if _, _, _, _, err := service.GetBusinessProductById(1); err == nil {
		//	已有商家，不必再进行初始化
		return
	}
	// 1个商家
	log.Println("初始化商家……")
	businessId, err := service.AddBusiness("无人零售柜", "数据初始化测试")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("商家ID	", businessId)
	//// 2个售货柜
	log.Println("初始化售货柜……")
	machineId, err := service.AddMachine(businessId, "哈工大荔园10栋", 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("1号售货柜ID  ", machineId)
	machineId2, err := service.AddMachine(businessId, "哈工大荔园9栋", 1)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("2号售货柜ID  ", machineId2)
	//// 100类货品，每个货品在每个售货柜中各有100件库存
	// 获取图像文件名list；图像按商品英文名命名
	imgDir := "./upload/img/"
	imgList, err := utils.ReadDir(imgDir)
	if err != nil {
		log.Fatal(err)
	}
	// 从excel文件的Sheet1中读取，列名：englishName, name, info, price
	log.Println("初始化商品信息……")
	file, err := excelize.OpenFile("./bev_name.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	// 读取Sheet1的所有行，每一行有两列：英文名和中文名
	rows := file.GetRows("Sheet1")
	for _, row := range rows {
		englishName := row[0]
		name := row[1]
		info := row[2]
		priceStr := row[3]
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			log.Fatal(err)
		}
		var imageUrl string
		// 从图像列表中找到对应图像路径
		for _, item := range imgList {
			fileFullName := item.Name()
			fileName, _ := utils.GetFileNameAndSuffix(fileFullName)
			if fileName == englishName {
				imageUrl = imgDir + fileFullName
				break
			}
		}
		productId, err := service.AddProduct(businessId, name, englishName, info, price, imageUrl)
		goodsId, err := service.AddGoods(machineId, productId, 100)   // 每个商品初始数量为100
		goodsId2, err := service.AddGoods(machineId2, productId, 100) // 每个商品初始数量为100
		log.Println(goodsId, "	 ", productId, "	", businessId, "	", machineId, "		", name, "	  ", englishName, "  	", info, "	  ", price, "	 ", imageUrl)
		log.Println(goodsId2, "	 ", productId, "	", businessId, "	", machineId2, "		", name, "	  ", englishName, "  	", info, "	  ", price, "	 ", imageUrl)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func main() {
	initDB()
	conf := config.GetConfig()
	if err := router.Router().Run(":" + conf.PORT); err != nil {
		log.Fatal(err)
	}
}

func ope() (err error) {
	db := database.GetDB()
	trans := db.Begin()
	defer func() {
		if err != nil {
			trans.Callback()
			err = fmt.Errorf("ope: %w", err)
		} else {
			trans.Commit()
		}
	}()
	rows, err := trans.Raw("select * from business").Rows()
	if err != nil {
		return err
	}
	var bus model.Business
	for rows.Next() {
		err := db.ScanRows(rows, &bus)
		if err != nil {
			return err
		}
		log.Println(bus.Name)
	}
	//err = trans.Raw("select * from businesses where name = ?", "丑八怪").First(&bus).Error
	//if err != nil {
	//	return err
	//}
	log.Println("first", bus.Name)
	return nil
}
