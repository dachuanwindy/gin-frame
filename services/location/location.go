package location

import (
	"context"
	"strconv"

	"github.com/why444216978/go-library/libraries/config"
	"github.com/why444216978/go-library/libraries/redis"
	"github.com/why444216978/go-library/libraries/util"
	"github.com/why444216978/go-library/libraries/util/conversion"

	redigo "github.com/gomodule/redigo/redis"
)

type Location struct{}

func (location Location) Init() {}

var location *Location

const (
	redisName         = "location"
	locationDetailKey = "location::id_detail:"
	locationNameKey   = "location::id_name:"
)

func (location *Location) getRedis() *redis.RedisDB {
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

	db, err := redis.Conn("location", hostCfg, passwordCfg, portCfg, dbCfg, maxActiveCfg, maxIdleCfg, logCfg, execTime)

	return db
}

func (location *Location) GetLocationDetail(ctx context.Context, id int) map[string]interface{} {
	db := location.getRedis()

	data, _ := redigo.String(db.Do(ctx, "GET", locationDetailKey+strconv.Itoa(id)))

	return conversion.JsonToMap(data)
}

func (location *Location) BatchLocationDetail(ctx context.Context, ids []int) []string {
	db := location.getRedis()

	var args []interface{}
	for _, v := range ids {
		args = append(args, locationDetailKey+strconv.Itoa(v))
	}

	data, _ := redigo.Strings(db.Do(ctx, "MGET", args...))

	return data
}
