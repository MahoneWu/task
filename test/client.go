package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io/ioutil"
    "log"
    "net"
    "strconv"
    "tlsdemo/task/common"
    "tlsdemo/task/protocol"
)

func main() {
    for i := 0; i < 100; i++ {
        //go creatClientConn(writeParam(strconv.Itoa(i),strconv.Itoa(i)))
        parma := login(common.JionStr("test",strconv.Itoa(10),""), "123")
        go creatClientConn(parma)
    }

    for {
        var str string
        fmt.Println("please input  exit to quit：\n")
        fmt.Scanln(&str)
        if str == "exit" {
            break
        }
    }
}

func writeParam(paramKey string,paramValue string) string{
    param := "WriteSecureMessage/"
    key := "key="
    key += paramKey

    value := "value="
    value += paramValue

    body := common.JionStr(key, value,"&")

    requestParam := common.JionStr(param, body,"")
    return requestParam
}

func login(nameParam string,passwordParam string) string{

    name := "Login/name="
    name += nameParam

    password := "password="
    password += passwordParam



    transportData := common.JionStr(name, password,"&")
    return transportData
}

func creatClientConn(reuestStr string)  {
    //load client certificate
    cert, err := tls.LoadX509KeyPair(common.FlatPath("/task/certs/client.pem"), common.FlatPath("/task/certs/client.key"))
    if err != nil {
        log.Println(err)
    }
    certBytes, err := ioutil.ReadFile(common.FlatPath("/task/certs/client.pem"))
    if err != nil {
        panic("Unable to read cert.pem")
    }
    clientCertPool := x509.NewCertPool()
    ok := clientCertPool.AppendCertsFromPEM(certBytes)
    if !ok {
        panic("failed to parse root certificate")
    }
    //InsecureSkipVerify field is to control client whether verification cert and host name
    conf := &tls.Config{
        RootCAs:            clientCertPool,
        Certificates:       []tls.Certificate{cert},
        InsecureSkipVerify: true,
    }
    conn, err := tls.Dial("tcp", "127.0.0.1:5000", conf)
    if err != nil {
        log.Println(err)
    }
    //at last to close connection
    defer conn.Close()



    readerChannel := make(chan []byte, 10000)
    writerChannel := make(chan []byte, 10000)

     readClientData(readerChannel, conn)
     writeClientData(writerChannel, conn)


    writerChannel <- []byte(reuestStr)

    //登录
    data := <- readerChannel
    token := string(data)
    fmt.Println("token= " , token)


    writeSecureMessage := "WriteSecureMessage/key=test&value=shopee&token="
    writeSecureMessage += token
    readSecureMessage := "ReadSecureMessage/key=test&token="
    readSecureMessage += token

    //fmt.Println(writeSecureMessage)
    //fmt.Println(readSecureMessage)

    writerChannel <- []byte(writeSecureMessage)
    writerChannel <- []byte(readSecureMessage)

    //var mutex sync.Mutex
    for i := 0; i < 10000; i++{
        writerChannel <- []byte(writeSecureMessage)


    }


    //time.Sleep(time.Second*5)
}


//read client response data
func readClientData(readerChannel chan []byte, conn net.Conn){
    //temp buffer，storing truncated data
    tmpBuffer := make([]byte, 0)
    buffer := make([]byte, 1024)
    for{
        n, err := conn.Read(buffer)
        if err != nil{
            log.Println(conn.RemoteAddr().String(), "connection error: ", err)
            return
        }
        //get read byte
        tmpBuffer = append(tmpBuffer, buffer[:n]...)
        for {
            flag := true
            packData := make([]byte, 0)
            packData, tmpBuffer, flag = protocol.Depack(tmpBuffer)
            //put the parsed data to channel
            readerChannel <- packData
            if !flag {
                break
            }
        }

    }
}

func writeClientData(writerChannel chan []byte, conn net.Conn) {
    for {
        data := <- writerChannel
        //write data to server
        n, err := conn.Write(protocol.Enpack(data))
        if err != nil{
            log.Println(conn.RemoteAddr().String(), "connection error: ", n,err)
            return
        }
    }
}

