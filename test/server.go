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
            continue
        }
        go handleClientRequest(conn)
    }
}



func handleClientRequest(conn net.Conn) {
    //define read channel to read client send data
    readerChannel := make(chan []byte, 10000)
    //define write channel to
    writerChannel := make(chan []byte, 10000)

    //start goroutines
    go readData(readerChannel, conn)
    go writeData(writerChannel, conn)

    //deal request data
    handleData(readerChannel, writerChannel)
    defer conn.Close()
}

//read client  request data
func readData(readerChannel chan []byte, conn net.Conn){
    //temp buffer，storing truncated data
    tmpBuffer := make([]byte, 0)
    //define buffer slice
    buffer := make([]byte, 1024)
    for{
        n, err := conn.Read(buffer)
        if err != nil{
            log.Println(conn.RemoteAddr().String(), "connection error: ", err)
            return
        }
        //temp buffer ,include last time truncated data
        tmpBuffer = append(tmpBuffer, buffer[:n]...)
        for {
            flag := true 
            packData := make([]byte, 0)
            //unpacking
            packData, tmpBuffer, flag = protocol.Depack(tmpBuffer)
            //put the parsed data to readerChannel
            readerChannel <- packData
            if !flag {
                break
            }
        }
         
    }
}

func writeData(writerChannel chan []byte, conn net.Conn) {
    for {
            //get write data from writerChannel
            data := <- writerChannel
            //write data to client
            n, err := conn.Write(protocol.Enpack(data))
            if err != nil{
                log.Println(conn.RemoteAddr().String(), "connection error: ", n, err)
                return
            }
        }
}

//deal channel data
func handleData(readerChannel chan []byte, writerChannel chan []byte) {
    is_authenticate := false
    for{
        select {
        case data := <-readerChannel:

            fmt.Println(string(data))

            if !is_authenticate {
                name,passwd := getLoginParam(string(data))
                if name != "" && passwd != "" {
                    token := business.SelectUser(name,passwd)
                    if(token != ""){
                        fmt.Println("token =" , token)
                            is_authenticate = true
                            //writerChannel <- protocol.IntToBytes(userId)
                            //writerChannel <- []byte(strconv.Itoa(userId))
                            //write token back
                            //writerChannel <- []byte(token)
                            writerChannel <- []byte(token)
                    }
                }
            }
            requestFunctionName, requestKey, requestValue,token := praseAPIData(string(data))
            switch
            {
            case requestFunctionName == "WriteSecureMessage":
                 writeResult := business.WriteMessage(token,requestKey,requestValue)
                 writerChannel <- []byte(writeResult)
            case requestFunctionName == "ReadSecureMessage":
                readResult :=business.ReadSecureMessage(token,requestKey)
                writerChannel <- []byte(readResult)
            }
        //case <-time.After(10*time.Second):
          //  return
        }
    }
}

// parse data
func getLoginParam(data string) (string,string){
    var resultName string
    var resultPassword string
    // parse data to get name data
    strSplit := strings.Split(data, "&")
    if(len(strSplit) == 2){
        name := strings.Split(strSplit[0], "=")
        if len(name) == 2 && name[0] == "name" {
            resultName = name[1]
        }
        //parse data to get password data
        password := strings.Split(strSplit[1], "=")
        if len(name) == 2 && password[0] == "password" {
            resultPassword = password[1]
        }
    }

    return resultName,resultPassword

}

//to parse
func praseAPIData(data string) (string, string, string,string){
    var method string
    var key    string
    var value  string
    var token string
    strSplit := strings.Split(data, "/")
    //fmt.Printf("len[%d] strSplit[0]=%s\n", len(strSplit), strSplit[0])
    if len(strSplit) == 2{
        if strSplit[0] == "WriteSecureMessage" {
            parament := strings.Split(strSplit[1], "&")
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
                        if kv[0] =="token"{
                            token = kv[1]

                        }
                    }
                }
                method = strSplit[0]
            }
        }
    }
    if len(strSplit) == 2 && strSplit[0] == "ReadSecureMessage" {
        kv := strings.Split(strSplit[1], "=")
        if len(kv) == 2 {
            if kv[0] == "key" {
                key = kv[1]
                method = strSplit[0]
            }
            if kv[0] =="token"{
                token = kv[1]

            }
        }
    }
    fmt.Println(method,key,value)
    return method, key, value,token
}