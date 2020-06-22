package price

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"gin-frame/libraries/config"
	"gin-frame/libraries/util"

	"gin-frame/controllers/base"
	"gin-frame/codes"
	"gin-frame/models/hangqing"
	"gin-frame/services"
)

type PriceController struct {
	base.BaseController
	PriceTitle				string
	IsAuthorize  			int
	AuthorizeUrl 			string
	LocationName			string
	AllBreeds				[]map[string]interface{}
	CustomerProductsNum		int
	CustomerProducts		[]map[string]interface{}
	Result  				map[string]interface{}
}

func (self *PriceController) Init(c *gin.Context, productName,moduleName string){
	self.BaseController.Init(c, productName, moduleName)
	self.BaseController.SetYmt()
	self.PriceTitle = "及时发表最新行情，获得更多客商关注"
	self.IsAuthorize = 1
	self.AuthorizeUrl = ""
	self.CustomerProductsNum = 0
	var data  = make(map[string]interface{})
	self.Result = data
}

func (self *PriceController) Do() {
	self.action()
	self.setData()
	self.ResultJson()
}

func (self *PriceController) action(){
	customerLocation := hangqing.GetCustomerLocationByCid(self.C.Request.Context(), self.Cid)
	if len(customerLocation) == 0 {
		self.setAuthUrl()
		self.Code = codes.NO_AUTHORIZE_SPY
		self.Msg = codes.ErrorMsg[codes.NO_AUTHORIZE_SPY]
		self.UserMsg = codes.ErrorUserMsg[codes.NO_AUTHORIZE_SPY]
		return
	}

	self.AllBreeds = hangqing.GetCustomerBreedsByCid(self.C.Request.Context(), self.Cid)
	if len(self.AllBreeds) == 0 {
		return
	}

	self.LocationName = services.GetLocationDetail(self.C.Request.Context(), customerLocation["location_id"].(int))
	self.formatLocation()
	self.formatProductAndBreedName()
	self.formatLastPrice()
}

func (self *PriceController) setData() {
	self.Data["price_title"] = self.PriceTitle
	self.Data["is_authorize"] = self.IsAuthorize
	self.Data["authorize_url"] = self.AuthorizeUrl
	self.Data["customer_products_num"] = len(self.AllBreeds)

	var customerProducts []map[string]interface{}
	for _,v := range self.AllBreeds {
		tmp := make(map[string]interface{})
		tmp["id"] = v["id"]
		tmp["price_id"] = v["price_id"]
		tmp["product_name"] = v["product_name"]
		tmp["breed_name"] = v["breed_name"]
		tmp["location"] = v["location"]
		tmp["type"] = v["type"]
		tmp["button"] = v["button"]
		tmp["price_list"] = v["price_list"]
		customerProducts = append(customerProducts, tmp)
	}
	self.Data["customer_products"] = customerProducts
}

func (self *PriceController) formatLastPrice() {
	allOriginPrices := hangqing.GetWithinThreeDaysOriginPriceByCustomerId(self.C.Request.Context(), self.Cid)

	//没有历史报价，初始化每一项为报价相关
	if len(allOriginPrices) == 0 {
		for _,v := range self.AllBreeds {
			v["price_id"] = 0
			v["type"] = 1
			v["button"] = "报价"
			v["price_list"] = make(map[string]interface{})
		}
	}else {
		var data = make(map[string]interface{})
		var keyProductId string
		var keyBreedtId string
		var key string
		var tmp = make(map[string]interface{})
		for _,v := range allOriginPrices {
			keyProductId = strconv.Itoa(v["product_id"].(int))
			keyBreedtId = strconv.Itoa(v["breed_id"].(int))
			key = fmt.Sprintf("%s_%s", keyProductId, keyBreedtId)
			tmp["id"] = v["id"]
			tmp["price_list"] = util.JsonToMapArray(v["price_list"].(string))
			tmp["desc_list"] = util.JsonToMapArray(v["desc_list"].(string))
			tmp["updated_time"] = v["updated_time"]
			data[key] = tmp
		}

		for _,v := range self.AllBreeds {
			v["price_id"] = 0
			v["type"] = 1
			v["type"] = 1
			v["button"] = "报价"
			v["price_list"] = make(map[string]interface{})
			keyProductId = strconv.Itoa(v["product_id"].(int))
			keyBreedtId = strconv.Itoa(v["breed_id"].(int))
			key = fmt.Sprintf("%s_%s", keyProductId, keyBreedtId)
			if data[key] != nil {
				tmp = data[key].(map[string]interface{})
				v["price_id"] = tmp["id"]
				if tmp["price_list"] != nil {
					item := tmp["price_list"].([]map[string]interface{})
					for k,v := range item {
						item[k]["unit_desc"] = "元/斤"
						change := 0
						if v["change"] != nil {
							change = int(v["change"].(float64))
						}
						item[k]["change"] = change
					}
				}
				v["price_list"] = tmp["price_list"]
				updated_time,err := strconv.Atoi(tmp["updated_time"].(string))
				v["last_time"] = util.TimeToHuman(updated_time)
				util.Must(err)
				if updated_time > int(time.Now().Add(-time.Hour * 6).Unix()) {
					v["type"] = 2;
					v["button"] = "修改报价";
				}
			}
		}
	}
}

func (self *PriceController) formatProductAndBreedName() {
	var ids []int
	for _,v := range self.AllBreeds {
		ids = append(ids, v["product_id"].(int))

		if v["breed_id"].(int) == 0 {
			v["breed_id"] = v["product_id"]
		}
		ids = append(ids, v["breed_id"].(int))
	}

	res := services.BatchProductDetail(self.C.Request.Context(), ids)
	var idsName = make(map[float64]string)
	var tmp = make(map[string]interface{})
	for _,v := range res {
		if len(v) == 0 {
			continue
		}
		tmp = util.JsonToMap(v)
		idsName[tmp["id"].(float64)] = tmp["name"].(string)
	}

	for _,v := range self.AllBreeds {
		v["product_name"] = idsName[float64(v["product_id"].(int))]
		v["breed_name"] = idsName[float64(v["product_id"].(int))]
	}
}

func (self *PriceController) formatLocation() {
	for _,v := range self.AllBreeds {
		v["location"] = self.LocationName
	}
}

func (self *PriceController) setAuthUrl(){
	urlCfg := config.GetConfig("url", "syp_auth")
	self.IsAuthorize = 0
	self.AuthorizeUrl = urlCfg.Key("auth_syp_url").String()
}
