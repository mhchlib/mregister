package register

import (
	"context"
	"crypto/tls"
	"time"
)

type Register_Model int

const (
	Register_Model_Reg = iota
	Register_Model_Disc
	Register_Model_RegDisc
)

type Options struct {
	NameSpaces string
	Addrs      []string
	Timeout    time.Duration
	Secure     bool
	TLSConfig  *tls.Config
	Model      Register_Model
	Context    context.Context
	TTL        int64
}
