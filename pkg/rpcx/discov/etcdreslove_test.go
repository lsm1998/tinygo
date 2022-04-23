package discov

import (
	"go.etcd.io/etcd/client/v3"
	"testing"
)

func TestRegister(t *testing.T) {
	client, err := clientv3.New(clientv3.Config{Endpoints: []string{"localhost:2379"}})
	if err != nil {
		t.Error(err)
		return
	}
	reg := NewRegister(client)
	discov := reg.Discov("project", "8088")
	for {
		select {
		case v := <-discov:
			t.Error(v)
		}
	}
}
