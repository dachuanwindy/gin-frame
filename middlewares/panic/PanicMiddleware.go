package panic

import (
	"gin-frame/controllers/base"
	"gin-frame/codes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"gin-frame/libraries/util"
)

func ThrowPanic(dir,moduleName string, baseController *base.BaseController, area int) gin.HandlerFunc{
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				baseController.HasReturn = true //用于处理重复返回JSON标志
				c.JSON(http.StatusInternalServerError, gin.H{
					"errno":	codes.SERVER_ERROR,
					"errmsg":	codes.ErrorMsg[codes.SERVER_ERROR],
					"data":	make(map[string]interface{}),
					"user_msg": codes.ErrorUserMsg[codes.SERVER_ERROR],
				})

				DebugStack := ""
				for _, v := range strings.Split(string(debug.Stack()), "\n") {
					DebugStack += v + "\n"
				}

				dateTime := time.Now().Format("2006-01-02 15:04:05")
				file := util.CreateDateDir(dir, moduleName + ".err." + util.HostName() + ".")
				file = file + "/" + strconv.Itoa(util.RandomN(area))
				util.WriteWithIo(file,"[" +dateTime+"]")
				util.WriteWithIo(file, fmt.Sprintf("%v\r\n", err))
				util.WriteWithIo(file, DebugStack)
				c.Done()
			}
		}()
		c.Next()

	}
}