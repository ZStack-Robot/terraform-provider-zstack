// # Copyright (c) ZStack.io, Inc.

package test

import (
	"testing"

	"zstack.io/zstack-sdk-go/pkg/param"
)

func TestQueryCluster(t *testing.T) {
	cluster, err := accountLoginCli.QueryCluster(param.NewQueryParam())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(cluster)
}

func TestQueryClusterByKey(t *testing.T) {
	cluster, err := accessKeyAuthCli.QueryCluster(param.NewQueryParam())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(cluster)
}
