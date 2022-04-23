package etcdx

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"time"
)

// defaultNotify Watch操作应该执行此事件来增删节点
func (r *Watcher) defaultNotify(event *clientv3.Event) {
	switch event.Type {
	case mvccpb.PUT:
		r.putAddress(string(event.Kv.Key), string(event.Kv.Value))
	case mvccpb.DELETE:
		r.delAddress(string(event.Kv.Key))
	}
}

func (r *Watcher) loadAddress() {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	response, err := r.client.Get(ctx, r.prefix, clientv3.WithPrefix())
	if err == nil {
		for _, kv := range response.Kvs {
			r.putAddress(string(kv.Key), string(kv.Value))
			addresses := r.getAddresses()
			for _, notify := range r.Notifies {
				notify(&clientv3.Event{Type: mvccpb.PUT, Kv: kv}, addresses)
			}
		}
	} else {
		log.Fatal("etcd get fail", err)
	}
}

func (r *Watcher) watch() {
	watch := r.client.Watch(context.Background(), r.prefix, clientv3.WithPrefix())
	for response := range watch {
		for _, event := range response.Events {
			r.defaultNotify(event)
			addresses := r.getAddresses()
			for _, notify := range r.Notifies {
				notify(event, addresses)
			}
		}
	}
}

// putAddress 新增节点
func (r *Watcher) putAddress(key, address string) {
	if strings.TrimSpace(address) == "" {
		return
	}
	r.lock.Lock()
	defer r.lock.Unlock()
	r.addresses[key] = address
}

// delAddress 删除节点
func (r *Watcher) delAddress(key string) {
	r.lock.Lock()
	defer r.lock.Unlock()
	delete(r.addresses, key)
}

// getAddresses 获取当前Watcher里面包含的ip列表
func (r *Watcher) getAddresses() []string {
	var addresses []string
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, address := range r.addresses {
		if strings.TrimSpace(address) == "" {
			continue
		}
		addresses = append(addresses, address)
	}
	return addresses
}
