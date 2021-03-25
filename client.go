package register

import (
	"errors"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/register/etcd"
	opts2 "github.com/mhchlib/register/registerOpts"
	"strings"
)

// InitRegister ...
func InitRegister(opts ...opts2.Option) (*RegisterClient, error) {
	options := &opts2.Options{}
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
	var registerType opts2.RegistryType
	switch options.RegisterType {
	case string(opts2.REGISTRYTYPE_ETCD):
		registerType = opts2.REGISTRYTYPE_ETCD
		srv, err = etcd.NewEtcdRegister(options)
	default:
		return nil, errors.New(string("registry type: " + options.RegisterType + " can not be supported, you can choose: etcd"))
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
