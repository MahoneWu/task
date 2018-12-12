// Package  provides database connection


package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

//define db related param
const(
	userName = "root"
	password = "admin"
	ip = "localhost"
	port = "3306"
	dbName = "demo"
)

//define db connection field
var DB *sql.DB

//initialize the database connection
func initDB()  {
	//splice the database's url param
	path := strings.Join([]string{userName,":",password,"@tcp(",ip,":",port,")/",dbName,"?charset=utf8"}, "")
	//get database connection
	DB, _ = sql.Open("mysql", path)
	DB.SetMaxOpenConns(2000)
	DB.SetMaxIdleConns(1000)
	//test the connecion whether available,if not ,then panic the error and return
	if err := DB.Ping();err != nil{
		panic(err)
		return
	}
	//print connection success
	fmt.Println("db connection success !!!")
}

//init method
func init()  {
	initDB()
}



