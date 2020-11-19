package main

import (
	"github.com/mhchlib/register"
	"github.com/mhchlib/register/etcd"
)

func main() {

	etcdRegister := etcd.EtcdRegister{}
	etcdRegister.Init(func(options *register.Options) {
		options.NameSpaces = "!!!11"
		options.Addrs = []string{"etcd.u.hcyang.top:31770"}
	})

	etcdRegister.Register("....")

	select {}

}
