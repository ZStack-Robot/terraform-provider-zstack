// Copyright (c) ZStack.io, Inc.

package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// PageLongJob Paginated query for long-running tasks
func (cli *ZSClient) PageLongJob(params param.QueryParam) ([]view.LongJobInventoryView, int, error) {
	var resp []view.LongJobInventoryView
	total, err := cli.Page("v1/longjobs", &params, &resp)
	return resp, total, err
}

// QueryLongJob Query long-running tasks
func (cli *ZSClient) QueryLongJob(queryParam param.QueryParam) ([]view.LongJobInventoryView, error) {
	var resp []view.LongJobInventoryView
	return resp, cli.List("v1/longjobs", &queryParam, &resp)
}

// GetLongJob Retrieve a long-running task
func (cli *ZSClient) GetLongJob(uuid string) (*view.LongJobInventoryView, error) {
	var resp view.LongJobInventoryView
	if err := cli.Get("v1/longjobs", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubmitLongJob Submit a long-running task
func (cli *ZSClient) SubmitLongJob(params *param.SubmitLongJobParam) (*view.LongJobInventoryView, error) {
	var resp view.LongJobInventoryView
	if err := cli.Post("v1/longjobs", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateLongJob Update a long-running task
func (cli *ZSClient) UpdateLongJob(uuid string, params *param.UpdateLongJobParam) (*view.LongJobInventoryView, error) {
	var resp view.LongJobInventoryView
	if err := cli.Put("v1/longjobs", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelLongJob Cancel a long-running task
func (cli *ZSClient) CancelLongJob(uuid string) error {
	params := map[string]struct{}{
		"cancelLongJob": {},
	}
	return cli.Put("v1/longjobs", uuid, params, nil)
}

// DeleteLongJob Delete a long-running task
func (cli *ZSClient) DeleteLongJob(uuid string) error {
	return cli.Delete("v1/longjobs", uuid, "")
}

// GetTaskProgress Retrieve task progress
func (cli *ZSClient) GetTaskProgress(apiId string) (*view.TaskProgressInventoryView, error) {
	var resp view.TaskProgressInventoryView
	if err := cli.Get("v1/task-progresses", apiId, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
