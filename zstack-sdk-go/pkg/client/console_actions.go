package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// RequestConsoleAccess 请求控制台访问地址
func (cli *ZSClient) RequestConsoleAccess(params param.RequestConsoleAccessParam) (view.ConsoleInventoryView, error) {
	var resp view.ConsoleInventoryView
	return resp, cli.Post("v1/consoles", &params, &resp)
}
