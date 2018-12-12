package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"tlsdemo/task/business"
	"tlsdemo/task/common"
	"tlsdemo/task/protocol"
)

func main() {
	cert, err := tls.LoadX509KeyPair(common.FlatPath("/task/certs/server.pem"), common.FlatPath("/task/certs/server.key"))
	if err != nil {
		log.Println(err)
		return
	}
	//load tls client pem
	certBytes, err := ioutil.ReadFile(common.FlatPath("/task/certs/client.pem"))
	if err != nil {
		panic("Unable to read cert.pem")
	}
	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		panic("failed to parse root certificate")
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}
	listener, err := tls.Listen("tcp", ":5000", config)
	//listener, err := net.Listen("tcp", "127.0.0.1:5000")
	if err != nil {
		fmt.Printf("Fatal error: %s", err.Error())
	}
	defer listener.Close()
	for {
		conn, err1 := listener.Accept()
		if err1 != nil {
			fmt.Println("accept error")
			continue
		}
		//fmt.Printf("new connection...%+v\n", conn)
		go handleClientRequest(conn)
	}
}

func handleClientRequest(conn net.Conn) {
	//deal request data
	for {
		if err := handleData(conn); err != nil {
			break
		}
	}
}

//read client  request data
func readData(conn net.Conn) ([]byte, error) {
	//temp buffer，storing truncated data
	tmpBuffer := make([]byte, 0)
	//define buffer slice
	buffer := make([]byte, 5000)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			//fmt.Println(conn.RemoteAddr().String(), "readData connection error = : ", err)
			return make([]byte, 0), err
		}
		//temp buffer ,include last time truncated data
		tmpBuffer = append(tmpBuffer, buffer[:n]...)
		var packData []byte
		for {
			flag := true
			packData = make([]byte, 0)
			//unpacking
			packData, tmpBuffer, flag = protocol.Depack(tmpBuffer)
			//put the parsed data to readerChannel
			if !flag {
				break
			}
		}
		fmt.Println(string(packData))

	}
	return nil, nil

}

func writeData(result string, conn net.Conn) {
	//write data to client
	//fmt.Println("RESULET",result)
	n, err := conn.Write(protocol.Enpack([]byte(result)))
	if err != nil {
		fmt.Println(conn.RemoteAddr().String(), "writeData connection error: ", n, err)
		return
	}
}

//deal channel data
func handleData(conn net.Conn) error {
	data, err := readData(conn)

	if err != nil {
		return err
	}
	method, params := getMethod(string(data))
	fmt.Println("request--", method, params)
	switch method {
	case "Login":
		name, passwd := getLoginParam(params)
		if name != "" && passwd != "" {
			token := business.SelectUser(name, passwd)
			if token != "" {
				fmt.Println("token =", token)
				writeData(token, conn)
			}
		}
	case "WriteSecureMessage":
		requestKey, requestValue, token := praseAPIData(method, params)

		//fmt.Println(fmt.Sprintf("before writing %v, %v", requestKey, requestValue))
		writeResult := business.WriteMessage(token, requestKey, requestValue)
		//fmt.Println(fmt.Sprintf("after writing %v, %v", requestKey, requestValue))
		writeData(writeResult, conn)
	case "ReadSecureMessage":
		requestKey, _, token := praseAPIData(method, params)
		readResult := business.ReadSecureMessage(token, requestKey)
		writeData(readResult, conn)
	default:
		fmt.Println("unrecognized method")
	}
	return nil
}

func getMethod(data string) (string, string) {
	strSplit := strings.Split(data, "/")
	if len(strSplit) == 2 {
		return strSplit[0], strSplit[1]
	}
	return "", ""
}

// parse data
func getLoginParam(data string) (string, string) {
	var resultName string
	var resultPassword string
	parament := strings.Split(data, "&")
	// parse data to get name data
	name := strings.Split(parament[0], "=")
	if len(name) == 2 && name[0] == "name" {
		resultName = name[1]
	}
	//parse data to get password data
	password := strings.Split(parament[1], "=")
	if len(name) == 2 && password[0] == "password" {
		resultPassword = password[1]
	}
	return resultName, resultPassword
}

//to parse
func praseAPIData(method string, data string) (string, string, string) {
	var key string
	var value string
	var token string
	if method == "WriteSecureMessage" {
		parament := strings.Split(data, "&")
		if len(parament) == 3 {
			for i := 0; i < 3; i++ {
				kv := strings.Split(parament[i], "=")
				if len(kv) == 2 {
					if kv[0] == "key" {
						key = kv[1]
					}
					if kv[0] == "value" {
						value = kv[1]
					}
					if kv[0] == "token" {
						token = kv[1]

					}
				}
			}
		}
	}

	if method == "ReadSecureMessage" {
		readParam := strings.Split(data, "&")
		if len(readParam) == 2 {
			for j := 0; j < 2; j++ {
				kv := strings.Split(readParam[j], "=")
				if len(kv) == 2 {
					if kv[0] == "key" {
						key = kv[1]
					}
					if kv[0] == "token" {
						token = kv[1]
					}
				}
			}
		}
	}
	fmt.Println(token, key, value)
	return key, value, token
}
