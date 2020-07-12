package library

import (
	"gin-frame/library/location"
	"gin-frame/library/product"
)

type LibraryFactory struct{}

func (factory *LibraryFactory) GetInstance(name string) map[string]interface{} {
	instances := make(map[string]interface{})

	switch name {
	case "Location":
		instances[name] = location.NewObj()
	case "Product":
		instances[name] = product.NewObj()
	default:
		panic("library name error")
	}
	return instances
}
