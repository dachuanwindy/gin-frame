package services

import (
	"context"
	"github.com/why444216978/go-library/libraries/redis"
	"github.com/why444216978/go-library/libraries/util"
	"strconv"

	redigo "github.com/gomodule/redigo/redis"
)

const (
	product_detail = "hangqing_category::id_detail:"
	product_name   = "product::id_name:"
)

func GetProductDetail(ctx context.Context, id int) string {
	db, err := redis.Conn("product")
	util.Must(err)

	data, err := redigo.String(db.Do(ctx, "GET", product_name+strconv.Itoa(id)))
	util.Must(err)
	return data
}

func BatchProductDetail(ctx context.Context, ids []int) []string {
	db, err := redis.Conn("product")
	util.Must(err)

	var args []interface{}
	for _, v := range ids {
		args = append(args, product_detail+strconv.Itoa(v))
	}

	data, err := redigo.Strings(db.Do(ctx, "MGET", args...))
	util.Must(err)

	return data
}
