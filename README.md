# register
feature register for mconfig

> support register center

* etcd:etcd://etcd.u.hcyang.top:31770
* select: select://xxxx-service::10.12.20.1:8080,10.12.20.2:8080,10.12.20.3:8080;yyyy-service::10.12.20.1:8080,10.12.20.2:8080,10.12.20.3:8080;

> start
```go
regClient, err := register.InitRegister(
    register.Namespace("test_register"),
    register.ResgisterAddress("etcd://etcd.u.hcyang.top:31770"),
    register.Metadata("key", "value"),
)
if err != nil {
    log.Fatal(err)
}
regClient.RegisterService("xxxx-service", map[string]interface{}{"key": "value"})
```