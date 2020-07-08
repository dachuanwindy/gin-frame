package product

import (
	"context"
	"strconv"

	"github.com/why444216978/go-library/libraries/config"
	"github.com/why444216978/go-library/libraries/redis"
	"github.com/why444216978/go-library/libraries/util/conversion"
	"github.com/why444216978/go-library/libraries/util/error"

	redigo "github.com/gomodule/redigo/redis"
)

type Product struct{}

func (product Product) Init() {}

var product *Product

const (
	redisName        = "product"
	productDetailKey = "product::id_detail:"
	productNameKey   = "product::id_name:"
)

func (product *Product) gerRedis() *redis.RedisDB {
	fileCfg := config.GetConfig("redis", redisName)

	hostCfg := fileCfg.Key("host").String()
	passwordCfg := fileCfg.Key("auth").String()
	portCfg, err := fileCfg.Key("port").Int()
	dbCfg, err := fileCfg.Key("db").Int()
	maxActiveCfg, err := fileCfg.Key("max_active").Int()
	maxIdleCfg, err := fileCfg.Key("max_idle").Int()
	logCfg, err := fileCfg.Key("is_log").Bool()
	execTime, err := fileCfg.Key("exec_timeout").Int64()
	error.Must(err)

	db, err := redis.Conn("product", hostCfg, passwordCfg, portCfg, dbCfg, maxActiveCfg, maxIdleCfg, logCfg, execTime)
	error.Must(err)

	return db
}

func (product *Product) GetProductDetail(ctx context.Context, id int) map[string]interface{} {
	db := product.gerRedis()

	data, _ := redigo.String(db.Do(ctx, "GET", productDetailKey+strconv.Itoa(id)))

	return conversion.JsonToMap(data)
}

func (product *Product) BatchProductDetail(ctx context.Context, ids []int) []string {
	db := product.gerRedis()

	var args []interface{}
	for _, v := range ids {
		args = append(args, productDetailKey+strconv.Itoa(v))
	}

	data, _ := redigo.Strings(db.Do(ctx, "MGET", args...))

	return data
}
