package configx

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var client *redis.Client

func get(key string) (string, error) {
	cmd := client.Get(context.Background(), key)
	if err := cmd.Err(); err != nil {
		return "", err
	}
	return cmd.Val(), cmd.Err()
}

func redisParse(c config, key string, obj interface{}) error {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Addr, c.Port),
		Password: c.Auth,
		DB:       c.Db,
	})
	s, err := get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(s), obj)
}
