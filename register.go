package register

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/register/common"
	"github.com/mhchlib/register/registerOpts"
)

// Register ...
type Register interface {
	RegisterService(serviceName string, metadata map[string]interface{}) (func(), error)
	UnRegisterService(serviceName string) error
	UnRegisterServiceAll() error
	GetService(serviceName string) (*common.ServiceVal, error)
	ListAllServices(serviceName string) ([]*common.ServiceVal, error)
	log.Logger
}

// RegisterClient ...
type RegisterClient struct {
	RegisterType registerOpts.RegistryType
	Srv          Register
}

func (r RegisterClient) RegisterService(serviceName string, metadata map[string]interface{}) (func(), error) {
	return r.Srv.RegisterService(serviceName, metadata)
}

func (r RegisterClient) UnRegisterService(serviceName string) error {
	return r.Srv.UnRegisterService(serviceName)
}

func (r RegisterClient) UnRegisterServiceAll() error {
	return r.Srv.UnRegisterServiceAll()
}

func (r RegisterClient) GetService(serviceName string) (*common.ServiceVal, error) {
	return r.Srv.GetService(serviceName)
}

func (r RegisterClient) ListAllServices(serviceName string) ([]*common.ServiceVal, error) {
	return r.Srv.ListAllServices(serviceName)
}
