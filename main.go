package main

import (
	"fmt"
	"tlsdemo/task/common"
	"tlsdemo/task/protocol"
)

func main() {


	token := "793df21ec43572357bdd3a2dc9711078"

	key := common.RandSeq(1000)
	value := common.RandSeq(1000)


	param := WriteParam(key,value,token)

	intParm := protocol.ByteToInt([]byte(param))

	intstr := protocol.IntToBytes(intParm)

	fmt.Println(len(intstr))



}

func WriteParam(paramKey string,paramValue string,paramToken string) string{
	param := "WriteSecureMessage/"
	key := "key="
	key += paramKey

	value := "value="
	value += paramValue

	token := "token="
	token += paramToken

	body1 := common.JionStr(key, value,"&")
	body2 := common.JionStr(body1, token,"&")
	requestParam := common.JionStr(param, body2,"")

	//fmt.Println("byte=",len([]byte(requestParam)))

	return requestParam
}