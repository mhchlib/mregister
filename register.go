package register

type Register interface {
	Init(opts ...Option)
	RegisterService(serviceName string)
	UnRegisterService(serviceName string)
	GetService(serviceName string) (string, error)
}
