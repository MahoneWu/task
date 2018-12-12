package bak

import (
    "log"
    "net"
    "tlsdemo/task/pool"
    "tlsdemo/task/protocol"
)


//write message interface
func WriteMessage(requestStr string)string{
    conn,_ := pool.NewCP.Acquire();
    writeClientData(conn,requestStr)

    response :=readClientData(conn)
    pool.NewCP.Release(conn)
    return response
}



// login interface
func Login(loginParam string) string {
    //load client certificate
    conn,_ := pool.NewCP.Acquire();
    writeClientData(conn,loginParam)

    token := readClientData(conn)
    pool.NewCP.Release(conn)
    return token

}


//read client response data
func readClientData(conn net.Conn) string{
    //temp buffer，storing truncated data
    tmpBuffer := make([]byte, 0)
    buffer := make([]byte, 1024)
    for{
        n, err := conn.Read(buffer)
        if err != nil{
            log.Println(conn.RemoteAddr().String(), "readClientData connection error: ", err)
            return ""
        }
        //get read byte
        tmpBuffer = append(tmpBuffer, buffer[:n]...)
        packData := make([]byte, 0)
        packData, tmpBuffer, _ = protocol.Depack(tmpBuffer)
        //put the parsed data to channel
        if(len(packData) != 0 ){
            return string(packData)
        }
    }

}

//write data
func writeClientData(conn net.Conn,data string) {
        //write data to server
        n, err := conn.Write(protocol.Enpack([]byte(data)))
        if err != nil{
            log.Println(conn.RemoteAddr().String(), "writeClientData connection error: ", n,err)
            return
        }
}

