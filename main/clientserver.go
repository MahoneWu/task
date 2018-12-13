package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"tlsdemo/task/common"
	_ "tlsdemo/task/db"
	_ "tlsdemo/task/pool"
	_ "tlsdemo/task/redis"
	"tlsdemo/task/tlsclient"
)

func main() {

	lparma := loginParam(common.JionStr("test",strconv.Itoa(5000),""), "123")
	token := tlsclient.Login(lparma)
	fmt.Println("---",token)


	//token := "793df21ec43572357bdd3a2dc9711078"

	t1 := time.Now()
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go toWriteData(token, i, &wg)
	}
	wg.Wait()
	t2 := time.Now()
	fmt.Println("-----time diff", t2.Sub(t1))
}

func toWriteData(token string, group int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 100; i++ {
		key := common.RandSeq(1000)
		value := common.RandSeq(1000)
		//fmt.Println(fmt.Sprintf("%d:%d", group, i))
		//parma := writeParam(fmt.Sprintf("%d:%d", group, i),fmt.Sprintf("%d:%d", group, i),token)
		parma := WriteParam(key, value, token)

		//fmt.Println(parma)
		//fmt.Println("--Len", len([]byte(parma)))

		writeResponse := tlsclient.WriteMessage(parma)
		fmt.Println("---Response", writeResponse)

		readParam := readParam(key,token)
		readResponse := tlsclient.WriteMessage(readParam)
		fmt.Println("---read",readResponse,"---",value)

	}
}

func WriteParam(paramKey string, paramValue string, paramToken string) string {
	param := "WriteSecureMessage/"
	key := "key="
	key += paramKey

	value := "value="
	value += paramValue

	token := "token="
	token += paramToken

	body1 := common.JionStr(key, value, "&")
	body2 := common.JionStr(body1, token, "&")
	requestParam := common.JionStr(param, body2, "")

	//fmt.Println("byte=",len([]byte(requestParam)))

	return requestParam
}

func readParam(paramKey string, paramToken string) string {
	param := "ReadSecureMessage/"
	key := "key="
	key += paramKey

	token := "token="
	token += paramToken

	body := common.JionStr(key, token, "&")
	requestParam := common.JionStr(param, body, "")

	return requestParam
}

func loginParam(nameParam string, passwordParam string) string {
	name := "Login/name="
	name += nameParam

	password := "password="
	password += passwordParam

	transportData := common.JionStr(name, password, "&")
	return transportData
}
