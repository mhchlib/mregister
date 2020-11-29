package mregister

import log "github.com/mhchlib/logger"

type Options struct {
	Address        []string
	NameSpace      string
	ServerInstance string
	logger         log.Logger
}

type Option func(*Options)
