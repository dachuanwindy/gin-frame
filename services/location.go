package services

import (
	"context"
	"strconv"
	"gin-frame/libraries/util"
	"gin-frame/libraries/redis"

	redigo "github.com/gomodule/redigo/redis"
)

const (
	location_detail  = 	"location::id_detail:"
	location_name	 =	"location::id_name:"
)

func GetLocationDetail(ctx context.Context, id int) string {
	db, err := redis.Conn("location")
	util.Must(err)

	data, err := redigo.String( db.Do(ctx, "GET", location_name + strconv.Itoa(id)) )
	util.Must(err)
	return data
}

func BatchLocationDetail(ctx context.Context, ids []int) []string {
	db, err := redis.Conn("location")
	util.Must(err)

	var args []interface{}
	for _,v := range ids {
		args = append(args, location_detail + strconv.Itoa(v))
	}

	data, err := redigo.Strings(db.Do(ctx,"MGET",  args...))
	util.Must(err)

	return data
}