package price

import (
	"gin-frame/controllers/base"
	"gin-frame/models/hangqing"
	"gin-frame/services"
	"gin-frame/services/location"
	"gin-frame/services/product"

	"github.com/gin-gonic/gin"
)

type FirstOriginPriceController struct {
	base.BaseController
	Result map[string]interface{}
}

func (self *FirstOriginPriceController) Init(c *gin.Context, productName, moduleName string) {
	self.BaseController.Init(c, productName, moduleName)
	self.BaseController.SetYmt()
	var data = make(map[string]interface{})
	self.Result = data
}

func (self *FirstOriginPriceController) Do() {
	self.action()
	self.setData()
	self.ResultJson()
}

func (self *FirstOriginPriceController) action() {
	self.getOrigin()
}

func (self *FirstOriginPriceController) setData() {
	self.Data["origin"] = self.getOrigin()

	self.Data["product"] = self.getProduct()

	self.Data["location"] = self.getLocation()
}

func (self *FirstOriginPriceController) getOrigin() []hangqing.OriginPrice {
	originPriceModel := hangqing.NewOriginPriceModel()
	return originPriceModel.GetFirst()
}

func (self *FirstOriginPriceController) getProduct() map[string]interface{} {
	servicesFactory := &services.ServicesFactory{}
	productFactory := servicesFactory.GetInstance("Product")
	productServices := productFactory["Product"].(*product.Product)
	return productServices.GetProductDetail(self.C, 8426)
}

func (self *FirstOriginPriceController) getLocation() map[string]interface{} {
	servicesFactory := &services.ServicesFactory{}
	locationFactory := servicesFactory.GetInstance("Location")
	locationServices := locationFactory["Location"].(*location.Location)
	return locationServices.GetLocationDetail(self.C, 2)
}
