package hangqing

import (
	"context"
	"github.com/jmoiron/sqlx"
	"gin-frame/libraries/util"
	"gin-frame/models/base"
)

type HqCustomer struct {
	Province_id     int
	City_id      	int
	County_id      	int
	Location_id     int
	Market_info_id  int
	Point_key      	string
	Point_key2      string
	Product_id      int
	Breed_id     	int
	Customer_id     int
}

func GetCustomerBreedsByCid(ctx context.Context, cid int) ( data []map[string]interface{} ) {
	query := "select province_id,city_id,county_id,location_id,market_info_id,point_key,point_key2,product_id,breed_id,customer_id from hq_customer where customer_id = ?"
	db := base.GetConn("hangqing")
	rows, err := db.SlaveDBQueryContext(ctx, query, cid)
	util.Must(err)

	var list = []*HqCustomer{}
	err = sqlx.StructScan(rows, &list)
	util.Must(err)

	if len(list) == 0 {
		return
	}

	for _, v := range list {
		tmp := util.StructToMap(*v)
		data = append(data, tmp)
	}

	return
}
