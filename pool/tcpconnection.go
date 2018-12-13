package pool

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"
	"tlsdemo/task/common"
)

var count int32
type Factory func() (net.Conn, error)

//define conn
type Conn struct {
	conn net.Conn;
	//define conn time
	time time.Time;
}

//conn pool
type ConnPool struct {
	mu sync.Mutex;
	conns chan *Conn;
	//factory to create connection resource
	factory Factory;
	//determine if the pool is closed
	closed bool;
	//connection timeout
	connTimeOut time.Duration;
}

func NewConnPool(factory Factory, cap int, connTimeOut time.Duration)  (*ConnPool, error){

	if cap <= 0 {
		return nil, errors.New("cap can not less than zero");
	}
	if connTimeOut <= 0 {
		return nil, errors.New("connTimeOut  not less than zero");
	}


	cp := &ConnPool{
		mu:          sync.Mutex{},
		conns:       make(chan *Conn, cap),
		factory:     factory,
		closed:      false,
		connTimeOut: connTimeOut,
	};

	for i := 0; i < cap; i++ {
		//create connection resources through the factory
		connRes, err := cp.factory();
		if err != nil {
			cp.Close();
			return nil, errors.New("factory error");
		}
		//insert the connection resource into the channel
		cp.conns <- &Conn{conn: connRes, time: time.Now()};
	}
	return cp, nil;
}

//get connection resources
func (cp *ConnPool) Acquire() (net.Conn, error) {
	if cp.closed {
		return nil, errors.New("connection pool is closed");
	}
	cp.mu.Lock();
	defer cp.mu.Unlock();
	fmt.Println("---get  ",len(cp.conns))
	for {
		timer := time.NewTimer(time.Second)
		select {
		case  <- timer.C:
			fmt.Println("timeout")
		//get connection resources from the channel
		case connRes, ok := <-cp.conns:
			{
				if !ok {
					return nil, errors.New("connection pool is closed");
				}
				//determine the time in the connection ,if it time out,then close
				if time.Now().Sub(connRes.time) > cp.connTimeOut {
					fmt.Println("get  close")
					connRes.conn.Close();
					continue;
				}
				/**
				atomic.AddInt32(&count, 1)
				fmt.Println("connection count ", count)
				*/
				return connRes.conn, nil;
			}
		}
	}
}

//connection resources back to the pool
func (cp *ConnPool) Release(conn net.Conn) error {
	if cp.closed {
		fmt.Println("connection pool is closed")
		return errors.New("connection pool is closed");
	}

	select {
	//add connection resources to the channel
	case cp.conns <- &Conn{conn: conn, time: time.Now()}:
		{
			fmt.Println("put----",len(cp.conns))
			return nil;
		}
	}
}


//关闭连接池
func (cp *ConnPool) Close() {
	if cp.closed {
		return;
	}
	cp.mu.Lock();
	defer cp.mu.Unlock();
	cp.closed = true;

	fmt.Println("关闭连接了")

	//关闭通道
	close(cp.conns);
	//循环关闭通道中的连接
	for conn := range cp.conns {
		conn.conn.Close();
	}

}






func CreateConnection()(net.Conn,error){
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
	return conn,nil;

	//at last to close connection
	//defer conn.Close()
}

func (cp *ConnPool) len() int {
	return len(cp.conns);
}




var CP *ConnPool

func init(){
	CP, _ = NewConnPool(func() (net.Conn, error) {
		return CreateConnection()
	}, 200, time.Second*60*50);
}


