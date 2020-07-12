package base

import (
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/why444216978/go-library/libraries/log"
	"github.com/why444216978/go-library/libraries/util"
)

var lock sync.RWMutex

type BaseController struct {
	HasError bool

	LogFormat *log.LogFormat
	C         *gin.Context
	XhopN     uint64

	Cid    int
	AppUid int
	AppId  int

	UserAppInfo map[string]interface{}

	Code    int
	Msg     string
	Data    map[string]interface{}
	UserMsg string
}

func (self *BaseController) Init(c *gin.Context, productName, moduleName string) {
	self.C = c
	self.XhopN = 0

	logFormat := log.NewLog()
	logFormat.Product = productName
	logFormat.Module = moduleName
	self.LogFormat = logFormat
	self.initResult()
}

func (self *BaseController) ResultJson() {
	self.C.JSON(http.StatusOK, gin.H{
		"errno":    self.Code,
		"errmsg":   self.Msg,
		"data":     self.Data,
		"user_msg": self.UserMsg,
	})
}

func (self *BaseController) Ping() {
	self.ResultJson()
}

func (self *BaseController) GetHeader(key string) string {
	return self.C.Request.Header.Get(key)
}

func (self *BaseController) SetYmt() {
	self.setCid()
	self.setAppId()
	self.setAppUid()
}

func (self *BaseController) initResult() {
	data := make(map[string]interface{})
	self.Code = 0
	self.Msg = "success"
	self.Data = data
	self.UserMsg = ""
}

func (self *BaseController) setCid() {
	var cid = 0
	res := self.GetHeader("X-Customer-Id")
	if res != "" {
		res, err := strconv.Atoi(res)
		util.Must(err)
		cid = res
	}

	self.Cid = cid
}

func (self *BaseController) setAppUid() {
	var appUid = 0
	res := self.GetHeader("X-User-Id")
	if res != "" {
		res, err := strconv.Atoi(res)
		util.Must(err)
		appUid = res
	}

	self.AppUid = appUid
}

func (self *BaseController) setAppId() {
	var AppId = 0
	res := self.GetHeader("X-User-Agent")
	if res != "" {
		res, err := strconv.Atoi(res)
		util.Must(err)
		AppId = res
	}

	self.AppId = AppId
}
