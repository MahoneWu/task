//this package is to deal bussiness

package business

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
	"tlsdemo/task/common"
	"tlsdemo/task/db"
	localRedis "tlsdemo/task/redis"
)

var loginUserSession = common.NewSafeMap()

var readRateLimitCtx *common.RateLimitCtx
var writeRateLimitCtx *common.RateLimitCtx

func InsertUser(user *User) int64 {
	//start transation
	tx, err := db.DB.Begin()
	var id int64
	if err != nil {
		fmt.Println("tx start fail")
		return 0
	}
	//prepare the insert sql
	stmt, err := db.DB.Prepare("insert t_user set name = ? , password = ?,createDate = ?")
	common.Checkerr(err)
	//replac the sql param and execute the sql
	ret, err := stmt.Exec(user.Name, user.Password, time.Now())
	//return lastid
	id, _ = ret.LastInsertId()
	//judge whether had error
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	//commit  the transation
	tx.Commit()
	//return the result
	return id
}

func SelectUser(name string, password string) string {
	//according to name and password to select meet the criteria
	rows, err := db.DB.Query("select id, name from t_user where name = ? and password = ? ", name, password)
	if err != nil {
		log.Println(err.Error())
	}
	token := common.Token()
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		common.Checkerr(err)
		//15 minute
		conn := localRedis.RedisClient.Get()
		defer conn.Close()
		//_,redisErr := conn.Do("SET",common.RedisKeyStr(LOGIN, token),strconv.Itoa(id),"EX",15*60)
		//common.Checkerr(redisErr)
		loginUserSession.WriteMap(token, id)
		return token
	}
	return ""
}

func WriteMessage(token string, key string, value string) string {
	//get user id from redis by token
	conn := localRedis.RedisClient.Get()
	defer conn.Close()
	//userId,_ := redis.Int(conn.Do("GET",common.RedisKeyStr(LOGIN, token)))
	userId := loginUserSession.ReadMap(token).(int)

	//if userId is zero then you need login first
	if userId == 0 {
		fmt.Println("you need login first")
		return "no login"
	}

	succ := writeRateLimitCtx.Acquire(userId)
	if succ {
		return "fail"
	}

	_, err := conn.Do("SET", common.RedisKey(userId, key), value)
	common.Checkerr(err)
	return "success"
}

func ReadSecureMessage(token string, key string) string {

	conn := localRedis.RedisClient.Get()
	defer conn.Close()

	//get user id from redis by token
	userId, _ := redis.Int(conn.Do("GET", common.RedisKeyStr(LOGIN, token)))
	if userId == 0 {
		fmt.Println("you need login first")
		return ""
	}
	succ := readRateLimitCtx.Acquire(userId)
	if succ {
		return ""
	}
	//get value by key from redis
	value, err := redis.String(conn.Do("GET", common.RedisKey(userId, key)))
	common.Checkerr(err)
	//return value
	return value
}

func init() {
	readRateLimitCtx = common.NewRateLimitCtx(ReadQuotaConfig)
	writeRateLimitCtx = common.NewRateLimitCtx(WriteQuotaConfig)
}
