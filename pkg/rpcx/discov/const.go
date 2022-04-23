package discov

import "os"

const (
	defaultSchema = "etcd-discov"
)

var (
	schema     string
	defaultTTL int64 = 5
)

func init() {
	schema = os.Getenv("ETCD_SCHEMA")
	if schema == "" {
		schema = defaultSchema
	}
}
