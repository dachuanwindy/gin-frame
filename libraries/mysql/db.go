package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"gin-frame/libraries/config"
	"gin-frame/libraries/log"
	"gin-frame/libraries/util"
	"gin-frame/libraries/xhop"
	"math/rand"
	"time"

	"github.com/opentracing/opentracing-go"

	_ "github.com/go-sql-driver/mysql"
)

type (
	DB struct {
		IsLog	bool
		masterDB *sql.DB
		slaveDB  []*sql.DB
		Config   *Config
	}
)

func New(c *Config) (db *DB, err error) {
	db = new(DB)
	db.Config = c
	db.IsLog = GetIsLog()
	db.masterDB, err = sql.Open("mysql", c.Master.DSN)
	if err != nil {
		err = errorsWrap(err, "init master db error")
		return
	}

	db.masterDB.SetMaxOpenConns(c.Master.MaxOpen)
	db.masterDB.SetMaxIdleConns(c.Master.MaxIdle)
	if err = db.masterDB.Ping(); err != nil {
		err = errorsWrap(err, "master db ping error")
		return
	}

	for i := 0; i < len(c.Slave); i++ {
		var mysqlDB *sql.DB
		mysqlDB, err = sql.Open("mysql", c.Slave[i].DSN)
		if err != nil {
			err = errorsWrap(err, "init slave db error")
			return
		}

		mysqlDB.SetMaxOpenConns(c.Slave[i].MaxOpen)
		mysqlDB.SetMaxIdleConns(c.Slave[i].MaxIdle)
		if err = mysqlDB.Ping(); err != nil {
			err = errorsWrap(err, "slave db ping error")
			return
		}

		db.slaveDB = append(db.slaveDB, mysqlDB)
	}
	return
}

func (db *DB) MasterDB() *sql.DB {
	return db.masterDB
}

func (db *DB) SlaveDB() *sql.DB {
	if len(db.slaveDB) == 0 {
		return db.masterDB
	}
	n := rand.Intn(len(db.slaveDB))
	return db.slaveDB[n]
}

// MasterDBClose 释放主库的资源
func (db *DB) MasterDBClose() error {
	if db.masterDB != nil {
		return db.masterDB.Close()
	}
	return nil
}

// SlaveDBClose 释放从库的资源
func (db *DB) SlaveDBClose() (err error) {
	for i := 0; i < len(db.slaveDB); i++ {
		err = db.slaveDB[i].Close()
		if err != nil {
			return err
		}
	}
	return nil
}

type operate int64

const (
	operateMasterExec operate = iota
	operateMasterQuery
	operateMasterQueryRow
	operateSlaveQuery
	operateSlaveQueryRow
)

var operationNames = map[operate]string{
	operateMasterExec:     "masterDBExec",
	operateMasterQuery:    "masterDBQuery",
	operateMasterQueryRow: "masterDBQueryRow",
	operateSlaveQuery:     "slaveDBQuery",
	operateSlaveQueryRow:  "slaveDBQueryRow",
}

func (db *DB) operate(ctx context.Context, op operate, query string, args ...interface{}) (i interface{}, err error) {
	var (
		parent        = opentracing.SpanFromContext(ctx)
		operationName = operationNames[op]
		span          = func() opentracing.Span {
			if parent == nil {
				return opentracing.StartSpan(operationName)
			}
			return opentracing.StartSpan(operationName, opentracing.ChildOf(parent.Context()))
		}()
		logFormat  = log.LogHeaderFromContext(ctx)
		startAt    = time.Now()
		endAt      time.Time
	)

	lastModule := logFormat.Module
	lastStartTime := logFormat.StartTime
	lastEndTime := logFormat.EndTime
	lastXHop := logFormat.XHop
	defer func() {
		logFormat.Module = lastModule
		logFormat.StartTime = lastStartTime
		logFormat.EndTime = lastEndTime
		logFormat.XHop = lastXHop
	}()

	defer span.Finish()
	defer func() {
		endAt = time.Now()

		logFormat.StartTime = startAt
		logFormat.EndTime = endAt
		latencyTime := logFormat.EndTime.Sub(logFormat.StartTime).Microseconds()// 执行时间
		logFormat.LatencyTime = latencyTime
		logFormat.XHop = xhop.NewXhopNull()

		span.SetTag("error", err != nil)
		span.SetTag("db.type", "sql")
		span.SetTag("db.statement", query)
		logFormat.Module = "databus/mysql"

		if err != nil {
			db.writeError(err.Error())
			panic(err.Error())
		}else if db.IsLog == true {
			log.Infof(logFormat, "%s:[%s], params:%s, used: %d milliseconds", operationName, query,
				args, endAt.Sub(startAt).Milliseconds())
		}
	}()

	switch op {
	case operateMasterQuery:
		i, err = db.MasterDB().QueryContext(ctx, query, args...)
	case operateMasterQueryRow:
		i = db.MasterDB().QueryRowContext(ctx, query, args...)
	case operateMasterExec:
		i, err = db.MasterDB().ExecContext(ctx, query, args...)
	case operateSlaveQuery:
		i, err = db.SlaveDB().QueryContext(ctx, query, args...)
	case operateSlaveQueryRow:
		i = db.SlaveDB().QueryRowContext(ctx, query, args...)
	}
	return
}

func (db *DB) MasterDBExecContext(ctx context.Context, query string, args ...interface{}) (result sql.Result, err error) {
	r, err := db.operate(ctx, operateMasterExec, query, args...)
	if err != nil {
		return nil, err
	}
	return r.(sql.Result), err
}

func (db *DB) MasterDBQueryContext(ctx context.Context, query string, args ...interface{}) (result *sql.Rows, err error) {
	r, err := db.operate(ctx, operateMasterQuery, query, args...)
	if err != nil {
		return nil, err
	}
	return r.(*sql.Rows), err
}

func (db *DB) MasterDBQueryRowContext(ctx context.Context, query string, args ...interface{}) (result *sql.Row) {
	r, _ := db.operate(ctx, operateMasterQueryRow, query, args...)
	return r.(*sql.Row)
}

func (db *DB) SlaveDBQueryContext(ctx context.Context, query string, args ...interface{}) (result *sql.Rows, err error) {
	r, err := db.operate(ctx, operateMasterQuery, query, args...)
	if err != nil {
		return nil, err
	}
	return r.(*sql.Rows), err
}

func (db *DB) SlaveDBQueryRowContext(ctx context.Context, query string, args ...interface{}) (result *sql.Row) {
	r, _ := db.operate(ctx, operateSlaveQueryRow, query, args...)
	return r.(*sql.Row)
}

func errorsWrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func (db *DB) writeError (errMsg string){
	errLogSection := "error"
	errorLogConfig := config.GetConfig("log", errLogSection)
	errorLogdir := errorLogConfig.Key("dir").String()

	date := time.Now().Format("2006-01-02")
	dateTime := time.Now().Format("2006-01-02 15:04:05")
	file := errorLogdir + "/mysql/" + date + ".err"
	util.WriteWithIo(file,"[" +dateTime+"]" + errMsg)
}

func GetIsLog() bool {
	cfg := config.GetConfig("log", "mysql_open")
	res,err := cfg.Key("turn").Bool()
	util.Must(err)
	return res
}