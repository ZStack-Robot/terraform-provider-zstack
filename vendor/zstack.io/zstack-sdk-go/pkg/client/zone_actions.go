package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// QueryZone 查询区域
func (cli *ZSClient) QueryZone(params param.QueryParam) ([]view.ZoneView, error) {
	resp := make([]view.ZoneView, 0)
	return resp, cli.List("v1/zones", &params, &resp)
}