package server

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

type pprof struct{}

func (p *pprof) Start() {
	fmt.Println(http.ListenAndServe(":6060", nil))
}

func (p *pprof) Shutdown(ctx context.Context) error {
	return nil
}
