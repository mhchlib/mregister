package etcd

import (
	"github.com/mhchlib/register/plugin"
)

func init() {
	plugin.RegisterRegisterPlugin("etcd", newEtcdRegister)
}
