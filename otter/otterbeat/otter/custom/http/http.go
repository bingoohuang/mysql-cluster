package http

import (
	"fmt"
	"github.com/bingoohuang/gg/pkg/jsoni"
	"github.com/bingoohuang/gg/pkg/rest"
)

func Get(url, basicAuth string, body interface{}) error {
	rsp, err := rest.Rest{Addr: url, BasicAuth: basicAuth, Result: &body}.Get()
	if err != nil {
		return fmt.Errorf("GET 请求出错，url: %s basicAuth: %s, err: %w", url, basicAuth, err)
	}
	if !isStatusOk(rsp.Status) {
		return fmt.Errorf("GET 请求出错，url: %s basicAuth: %s, status: %d ", url, basicAuth, rsp.Status)
	}

	return nil
}

func Post(url, basicAuth string, requestBody, responseBody interface{}) error {
	bytes, err := jsoni.Marshal(requestBody)
	if err != nil {
		return err
	}

	rsp, err := rest.Rest{
		Body:      bytes,
		Addr:      url,
		Headers:   jsonContentType(),
		Result:    &responseBody,
		Timeout:   0,
		BasicAuth: basicAuth,
	}.Post()

	if err != nil {
		return fmt.Errorf("POST 请求出错，url: %s basicAuth: %s, err: %w", url, basicAuth, err)
	}

	if !isStatusOk(rsp.Status) {
		return fmt.Errorf("POST 请求出错，url: %s basicAuth: %s, status: %d ", url, basicAuth, rsp.Status)
	}

	return nil
}

// http 状态被认为正确 200 <= status < 300
func isStatusOk(status int) bool {
	return status >= 200 && status < 300
}

// 请求和响应都是 application/json 格式
func jsonContentType() map[string]string {
	return map[string]string{"Accept": "application/json", "Content-Type": "application/json"}
}
