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
	body := string(r.client.body)
	if r.client.request == nil {
		return ""
	}
	if len(body) == 0 && len(r.client.request.PostForm) > 0 {
		body = r.client.request.PostForm.Encode()
	}
	return buildCurl(r.client.url, r.client.method, body, r.client.request.Header, r.client.request.Cookies())
}

func (r *muteHttpResponse) UseTime() int64 {
	return r.client.useTime
}

func (r *muteHttpResponse) Unmarshal(resp interface{}) error {
	return json.Unmarshal(r.body, resp)
}
