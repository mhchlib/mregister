package etcd_kit

import (
	"github.com/mhchlib/register"
	"testing"
)

func TestA01(t *testing.T) {
	reg := &EtcdRegister{}
	reg.Init(func(options *register.Options) {
		options.Address = []string{"etcd_custom.u.hcyang.top:31770"}
		options.NameSpace = "/com.github.mhchlib.pkg"
		options.ServerInstance = "127.0.0.1:8080"
	})
	reg.RegisterService("Test")

}

func TestA02(t *testing.T) {
	reg := &EtcdRegister{}
	reg.Init(func(options *register.Options) {
		options.Address = []string{"etcd_custom.u.hcyang.top:31770"}
		options.NameSpace = "com.github.mhchlib.pkg"
		options.ServerInstance = "127.0.0.1:8080"
	})
	reg.RegisterService("Test")
	service, err := reg.GetService("Test")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(service)
}
