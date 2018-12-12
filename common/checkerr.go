//this package is to check common error log print


package common

import "log"

func Checkerr(err error)  {
	if err != nil{
		log.Println(err.Error())
	}
}
