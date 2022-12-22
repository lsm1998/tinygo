package example

import (
	"github.com/lsm1998/tinygo/pkg/configx"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type Config struct {
		Name string `json:"name" yaml:"name"`
	}
	var C Config
	if err := configx.Must(configx.WithLocal("config/config.yaml", &C)); err != nil {
		t.Fatal(err)
	}
	t.Log(C)
}
