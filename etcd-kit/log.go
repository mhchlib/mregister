package etcd_kit

import log "github.com/mhchlib/logger"

type RegisterLog struct {
}

func (A RegisterLog) Log(vals ...interface{}) error {
	log.Info(vals)
	return nil
}
