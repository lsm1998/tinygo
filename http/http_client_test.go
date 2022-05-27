package http

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()
	response, err := New("http://www.baidu.com").
		SetQuery("name", "lsm").
		SetQuery("age", "25").
		SetQuery("sex", "1").
		SetPostForm(url.Values{"username": []string{"lsm"}}).
		Post(ctx)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Println(string(response.GetBody()))
	fmt.Println(response.Curl())
	fmt.Println(response.UseTime())
}
