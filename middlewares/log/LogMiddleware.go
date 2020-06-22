package log

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"gin-frame/libraries/config"
	"gin-frame/libraries/log"
	"gin-frame/libraries/util/dir"
	"gin-frame/libraries/util/conversion"
	"gin-frame/libraries/util/url"
	"gin-frame/libraries/util/random"
	"gin-frame/libraries/util/sys"
	"gin-frame/libraries/xhop"
	"net/http"
	"strconv"
	"time"
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

func LoggerMiddleware(port int, productName, moduleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logHeader := &log.LogFormat{}

		runLogSection := "run"
		runLogConfig := config.GetConfig("log", runLogSection)
		runLogDir := runLogConfig.Key("dir").String()
		area, _ := runLogConfig.Key("area").Int()

		file := dir.CreateHourLogFile(runLogDir, moduleName+".log."+sys.HostName()+".")
		file = file + "/" + strconv.Itoa(random.RandomN(area))

		log.Init(&log.LogConfig{
			File:           file,
			Path:           runLogDir,
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
			logID = log.NewObjectId().Hex()
		}

		ctx := c.Request.Context()
		dst := new(log.LogFormat)
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

		ctx = log.ContextWithLogHeader(ctx, dst)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set(HeaderXLogID, dst.LogId)
		c.Writer.Header().Set(HeaderXHop, dst.XHop.String())

		reqBody := []byte{}
		if c.Request.Body != nil { // Read
			reqBody, _ = ioutil.ReadAll(c.Request.Body)
		}
		strReqBody := string(reqBody)

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody)) // Reset
		responseWriter := &bodyLogWriter{body: bytes.NewBuffer(nil), ResponseWriter: c.Writer}
		c.Writer = responseWriter

		dst.StartTime = time.Now()

		c.Next() // 处理请求

		responseBody := responseWriter.body.String()

		log.Info(dst, map[string]interface{}{
			"requestHeader": c.Request.Header,
			"requestBody":   conversion.JsonToMap(strReqBody),
			"responseBody":  conversion.JsonToMap(responseBody),
			"uriQuery":      url.ParseUriQueryToMap(c.Request.URL.RawQuery),
		})
	}
}
