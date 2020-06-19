package hangqing

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gin-frame/libraries/util"
	"gin-frame/models/base"
)

type HqCustomerLocation struct {
	Province_id     int
	City_id      	int
	County_id      	int
	Location_id     int
	Market_info_id  int
	Customer_id     int
	Customer_name   string
}

func GetCustomerLocationByCid(ctx context.Context, cid int) (data map[string]interface{}) {
	db := base.GetConn("hangqing")

	query := "select province_id,city_id,county_id,location_id,market_info_id,customer_id,customer_name from hq_customer_location where customer_id = ? order by id desc limit 1"
	rows,err := db.SlaveDBQueryContext(ctx, query, cid)
	util.Must(err)

	var list = []*HqCustomerLocation{}
	err = sqlx.StructScan(rows, &list)
	util.Must(err)

	if len(list) == 0 {
		return
	}
	data = util.StructToMap(*list[0])

	return

}
