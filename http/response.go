package http

import (
	"encoding/json"
	"net/http"
)

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
	return buildCurl(r.client.url, r.client.method, string(r.client.body), r.client.request.Header, r.client.request.Cookies())
}

func (r *muteHttpResponse) UseTime() int64 {
	return r.client.useTime
}

func (r *muteHttpResponse) Unmarshal(resp interface{}) error {
	return json.Unmarshal(r.body, resp)
}
