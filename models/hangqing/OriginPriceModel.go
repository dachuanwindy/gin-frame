package hangqing

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gin-frame/libraries/util"
	"gin-frame/models/base"
)

type OriginPrice struct {
	Id 				int
	Customer_id     int
	Province_id     int
	City_id      	int
	County_id      	int
	Location_id     int
	Product_id		int
	Breed_id		int
	Point_key		string
	Day_time		string
	Price_list		string
	Desc_list		string
	Updated_time	string
}

func GetWithinThreeDaysOriginPriceByCustomerId(ctx context.Context, cid int) (data []map[string]interface{}) {
	db := base.GetConn("hangqing")

	dateTime := util.GetDaysAgoZeroTime(-3)
	sql := "select id,customer_id, province_id,city_id,county_id,location_id,product_id,breed_id,point_key,day_time,price_list,desc_list,updated_time from origin_price " +
		"where customer_id = ? and created_time >= ? order by updated_time desc"
	rows, err := db.MasterDBQueryContext(ctx, sql, cid, dateTime)
	util.Must(err)

	var list = []*OriginPrice{}
	err = sqlx.StructScan(rows, &list)
	util.Must(err)

	for _, v := range list {
		tmp := util.StructToMap(*v)
		data = append(data, tmp)
	}

	return
}



