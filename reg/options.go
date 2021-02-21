package reg

import (
	log "github.com/mhchlib/logger"
)

type Options struct {
	RegisterStr    string
	RegisterType   RegistryType
	Address        []string
	NameSpace      string
	ServerInstance string
	Metadata       map[string]interface{}
	logger         log.Logger
}

type Option func(*Options)

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
