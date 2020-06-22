package main

import (
	"fmt"
	"gin-frame/models/hangqing"
)

func main(){
	originPriceModel := hangqing.NewOriginPriceModel()
	resList := originPriceModel.GetFirst()
	fmt.Println(resList)
}