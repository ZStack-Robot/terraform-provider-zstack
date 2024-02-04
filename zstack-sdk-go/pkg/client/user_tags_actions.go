package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// CreateUserTag 创建用户标签
func (cli *ZSClient) CreateUserTag(params param.CreateTagParam) (view.UserTagInventoryView, error) {
	var resp view.UserTagInventoryView
	return resp, cli.Post("v1/user-tags", params, &resp)
}

// QueryUserTag 查询用户标签
func (cli *ZSClient) QueryUserTag(params param.QueryParam) ([]view.UserTagInventoryView, error) {
	var tags []view.UserTagInventoryView
	return tags, cli.List("v1/user-tags", &params, &tags)
}

// QueryUserTag 查询全部用户标签
func (cli *ZSClient) ListAllUserTags() ([]view.UserTagInventoryView, error) {
	params := param.NewQueryParam()
	var tags []view.UserTagInventoryView
	return tags, cli.ListAll("v1/user-tags", &params, &tags)
}
