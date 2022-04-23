package http

import (
	"fmt"
	"net/http"
)

func buildCurl(uri string, method string, data string, header http.Header, cookies []*http.Cookie) string {
	c := fmt.Sprintf("curl  -X %s '%s'", method, uri)
	if header != nil {
		for k, v := range header {
			if len(v) > 0 {
				c += fmt.Sprintf(" -H '%s:%s'", k, v[0])
			}
		}
	}
	if cookies != nil && len(cookies) > 0 {
		c += " -H 'Cookie: "
		for _, v := range cookies {
			c += fmt.Sprintf("%s=%s;", v.Name, v.Value)
		}
		c += "'"
	}
	if data != "" {
		c += fmt.Sprintf(" -d '%s'", data)
	}
	return c + ` -w '\n\ntime_connect %{time_connect}\ntime_starttransfer %{time_starttransfer}\ntime_total %{time_total}\n'`
}
