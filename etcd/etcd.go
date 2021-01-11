package etcd

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/etcdv3"
	"github.com/go-kit/kit/sd/lb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/register/mregister"
	"github.com/pborman/uuid"
	"io"
	"sync"
	"time"
)

type EtcdRegister struct {
	Opts     *mregister.Options
	services *ServiceMap
	log.Logger
}

type ServiceMap struct {
	data map[string]*Service
	sync.RWMutex
}

type Service struct {
	balancer lb.Balancer
	client   etcdv3.Client
	key      string
}

func NewEtcdRegister(opts []mregister.Option) mregister.Register {
	options := &mregister.Options{}
	for _, o := range opts {
		o(options)
	}
	reg := &EtcdRegister{}
	reg.Opts = options
	if reg.Logger == nil {
		reg.Logger = log.NewLogger()
	}
	reg.services = &ServiceMap{
		data: map[string]*Service{},
	}
	return reg
}

func newEtcdClient(er *EtcdRegister) etcdv3.Client {
	ctx := context.Background()
	option := etcdv3.ClientOptions{DialTimeout: time.Second * 3, DialKeepAlive: time.Second * 3}
	client, err := etcdv3.NewClient(ctx, er.Opts.Address, option)
	if err != nil {
		panic(err)
	}
	return client
}

func (er *EtcdRegister) UnRegisterService(serviceName string) error {
	services := er.services
	if services == nil || services.data == nil {
		return errors.New("service not found")
	}
	services.Lock()
	defer services.Unlock()
	service, ok := services.data[serviceName]
	if !ok {
		return errors.New("service not found")
	}
	err := service.client.Deregister(etcdv3.Service{
		Key: service.key,
	})
	er.Logger.Info(service.key)
	if err == nil {
		service.client = nil
		service.balancer = nil
		delete(services.data, serviceName)
	}
	er.Logger.Info("SUCCES UnRegister Service ", serviceName)
	return err
}

func (er *EtcdRegister) UnRegisterServiceAll() error {

	return nil
}

func (er *EtcdRegister) RegisterService(serviceName string) error {
	client := newEtcdClient(er)
	key := getEtcdKey(er.Opts.NameSpace, serviceName, uuid.New())
	registrar := etcdv3.NewRegistrar(client, etcdv3.Service{Key: key, Value: er.Opts.ServerInstance}, er.Logger)
	registrar.Register()
	services := er.services
	if services == nil {
		services = &ServiceMap{
			data: map[string]*Service{},
		}
		er.services = services
	}
	services.Lock()
	services.data[serviceName] = &Service{
		client: client,
		key:    key,
	}
	services.Unlock()
	er.Logger.Info("register service ", serviceName, " success")
	return nil
}

func (er *EtcdRegister) GetService(serviceName string) (string, error) {
	prefix := getEtcdKey(er.Opts.NameSpace, serviceName, "")
	services := er.services
	exist := false
	var bl lb.Balancer

	if services == nil {
		services = &ServiceMap{}
		er.services = services
		exist = false
	} else {
		services.RLock()
		service, ok := services.data[serviceName]
		if !ok {
			exist = false
		} else {
			bl = service.balancer
		}
		services.RUnlock()
	}

	if exist == false {
		client := newEtcdClient(er)
		instancer, err := etcdv3.NewInstancer(client, prefix, er.Logger)
		if err != nil {
			panic(err)
		}
		endpointer := sd.NewEndpointer(instancer, func(instance string) (endpoint.Endpoint, io.Closer, error) {
			return func(ctx context.Context, request interface{}) (response interface{}, err error) {
				return instance, nil
			}, nil, nil
		}, er.Logger)
		balancer := lb.NewRoundRobin(endpointer)
		er.services.Lock()
		er.services.data[serviceName] = &Service{
			balancer: balancer,
			client:   client,
		}
		er.services.Unlock()
		bl = balancer
	}
	if bl != nil {
		reqEndPoint, err := bl.Endpoint()
		if err != nil {
			log.Error(err)
			return "", err
		}
		ctx := context.Background()
		data, err := reqEndPoint(ctx, nil)
		if err != nil {
			log.Error(err)
			return "", err
		}
		return data.(string), nil
	} else {
		return "", nil
	}
}

func getEtcdKey(namespace string, serviceName string, salt string) string {
	return "/" + namespace + "/" + serviceName + "/" + salt
}
