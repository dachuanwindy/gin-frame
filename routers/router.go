package routers

import (
	"github.com/gin-gonic/gin"

	"gin-frame/libraries/config"
	"gin-frame/libraries/util"

	"gin-frame/controllers/base"
	"gin-frame/controllers/price"

	"gin-frame/middlewares/log"
	"gin-frame/middlewares/trace"
	"gin-frame/middlewares/panic"
)

func InitRouter(Port int, productName, moduleName string) *gin.Engine {
	server := gin.New()

	baseController 	 := &base.BaseController{}
	firstOriginPriceController := &price.FirstOriginPriceController{}

	server.Use(gin.Recovery())
	server.Use(trace.OpenTracing(productName))
	server.Use(log.LoggerMiddleware(Port, productName, moduleName))
	errLogSection := "error"
	errorLogConfig := config.GetConfig("log", errLogSection)
	errorLogDir := errorLogConfig.Key("dir").String()
	errorLogArea, err := errorLogConfig.Key("area").Int()
	util.Must(err)
	server.Use(panic.ThrowPanic(errorLogDir, moduleName, baseController, errorLogArea))
	//server.Use(dump.BodyDump())

	group := server.Group("")
	group.GET("/ping", func(c *gin.Context) {
		baseController.Init(c, productName, "ping")
		baseController.Ping()
	})

	group.GET("/origin/first_origin_price", func(c *gin.Context) {
		firstOriginPriceController.Init(c, productName, moduleName)
		firstOriginPriceController.Do()
	})
	return server
}
