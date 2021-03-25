package memory

import (
	"encoding/json"
	"github.com/mhchlib/go-kit/sd"
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/register/common"
)

type Instancer struct {
	Services []common.ServiceVal
}

func (i Instancer) Register(events chan<- sd.Event) {
	services := make([]string, 0)
	for _, service := range i.Services {
		data, _ := json.Marshal(service)
		services = append(services, string(data))

	}
	events <- sd.Event{
		Instances: services,
		Err:       nil,
	}
}

func (i Instancer) Deregister(events chan<- sd.Event) {
	log.Info(" Instancer Deregister")
}

func (i Instancer) Stop() {
	log.Info(" Instancer Stop")
}
