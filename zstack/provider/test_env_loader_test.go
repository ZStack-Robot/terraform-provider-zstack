// Copyright (c) ZStack.io, Inc.

package provider

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"sync"
	"testing"
)

// EnvData holds all environment resource data loaded from testdata/env.json
type EnvData struct {
	Zones                  []map[string]interface{} `json:"zones"`
	Clusters               []map[string]interface{} `json:"clusters"`
	Hosts                  []map[string]interface{} `json:"hosts"`
	Images                 []map[string]interface{} `json:"images"`
	BackupStorages         []map[string]interface{} `json:"backup_storages"`
	PrimaryStorages        []map[string]interface{} `json:"primary_storages"`
	InstanceOfferings      []map[string]interface{} `json:"instance_offerings"`
	DiskOfferings          []map[string]interface{} `json:"disk_offerings"`
	L2Networks             []map[string]interface{} `json:"l2_networks"`
	L3Networks             []map[string]interface{} `json:"l3_networks"`
	VmInstances            []map[string]interface{} `json:"vm_instances"`
	SecurityGroups         []map[string]interface{} `json:"security_groups"`
	SecurityGroupRules     []map[string]interface{} `json:"security_group_rules"`
	VirtualRouterOfferings []map[string]interface{} `json:"virtual_router_offerings"`
	VirtualRouters         []map[string]interface{} `json:"virtual_routers"`
	SdnControllers         []map[string]interface{} `json:"sdn_controllers"`
	InstanceScripts        []map[string]interface{} `json:"instance_scripts"`
	ScriptExecutions       []map[string]interface{} `json:"script_executions"`
	MnNodes                []map[string]interface{} `json:"mn_nodes"`
	IpRanges               []map[string]interface{} `json:"ip_ranges"`
	VmNics                 []map[string]interface{} `json:"vm_nics"`
}

var (
	envData     *EnvData
	envDataOnce sync.Once
	envDataErr  error
)

// loadEnvData loads testdata/env.json. It is safe for concurrent use.
// If the file does not exist, it calls t.Skip.
func loadEnvData(t *testing.T) *EnvData {
	t.Helper()
	envDataOnce.Do(func() {
		_, filename, _, _ := runtime.Caller(0)
		dir := filepath.Dir(filename)
		path := filepath.Join(dir, "testdata", "env.json")

		data, err := os.ReadFile(path)
		if err != nil {
			envDataErr = fmt.Errorf("cannot read testdata/env.json: %w (run: source .env.test && go run ./zstack/provider/testdata/generate_env.go)", err)
			return
		}
		var env EnvData
		if err := json.Unmarshal(data, &env); err != nil {
			envDataErr = fmt.Errorf("cannot parse testdata/env.json: %w", err)
			return
		}
		envData = &env
	})
	if envDataErr != nil {
		t.Skip(envDataErr.Error())
	}
	return envData
}

// Helper functions to extract string values from map[string]interface{}

func envStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case string:
			return val
		case float64:
			if val == float64(int64(val)) {
				return strconv.FormatInt(int64(val), 10)
			}
			return strconv.FormatFloat(val, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(val)
		case nil:
			return ""
		default:
			return fmt.Sprintf("%v", val)
		}
	}
	return ""
}

func envInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		}
	}
	return 0
}
