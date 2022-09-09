package discover

import (
	"context"
	etcdapiv3mvccpb "go.etcd.io/etcd/api/v3/mvccpb"
	etcdclientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

type Instancer struct {
	mu        sync.RWMutex
	client    *etcdclientv3.Client
	prefix    string
	addresses map[string]string
}

func (ins *Instancer) watch() {
	rch := ins.client.Watch(context.Background(), ins.prefix, etcdclientv3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			ins.mu.Lock()
			if ev.Type == etcdapiv3mvccpb.DELETE {
				delete(ins.addresses, string(ev.Kv.Key))
			} else if ev.Type == etcdapiv3mvccpb.PUT {
				ins.addresses[string(ev.Kv.Key)] = string(ev.Kv.Value)
			}
			ins.mu.Unlock()
		}
	}
}

func (ins *Instancer) Addresses() []string {
	ins.mu.RLock()
	defer ins.mu.RUnlock()

	var a []string
	for _, v := range ins.addresses {
		a = append(a, v)
	}

	return a
}

func NewInstacer(client *etcdclientv3.Client, prefix string) (*Instancer, error) {
	ins := &Instancer{
		client:    client,
		prefix:    prefix,
		addresses: map[string]string{},
	}

	resp, err := client.Get(context.Background(), prefix, etcdclientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	for _, ev := range resp.Kvs {
		ins.addresses[string(ev.Key)] = string(ev.Value)
	}

	go ins.watch()

	return ins, nil
}
