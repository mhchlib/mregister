package mregister

import log "github.com/mhchlib/logger"

type Register interface {
	RegisterService(serviceName string) error
	UnRegisterService(serviceName string) error
	UnRegisterServiceAll() error
	GetService(serviceName string) (string, error)
	log.Logger
}
