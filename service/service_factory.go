package service

import (
	"gin-frame/service/origin_price_service"
)

type ServiceFactory struct{}

func (factory *ServiceFactory) GetInstance(name string) map[string]interface{} {
	instances := make(map[string]interface{})

	switch name {
	case "OriginPriceService":
		instances[name] = origin_price_service.NewObj()
	default:
		panic("service name error")
	}
	return instances
}
