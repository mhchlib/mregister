package register

import (
	log "github.com/mhchlib/logger"
)

const DEFAULT_PORT = ":8080"

type Register interface {
	RegisterService(serviceName string, metadata map[string]interface{}) (func(), error)
	UnRegisterService(serviceName string) error
	UnRegisterServiceAll() error
	GetService(serviceName string) (*ServiceVal, error)
	ListAllServices(serviceName string) ([]*ServiceVal, error)
	log.Logger
}

type ServiceVal struct {
	Address  string                 `json:"address"`
	Metadata map[string]interface{} `json:"metadata"`
}
