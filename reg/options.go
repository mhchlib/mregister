package reg

import log "github.com/mhchlib/logger"

type Options struct {
	Address        []string
	NameSpace      string
	ServerInstance string
	Metadata       map[string]interface{}
	logger         log.Logger
}

type Option func(*Options)
