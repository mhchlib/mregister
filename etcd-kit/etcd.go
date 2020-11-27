package etcd_kit

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/register"
	"github.com/pborman/uuid"
	"io"
	"sync"
	"time"
)

type EtcdRegister struct {
	opts     *register.Options
	client   etcdv3.Client
	balancer *BalancerMap
}

func (er *EtcdRegister) UnRegisterService(serviceName string) {
	panic("implement me")
}

type BalancerMap struct {
	data map[string]*lb.Balancer
	sync.RWMutex
}

func (er *EtcdRegister) RegisterService(serviceName string) {
	client := er.client
	key := getEtcdKey(er.opts.NameSpace, serviceName, uuid.New())
	registrar := etcdv3.NewRegistrar(client, etcdv3.Service{Key: key, Value: er.opts.ServerInstance}, &RegisterLog{})
	registrar.Register()
}

func (er *EtcdRegister) Init(opts ...register.Option) {
	options := &register.Options{}
	for _, o := range opts {
		o(options)
	}
	er.opts = options
	ctx := context.Background()
	option := etcdv3.ClientOptions{DialTimeout: time.Second * 3, DialKeepAlive: time.Second * 3}
	client, err := etcdv3.NewClient(ctx, er.opts.Address, option)
	if err != nil {
		panic(err)
	}
	er.client = client
	er.balancer = &BalancerMap{
		data:    make(map[string]*lb.Balancer),
		RWMutex: sync.RWMutex{},
	}
}

func (er *EtcdRegister) GetService(serviceName string) (string, error) {
	prefix := getEtcdKey(er.opts.NameSpace, serviceName, "")
	er.balancer.RLock()
	ctx := context.Background()
	bl, ok := er.balancer.data[prefix]
	er.balancer.RUnlock()
	if ok == false {
		logger := &RegisterLog{}
		instancer, err := etcdv3.NewInstancer(er.client, prefix, logger)
		if err != nil {
			panic(err)
		}
		endpointer := sd.NewEndpointer(instancer, func(instance string) (endpoint.Endpoint, io.Closer, error) {
			return func(ctx context.Context, request interface{}) (response interface{}, err error) {
				return instance, nil
			}, nil, nil
		}, logger)
		balancer := lb.NewRoundRobin(endpointer)
		er.balancer.Lock()
		er.balancer.data[prefix] = &balancer
		er.balancer.Unlock()
		bl = &balancer
	}
	reqEndPoint, err := (*bl).Endpoint()
	if err != nil {
		log.Error(err)
		return "", err
	}
	data, err := reqEndPoint(ctx, nil)
	if err != nil {
		log.Error(err)
		return "", err
	}
	return data.(string), nil
}

func getEtcdKey(namespace string, serviceName string, salt string) string {
	return "/" + namespace + "/" + serviceName + "/" + salt
}
