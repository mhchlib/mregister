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
	if err != nil {
		log.Fatal(err)
	}
	service, err := regClient.GetService("test")
	if err != nil {
		log.Fatal(err)
	}
	log.Info(service.Address)
	log.Info(service.Metadata["key"])
}
