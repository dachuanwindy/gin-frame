package price

import (
	"gin-frame/controllers/base"
	"gin-frame/models/hangqing"
	"gin-frame/services/product"

	"github.com/gin-gonic/gin"
)

type FirstOriginPriceController struct {
	base.BaseController
	OriginPrice []hangqing.OriginPrice
	Result      map[string]interface{}
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
	self.Data["origin"] = self.OriginPrice
	self.Data["redis"] = product.GetProductDetail(self.C, 8426)
}

func (self *FirstOriginPriceController) getOrigin() {
	originPriceModel := hangqing.NewOriginPriceModel()
	self.OriginPrice = originPriceModel.GetFirst()
}
