package origin_price_service

import (
	"context"
	"gin-frame/dao"
	"gin-frame/dao/origin_price_dao"
	"gin-frame/library"
	"gin-frame/library/location"
	"gin-frame/library/product"
	"log"
	"sync"
)

type OriginPriceService struct {
	productService  *product.ProductLibrary
	locationService *location.LocationLibrary
	originPriceDao  *origin_price_dao.OriginPriceDao
}

var onceOriginPriceService sync.Once
var originPriceService *OriginPriceService

func NewObj() *OriginPriceService {
	onceOriginPriceService.Do(func() {
		originPriceService = &OriginPriceService{}

		libraryFactory := &library.LibraryFactory{}
		locationInterface := libraryFactory.GetInstance("Location")
		originPriceService.locationService = locationInterface["Location"].(*location.LocationLibrary)
		productInterface := libraryFactory.GetInstance("Product")
		originPriceService.productService = productInterface["Product"].(*product.ProductLibrary)

		daoFactory := dao.DaoFactory{}
		originPriceInterface := daoFactory.GetInstance("OriginPriceDao")
		originPriceService.originPriceDao = originPriceInterface["OriginPriceDao"].(*origin_price_dao.OriginPriceDao)

		log.Printf("new origin_price_service")
	})

	return originPriceService
}

func (self *OriginPriceService) GetFirstRow(noCache bool) map[string]interface{} {
	return self.originPriceDao.GetFirstRow(true)
}

func (self *OriginPriceService) GetOriginPriceLocation(ctx context.Context, locationId int) map[string]interface{} {
	return self.locationService.GetLocationDetail(ctx, locationId)
}

func (self *OriginPriceService) GetOriginPriceProduct(ctx context.Context, productId int) map[string]interface{} {
	return self.productService.GetProductDetail(ctx, productId)

}
