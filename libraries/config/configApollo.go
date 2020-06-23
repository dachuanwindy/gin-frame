package config

import (
	"fmt"
	util_err "gin-frame/libraries/util/error"
	"github.com/astaxie/beego/config"
	agollo "github.com/zouyx/agollo/v3"
	env "gin-frame/libraries/config/env"
)

var AppConfig config.Configer

func init() {
	sysEnvConfig := &env.SysEnvConfig{AppID: "moments-server"}
	agollo.InitCustomConfig(sysEnvConfig.LoadSysConfig)
	if err := agollo.Start(); err != nil {
		panic(err)
	}
	base := agollo.GetStringValue("BASE", "")
	if base == "" {
		panic("apollo 配置为空")
	}

	var err error
	AppConfig, err = config.NewConfigData("ini", []byte(base))
	util_err.Must(err)
}

func Test() {
	section, err := AppConfig.GetSection("db")
	util_err.Must(err)
	fmt.Println(section)
	fmt.Println(section["user_dynamic"])
}
