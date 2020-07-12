package price

import (
	"gin-frame/controllers/base"
	"gin-frame/service"
	"gin-frame/service/origin_price_service"

	"github.com/gin-gonic/gin"
)

type FirstOriginPriceController struct {
	base.BaseController
	OriginPriceService *origin_price_service.OriginPriceService
	Result             map[string]interface{}
}

func (self *FirstOriginPriceController) Init(c *gin.Context, productName, moduleName string) {
	self.BaseController.Init(c, productName, moduleName)
	self.BaseController.SetYmt()
	var data = make(map[string]interface{})
	self.Result = data
}

func (self *FirstOriginPriceController) Do() {
	self.load()
	self.action()
	self.setData()
	self.ResultJson()
}

func (self *FirstOriginPriceController) load() {
	serviceFactory := &service.ServiceFactory{}
	originPriceFactory := serviceFactory.GetInstance("OriginPriceService")
	self.OriginPriceService = originPriceFactory["OriginPriceService"].(*origin_price_service.OriginPriceService)
}

func (self *FirstOriginPriceController) action() {
	origin := self.OriginPriceService.GetFirstRow(true)
	self.Data["origin"] = origin

	productId := 0
	locationId := 0

	if origin != nil {
		if origin["product_id"] != nil {
			productId = origin["product_id"].(int)
			if origin["breed_id"] != nil {
				productId = origin["breed_id"].(int)
			}
		}

		if origin["location_id"] != nil {
			locationId = origin["location_id"].(int)
		}
	}

	self.Data["product"] = self.OriginPriceService.GetOriginPriceProduct(self.C, productId)

	self.Data["location"] = self.OriginPriceService.GetOriginPriceLocation(self.C, locationId)
}

func (self *FirstOriginPriceController) setData() {

}
