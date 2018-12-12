package tlsclient

import (
    "fmt"
    "net"
    "tlsdemo/task/pool"
    "tlsdemo/task/protocol"
)


//write message interface
func WriteMessage(requestStr string)string{
    //conn,_ := pool.CP.Get();
    conn,_ := pool.CP.Acquire()
    if err := writeClientData(conn,requestStr); err != nil{
        fmt.Println("err to write ", err)
        return ""
    }else{
        response :=readServerData(conn)
        //pool.CP.Put(conn)
        pool.CP.Release(conn)
        return response
    }
}


// login interface
func Login(loginParam string) string {
    //load client certificate
    //conn,_ := pool.CP.Get();
    conn,_ := pool.CP.Acquire();
    writeClientData(conn,loginParam)

    token := readServerData(conn)
    //pool.CP.Put(conn)
    pool.CP.Release(conn)
    return token

}


//read client response data
func readServerData(conn net.Conn) string{
    //temp buffer，storing truncated data
    tmpBuffer := make([]byte, 0)
    buffer := make([]byte, 1024)
    for{
        //conn.SetReadDeadline(time.Now().Add(time.Second))
        //fmt.Println("before Read")
        n, err := conn.Read(buffer)
        fmt.Println("ater Read , n =",n)
        if err != nil{
            fmt.Println(conn.RemoteAddr().String(), "readClientData connection error: ", err)
            return ""
        }
        //get read byte
        tmpBuffer = append(tmpBuffer, buffer[:n]...)
        packData := make([]byte, 0)
        packData, tmpBuffer, _ = protocol.Depack(tmpBuffer)
        //fmt.Println("patch", len(packData))
        //put the parsed data to channel
        if(len(packData) != 0 ){
            return string(packData)
        }else{
            return "fail"
        }
    }
}

//write data
func writeClientData(conn net.Conn,data string) error {
        //write data to server
        n, err := conn.Write(protocol.Enpack([]byte(data)))
       // fmt.Println(fmt.Sprintf("%s, len %d", data, n))
        if err != nil{
            fmt.Println(conn.RemoteAddr().String(), "writeClientData connection error: ", n,err)
            return err
        }
        return nil
}

