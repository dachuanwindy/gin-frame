package main

import (
    "fmt"
    "gin-frame/libraries/endless"
    "gin-frame/routers"
    "log"
    "runtime"
    "strconv"
    "syscall"
)

const Port = 777
const productName = "zhuayu"
const moduleName = "hangqing_chandi"

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
    server := routers.InitRouter(Port, productName, moduleName)
    
    tmpServer := endless.NewServer(fmt.Sprintf(":%s", strconv.Itoa(Port)), server)
    tmpServer.BeforeBegin = func(add string) {
        log.Printf("Actual pid is %d", syscall.Getpid())
    }
    err := tmpServer.ListenAndServe()
    if err != nil {
        log.Printf("Server err: %v", err)
    }
}

