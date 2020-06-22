package base

import (
	"gin-frame/libraries/config"
	"gin-frame/libraries/mysql"
	"gin-frame/libraries/util"
	"gopkg.in/ini.v1"
	util_err "gin-frame/libraries/util/error"
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
	write := getDSN(conn + "_write")
	read := getDSN(conn + "_read")

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

	cfg := &mysql.Config{
		Master: writeDSN,
		Slave:  readDSN,
	}

	db, err := mysql.New(cfg)
	util.Must(err)

	return db
}

func getMaxOpen(conn string) int {
	cfg := getCfg(conn)
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
