package base

import (
	"gin-frame/libraries/config"
	"gin-frame/libraries/util"
	"gin-frame/libraries/mysql"
)

var priceInstance *mysql.DB
var ymtInstance *mysql.DB

func GetConn(conn string) *mysql.DB{
	db := &mysql.DB{}

	if conn == "hangqing" {
		if priceInstance == nil {
			db = getPriceConn()
		}else {
			db = priceInstance
		}
	}else if conn == "ymt360" {
		if ymtInstance == nil {
			db = getYmtConn()
		}else {
			db = ymtInstance
		}
	}else{
		panic("err conn string")
	}

	return db
}

func getYmtConn() *mysql.DB{
	write := GetDSN("ymt360_write")
	read := GetDSN("ymt360_read")

	writeDSN := mysql.Conn{
		DSN:     write,
		MaxOpen: 5,
		MaxIdle: 5,
	}
	readDSN := mysql.Conn{
		DSN:     read,
		MaxOpen: 5,
		MaxIdle: 5,
	}

	arrDSN := []mysql.Conn{}
	arrDSN = append(arrDSN, readDSN)

	cfg := &mysql.Config{
		Master:      writeDSN,
		Slave:       arrDSN,
	}

	db, err := mysql.New(cfg)
	util.Must(err)

	return db
}

func getPriceConn() *mysql.DB{
	write := GetDSN("hangqing_write")
	read := GetDSN("hangqing_read")

	writeDSN := mysql.Conn{
		DSN:     write,
		MaxOpen: 5,
		MaxIdle: 5,
	}
	readDSN := mysql.Conn{
		DSN:     read,
		MaxOpen: 5,
		MaxIdle: 5,
	}

	arrDSN := []mysql.Conn{}
	arrDSN = append(arrDSN, readDSN)

	cfg := &mysql.Config{
		Master:      writeDSN,
		Slave:       arrDSN,
	}

	db, err := mysql.New(cfg)
	util.Must(err)

	return db
}

func GetDSN(conn string) string {
	cfg := config.GetConfig("mysql", conn)
	dsn := cfg.Key("user").String() + ":" + cfg.Key("password").String() + "@tcp(" + cfg.Key("host").String() + ":" + cfg.Key("port").String() + ")/" + cfg.Key("db").String() + "?charset=" + cfg.Key("charset").String()
	return dsn
}