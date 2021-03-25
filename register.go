package mregister

import (
	"errors"
	"fmt"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mregister/plugin"
	"github.com/mhchlib/mregister/register"
	"strings"
)

// RegistryType ...
type RegistryType string

// InitRegister ...
func InitRegister(opts ...register.Option) (*RegisterClient, error) {
	options := &register.Options{}
	for _, o := range opts {
		o(options)
	}
	if options.RegisterStr != "" {
		t, address, err := parseAddressStr(options.RegisterStr)
		if err != nil {
			return nil, err
		}
		options.RegisterType = t
		options.Address = address
	}
	if options.ServerInstance == "" {
		//return nil, errors.New("server instance can not be empty")
		options.ServerInstance = "127.0.0.1:8080"
	}
	p, ok := plugin.RegisterPluginMap[options.RegisterType]
	if !ok {
		return nil, errors.New(fmt.Sprintf("register type: %s does not be supported, you can choose: %s", options.RegisterType, plugin.RegisterPluginNames))
	}
	srv, err := p.New(options)
	if err != nil {
		return nil, err
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &RegisterClient{
		RegisterType: RegistryType(options.RegisterType),
		Srv:          srv,
	}, nil
}

// ConfigSeparateSymbol ...
const ConfigSeparateSymbol = "://"

func parseAddressStr(str string) (string, string, error) {
	splits := strings.Split(str, ConfigSeparateSymbol)
	if len(splits) != 2 {
		return "", "", errors.New(str + " is invalid Address")
	}
	return splits[0], splits[1], nil
}

// RegisterClient ...
type RegisterClient struct {
	RegisterType RegistryType
	Srv          register.Register
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

func (r RegisterClient) GetService(serviceName string) (*register.ServiceVal, error) {
	return r.Srv.GetService(serviceName)
}

func (r RegisterClient) ListAllServices(serviceName string) ([]*register.ServiceVal, error) {
	return r.Srv.ListAllServices(serviceName)
}
