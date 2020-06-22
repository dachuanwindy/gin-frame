package hangqing

import (
	"gin-frame/libraries/mysql"
	"gin-frame/models/base"
	"github.com/jinzhu/gorm"
)

type OriginPrice struct {
	//gorm.Model
	Id           int `gorm:"primary_key"`
	CustomerId   int
	ProvinceId   int
	CityId       int
	CountyId     int
	LocationId   int
	ProductId    int
	BreedId      int
	PointKey     string
	DayTime      string
	PriceList    string
	DescList     string
	Status       int
	CreatedTime  int
	UpdatedTime  int
	RefuseReason string
	IsSync       int
}

func (OriginPrice) TableName() string {
	return "origin_price"
}

type OriginPriceModel struct {
	Db          *mysql.DB
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
