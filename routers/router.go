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
	chandiController := &price.PriceController{}

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

	group.POST("/hangqing_chandi/api/customer_origin_list", func(c *gin.Context) {
		chandiController.Init(c, productName, moduleName)
		chandiController.Do()
	})
	return server
}
