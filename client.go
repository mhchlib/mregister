package register

import (
	"errors"
	log "github.com/mhchlib/logger"
	"strings"
)

func InitRegister(opts ...Option) (*RegisterClient, error) {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	if options.registerStr != "" {
		t, address, err := parseAddressStr(options.registerStr)
		if err != nil {
			return nil, err
		}
		options.registerType = t
		options.address = strings.Split(address, ",")
	}

	if options.serverInstance == "" {
		return nil, errors.New("server instance can not be empty")
	}
	var srv Register
	var err error
	var registerType RegistryType
	switch options.registerType {
	case string(REGISTRYTYPE_ETCD):
		registerType = REGISTRYTYPE_ETCD
		srv, err = NewEtcdRegister(options)
	default:
		return nil, errors.New(string("registry type: " + options.registerType + " can not be supported, you can choose: etcd"))
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

const ConfigSeparateSymbol = "://"

func parseAddressStr(str string) (string, string, error) {
	splits := strings.Split(str, ConfigSeparateSymbol)
	if len(splits) != 2 {
		return "", "", errors.New(str + " is invalid Address")
	}
	return splits[0], splits[1], nil
}
