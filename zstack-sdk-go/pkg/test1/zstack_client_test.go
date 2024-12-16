// # Copyright (c) ZStack.io, Inc.

package test

import (
	"reflect"
	"testing"

	"github.com/kataras/golog"
	"zstack.io/zstack-sdk-go/pkg/client"
	"zstack.io/zstack-sdk-go/pkg/util/jsonutils"
)

func TestZSClient_ValidateSession(t *testing.T) {
	tests := []struct {
		name    string
		cli     *client.ZSClient
		want    map[string]bool
		wantErr bool
	}{
		{accountLogin, accountLoginCli, map[string]bool{"valid": true}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cli.ValidateSession()
			if (err != nil) != tt.wantErr {
				t.Errorf("ZSClient.ValidateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZSClient.ValidateSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestZSClient_Zql(t *testing.T) {
	if _, err := accountLoginCli.Zql("query Image where __systemTag__='applianceType::vrouter'", nil); err != nil {
		t.Fatalf("TestZSClient_Zql error: %v", err)
	}
}

func TestZSClient_WebLogin(t *testing.T) {
	data, err := accountLoginCli.WebLogin()
	if err != nil {
		t.Fatalf("TestZSClient_WebLogin error: %v", err)
	}
	golog.Infof("%v", jsonutils.Marshal(data))

	result, err := accountLoginCli.ValidateSessionId(data.SessionId)
	if err != nil {
		t.Fatalf("TestZSClient_WebLogin ValidateSessionId error: %v", err)
	}
	golog.Infof("%v", jsonutils.Marshal(result))
}
