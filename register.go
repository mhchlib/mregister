package register

import (
	"errors"
	registerEtcd "github.com/mhchlib/register/etcd"
	"github.com/mhchlib/register/reg"
)

func InitRegister(registryType string, opts ...reg.Option) (reg.Register, error) {
	if registryType == "etcd" {
		reg, err := registerEtcd.NewEtcdRegister(opts)
		return reg, err
	}
	return nil, errors.New("registry type: " + registryType + " can not be supported, you can choose: etcd")
}
