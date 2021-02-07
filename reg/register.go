package reg

import log "github.com/mhchlib/logger"

type RegistryType string

const (
	RegistryType_Etcd RegistryType = "etcd"
)

type Register interface {
	RegisterService(serviceName string, metadata map[string]interface{}) error
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
