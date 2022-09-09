package discover

import (
	"context"
	"github.com/go-kit/kit/sd/etcdv3"
	kitlog "github.com/go-kit/log"
	"time"
)

type discover struct {
	client etcdv3.Client
}

func (d *discover) HttpRegister(serviceName string, value string) *etcdv3.Registrar {
	register := etcdv3.NewRegistrar(d.client, etcdv3.Service{
		Key:   "/services/" + serviceName + "/http/" + value,
		Value: value,
	}, kitlog.NewNopLogger())

	register.Register()
	return register
}

func (d *discover) RPCRegister(serviceName string, value string) *etcdv3.Registrar {
	register := etcdv3.NewRegistrar(d.client, etcdv3.Service{
		Key:   "/services/" + serviceName + "/rpc/" + value,
		Value: value,
	}, kitlog.NewNopLogger())

	register.Register()
	return register
}

func NewDiscover(ctx context.Context, addr string) *discover {
	options := etcdv3.ClientOptions{
		DialTimeout:   time.Second * 3,
		DialKeepAlive: time.Second * 3,
	}

	client, err := etcdv3.NewClient(ctx, []string{addr}, options)
	if err != nil {
		panic(err)
	}

	return &discover{
		client: client,
	}
}
