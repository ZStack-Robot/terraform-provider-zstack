// Copyright (c) HashiCorp, Inc.

package client

import (
	"fmt"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/view"
)

type FencerStrategy string

const (
	Force      FencerStrategy = "Force"      //激进
	Permissive FencerStrategy = "Permissive" //保守
)

func (cli *ZSClient) QueryGlobalConfig(params param.QueryParam) ([]view.GlobalConfigView, error) {
	var configurations []view.GlobalConfigView
	return configurations, cli.List("v1/global-configurations", &params, &configurations)
}

func (cli *ZSClient) QueryResourceConfig(params param.QueryParam) ([]view.ResourceConfigView, error) {
	var configurations []view.ResourceConfigView
	return configurations, cli.List("v1/resource-configurations", &params, &configurations)
}

// GetResourceConfig 获取资源的资源高级设置
func (cli *ZSClient) GetResourceConfig(resourceUuid, category, name string) ([]view.ResourceConfigView, error) {
	resp := new([]view.ResourceConfigView)

	return *resp, cli.GetWithSpec("v1/resource-configurations", resourceUuid, fmt.Sprintf("%s/%s", category, name), "effectiveConfigs", nil, resp)
}

// UpdateGlobalConfig 更新资源高级设置
func (cli *ZSClient) UpdateGlobalConfig(category, name string, params param.UpdateGlobalConfigParam) (view.GlobalConfigView, error) {
	resp := new(view.GlobalConfigView)

	return *resp, cli.Put("v1/global-configurations", fmt.Sprintf("%s/%s", category, name), params, resp)
}

func (cli *ZSClient) UpdateResourceConfig(category, name, resourceUuid string, params param.UpdateResourceConfigParam) (view.ResourceConfigView, error) {
	resp := new(view.ResourceConfigView)

	return *resp, cli.Put("v1/resource-configurations", fmt.Sprintf("%s/%s/%s", category, name, resourceUuid), params, resp)
}
