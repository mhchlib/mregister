package etcd

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mhchlib/go-kit/endpoint"
	"github.com/mhchlib/go-kit/sd"
	"github.com/mhchlib/go-kit/sd/etcdv3"
	"github.com/mhchlib/go-kit/sd/lb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/register/common"
	"github.com/mhchlib/register/registerOpts"
	"github.com/mhchlib/register/robin"
	"io"
	"strings"
	"sync"
	"time"
)

// EtcdRegister ...
type EtcdRegister struct {
	Opts     *registerOpts.Options
	services *EtcdServiceMap
	log.Logger
}

// EtcdServiceMap ...
type EtcdServiceMap struct {
	data map[string]*EtcdService
	sync.RWMutex
}

// EtcdService ...
type EtcdService struct {
	balancer   lb.Balancer
	endpointer sd.Endpointer
	client     etcdv3.Client
	key        string
}

//NewEtcdRegister ...
func NewEtcdRegister(options *registerOpts.Options) (*EtcdRegister, error) {
	reg := &EtcdRegister{}
	reg.Opts = options
	if reg.Logger == nil {
		reg.Logger = log.NewLogger()
	}
	reg.services = &EtcdServiceMap{
		data: map[string]*EtcdService{},
	}
	return reg, nil
}

func newEtcdClient(er *EtcdRegister) etcdv3.Client {
	ctx := context.Background()
	option := etcdv3.ClientOptions{DialTimeout: time.Second * 3, DialKeepAlive: time.Second * 3}
	client, err := etcdv3.NewClient(ctx, strings.Split(er.Opts.Address, ","), option)
	if err != nil {
		panic(err)
	}
	return client
}

// UnRegisterService ...
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
	er.Logger.Info("success unregister service ", serviceName)
	return err
}

// UnRegisterServiceAll ...
func (er *EtcdRegister) UnRegisterServiceAll() error {

	return nil
}

// RegisterService ...
func (er *EtcdRegister) RegisterService(serviceName string, metadata map[string]interface{}) (func(), error) {
	globalMetadata := er.Opts.Metadata
	for key, value := range metadata {
		globalMetadata[key] = value
	}

	serviceVal := &common.ServiceVal{
		Address:  er.Opts.ServerInstance,
		Metadata: globalMetadata,
	}
	serviceValStr, err := json.Marshal(serviceVal)
	if err != nil {
		return nil, err
	}
	client := newEtcdClient(er)
	key := getEtcdKey(er.Opts.Namespace, serviceName, er.Opts.ServerInstance)
	registrar := etcdv3.NewRegistrar(client, etcdv3.Service{Key: key, Value: string(serviceValStr)}, er.Logger)
	registrar.Register()
	services := er.services
	if services == nil {
		services = &EtcdServiceMap{
			data: map[string]*EtcdService{},
		}
		er.services = services
	}
	services.Lock()
	services.data[serviceName] = &EtcdService{
		client: client,
		key:    key,
	}
	services.Unlock()
	er.Logger.Info("register service", serviceName, "success")
	return func() {
		err := er.UnRegisterService(serviceName)
		if err != nil {
			log.Error(err)
		}
	}, nil
}

// GetService ...
func (er *EtcdRegister) GetService(serviceName string) (*common.ServiceVal, error) {
	prefix := getEtcdKey(er.Opts.Namespace, serviceName, "")
	services := er.services
	exist := false
	var bl lb.Balancer
	if services == nil {
		services = &EtcdServiceMap{}
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
			return func(ctx context.Context, request interface{}) (Response interface{}, err error) {
				return instance, nil
			}, nil, nil
		}, er.Logger)
		balancer := lb.NewRoundRobin(endpointer)
		er.services.Lock()
		er.services.data[serviceName] = &EtcdService{
			balancer:   balancer,
			client:     client,
			endpointer: endpointer,
		}
		er.services.Unlock()
		bl = balancer
	}
	if bl != nil {
		reqEndPoint, err := bl.Endpoint()
		if err != nil {
			log.Error(err)
			return nil, err
		}
		ctx := context.Background()
		data, err := reqEndPoint(ctx, nil)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		serviceVal := &common.ServiceVal{}
		err = json.Unmarshal([]byte(data.(string)), &serviceVal)
		if err != nil {
			return nil, err
		}
		return serviceVal, nil
	} else {
		return nil, errors.New("no registration information was found")
	}
}

// ListAllServices ...
func (er *EtcdRegister) ListAllServices(serviceName string) ([]*common.ServiceVal, error) {
	var endpointer sd.Endpointer
	er.services.RLock()
	service, ok := er.services.data[serviceName]
	if ok {
		endpointer = service.endpointer
	}
	er.services.RUnlock()
	if endpointer == nil {
		prefix := getEtcdKey(er.Opts.Namespace, serviceName, "")
		client := newEtcdClient(er)
		instancer, err := etcdv3.NewInstancer(client, prefix, er.Logger)
		defer func() {
			instancer.Stop()
		}()
		if err != nil {
			panic(err)
		}
		endpointer = sd.NewEndpointer(instancer, func(instance string) (endpoint.Endpoint, io.Closer, error) {
			return func(ctx context.Context, request interface{}) (Response interface{}, err error) {
				return instance, nil
			}, nil, nil
		}, er.Logger)
	}
	balancer := robin.NewListRobin(endpointer)
	r := balancer.(*robin.ListRobin)
	reqEndPoints, err := r.Endpoints()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ctx := context.Background()
	serviceVals := make([]*common.ServiceVal, 0)
	for _, reqEndPoint := range reqEndPoints {
		data, err := reqEndPoint(ctx, nil)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		serviceVal := &common.ServiceVal{}
		err = json.Unmarshal([]byte(data.(string)), &serviceVal)
		if err != nil {
			return nil, err
		}
		serviceVals = append(serviceVals, serviceVal)
	}
	return serviceVals, nil
}

func getEtcdKey(namespace string, serviceName string, salt string) string {
	return "/" + namespace + "/" + serviceName + "/" + salt
}
