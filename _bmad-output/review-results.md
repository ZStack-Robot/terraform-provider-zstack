# 分支审查 & 测试结果汇总

> 审查日期: 2026-04-17
> 审查环境: ZStack @ 172.24.248.129:8080 (部分测试因 503 不可用)

---

## 总览

| # | 分支 | 结论 | Critical | Major | Medium | Minor | 验收测试 |
|---|------|------|----------|-------|--------|-------|---------|
| 1 | `fix/policy-empty-statements` | **Needs Changes** | 1 | 1 | 1 | 0 | FAIL (description unknown) |
| 2 | `test/composite-id-import-steps` | **Approve (comment)** | 0 | 0 | 1 | 2 | scheduler_job PASS / global_config FAIL (预存) |
| 3 | `fix/port-forwarding-rule` | **Approve (comment)** | 0 | 0 | 1 | 1 | FAIL (503 基础设施问题) |
| 4a | `fix/monitor-group-requires-replace` | **Pass** | 0 | 0 | 0 | 0 | 无验收测试 |
| 4b | `fix/monitor-template-requires-replace` | **Pass** | 0 | 0 | 0 | 0 | 无验收测试 |
| 4c | `fix/preconfiguration-template-requires-replace` | **Pass** | 0 | 0 | 0 | 0 | 无验收测试 |
| 4d | `fix/price-table-requires-replace` | **Pass** | 0 | 0 | 0 | 0 | 无验收测试 |
| 5 | `feat/l2vxlan-network-update` | **Needs Changes** | 1 | 0 | 1 | 1 | 无验收测试 |
| 6 | `refactor/read-error-handling` | **Needs Changes** | 1 | 0 | 1 | 0 | BLOCKED (编译失败) |
| 7 | `fix/alarm-update-and-test` | **Approve (comment)** | 0 | 0 | 1 | 1 | 与 #6 共享 |
| 8 | `refactor/provider-quality-hardening` | **Needs Changes** | 1 | 0 | 2 | 1 | FAIL (503 预存) |

**统计**: 4 Pass / 3 Approve with comments / 4 Needs Changes

---

## 详细发现

### #1 `fix/policy-empty-statements` — Needs Changes

**[CRITICAL] description 字段 Create 后仍为 Unknown**
- 位置: `resource_zstack_policy.go` Create 方法
- 问题: `description` 是 `Optional+Computed` + `UseStateForUnknown()`，用户省略时 plan 阶段为 Unknown。Create 方法未设置具体值，Terraform 报错 "provider still indicated an unknown value"
- 验收测试确认: API 调用成功（policy 创建+清理），但 state 处理失败
- 修复:
  ```go
  if plan.Description.IsUnknown() || plan.Description.IsNull() {
      plan.Description = types.StringValue("")
  }
  ```

**[HIGH] 默认 statement Allow+** 是全权限策略**
- 位置: `resource_zstack_policy.go:120-126`
- 问题: 每个通过此资源创建的 policy 都获得 `Effect: "Allow", Actions: ["**"]`（god-mode），用户无法控制
- 最低要求: 在 schema Description 中文档化此默认行为

**[MEDIUM] Update 返回 error 未说明原因**
- 所有字段都有 RequiresReplace，Update 理论上不会被调用，但缺少注释说明

**单元测试**: PASS (`TestPolicyResource_Metadata`)
**验收测试**: FAIL (`TestAccPolicyResource` — description unknown 错误)

---

### #2 `test/composite-id-import-steps` — Approve (comment)

**[MEDIUM] 函数命名不符合项目惯例**
- `importStateIdSchedulerJob` / `importStateIdGlobalConfig` 应改为 `importStateIdFrom<Source>` 格式
- 项目惯例: `importStateIdFromUUID`（描述数据来源，而非资源名）
- 建议: `importStateIdFromUUIDType` / `importStateIdFromCategoryName`

**[LOW] helper 函数应移至 provider_test.go 共享**
- 当前各自定义在 resource test 文件中，与 `importStateIdFromUUID` 的共享模式不一致

