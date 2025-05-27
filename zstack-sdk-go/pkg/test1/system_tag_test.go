// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"github.com/kataras/golog"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestQuerySystemTags(t *testing.T) {
	queryParam := param.NewQueryParam()
	//queryParam.AddQ("resourceUuid=e29e11da127e42fb8844d4a0e421c84a")
	tags, err := accountLoginCli.QuerySystemTags(queryParam)
	if err != nil {
		t.Errorf("TestQuerySystemTags %v", err)
	}
	golog.Info(tags)
}
