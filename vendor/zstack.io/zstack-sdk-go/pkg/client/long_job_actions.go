package client

import (
	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

// PageLongJob 分页查询长时任务
func (cli *ZSClient) PageLongJob(params param.QueryParam) ([]view.LongJobInventoryView, int, error) {
	var resp []view.LongJobInventoryView
	total, err := cli.Page("v1/longjobs", &params, &resp)
	return resp, total, err
}

// QueryLongJob 查询长时任务
func (cli *ZSClient) QueryLongJob(queryParam param.QueryParam) ([]view.LongJobInventoryView, error) {
	var resp []view.LongJobInventoryView
	return resp, cli.List("v1/longjobs", &queryParam, &resp)
}

// GetLongJob 获取长时任务
func (cli *ZSClient) GetLongJob(uuid string) (*view.LongJobInventoryView, error) {
	var resp view.LongJobInventoryView
	if err := cli.Get("v1/longjobs", uuid, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SubmitLongJob 提交长时任务
func (cli *ZSClient) SubmitLongJob(params *param.SubmitLongJobParam) (*view.LongJobInventoryView, error) {
	var resp view.LongJobInventoryView
	if err := cli.Post("v1/longjobs", params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateLongJob 更新长时任务
func (cli *ZSClient) UpdateLongJob(uuid string, params *param.UpdateLongJobParam) (*view.LongJobInventoryView, error) {
	var resp view.LongJobInventoryView
	if err := cli.Put("v1/longjobs", uuid, params, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// CancelLongJob 取消长时任务
func (cli *ZSClient) CancelLongJob(uuid string) error {
	params := map[string]struct{}{
		"cancelLongJob": {},
	}
	return cli.Put("v1/longjobs", uuid, params, nil)
}

// DeleteLongJob 删除长时任务
func (cli *ZSClient) DeleteLongJob(uuid string) error {
	return cli.Delete("v1/longjobs", uuid, "")
}

// GetTaskProgress 获取任务进度
func (cli *ZSClient) GetTaskProgress(apiId string) (*view.TaskProgressInventoryView, error) {
	var resp view.TaskProgressInventoryView
	if err := cli.Get("v1/task-progresses", apiId, nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
