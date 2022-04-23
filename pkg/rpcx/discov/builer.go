package discov

import (
	"fmt"
	"github.com/lsm1998/tinygo/pkg/etcdx"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
)

type Builder struct {
	watcher *etcdx.Watcher
	domain  string
}

// NewBuilder TODO 将 client，domain 用 options 包裹起来，以统一写法
func NewBuilder(client *clientv3.Client, domain string) *Builder {
	if domain == "" {
		logrus.Warn("domain cannot be empty")
	}
	r := &Builder{
		watcher: etcdx.NewWatcher(client),
		domain:  domain,
	}
	resolver.Register(r)
	return r
}

func (b *Builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		cc: cc,
	}
	prefix := fmt.Sprintf("/%s/%s/", "docer_discov", target.Endpoint)
	notify := etcdx.Notify(func(_ *clientv3.Event, address []string) {
		if len(address) == 0 {
			address = append(address, b.domain)
			logrus.Warnf("target: %s, grpc address is empty , add the application address  %+v", prefix, address)
		} else {
			logrus.Infof("target: %s, grpc address change to %+v", prefix, address)
		}
		r.Update(address)
	})
	b.watcher.SetPrefix(prefix)
	b.watcher.AddEvent(notify)
	go b.watcher.Run()
	return r, nil
}

func (b *Builder) Scheme() string {
	return schema
}
