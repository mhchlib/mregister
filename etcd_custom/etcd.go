package etcd_custom

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/register/bak"
	"time"
)

type EtcdRegister struct {
	opts bak.Options
	cli  *clientv3.Client
	id   string
}

func (e *EtcdRegister) Init(options ...bak.Option) {
	opts := &bak.Options{}
	for _, o := range options {
		o(opts)
	}
	e.opts = *opts

	if e.opts.TTL == 0 {
		e.opts.TTL = 20
	}

	//link to etcd_custom
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   e.opts.Addrs,
		DialTimeout: e.opts.Timeout,
	})
	if err != nil {
		log.Fatal(err)
	}

	e.cli = cli
	e.id = "....."

	return
}

func (e *EtcdRegister) Options() bak.Options {
	return e.opts
}

func (e *EtcdRegister) Register(service string) error {
	namespaces := e.opts.NameSpaces
	log.Info(e.cli)
	log.Info(e.opts.NameSpaces)
	kv := clientv3.NewKV(e.cli)
	key := "/register/services/" + namespaces + "/" + e.id
	ctx := context.Background()
	lease := clientv3.NewLease(e.cli)
	leaseRes, err := clientv3.NewLease(e.cli).Grant(ctx, e.opts.TTL)
	if err != nil {
		return err
	}
	var server string = "1111"
	_, err = kv.Put(context.Background(), key, server, clientv3.WithLease(leaseRes.ID))
	if err != nil {
		return err
	}
	keepaliveRes, err := lease.KeepAlive(context.TODO(), leaseRes.ID)
	if err != nil {
		return err
	}
	go lisKeepAlive(keepaliveRes)
	return err
}

func lisKeepAlive(keepaliveRes <-chan *clientv3.LeaseKeepAliveResponse) {
	for {
		select {
		case ret := <-keepaliveRes:
			if ret != nil {
				log.Info("续租成功", time.Now())
			}
		}
	}
}

func (e *EtcdRegister) Deregister() error {
	panic("implement me")
}

func (e *EtcdRegister) GetService(s string) ([]*bak.Service, error) {
	panic("implement me")
}

func (e *EtcdRegister) GetServiceAddress(s string) ([]*bak.Service, error) {
	panic("implement me")
}

func (e *EtcdRegister) ListServices() ([]*bak.Service, error) {
	panic("implement me")
}

func (e *EtcdRegister) String() string {
	panic("implement me")
}