**[LOW] 分隔符不一致 `/` vs `:`**
- 确认是资源 ImportState 实现层面的设计，非本分支问题

**验收测试**:
- `TestAccSchedulerJobResource` (含 Import Step): **PASS**
- `TestAccGlobalConfigResource`: **FAIL** — "Provider produced inconsistent result after apply"（预存问题，非本分支引入，Import Step 未被执行到）

---

### #3 `fix/port-forwarding-rule` — Approve (comment)

**[MEDIUM] 同一 bug 模式存在于其他资源**
- `resource_zstack_l3network.go`: `IpVersion` (Int64), `System` (Bool) — Optional+Computed 只检查 IsNull()
- `resource_zstack_instance_scripts.go`: `ScriptTimeout` (Int64) — 同上
- 建议: 后续全局修复 `fix/optional-computed-isnull-guard`

**[LOW] VmNicUuid 的 IsNull() 检查不一致**
- Optional+Computed String 用 `!IsNull() && ValueString() != ""` 但没检查 `IsUnknown()`
- 偶然安全（Unknown 时 ValueString() 返回 ""），但与本分支建立的模式不一致

**本分支修复正确**: 3 个 Int64 字段全部覆盖，测试断言类型修正 (String→Int64) 也正确
**单元测试**: PASS
**���收测试**: FAIL (503 — ZStack API 不可达，基础设施问题)

---

### #4a-d RequiresReplace 清理 — 全部 Pass

| 分支 | 文件 | 移除字段 | Update 方法验证 |
|------|------|---------|----------------|
| 4a monitor_group | resource_zstack_monitor_group.go | name, description | UpdateMonitorGroupParamDetail 包含两字段 ✅ |
| 4b monitor_template | resource_zstack_monitor_template.go | name, description | UpdateMonitorTemplateParamDetail 包含两字段 ✅ |
| 4c preconfiguration_template | resource_zstack_preconfiguration_template.go | name, desc, distribution, type, content | UpdatePreconfigurationTemplateParamDetail 包含全部 5 字段 ✅ |
| 4d price_table | resource_zstack_price_table.go | name, description | UpdatePriceTableParamDetail 包含两字段 ✅ |

每个被移除 RequiresReplace 的字段都在 Update API 调用参数和 state 回写中得到完整处理。纯删除变更，无风险。

---

### #5 `feat/l2vxlan-network-update` — Needs Changes

**[CRITICAL] vni 属性缺少 RequiresReplace — 静默数据丢失**
- 位置: `resource_zstack_l2vxlan_network.go:90-94` (vni schema)
- 问题: `vni` 是 Optional+Computed 但无 RequiresReplace。Update 只发送 name/description。用户修改 vni 时:
  1. Terraform 调用 Update（不是 replace）
  2. API 忽略 vni 变更
  3. Read-back 用旧值覆盖 plan
  4. Plan 显示成功但实际 vni 未变
- 对比: 兄弟资源 `l2vlan_network.vlan` 正确标记了 `RequiresReplace`；`pool_uuid`、`zone_uuid` 也有
- 修复:
  ```go
  "vni": schema.Int64Attribute{
      Optional:    true,
      Computed:    true,
      Description: "The VXLAN Network Identifier (VNI).",
      PlanModifiers: []planmodifier.Int64{
          int64planmodifier.RequiresReplace(),
      },
  },
  ```

**[MEDIUM] Update 中读取完整 state 但只用 UUID**

**[LOW] 无验收测试覆盖 Update 路径**

**正面**: API 选择正确 (UpdateL2Network)、read-back 字段完整匹配 Read 方法

---

### #6 `refactor/read-error-handling` — Needs Changes

**[CRITICAL] stateCheckAlarmDisappears 未定义 — 整个 test package 编译失败**
- 位置: `resource_zstack_alarm_test.go:77` 调用 `stateCheckAlarmDisappears`
- 问题: 该函数定义在 #7 分支的 `check_disappears_test.go` 中，但 #6 分支没有
- 后果: **所有 provider 测试都无法编译和运行**
- 修复: 将 `stateCheckAlarmDisappears` 定义加入 `check_disappears_test.go`
- 或: 调整合入顺序 — #7 必须先于 #6 合入

