package log

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"gin-frame/libraries/config"
	"gin-frame/libraries/util"
	"gin-frame/libraries/xhop"
)

const (
	QueryLogID   = "logid"
	HeaderXLogID = "x-logid"
	HeaderXHop   = "x-hop"
)


type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func LoggerMiddleware(port int, productName, moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logHeader := &LogFormat{}

		runLogSection := "run"
		runLogConfig := config.GetConfig("log", runLogSection)
		runLogdir := runLogConfig.Key("dir").String()

		file := util.CreateDateDir(runLogdir, moduleName + ".log." + util.HostName() + ".")
		file = file + "/" + strconv.Itoa(util.RandomN(9))

		Init(&LogConfig{
			File:			file,
			Path:			runLogdir,
			Mode:           1,
			AsyncFormatter: false,
			Debug:          true,
		}, file)

		var logID string
		switch {
		case c.Query(QueryLogID) != "":
			logID = c.Query(QueryLogID)
		case c.Request.Header.Get(HeaderXLogID) != "":
			logID = c.Request.Header.Get(HeaderXLogID)
		default:
			logID = NewObjectId().Hex()
		}

		ctx := c.Request.Context()
		dst := new(LogFormat)
		*dst = *logHeader

		dst.HttpCode = c.Writer.Status()
		dst.Port = port
		dst.LogId = logID
		dst.Method = c.Request.Method
		dst.CallerIp = c.ClientIP()
		dst.UriPath = c.Request.RequestURI
		dst.XHop = NextXhop(c.Request.Header)
		dst.Product = productName
		dst.Module = moduleName

		//TODO
		dst.Env = "development"

		ctx = ContextWithLogHeader(ctx, dst)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(HeaderXLogID, dst.LogId)
		c.Writer.Header().Set(HeaderXHop, dst.XHop.String())

		reqBody := []byte{}
		if c.Request.Body != nil { // Read
			reqBody, _ = ioutil.ReadAll(c.Request.Body)
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset
		blw := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = blw

		dst.StartTime = time.Now()

		c.Next() // 处理请求

		Info(dst, map[string]interface{}{
			"reqHeader":  c.Request.Header,
			"reqBody":    string(reqBody),
			"respBody":   string(blw.body.Bytes()),
			"query":      c.Request.URL.Query(),
		})
	}
}

func NextByHeaderXhop(xhopHex string, n uint64) *xhop.XHop {
	var xHopInfo *xhop.XHop
	var err error
	if xhopHex == "" {
		xHopInfo = xhop.NewXHop()
	} else if xHopInfo, err = xhop.NewFromHex(xhopHex); err != nil {
		xHopInfo = xhop.NewXHop()
	} else {
		xHopInfo = xHopInfo.NextN(n)
	}

	return xHopInfo
}

func NextXhop(header http.Header) *xhop.XHop {
	var xhopHex = header.Get(HeaderXHop)
	var xHopInfo *xhop.XHop
	var err error
	if xhopHex == "" {
		xHopInfo = xhop.NewXHop()
	} else if xHopInfo, err = xhop.NewFromHex(xhopHex); err != nil {
		xHopInfo = xhop.NewXHop()
	} else {
		xHopInfo = xHopInfo.Next()
	}

	return xHopInfo
}
