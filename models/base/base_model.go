package base

import (
	"fmt"
	"github.com/why444216978/go-library/libraries/config"
	"github.com/why444216978/go-library/libraries/mysql"
	"github.com/why444216978/go-library/libraries/util"
	util_err "github.com/why444216978/go-library/libraries/util/error"
	"gopkg.in/ini.v1"
)

var cfgs map[string]*ini.Section
var dbInstance map[string]*mysql.DB

func GetInstance(conn string) *mysql.DB {
	if len(dbInstance) == 0 {
		dbInstance = make(map[string]*mysql.DB, 30)
	}

	if dbInstance[conn] == nil {
		dbInstance[conn] = getConn(conn)
	}

	return dbInstance[conn]
}

func getConn(conn string) *mysql.DB {
	write := conn + "_write"
	read := conn + "_read"
	writeDsn := getDSN(conn + "_write")
	readDsn := getDSN(conn + "_read")

	writeObj := mysql.Conn{
		DSN:     writeDsn,
		MaxOpen: getMaxOpen(write),
		MaxIdle: getMaxIdle(write),
	}

	readObj := mysql.Conn{
		DSN:     readDsn,
		MaxOpen: getMaxOpen(read),
		MaxIdle: getMaxOpen(read),
	}

	cfg := &mysql.Config{
		Master: writeObj,
		Slave:  readObj,
	}

	db, err := mysql.New(cfg)
	util.Must(err)

	return db
}

func getMaxOpen(conn string) int {
	cfg := getCfg(conn)
	fmt.Println(cfg)
	masterNum, err := cfg.Key("max_open").Int()
	util_err.Must(err)
	return masterNum
}

func getMaxIdle(conn string) int {
	cfg := getCfg(conn)
	masterNum, err := cfg.Key("max_idle").Int()
	util_err.Must(err)
	return masterNum
}

func getDSN(conn string) string {
	cfg := getCfg(conn)
	dsn := cfg.Key("user").String() + ":" + cfg.Key("password").String() + "@tcp(" + cfg.Key("host").String() + ":" + cfg.Key("port").String() + ")/" + cfg.Key("db").String() + "?charset=" + cfg.Key("charset").String()
	return dsn
}

func getCfg(conn string) *ini.Section {
	if cfgs == nil {
		cfgs = make(map[string]*ini.Section, 30)
	}
	if cfgs[conn] == nil {
		cfgs[conn] = config.GetConfig("mysql", conn)
	}
	return cfgs[conn]
}
