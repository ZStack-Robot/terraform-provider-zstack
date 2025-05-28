// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestZSClient_QueryUserTag(t *testing.T) {
	queryParam := param.NewQueryParam()
	queryParam.AddQ("tagPatternUuid=fbe3c0f934ce417995cebbe79d415cc0")
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

func TestCreateTag(t *testing.T) {
	tag, err := accountLoginCli.CreateTag(param.CreateResourceTagParam{
		BaseParam: param.BaseParam{},
		Params: param.CreateResourceTagDetailParam{
			Name:        "performance88",
			Value:       "performance::{performance1}",
			Description: "tag2 for test",
			Color:       "#000000",
			Type:        "withToken",
		},
	})
	if err != nil {
		t.Errorf("TestCreateTag %v", err)
		return
	}
	golog.Info(tag)
}

func TestUpdateTag(t *testing.T) {
	tag, err := accountLoginCli.UpdateTag("0d2dfada733746699cbb2a11276ada17", param.UpdateResourceTagParam{
		BaseParam: param.BaseParam{},
		UpdateResourceTag: param.UpdateResourceTagDetailParam{
			Name: "performance50",
			//	Value:       "performance::{performance222}",
			Description: "tag for test yz5",
			Color:       "#000000",
		},
	})
	if err != nil {
		t.Errorf("TestUpdateTag %v", err)
		return
	}
	golog.Info(tag)
}

func TestDeleteTag(t *testing.T) {
	err := accountLoginCli.DeleteTag("eed8f45ede5348f3960b9dff1cd55f7c", param.DeleteModePermissive)
	if err != nil {
		t.Errorf("TestDeleteTag %v", err)
		return
	}
	golog.Info("delete tag success")
}

func TestQueryTag(t *testing.T) {
	queryParam := param.NewQueryParam()
	queryParam.AddQ("name=test::tag1")
	tags, err := accountLoginCli.QueryTag(queryParam)
	if err != nil {
		t.Errorf("TestQueryResourceTags %v", err)
	}
	golog.Info(tags)
}

func TestListAllTags(t *testing.T) {
	tags, err := accountLoginCli.ListAllTags()
	if err != nil {
		t.Errorf("TestListAllTags %v", err)
	}
	golog.Info(tags)
}

func TestGetTag(t *testing.T) {
	tag, err := accountLoginCli.GetTag("15b4769130f542b18e0296da4ef2ef80")
	if err != nil {
		t.Errorf("TestGetTag %v", err)
	}
	golog.Info(tag)
}

func TestGetUserTag(t *testing.T) {
	tag, err := accountLoginCli.GetUserTag("15b4769130f542b18e0296da4ef2ef80")
	if err != nil {
		t.Errorf("TestGetUserTag %v", err)
	}
	golog.Info(tag)
}

func TestAttachTagToResource(t *testing.T) {
	tag, err := accountLoginCli.AttachTagToResource("3d7ae53107994d36a202bdf704c007b5", []string{"1e7e24218e314b2b80f3c0e655efc3cf"})
	if err != nil {
		t.Errorf("TestAttachTagToResource %v", err)
	}
	golog.Info(tag)
}

func TestAttachTagToResourceWithToken(t *testing.T) {
	tag, err := accountLoginCli.AttachTagToResource("fbe3c0f934ce417995cebbe79d415cc0", []string{"1bfd7c3dd7ba48b29b9c13cf26ec7037"}, "withToken", `{"performance1":"low"}`)
	if err != nil {
		t.Errorf("TestAttachTagToResourceWithToken %v", err)
	}
	golog.Info(tag)
}

func TestDetachTagFromResource(t *testing.T) {
	err := accountLoginCli.DetachTagFromResource("0b0a50adfcba457db09629b3c365d66f", []string{"1e7e24218e314b2b80f3c0e655efc3cf"})
	if err != nil {
		t.Errorf("TestDetachTagFromResource %v", err)
	}
	golog.Info("detach tag success")
}
