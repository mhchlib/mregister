package reg

import log "github.com/mhchlib/logger"

type Register interface {
	RegisterService(serviceName string, metadata map[string]interface{}) error
	UnRegisterService(serviceName string) error
	UnRegisterServiceAll() error
	GetService(serviceName string) (*ServiceVal, error)
	ListAllServices(serviceName string) ([]*ServiceVal, error)
	log.Logger
}

type ServiceVal struct {
	Address  string
	Metadata map[string]interface{}
}
