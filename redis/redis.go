package redis

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

var RedisClient *redis.Pool

//const param
const (
	HOST = "127.0.0.1:6379"
	TIMEOUT = 10
	IDLETIMEOUT = 60 //idle timeout data
	MAXIDLE = 100 //max idle  amount
	MXAACTIVE = 200 //max active amount
)

//init redis connection
func InitRedis(){
	// create redis pool
	RedisClient = &redis.Pool{
		MaxIdle:     MAXIDLE,
		MaxActive:   MXAACTIVE,
		IdleTimeout: IDLETIMEOUT * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", HOST,
				redis.DialConnectTimeout(TIMEOUT*time.Second),
				redis.DialReadTimeout(TIMEOUT*time.Second),
				redis.DialWriteTimeout(TIMEOUT*time.Second))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}

}

//init method
func init()  {
	InitRedis()

}