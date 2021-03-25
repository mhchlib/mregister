package register

import (
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/register/etcd"
	"github.com/mhchlib/register/memory"
	"github.com/mhchlib/register/registerOpts"
	"strings"
)

// InitRegister ...
func InitRegister(opts ...registerOpts.Option) (*RegisterClient, error) {
	options := &registerOpts.Options{}
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
	var srv Register
	var err error
	var registerType registerOpts.RegistryType
	switch options.RegisterType {
	case string(registerOpts.REGISTRYTYPE_ETCD):
		registerType = registerOpts.REGISTRYTYPE_ETCD
		srv, err = etcd.NewEtcdRegister(options)
	case string(registerOpts.REGISTRYTYPE_MEMORY):
		registerType = registerOpts.REGISTRYTYPE_MEMORY
		srv, err = memory.NewMemoryRegister(options)
	default:
		return nil, errors.New(string("registry type: " + options.RegisterType + " can not be supported, you can choose: etcd,memory"))
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &RegisterClient{
		RegisterType: registerType,
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
