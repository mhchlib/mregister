package register

import (
	"errors"
	registerEtcd "github.com/mhchlib/register/etcd"
	"github.com/mhchlib/register/mregister"
)

func InitRegister(registryType string, opts ...mregister.Option) (mregister.Register, error) {
	if registryType == "etcd" {
		reg := registerEtcd.NewEtcdRegister(opts)
		return reg, nil
	}
	return nil, errors.New("registry type: " + registryType + " can not be supported, you can choose: etcd")
}
