package utils

import (
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

/**
获得TimeFormat
*/
func GetTimeFormat(haveDay bool) string {
	if haveDay {
		return "2006-01-02"
	} else {
		return "2006-01"
	}
}

/**
获得redis月排行榜zSet的key值：machinId-year-month
*/
func GetRankingZSetKey(machineId uint, date string) string {
	return strconv.Itoa(int(machineId)) + "-" + date
}

func GetUUID() string {
	return uuid.NewV4().String()
}

/**
读取指定文件夹下按名称排序的文件名
*/
func ReadDir(dirPath string) (list []os.FileInfo, err error) {
	list, err = ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Name() < list[j].Name()
	})
	//for _, item := range list {
	//	log.Println(item.Name())
	//}
	return
}

/**
输入文件名带后缀，获取文件名和后缀
*/
func GetFileNameAndSuffix(fullName string) (fileName, suffix string) {
	suffix = path.Ext(fullName)
	fileName = strings.TrimSuffix(fullName, suffix)
	return
}
