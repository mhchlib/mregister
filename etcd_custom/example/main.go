package main

import (
	"github.com/mhchlib/register/bak"
	"github.com/mhchlib/register/etcd_custom"
)

func main() {

	etcdRegister := etcd_custom.EtcdRegister{}
	etcdRegister.Init(func(options *bak.Options) {
		options.NameSpaces = "!!!11"
		options.Addrs = []string{"etcd_custom.u.hcyang.top:31770"}
	})

	etcdRegister.Register("....")

	select {}

}
