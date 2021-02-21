module github.com/mhchlib/register

go 1.14

require (
	github.com/go-kit/kit v0.10.1
	github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc
	github.com/pborman/uuid v1.2.0
)

replace github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc => ../logger

replace github.com/go-kit/kit v0.10.1 => ../kit
