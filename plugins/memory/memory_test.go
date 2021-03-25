package memory

import (
	"github.com/mhchlib/mregister/register"
	"testing"
)

func TestNewMemoryRegister(t *testing.T) {
	client, err := newMemoryRegister(&register.Options{
		RegisterStr:    "",
		RegisterType:   "memory",
		Address:        "xxxx-service::10.12.20.1:8080(xx:111#dada:dsdas#dsada:dasdas),10.12.20.3:8080;yyyy-service::10.12.20.1:8080,10.12.20.2:8080,10.12.20.3:8080",
		Namespace:      "xxxx",
		ServerInstance: "xxxxx",
		Metadata:       nil,
		Logger:         nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(client.GetService("xxxx-service"))
	services, err := client.ListAllServices("xxxx-service")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("print services ")
	for _, service := range services {
		t.Log(service)
	}
}

func Test_checkAddressLegal(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1", args: args{address: "10.12.20.1:8080"}, wantErr: true,
		}, {
			name: "2", args: args{address: "127.0.0.1:880"}, wantErr: true,
		}, {
			name: "3", args: args{address: "127.0.0.1:880"}, wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkAddressLegal(tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("checkAddressLegal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
