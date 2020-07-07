package product

import (
	"context"
	"strconv"

	"github.com/why444216978/go-library/libraries/config"
	"github.com/why444216978/go-library/libraries/redis"
	"github.com/why444216978/go-library/libraries/util"

	redigo "github.com/gomodule/redigo/redis"
)

const (
	redisName         = "location"
	locationDetailKey = "location::id_detail:"
	locationNameKey   = "location::id_name:"
)

func getDb() *redis.RedisDB {
	fileCfg := config.GetConfig("redis", redisName)

	hostCfg := fileCfg.Key("host").String()
	passwordCfg := fileCfg.Key("auth").String()
	portCfg, err := fileCfg.Key("port").Int()
	dbCfg, err := fileCfg.Key("db").Int()
	maxActiveCfg, err := fileCfg.Key("max_active").Int()
	maxIdleCfg, err := fileCfg.Key("max_idle").Int()
	logCfg, err := fileCfg.Key("is_log").Bool()
	execTime, err := fileCfg.Key("exec_timeout").Int64()
	util.Must(err)

	db, err := redis.Conn("product", hostCfg, passwordCfg, portCfg, dbCfg, maxActiveCfg, maxIdleCfg, logCfg, execTime)
	util.Must(err)

	return db
}

func GetLocationDetail(ctx context.Context, id int) string {
	db := getDb()

	data, err := redigo.String(db.Do(ctx, "GET", locationNameKey+strconv.Itoa(id)))
	util.Must(err)
	return data
}

func BatchLocationDetail(ctx context.Context, ids []int) []string {
	db := getDb()

	var args []interface{}
	for _, v := range ids {
		args = append(args, locationDetailKey+strconv.Itoa(v))
	}

	data, err := redigo.Strings(db.Do(ctx, "MGET", args...))
	util.Must(err)

	return data
}
