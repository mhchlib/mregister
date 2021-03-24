package register

import (
	log "github.com/mhchlib/logger"
)

// Register ...
type Register interface {
	RegisterService(serviceName string, metadata map[string]interface{}) (func(), error)
	UnRegisterService(serviceName string) error
	UnRegisterServiceAll() error
	GetService(serviceName string) (*ServiceVal, error)
	ListAllServices(serviceName string) ([]*ServiceVal, error)
	log.Logger
}

// RegisterClient ...
type RegisterClient struct {
	RegisterType RegistryType
	Srv          Register
}

// ServiceVal ...
type ServiceVal struct {
	Address  string                 `json:"address"`
	Metadata map[string]interface{} `json:"metadata"`
}
