package model

import (
	"fmt"
	"log"
)

/**
从redis zSet中降序获得[start...end]的value和score
*/
func ZsetRevRange(zSet string, start int64, stop int64) (values []string, scores []int) {
	values = RedisDB.ZRevRange(zSet, start, stop).Val() // 值
	for _, value := range values {
		// 对应分数
		score := RedisDB.ZScore(zSet, value).Val()
		scores = append(scores, int(score))
	}
	log.Println("zSet values: ", values)
	log.Println("zSet score: ", scores)
	return
}

/**
按value查找并更新redis zSet的score
*/
func ZsetAdd(zSet string, value string, newSale int) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("model.ZsetAdd: %w", err)
		}
	}()
	var score float64
	// 判断是否存在该成员
	isExist := RedisDB.ZLexCount(zSet, "["+value, "["+value).Val()
	if isExist != 0 {
		// 存在，取出原销量
		score, err = RedisDB.ZScore(zSet, value).Result()
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		score = 0
	}
	newScore := score + float64(newSale)
	RedisDB.ZAdd(zSet, struct {
		Score  float64
		Member interface{}
	}{Score: newScore, Member: value})
	log.Println("zSet value: ", value, " score: ", score, " newScore: ", newScore)
	return nil
}
