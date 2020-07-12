package origin_price_model

import (
	"gin-frame/models/base"

	"github.com/jinzhu/gorm"
	"github.com/why444216978/go-library/libraries/mysql"
)

type OriginPrice struct {
	//gorm.Model
	Id            int `gorm:"primary_key"`
	Customer_id   int
	Province_id   int
	City_id       int
	County_id     int
	Location_id   int
	Product_id    int
	Breed_id      int
	Point_key     string
	Day_time      string
	Price_list    string
	Desc_list     string
	Status        int
	Created_time  int
	Updated_time  int
	Refuse_reason string
	Is_sync       int
}

func (OriginPrice) TableName() string {
	return "origin_price"
}

type OriginPriceModel struct {
	Db *mysql.DB
}

var instance *OriginPriceModel

func NewOriginPriceModel() *OriginPriceModel {
	if instance == nil {
		instance = &OriginPriceModel{}
		instance.Db = base.GetInstance("hangqing")
	}
	return instance
}

func (instance *OriginPriceModel) GetFirst() []OriginPrice {
	originPrices := []OriginPrice{}
	orm := instance.Db.SlaveOrm()
	dbRes := orm.First(&originPrices)
	instance.checkRes(dbRes)
	return originPrices
}

func (instance *OriginPriceModel) checkRes(dbRes *gorm.DB) {
	if dbRes.Error != nil {
		panic(dbRes.Error)
	}
}
