package memory

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mhchlib/go-kit/endpoint"
	"github.com/mhchlib/go-kit/sd"
	"github.com/mhchlib/go-kit/sd/lb"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mregister/register"
	"github.com/mhchlib/mregister/robin"
	"io"
	net "net"
	"strconv"
	"strings"
	"sync"
)

// MemoryRegister ...
type MemoryRegister struct {
	Opts     *register.Options
	services *MemoryServiceMap
	log.Logger
}

// MemoryServiceMap ...
type MemoryServiceMap struct {
	data map[string]*MemoryService
	sync.RWMutex
}

// MemoryService ...
type MemoryService struct {
	endpointer sd.Endpointer
	balancer   lb.Balancer
	services   []register.ServiceVal
	key        string
}

func newMemoryRegister(options *register.Options) (register.Register, error) {
	reg := &MemoryRegister{}
	reg.Opts = options
	if reg.Logger == nil {
		reg.Logger = log.NewLogger()
	}
	reg.services = &MemoryServiceMap{
		data: map[string]*MemoryService{},
	}
	err := parseMemoryAddressStr(reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func parseMemoryAddressStr(reg *MemoryRegister) error {
	//xxxx-service::10.12.20.1:8080(xx:111,dada:dsdas,dsada,dasdas),10.12.20.3:8080;yyyy-service::10.12.20.1:8080,10.12.20.2:8080,10.12.20.3:8080;
	addressStr := reg.Opts.Address
	serviceArray := strings.Split(addressStr, ";")
	for _, serviceItem := range serviceArray {
		services := make([]register.ServiceVal, 0)
		serviceItemSplits := strings.Split(serviceItem, "::")
		serviceName := "default"
		serviceEntitys := serviceItem
		if len(serviceItemSplits) == 2 {
			serviceName = serviceItemSplits[0]
			serviceEntitys = serviceItemSplits[1]
		}
		if len(serviceItemSplits) > 2 {
			return errors.New("Register Memory mode address string format is error")
		}
		serviceEntityArray := strings.Split(serviceEntitys, ",")
		for _, serviceEntity := range serviceEntityArray {
			service := &register.ServiceVal{}
			point := strings.Index(serviceEntity, "(")
			serviceAddress := ""
			if point == -1 {
				serviceAddress = serviceEntity
			} else {
				b := strings.HasSuffix(serviceEntity, ")")
				log.Info(serviceEntity)
				if !b {
					return errors.New("Register Memory mode address string format is error")
				}
				metaDataStr := serviceEntity[point+1 : len(serviceEntity)-1]
				metaDataArr := strings.Split(metaDataStr, "#")
				metaData := make(map[string]interface{})
				for _, metaDataItem := range metaDataArr {
					metaDataItemSplitArr := strings.Split(metaDataItem, ":")
					if len(metaDataItemSplitArr) == 2 {
						metaData[metaDataItemSplitArr[0]] = metaDataItemSplitArr[1]
					} else {
						return errors.New("Register Memory mode address string format is error")
					}
				}
				service.Metadata = metaData
				serviceAddress = serviceEntity[0:point]
			}
			err := checkAddressLegal(serviceAddress)
			if err != nil {
				return err
			}
			service.Address = serviceAddress
			services = append(services, *service)
		}
		reg.services.Lock()
		instancer := &Instancer{
			Services: services,
		}
		endpointer := sd.NewEndpointer(instancer, func(instance string) (endpoint.Endpoint, io.Closer, error) {
			return func(ctx context.Context, request interface{}) (Response interface{}, err error) {
				return instance, nil
			}, nil, nil
		}, reg.Logger)
		balancer := lb.NewRoundRobin(endpointer)
		reg.services.data[serviceName] = &MemoryService{
			balancer:   balancer,
			services:   services,
			endpointer: endpointer,
			key:        serviceName,
		}
		reg.services.Unlock()
	}
	return nil
}

func checkAddressLegal(address string) error {
	splits := strings.Split(address, ":")
	if len(splits) != 2 {
		return fmt.Errorf("Register memory mode address %s is illegal，should be xx.xx.xx.xx:xxxx", address)
	}
	ipStr := splits[0]
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return fmt.Errorf("Register memory mode address %s is illegal，should be xx.xx.xx.xx:xxxx", address)
	}
	port := splits[1]
	_, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("Register memory mode address %s is illegal，should be xx.xx.xx.xx:xxxx", address)
	}
	return nil
}

// UnRegisterService ...
func (er *MemoryRegister) UnRegisterService(serviceName string) error {
	return nil
}

// UnRegisterServiceAll ...
func (er *MemoryRegister) UnRegisterServiceAll() error {
	return nil
}

// RegisterService ...
func (er *MemoryRegister) RegisterService(serviceName string, metadata map[string]interface{}) (func(), error) {
	return func() {}, nil
}

// GetService ...
func (er *MemoryRegister) GetService(serviceName string) (*register.ServiceVal, error) {
	//prefix := getMemoryKey(er.Opts.namespace, serviceName, "")
	services := er.services
	var bl lb.Balancer
	if services == nil {
		services = &MemoryServiceMap{}
		er.services = services
		return nil, register.SERVICES_NOT_FOUND
	} else {
		services.RLock()
		service, ok := services.data[serviceName]
		if !ok {
			services.RUnlock()
			return nil, register.SERVICES_NOT_FOUND
		} else {
			bl = service.balancer
		}
		services.RUnlock()
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
		serviceVal := &register.ServiceVal{}
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
func (er *MemoryRegister) ListAllServices(serviceName string) ([]*register.ServiceVal, error) {
	er.services.RLock()
	service, ok := er.services.data[serviceName]
	if !ok {
		return nil, register.SERVICES_NOT_FOUND
	}
	er.services.RUnlock()

	balancer := robin.NewListRobin(service.endpointer)
	r := balancer.(*robin.ListRobin)
	reqEndPoints, err := r.Endpoints()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ctx := context.Background()
	serviceVals := make([]*register.ServiceVal, 0)
	for _, reqEndPoint := range reqEndPoints {
		data, err := reqEndPoint(ctx, nil)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		serviceVal := &register.ServiceVal{}
		err = json.Unmarshal([]byte(data.(string)), &serviceVal)
		if err != nil {
			return nil, err
		}
		serviceVals = append(serviceVals, serviceVal)
	}
	return serviceVals, nil
}

func getMemoryKey(namespace string, serviceName string, salt string) string {
	return "/" + namespace + "/" + serviceName + "/" + salt
}
