module gin-frame

go 1.13

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
	github.com/jinzhu/gorm v1.9.14
	github.com/opentracing/opentracing-go v1.2.0
	github.com/why444216978/go-library v0.0.0-20200707083756-fe5dd13ddf61
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	gopkg.in/ini.v1 v1.57.0
)

replace github.com/gomodule/redigo v2.0.0+incompatible => github.com/gomodule/redigo/redis v0.0.0-20200429221454-e14091dffc1b
