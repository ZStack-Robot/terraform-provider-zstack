// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// DeleteTag Delete a tag
func (cli *ZSClient) DeleteTag(uuid string, mode param.DeleteMode) error {
	return cli.Delete("v1/tags", uuid, string(mode))
}

// CreateSystemTag Create a system tag
func (cli *ZSClient) CreateSystemTag(params param.CreateTagParam) (view.SystemTagView, error) {
	var resp view.SystemTagView
	return resp, cli.Post("v1/system-tags", params, &resp)
}

// UpdateSystemTag Update a system tag
func (cli *ZSClient) UpdateSystemTag(uuid string, params param.UpdateSystemTagParam) (view.SystemTagView, error) {
	var resp view.SystemTagView
	return resp, cli.Put("v1/system-tags", uuid, params, &resp)
}

// QuerySystemTags Query system tags
func (cli *ZSClient) QuerySystemTags(params param.QueryParam) ([]view.SystemTagView, error) {
	var tags []view.SystemTagView
	return tags, cli.List("v1/system-tags", &params, &tags)
}
