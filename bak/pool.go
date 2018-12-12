package bak

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"
	"tlsdemo/task/common"
)


var (
	ErrInvalidConfig = errors.New("invalid pool config")
	ErrPoolClosed    = errors.New("pool closed")
)


type factory func() (net.Conn, error)


type Pool interface {
	Acquire() (net.Conn, error) // 获取资源
	Release(net.Conn) error     // 释放资源
	Close(net.Conn) error       // 关闭资源
	Shutdown() error            // 关闭池
}

type GenericPool struct {
	sync.Mutex
	pool        chan net.Conn
	maxOpen     int  // 池中最大资源数
	numOpen     int  // 当前池中资源数
	minOpen     int  // 池中最少资源数
	closed      bool // 池是否已关闭
	maxLifetime time.Duration
	factory     factory // 创建连接的方法
}


func NewGenericPool(minOpen, maxOpen int, maxLifetime time.Duration, factory factory) (*GenericPool, error) {
	if maxOpen <= 0 || minOpen > maxOpen {
		return nil, ErrInvalidConfig
	}
	p := &GenericPool{
		maxOpen:     maxOpen,
		minOpen:     minOpen,
		maxLifetime: maxLifetime,
		factory:     factory,
		pool:        make(chan net.Conn, maxOpen),
	}

	for i := 0; i < minOpen; i++ {
		netCoon, err := factory()
		if err != nil {
			continue
		}
		p.numOpen++
		p.pool <- netCoon
	}
	return p, nil
}


func (p *GenericPool) Acquire() (net.Conn, error) {
	if p.closed {
		return nil, ErrPoolClosed
	}
	for {
		netCoon, err := p.getOrCreate()
		if err != nil {
			return nil, err
		}
		return netCoon, nil
	}
}

func (p *GenericPool) getOrCreate() (net.Conn, error) {
	select {
	case netCoon := <-p.pool:
		return netCoon, nil
	default:
	}
	p.Lock()
	if p.numOpen >= p.maxOpen {
		netCoon := <-p.pool
		p.Unlock()
		return netCoon, nil
	}
	// 新建连接
	netCoon, err := p.factory()
	if err != nil {
		p.Unlock()
		return nil, err
	}
	p.numOpen++
	p.Unlock()
	return netCoon, nil
}

// 释放单个资源到连接池
func (p *GenericPool) Release(closer net.Conn) error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	p.pool <- closer
	p.Unlock()
	return nil
}


// 关闭单个资源
func (p *GenericPool) Close(closer net.Conn) error {
	p.Lock()
	closer.Close()
	p.numOpen--
	p.Unlock()
	return nil
}


// 关闭连接池，释放所有资源
func (p *GenericPool) Shutdown() error {
	if p.closed {
		return ErrPoolClosed
	}
	p.Lock()
	close(p.pool)
	for netCoon := range p.pool {
		netCoon.Close()
		p.numOpen--
	}
	p.closed = true
	p.Unlock()
	return nil
}

var NewCP *GenericPool

//(minOpen, maxOpen int, maxLifetime time.Duration, factory factory) (*GenericPool, error) {
func init(){
	NewCP, _ = NewGenericPool(100,500,time.Second*60*5, func() (net.Conn, error) {
		return  CreateConnections()
	});
}


func CreateConnections()(net.Conn,error){
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
}





