package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// DeleteTag 删除标签
func (cli *ZSClient) DeleteTag(uuid string, mode param.DeleteMode) error {
	return cli.Delete("v1/tags", uuid, string(mode))
}

// CreateSystemTag 创建系统标签
func (cli *ZSClient) CreateSystemTag(params param.CreateTagParam) (view.SystemTagView, error) {
	var resp view.SystemTagView
	return resp, cli.Post("v1/system-tags", params, &resp)
}

// UpdateSystemTag 更新系统标签
func (cli *ZSClient) UpdateSystemTag(uuid string, params param.UpdateSystemTagParam) (view.SystemTagView, error) {
	var resp view.SystemTagView
	return resp, cli.Put("v1/system-tags", uuid, params, &resp)
}

// QuerySystemTags 查询系统标签
func (cli *ZSClient) QuerySystemTags(params param.QueryParam) ([]view.SystemTagView, error) {
	var tags []view.SystemTagView
	return tags, cli.List("v1/system-tags", &params, &tags)
}
