package services

import (
	"gin-frame/services/location"
	"gin-frame/services/product"
)

type ServicesFactory struct{}

func (factory *ServicesFactory) GetInstance(name string) map[string]interface{} {
	instances := make(map[string]interface{})
	switch name {
	case "Location":
		instances[name] = &location.Location{}
	case "Product":
		instances[name] = &product.Product{}
	default:
		panic("services name error")
	}
	return instances
}
