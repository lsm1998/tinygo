package redisx

import (
	"context"
	"fmt"
	"testing"
)

func TestMust(t *testing.T) {
	client := Must(WithConfig(Config{
		Addr: "120.79.132.241",
		Port: 6379,
		Auth: "redisyyds123",
	}))
	stringCmd := client.Get(context.Background(), "yidu-book-config")
	fmt.Println(stringCmd.Val(), stringCmd.Err())
}
