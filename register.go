package register

import (
	"errors"
	"github.com/mhchlib/register/common"
	registerEtcd "github.com/mhchlib/register/etcd"
	"github.com/mhchlib/register/reg"
	"strings"
)

const DEFAULT_PORT = ":8080"

func InitRegister(opts ...reg.Option) (reg.Register, error) {
	options := &reg.Options{}
	for _, o := range opts {
		o(options)
	}
	if options.RegisterStr != "" {
		t, address, err := parseAddressStr(options.RegisterStr)
		if err != nil {
			return nil, err
		}
		options.RegisterType = reg.RegistryType(t)
		options.Address = strings.Split(address, ",")
	}

	if options.ServerInstance == "" {
		ip, err := common.GetClientIp()
		if err != nil {
			return nil, err
		}
		options.ServerInstance = ip + DEFAULT_PORT
	}
	if options.RegisterType == reg.RegistryType_Etcd {
		regClient, err := registerEtcd.NewEtcdRegister(options)
		return regClient, err
	}
	return nil, errors.New(string("registry type: " + options.RegisterType + " can not be supported, you can choose: etcd"))
}

const ConfigSeparateSymbol = "://"

func parseAddressStr(str string) (string, string, error) {
	splits := strings.Split(str, ConfigSeparateSymbol)
	if len(splits) != 2 {
		return "", "", errors.New(str + " is invalid Address")
	}
	return splits[0], splits[1], nil
}
