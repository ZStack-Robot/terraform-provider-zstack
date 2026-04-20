# Story 10: Sprint A — 自包含资源验收测试 (8 个)

> **分支**: 待创建 (`test/sprint-a-self-contained`)  
> **状态**: 未开始  
> **优先级**: P1 — Phase 1 Sprint A  
> **前置依赖**: story-07 (RequiresReplace 清理，确保 monitor_template/price_table/preconfiguration_template Update 可用)  
> **预计工作量**: 2 天

---

## Story

作为 Provider 维护者，我需要为 8 个无外部依赖的核心资源新增验收测试，直接达到"标准"级 (C+U+I+D) 或"基础"级 (C+I+D, U=N/A)。

## 资源清单

| # | Resource | Update 状态 | 目标级别 | 环境依赖 |
|---|----------|-----------|---------|---------|
| 1 | `access_key` | N/A (no-op) | 基础 | account_uuid (env.json) |
| 2 | `monitor_template` | 有 (UpdateMonitorTemplate) | 标准 | 无 |
| 3 | `log_server` | 有 (UpdateLogServer) | 标准 | 无 |
| 4 | `cdp_policy` | 有 (UpdateCdpPolicy) | 标准 | 无 |
| 5 | `price_table` | 有 (UpdatePriceTable) | 标准 | 无 |
| 6 | `stack_template` | 有 (UpdateStackTemplate) | 标准 | 无 |
| 7 | `preconfiguration_template` | 有 (UpdatePreconfigurationTemplate) | 标准 | 无 |
| 8 | `instance_scripts` | 有 (UpdateGuestVmScript) | 标准 | 无 |

## 验收标准

- [ ] AC1: 8 个资源各有 `TestAcc<X>Resource` 函数，包含 Create + Import + Destroy (基础)
- [ ] AC2: 7 个有 Update 的资源包含 Update Step (标准)
- [ ] AC3: 8 个资源各有 `testAccCheck<X>Destroy` 函数
- [ ] AC4: `check_destroy_test.go` 新增对应的 Destroy 检查函数
- [ ] AC5: 编译通过，所有新测试编译通过
- [ ] AC6: 抽样运行 3 个验收测试确认通过

## Tasks

### 基础设施

- [ ] 1. 确认 `env.json` 中包含 accounts 数据（access_key 需要 account_uuid）
- [ ] 2. 在 `check_destroy_test.go` 新增 8 个 Destroy 检查函数

### 验收测试（每个资源一个 task）

- [ ] 3. `resource_zstack_access_key_test.go` — Create (需 account_uuid) + Import + Destroy
- [ ] 4. `resource_zstack_monitor_template_test.go` — Create + Update (name) + Import + Destroy
- [ ] 5. `resource_zstack_log_server_test.go` — Create + Update (name) + Import + Destroy
- [ ] 6. `resource_zstack_cdp_policy_test.go` — Create + Update (name) + Import + Destroy
- [ ] 7. `resource_zstack_price_table_test.go` — Create + Update (name) + Import + Destroy
- [ ] 8. `resource_zstack_stack_template_test.go` — Create + Update (name) + Import + Destroy
- [ ] 9. `resource_zstack_preconfiguration_template_test.go` — Create + Update (name) + Import + Destroy
- [ ] 10. `resource_zstack_instance_scripts_test.go` — Create + Update (name) + Import + Destroy

### 验证

- [ ] 11. 编译确认 (`go build ./...`)
- [ ] 12. Schema 单元测试通过
- [ ] 13. 抽样运行 3 个验收测试

## 审查要点

1. 每个资源的 HCL Config 是否包含所有 Required 属性
2. Update Step 修改的属性是否为实际可变属性（不是 ForceNew 的）
3. Import Step 的 ImportStateIdFunc 是否匹配 ImportState 方法的 ID 格式
4. Destroy 检查函数使用 Get-style 还是 Query-style
