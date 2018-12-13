//this package is to deal bussiness

package business

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
	"time"
	"tlsdemo/task/common"
	"tlsdemo/task/db"
	localRedis "tlsdemo/task/redis"
)

func InsertUser(user *User)(int64){
	//start transation
	tx, err := db.DB.Begin()
	var id int64
	if err !=  nil {
		fmt.Println("tx start fail")
		return 0
	}
	//prepare the insert sql
	stmt, err := db.DB.Prepare("insert t_user set name = ? , password = ?,createDate = ?")
	common.Checkerr(err)
	//replac the sql param and execute the sql
	ret, err := stmt.Exec(user.Name, user.Password, time.Now())
	//return lastid
	id ,_= ret.LastInsertId()
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

func SelectUser(name string,password string) string {
	//according to name and password to select meet the criteria
	rows, err := db.DB.Query("select id, name from t_user where name = ? and password = ? ",name,password)
	if err != nil {
		log.Println(err.Error())
	}
	token := common.Token()
	for rows.Next()  {
		var id int
		var name string
		err := rows.Scan(&id,&name)
		common.Checkerr(err)
		//15 minute
		conn := localRedis.RedisClient.Get()
		defer conn.Close()
		_,redisErr := conn.Do("SET",common.RedisKeyStr(LOGIN, token),strconv.Itoa(id),"EX",15*60)
		common.Checkerr(redisErr)
		return token
	}
	return ""
}

func WriteMessage(token string ,key string, value string)(string){
	//get user id from redis by token
	conn := localRedis.RedisClient.Get()
	defer conn.Close()
	userId,_ := redis.Int(conn.Do("GET",common.RedisKeyStr(LOGIN, token)))


	//if userId is zero then you need login first
	if userId == 0{
		fmt.Println("you need login first")
		return "no login"
	}
	//get user write speed
	writeSpeed,_ := redis.Int(conn.Do("GET",common.RedisKey(userId, WRITE)))
	//write limit key，composed of userId and time
	writeLimitKey := common.AddStringWithBuff(WRITE,strconv.Itoa(userId),strconv.FormatInt(time.Now().Unix(), 10),"_")

	fmt.Println("writeSpeed = ",writeSpeed)

	//limitStr := "10_1544509508"
	count,_ := redis.Int(conn.Do("GET", writeLimitKey))
	if(count > writeSpeed){
		log.Printf("超过最大写次数了,writeSpeed = %d,count = %d ",writeSpeed,count)
		return "fail"
	}
	conn.Do("incr",writeLimitKey)

	_,err := conn.Do("SET",common.RedisKey(userId, key), value)
	common.Checkerr(err)
	return "success"
}


func ReadSecureMessage(token string,key string) string{

	conn := localRedis.RedisClient.Get()
	defer conn.Close()

	//get user id from redis by token
	userId,_ := redis.Int(conn.Do("GET",common.RedisKeyStr(LOGIN, token)))
	if userId == 0{
		fmt.Println("you need login first")
		return ""
	}
	//get current user read speed
	readSpeed,_ := redis.Int(conn.Do("GET",common.RedisKey(userId, READ)))
	//current user key，composed of userId and
	readLimitKey := common.AddStringWithBuff(READ,strconv.Itoa(userId),strconv.FormatInt(time.Now().Unix(), 10),"_")


	//get number of times already read
	count,_ := redis.Int(conn.Do("GET", readLimitKey))
	//if count more than the current user read speed ,then return
	if(count > readSpeed){
		fmt.Printf("超过最大读次数了,readSpeed = %d,count = %d ",readSpeed,count)
		fmt.Println()
		return ""
	}

	//increase one
	conn.Do("incr",readLimitKey)
	//get value by key from redis
	value,err := redis.String(conn.Do("GET",common.RedisKey(userId, key)))
	common.Checkerr(err)
	//return value
	return value
}




