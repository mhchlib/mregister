package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mregister"
	"github.com/mhchlib/mregister/register"
)

func main() {
	regClient, err := mregister.InitRegister(
		register.Namespace("test_register"),
		register.ResgisterAddress("etcd://etcd.u.hcyang.top:31770"),
		register.Metadata("key", "value"),
	)
	if err != nil {
		log.Fatal(err)
	}
	regClient.RegisterService("test", map[string]interface{}{"key": "value2"})
	regClient.RegisterService("test", map[string]interface{}{"key": "value3"})
	regClient.RegisterService("test", map[string]interface{}{"key": "value4"})
	regClient.RegisterService("test", map[string]interface{}{"key": "value5"})
	regClient.RegisterService("test", map[string]interface{}{"key": "value6"})
	regClient.RegisterService("test", map[string]interface{}{"key": "value7"})
	regClient.RegisterService("test", map[string]interface{}{"key": "value8"})
	if err != nil {
		log.Fatal(err)
	}
	service, err := regClient.GetService("test")
	if err != nil {
		log.Fatal(err)
	}
	log.Info(service.Address)
	log.Info(service.Metadata["key"])

	///list
	allService, err := regClient.ListAllServices("test")
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range allService {
		log.Info(val.Metadata["key"])
	}
}
