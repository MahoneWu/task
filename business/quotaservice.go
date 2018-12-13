//package is to quota business

package business

import (
	"fmt"
	"log"
	"time"
	"tlsdemo/task/common"
	"tlsdemo/task/db"
	"tlsdemo/task/redis"
)

var WriteQuotaConfig = common.NewSafeMap()
var ReadQuotaConfig = common.NewSafeMap()

func SelectQuotaInfo(userId int) (quota UserQuota) {
	var writeSpeed int
	var readSpeed int
	//according to userId to select meet the criteria
	err := db.DB.QueryRow("select writeSpeed,readSpeed from t_user_quota where userId = ?  ", userId).Scan(&writeSpeed, &readSpeed)
	common.Checkerr(err)
	result := UserQuota{
		WriteSpeed: writeSpeed,
		ReadSpeed:  readSpeed,
	}
	return result
}

func SelectQuotaList() {
	var writeSpeed int
	var readSpeed int
	var userId int

	//according to userId to select meet the criteria
	rows, err := db.DB.Query("select writeSpeed,readSpeed,userId from t_user_quota")
	common.Checkerr(err)
	conn := redis.RedisClient.Get()
	defer conn.Close()
	for rows.Next() {
		err = rows.Scan(&writeSpeed, &readSpeed, &userId)
		// _,writeErr := conn.Do("SET",common.RedisKey(userId, WRITE), writeSpeed,"EX",10*60)
		// common.Checkerr(writeErr)
		// _,readErr := conn.Do("SET",common.RedisKey(userId, READ), readSpeed,"EX",10*60)
		// common.Checkerr(readErr)
		WriteQuotaConfig.WriteMap(userId, writeSpeed)
		ReadQuotaConfig.WriteMap(userId, readSpeed)
	}
}

func InsertQuota(quota UserQuota) bool {
	//define return field
	var flag = false
	//start transation
	tx, err := db.DB.Begin()
	if err != nil {
		fmt.Println("tx start fail")
		return flag
	}
	//prepare the insert sql
	stmt, err := db.DB.Prepare("insert t_user_quota set userId = ? , writeSpeed = ?,readSpeed = ? ,createDate = ?")
	common.Checkerr(err)
	//replac the sql param and execute the sql
	_, err = stmt.Exec(quota.UserId, quota.WriteSpeed, quota.ReadSpeed, time.Now())
	//judge whether had error
	if err != nil {
		log.Println(err.Error())
		return false
	}
	//commit  the transation
	tx.Commit()
	flag = true
	//return the result
	return flag
}

func init() {
	SelectQuotaList()
}
