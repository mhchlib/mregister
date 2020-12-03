module github.com/mhchlib/register

go 1.14

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/go-kit/kit v0.10.0
	github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc
	github.com/pborman/uuid v1.2.0
)

replace github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc => ../logger
