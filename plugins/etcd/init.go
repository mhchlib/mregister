package etcd

import (
	"github.com/mhchlib/mregister/plugin"
)

func init() {
	plugin.RegisterRegisterPlugin("etcd", newEtcdRegister)
}