**[MEDIUM] alarm Update 有 nil-client guard，���他 4 个资源没有**
- 预存不一致，非阻塞

**核心重构正确**:
- 5 个资源全部遵循 Update → 忽略返回 → Get/Query → 写 state 模式
- alarm 正确使用 `findResourceByQuery`（SDK 无 GetAlarm）
- 字段赋值列表与 master 完全一致
- metric_name 修正合理

**验收测试**: 全部 BLOCKED（编译失败）

---

### #7 `fix/alarm-update-and-test` — Approve (comment)

**[MEDIUM] SDK 版本 pinning 缺失**
- 文档引用 SDK v0.0.4 但无 go.mod 约束或 TODO 注释
- 建议: 在 `resource_zstack_alarm.go` workaround 处加 `// TODO(sdk-bug)` 注释

**[LOW] 文档 grep 命令路径含省略号，不可直接复制使用**

**正面**:
- SDK Bug 文档质量高（HTTP trace + 源码定位 + 多环境验证 + 补救方案）
- `stateCheckAlarmDisappears` 模式与现有 6 个函数完全一致

---

### #8 `refactor/provider-quality-hardening` — Needs Changes

**[CRITICAL] Expunge 失败��止 Terraform 清除 state**
- 位置: `resource_zstack_iam2_project.go:236-242`
- 问题: Delete 成功 → Expunge 失败 → `return` 阻止 state 清理 → 资源在 ZStack 已删除但 TF state 中残留 → 后续 destroy 也失败（资源不存在）
- `ExpungeIAM2Project` 非幂等（UUID 不存在时 PUT 会报错）
- 修复:
  ```go
  if err := r.client.ExpungeIAM2Project(state.Uuid.ValueString()); err != nil {
      resp.Diagnostics.AddWarning(
          "Warning expunging IAM2 Project",
          "Delete succeeded but expunge failed (name may not be immediately reusable): "+err.Error(),
      )
      // 不 return — Delete 已成功，让 Terraform 正常清除 state
  }
  ```

**[MEDIUM] 与 Instance/Image 的 expunge 模式不一致**
- Instance/Image 用 `expunge` boolean 让用户选择；IAM2Project 无条件 expunge
- 建议加注释解释: name 唯一性约束要求必须 expunge

**[MEDIUM] Chore 变更混入 fix 分支**
- 已在独立 commit，可接受

**验收测试**: FAIL (503，master 上同样失败)

---

## 合入顺序建议

```
#4a → #4b → #4c → #4d    (独立，可并行合入)
       ↓
#7 → #6                   (#6 依赖 #7 的 stateCheckAlarmDisappears)
       ���
#1 (修复 description 后)
#2 (命名可选修复)
#3 (独立)
#5 (修复 vni RequiresReplace 后)
#8 (修复 Expunge error handling 后)
```

## 待修复项清单

| 分支 | 严重度 | 修复项 | 工作量 |
|------|--------|--------|--------|
| #1 | CRITICAL | Create 后设置 description 为空字符串 | 1 行 |
| #5 | CRITICAL | vni 加 RequiresReplace | 5 行 |
| #6 | CRITICAL | 加 stateCheckAlarmDisappears 或调整合入顺序 | 6 行 |
| #8 | CRITICAL | Expunge 失败改 AddWarning 不 return | 5 行 |
| #1 | HIGH | Schema description 文档化 god-mode 默认 | 1 行 |
| #3 | MEDIUM | 后续全局修复 Optional+Computed IsNull guard | 新分支 |
| #2 | MEDIUM | 函数重命名 importStateIdFrom* | 可选 |

## 基础设施问题

多个验收测试因 ZStack API 503 失败（非代码问题）:
- `TestAccPortForwardingRuleResource`
- `TestAccIAM2ProjectResource`
- `TestAccIAM2ProjectResource_disappears`

`TestAccGlobalConfigResource` 的 "inconsistent result after apply" 是 master 上的预存 bug。
