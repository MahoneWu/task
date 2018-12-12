//this package is for file packaging

package common

import (
	"bytes"
	"os"
	"strconv"
	"strings"
)

//package the file path
func FlatPath(sufferPath string) string{
	dir, _ := os.Getwd()
	path := strings.Join([]string{dir,sufferPath},"")
	return path
}

//flat int and string
func RedisKey(userId int,key string) string{
	redisKey := strings.Join([]string{strconv.Itoa(userId) ,"_",key},"")
	return redisKey
}

//flat two string
func RedisKeyStr(prefix string,suffix string) string{
	redisKey := strings.Join([]string{prefix ,"_",suffix},"")
	return redisKey
}


//flat two string
func JionStr(prefix string,suffix string,joinType string) string{
	redisKey := strings.Join([]string{prefix ,joinType,suffix},"")
	return redisKey
}

func AddStringWithBuff(str1 string,str2 string,str3 string,symbol string)string{
	var buffer bytes.Buffer
	buffer.WriteString(str1)
	buffer.WriteString(symbol)
	buffer.WriteString(str2)
	buffer.WriteString(symbol)
	buffer.WriteString(str3)
	return buffer.String()
}

