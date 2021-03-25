package register

import (
	log "github.com/mhchlib/logger"
)

// Options ...
type Options struct {
	RegisterStr    string
	RegisterType   string
	Address        string
	Namespace      string
	ServerInstance string
	Metadata       map[string]interface{}
	Logger         log.Logger
}

// Option ...
type Option func(*Options)

// Namespace ...
func Namespace(namespace string) Option {
	return func(options *Options) {
		options.Namespace = namespace
	}
}

// Instance ...
func Instance(address string) Option {
	return func(options *Options) {
		options.ServerInstance = address
	}
}

// MetadataMap ...
func MetadataMap(metadata map[string]interface{}) Option {
	return func(options *Options) {
		options.Metadata = metadata
	}
}

// Metadata ...
func Metadata(key string, value interface{}) Option {
	return func(options *Options) {
		if options.Metadata == nil {
			options.Metadata = make(map[string]interface{})
		}
		options.Metadata[key] = value
	}
}

// ResgisterAddress ...
func ResgisterAddress(address string) Option {
	return func(options *Options) {
		options.RegisterStr = address
	}
}
