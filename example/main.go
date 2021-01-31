package main

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/register"
	"github.com/mhchlib/register/reg"
)

func main() {
	regClient, err := register.InitRegister("etcd", func(options *reg.Options) {
		options.NameSpace = "test_register"
		options.Address = []string{"etcd.u.hcyang.top:31770"}
		options.Metadata = make(map[string]interface{})
		options.Metadata["key"] = "value1"
	})
	if err != nil {
		log.Fatal(err)
	}
	err = regClient.RegisterService("test", map[string]interface{}{"key": "value2"})
	err = regClient.RegisterService("test", map[string]interface{}{"key": "value3"})
	err = regClient.RegisterService("test", map[string]interface{}{"key": "value4"})
	err = regClient.RegisterService("test", map[string]interface{}{"key": "value5"})
	err = regClient.RegisterService("test", map[string]interface{}{"key": "value6"})
	err = regClient.RegisterService("test", map[string]interface{}{"key": "value7"})
	err = regClient.RegisterService("test", map[string]interface{}{"key": "value8"})
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
	allService, err := regClient.ListAllService("test")
	if err != nil {
		log.Fatal(err)
	}
	for _, val := range allService {
		log.Info(val.Metadata["key"])
	}
}
