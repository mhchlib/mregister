package register

import (
	log "github.com/mhchlib/logger"
)

// Options ...
type Options struct {
	registerStr    string
	registerType   string
	address        []string
	namespace      string
	serverInstance string
	metadata       map[string]interface{}
	logger         log.Logger
}

// Option ...
type Option func(*Options)

// RegistryType ...
type RegistryType string

const (
	// REGISTRYTYPE_ETCD ...
	REGISTRYTYPE_ETCD RegistryType = "etcd"
)

// Namespace ...
func Namespace(namespace string) Option {
	return func(options *Options) {
		options.namespace = namespace
	}
}

// Instance ...
func Instance(address string) Option {
	return func(options *Options) {
		options.serverInstance = address
	}
}

// MetadataMap ...
func MetadataMap(metadata map[string]interface{}) Option {
	return func(options *Options) {
		options.metadata = metadata
	}
}

// Metadata ...
func Metadata(key string, value interface{}) Option {
	return func(options *Options) {
		if options.metadata == nil {
			options.metadata = make(map[string]interface{})
		}
		options.metadata[key] = value
	}
}

// ResgisterAddress ...
func ResgisterAddress(address string) Option {
	return func(options *Options) {
		options.registerStr = address
	}
}
