package plugin

import (
	log "github.com/mhchlib/logger"
	"github.com/mhchlib/mregister/register"
)

// StorePlugin ...
type RegisterPlugin struct {
	Name string
	New  func(options *register.Options) (register.Register, error)
	//...
}

var RegisterPluginMap map[string]*RegisterPlugin

var RegisterPluginNames []string

// NewStorePlugin ...
func NewRegisterPlugin(name string, new func(options *register.Options) (register.Register, error)) *RegisterPlugin {
	return &RegisterPlugin{Name: name, New: new}
}

// RegisterStorePlugin ...
func RegisterRegisterPlugin(name string, new func(options *register.Options) (register.Register, error)) error {
	if RegisterPluginMap == nil {
		RegisterPluginMap = make(map[string]*RegisterPlugin)
	}
	if RegisterPluginNames == nil {
		RegisterPluginNames = []string{}
	}

	if _, ok := RegisterPluginMap[name]; ok {
		log.Fatal("repeated register same name register plugin ...")
	}
	RegisterPluginMap[name] = NewRegisterPlugin(name, new)
	RegisterPluginNames = append(RegisterPluginNames, name)
	return nil
}
