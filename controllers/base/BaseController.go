package base

import (
	"github.com/gin-gonic/gin"
	"gin-frame/models/user"
	"net/http"
	"strconv"
	"gin-frame/libraries/log"
	"gin-frame/libraries/util"
)

type BaseController struct {
	HasReturn	bool

	LogFormat  *log.LogFormat
	C			*gin.Context
	XhopN		uint64

	Cid			int
	AppUid		int
	AppId		int

	UserAppInfo	map[string]interface{}

	Code 		int
	Msg  		string
	Data 		map[string]interface{}
	UserMsg		string
}

func (self *BaseController) Init(c *gin.Context, productName,moduleName string){
	self.C = c
	self.XhopN = 0

	logFormat := log.NewLog()
	logFormat.Product = productName
	logFormat.Module = moduleName
	self.LogFormat = logFormat
	self.initResult()
}

func (self *BaseController) ResultJson(){
	if self.HasReturn == false {
		self.C.JSON(http.StatusOK, gin.H{
			"errno":	self.Code,
			"errmsg":	self.Msg,
			"data":	self.Data,
			"user_msg": self.UserMsg,
		})
	}
}

func (self *BaseController) Ping() {
	//self.Data["why"] = "111"
	self.ResultJson()
}

func (self *BaseController) GetHeader(key string) string {
	return self.C.Request.Header.Get(key)
}

func (self *BaseController) SetYmt() {
	self.setCid()
	self.setAppId()
	self.setAppUid()

	self.setUserAppInfo()
}

func (self *BaseController) setUserAppInfo() {
	var userAppInfo = make(map[string]interface{})

	res := user.GetInfoByAppUid(self.C, self.AppUid)
	if len(res) != 0 {
		userAppInfo = res
	}

	self.UserAppInfo = userAppInfo
}

func (self *BaseController) initResult(){
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

