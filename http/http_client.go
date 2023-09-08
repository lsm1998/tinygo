package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	codeMustErr = errors.New("HTTP code mismatch")

	defaultClient = &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				var d net.Dialer
				c, err := d.DialContext(ctx, network, addr)
				// c, err := net.DialTimeout(network, addr, timeout) //设置建立连接超时
				if err != nil {
					return nil, err
				}
				// _ = c.SetDeadline(time.Now().Add(timeout)) //设置发送接收数据超时
				return c, nil
			},
		},
	}
)

type MuteHttpClient interface {
	AddCookie(cookies ...*http.Cookie) MuteHttpClient

	SetHeader(key, value string) MuteHttpClient

	SetQuery(key, value string) MuteHttpClient

	Header(header http.Header) MuteHttpClient

	SetPostForm(value url.Values) MuteHttpClient

	SetBodyJSON(obj interface{}) MuteHttpClient

	MustCode(code int) MuteHttpClient

	Post(ctx context.Context) (MuteHttpResponse, error)

	Get(ctx context.Context) (MuteHttpResponse, error)

	Put(ctx context.Context) (MuteHttpResponse, error)

	Delete(ctx context.Context) (MuteHttpResponse, error)

	Options(ctx context.Context) (MuteHttpResponse, error)

	Patch(ctx context.Context) (MuteHttpResponse, error)
}

type muteHttpClient struct {
	url      string
	mustCode int
	method   string
	bodyByte []byte
	body     io.Reader
	request  *http.Request
	client   *http.Client
	useTime  int64
	header   http.Header
	cookies  []*http.Cookie
	postForm url.Values
	query    url.Values
}

func New(url string) MuteHttpClient {
	return &muteHttpClient{url: url, client: defaultClient, header: make(http.Header)}
}

func (c *muteHttpClient) SetBodyJSON(obj interface{}) MuteHttpClient {
	c.bodyByte, _ = json.Marshal(obj)

	c.body = strings.NewReader(string(c.bodyByte))

	c.header.Set("Content-Type", "application/json")
	return c
}

func (c *muteHttpClient) AddCookie(cookies ...*http.Cookie) MuteHttpClient {
	for _, cookie := range cookies {
		c.request.AddCookie(cookie)
	}
	return c
}

func (c *muteHttpClient) SetHeader(key, value string) MuteHttpClient {
	if _, ok := c.header[key]; ok {
		c.header[key] = append(c.header[key], value)
	} else {
		c.header[key] = []string{value}
	}
	return c
}

func (c *muteHttpClient) SetQuery(key, value string) MuteHttpClient {
	if c.query == nil {
		c.query = make(url.Values)
	}
	c.query.Set(key, value)
	return c
}

func (c *muteHttpClient) Header(header http.Header) MuteHttpClient {
	c.header = header
	return c
}

func (c *muteHttpClient) SetPostForm(value url.Values) MuteHttpClient {
	c.header.Set("Content-Type", "application/x-www-form-urlencoded")
	c.postForm = value
	return c
}

func (c *muteHttpClient) MustCode(code int) MuteHttpClient {
	c.mustCode = code
	return c
}

func (c *muteHttpClient) Post(ctx context.Context) (MuteHttpResponse, error) {
	return c.do(http.MethodPost, ctx)
}

func (c *muteHttpClient) Get(ctx context.Context) (MuteHttpResponse, error) {
	return c.do(http.MethodGet, ctx)
}

func (c *muteHttpClient) Put(ctx context.Context) (MuteHttpResponse, error) {
	return c.do(http.MethodPut, ctx)
}

func (c *muteHttpClient) Delete(ctx context.Context) (MuteHttpResponse, error) {
	return c.do(http.MethodDelete, ctx)
}

func (c *muteHttpClient) Options(ctx context.Context) (MuteHttpResponse, error) {
	return c.do(http.MethodOptions, ctx)
}

func (c *muteHttpClient) Patch(ctx context.Context) (MuteHttpResponse, error) {
	return c.do(http.MethodPatch, ctx)
}

func (c *muteHttpClient) do(method string, ctx context.Context) (MuteHttpResponse, error) {
	now := time.Now().UnixMilli()
	var err error
	var response *http.Response
	var result = new(muteHttpResponse)
	var query url.Values
	c.method = method
	c.request, err = http.NewRequest(method, c.url, c.body)
	if err != nil {
		goto RESULT
	}
	c.request = c.request.WithContext(ctx)
	c.request.Header = c.header
	c.request.PostForm = c.postForm
	c.request.URL, err = url.ParseRequestURI(c.url)
	query = c.request.URL.Query()
	if err != nil {
		goto RESULT
	}

	for key, values := range c.query {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	c.request.URL.RawQuery = query.Encode()

	response, err = c.client.Do(c.request)
	if err != nil {
		goto RESULT
	}
	result.body, err = io.ReadAll(response.Body)
	if err != nil {
		goto RESULT
	}
	defer response.Body.Close()
	if c.mustCode > 0 && response.StatusCode != c.mustCode {
		err = codeMustErr
		goto RESULT
	}
	c.useTime = time.Now().UnixMilli() - now
RESULT:
	result.response = response
	result.client = *c
	return result, err
}
