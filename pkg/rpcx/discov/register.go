package discov

import (
	"context"
	"fmt"
	"github.com/lsm1998/tinygo"
	"github.com/prometheus/common/log"
	"go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

const (
	PAUSE = "pause"
	START = "start"
	EXIT  = "exit"
)

type Register struct {
	client      *clientv3.Client
	state       string
	serviceName string
	serviceAddr string
	ttl         int64
}

func NewRegister(client *clientv3.Client) *Register {
	r := &Register{client: client, state: START}
	return r
}

//Discov 服务注册
func (r *Register) Discov(serviceName string, port string) <-chan error {
	e := make(chan error)
	split := strings.Split(port, ":")
	if len(split) > 0 {
		port = split[len(split)-1]
	}
	var err error
	var ip string

	for {
		ip, err = tinygo.IpAddr()
		if err != nil {
			e <- err
			return e
		} else if ip == "" {
			log.Warn("re acquire IP...")
			time.Sleep(time.Second * 1)
		}
		break
	}

	r.serviceAddr = ip + ":" + port
	r.serviceName = serviceName

	ch := make(chan struct{}, 1)
	go r.async(ch, e)
	return e
}

func (r *Register) async(ch chan struct{}, e chan error) {
	ch <- struct{}{}
	for {
		switch r.state {
		case PAUSE:
			time.Sleep(time.Second)
			continue
		case EXIT:
			return
		}

		select {
		case <-ch:
			err := r.keepAlive(ch)
			if err != nil {
				e <- err
				time.Sleep(time.Second * time.Duration(defaultTTL))
				if len(ch) == 0 {
					ch <- struct{}{}
				}
			}
		}
	}
}

func (r *Register) keepAlive(restart chan struct{}) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer func() {
		if err != nil {
			cancel()
		}
	}()
	// 创建租约
	lease, err := r.client.Grant(ctx, defaultTTL)
	if err != nil {
		return err
	}
	// 设置键值对
	_, err = r.client.Put(ctx, r.key(), r.serviceAddr, clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}
	// 续签
	ch, err := r.client.KeepAlive(context.Background(), lease.ID)
	if err != nil {
		return err
	}
	go func() {
		for {
			if v := <-ch; v == nil {
				restart <- struct{}{}
				return
			}
		}
	}()
	return nil
}

func (r *Register) key() string {
	return fmt.Sprintf("%s/%s", GetPrefix(r.serviceName), r.serviceAddr)
}

func GetPrefix(serviceName string) string {
	return fmt.Sprintf("/%s/%s/%s", "docer_discov", tinygo.AppMod, serviceName)
}

func (r *Register) Pause() {
	r.state = PAUSE
}
func (r *Register) Start() {
	r.state = START
}

func (r *Register) Exit() {
	_, _ = r.client.Delete(context.Background(), r.key())
	go func() {
		for i := 0; i < 20; i++ {
			_, _ = r.client.Delete(context.Background(), r.key())
			time.Sleep(time.Millisecond * 300)
		}
	}()
	r.state = EXIT
}
