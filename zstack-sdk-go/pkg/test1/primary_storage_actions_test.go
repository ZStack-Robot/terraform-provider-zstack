// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"zstack.io/zstack-sdk-go/pkg/param"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestQueryPrimaryStorage(t *testing.T) {
	storage, err := accountLoginCli.QueryPrimaryStorage(param.NewQueryParam())
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(jsonutils.Marshal(storage))
}

func TestPagePrimaryStorage(t *testing.T) {
	p := param.NewQueryParam()
	p.Start(0)
	p.Limit(10)
	storage, total, err := accountLoginCli.PagePrimaryStorage(p)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(jsonutils.Marshal(storage))
	t.Log(total)
}

func TestPageCephPrimaryStoragePool(t *testing.T) {
	p := param.NewQueryParam()
	p.AddQ("primaryStorageUuid=85a7510f50d144b6a664ff55b05a9bd2")
	p.Start(0)
	p.Limit(10)
	storage, total, err := accountLoginCli.PageCephPrimaryStoragePool(p)
	if err != nil {
		t.Error(err)
		return
	}

	t.Log(jsonutils.Marshal(storage))
	t.Log(total)
}
