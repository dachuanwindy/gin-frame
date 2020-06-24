package routers

import (
	"github.com/gin-gonic/gin"
	"gin-frame/libraries/config"
	"gin-frame/libraries/util"
	"gin-frame/middlewares/log"
	"gin-frame/middlewares/trace"
	"gin-frame/middlewares/panic"
	"gin-frame/controllers/base"
	"gin-frame/controllers/price"

)

func InitRouter(port int, productName, moduleName,env string) *gin.Engine {
	server := gin.New()

	server.Use(gin.Recovery())

	server.Use(trace.OpenTracing(productName))

	logFields := make(map[string]string, 3)
	logFieldsSection := "log_fields"
	logFieldsConfig := config.GetConfig("log", logFieldsSection)
	logFields["query_id"] = logFieldsConfig.Key("query_id").String()
	logFields["header_id"] = logFieldsConfig.Key("header_id").String()
	logFields["header_hop"] = logFieldsConfig.Key("header_hop").String()

	runLogSection := "run"
	runLogConfig := config.GetConfig("log", runLogSection)
	runLogDir := runLogConfig.Key("dir").String()
	runLogArea, _ := runLogConfig.Key("area").Int()
	server.Use(log.LoggerMiddleware(port, logFields, runLogDir, runLogArea, productName, moduleName, env))

	errLogSection := "error"
	errorLogConfig := config.GetConfig("log", errLogSection)
	errorLogDir := errorLogConfig.Key("dir").String()
	errorLogArea, err := errorLogConfig.Key("area").Int()
	util.Must(err)
	server.Use(panic.ThrowPanic(port, logFields, errorLogDir, errorLogArea, productName, moduleName, env))
	//server.Use(dump.BodyDump())

	baseController 	 := &base.BaseController{}
	firstOriginPriceController := &price.FirstOriginPriceController{}

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
