// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestZSClient_GetUserTag(t *testing.T) {
	queryParam := param.NewQueryParam()
	queryParam.AddQ("resourceUuid=2758914006f244879ec642a82406f8f3")
	tags, err := accountLoginCli.QueryUserTag(queryParam)
	if err != nil {
		t.Errorf("TestQuerySystemTags %v", err)
	}
	golog.Info(tags)
}

func TestCreateUserTag(t *testing.T) {
	tag, err := accountLoginCli.CreateUserTag(param.CreateTagParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateTagDetailParam{
			ResourceType: param.ResourceTypeVolumeVo,
			ResourceUuid: "5a7f72aa7f8041ea984f3cdabc3e9840",
			Tag:          "userID::1",
		},
	})
	if err != nil {
		t.Errorf("TestCreateUserTag %v", err)
		return
	}
	golog.Info(tag)
}
