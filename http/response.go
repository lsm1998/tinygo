package http

import (
	"encoding/json"
	"net/http"
)

type MuteHttpResponse interface {
	Code() int

	GetBody() []byte

	Curl() string

	UseTime() int64

	Unmarshal(resp interface{}) error

	GetHeader(key string) string

	Header() http.Header
}

type muteHttpResponse struct {
	response *http.Response
	client   muteHttpClient
	body     []byte
	url      string
}

func (r *muteHttpResponse) Code() int {
	if r.response == nil {
		return 0
	}
	return r.response.StatusCode
}

func (r *muteHttpResponse) GetBody() []byte {
	return r.body
}

func (r *muteHttpResponse) Curl() string {
	if r.client.request == nil {
		return ""
	}
	return buildCurl(r.client.url, r.client.method, string(r.client.bodyByte), r.client.request.Header, r.client.request.Cookies())
}

func (r *muteHttpResponse) UseTime() int64 {
	return r.client.useTime
}

func (r *muteHttpResponse) Unmarshal(resp interface{}) error {
	return json.Unmarshal(r.body, resp)
}

func (r *muteHttpResponse) GetHeader(key string) string {
	if r.response == nil {
		return ""
	}
	return r.response.Header.Get(key)
}

func (r *muteHttpResponse) Header() http.Header {
	if r.response == nil {
		return nil
	}
	return r.response.Header
}
