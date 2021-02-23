package register

import (
	"errors"
	"github.com/mhchlib/register/regutils"
	"strings"
)

func InitRegister(opts ...Option) (Register, error) {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	if options.registerStr != "" {
		t, address, err := parseAddressStr(options.registerStr)
		if err != nil {
			return nil, err
		}
		options.registerType = RegistryType(t)
		options.address = strings.Split(address, ",")
	}

	if options.serverInstance == "" {
		ip, err := regutils.GetClientIp()
		if err != nil {
			return nil, err
		}
		options.serverInstance = ip + DEFAULT_PORT
	}
	if options.registerType == REGISTRYTYPE_ETCD {
		regClient, err := NewEtcdRegister(options)
		return regClient, err
	}
	return nil, errors.New(string("registry type: " + options.registerType + " can not be supported, you can choose: etcd"))
}

const ConfigSeparateSymbol = "://"

func parseAddressStr(str string) (string, string, error) {
	splits := strings.Split(str, ConfigSeparateSymbol)
	if len(splits) != 2 {
		return "", "", errors.New(str + " is invalid Address")
	}
	return splits[0], splits[1], nil
}
