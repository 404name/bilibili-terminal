package utils

import (
	"strconv"

	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

// 代码来自：https://github.com/CuteReimu/bilibili/blob/90408c0f9534f5689290a386467fbdb75403118b/client.go#L126

func formatError(prefix string, code int64, message ...string) error {
	for _, m := range message {
		if len(m) > 0 {
			return errors.New(prefix + "失败，返回值：" + strconv.FormatInt(code, 10) + "，返回信息：" + m)
		}
	}
	return errors.New(prefix + "失败，返回值：" + strconv.FormatInt(code, 10))
}

func GetRespData(resp *resty.Response, prefix string) ([]byte, error) {
	if resp.StatusCode() != 200 {
		return nil, errors.Errorf(prefix+"失败，status code: %d", resp.StatusCode())
	}
	if !gjson.ValidBytes(resp.Body()) {
		return nil, errors.New("json解析失败：" + resp.String())
	}
	res := gjson.ParseBytes(resp.Body())
	code := res.Get("code").Int()
	if code != 0 {
		return nil, formatError(prefix, code, res.Get("message").String(), res.Get("msg").String())
	}

	return []byte(res.Get("data").Raw), nil
}
