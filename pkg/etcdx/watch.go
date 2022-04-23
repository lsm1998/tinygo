package etcdx

import (
	"go.etcd.io/etcd/client/v3"
	"sync"
)

type Watcher struct {
	lock      sync.RWMutex
	client    *clientv3.Client
	prefix    string
	addresses map[string]string
	Notifies  []Notify
}

func NewWatcher(client *clientv3.Client) *Watcher {
	return &Watcher{
		lock:      sync.RWMutex{},
		client:    client,
		addresses: make(map[string]string),
	}
}

type Notify func(ev *clientv3.Event, address []string)

func (r *Watcher) AddEvent(notify Notify) *Watcher {
	r.Notifies = append(r.Notifies, notify)
	return r
}

func (r *Watcher) SetPrefix(prefix string) *Watcher {
	r.prefix = prefix
	return r
}

// Run 启动节点监听机制
// 在刚进入函数，会先执行一遍notify函数，来处理当地址列表为空时的操作，这样我们就能够填充应用域名保证连接有效
func (r *Watcher) Run() {
	for _, notify := range r.Notifies {
		notify(nil, []string{})
	}
	r.loadAddress()
	r.watch()
}

func (r *Watcher) GetAddresses() []string {
	return r.getAddresses()
}
