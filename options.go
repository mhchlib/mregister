package register

import (
	log "github.com/mhchlib/logger"
)

type Options struct {
	registerStr    string
	registerType   RegistryType
	address        []string
	namespace      string
	serverInstance string
	metadata       map[string]interface{}
	logger         log.Logger
}

type Option func(*Options)

type RegistryType string

const (
	REGISTRYTYPE_ETCD RegistryType = "etcd"
)

func RegisterStr(registerStr string) Option {
	return func(options *Options) {
		options.registerStr = registerStr
	}
}

func SelectEtcdRegister() Option {
	return func(options *Options) {
		options.registerType = REGISTRYTYPE_ETCD
	}
}

func Namespace(namespace string) Option {
	return func(options *Options) {
		options.namespace = namespace
	}
}

func Instance(address string) Option {
	return func(options *Options) {
		options.serverInstance = address
	}
}

func MetadataMap(metadata map[string]interface{}) Option {
	return func(options *Options) {
		options.metadata = metadata
	}
}

func Metadata(key string, value interface{}) Option {
	return func(options *Options) {
		if options.metadata == nil {
			options.metadata = make(map[string]interface{})
		}
		options.metadata[key] = value
	}
}

func ResgisterAddress(address []string) Option {
	return func(options *Options) {
		options.address = address
	}
}
