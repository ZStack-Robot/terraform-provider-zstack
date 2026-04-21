# Bug 状态跟踪

> **项目**: Terraform Provider ZStack (zstack-sdk-go-v2)
> **最后更新**: 2026-04-21
> **统计**: OPEN 10 | FIXING 3 | FIXED 2 | WONTFIX 0

---

## 总览

| ID | 严重度 | 资源 | 资源生命周期 | 一句话描述 | 状态 | 负责人 | 关联 |
|----|--------|------|-------------|-----------|------|--------|------|
| [BUG-1](#bug-1) | **CRITICAL** | policy | Create | Create 后 description 为 Unknown | OPEN | 未分配 | Story-05 |
| [BUG-2](#bug-2) | **CRITICAL** | l2vxlan_network | Update | vni 缺少 RequiresReplace，静默数据丢失 | OPEN | 未分配 | Story-07 |
| [BUG-3](#bug-3) | **CRITICAL** | iam2_project | Delete | Expunge 逻辑缺失，name 不可复用 | FIXING | 未分配 | Story-04 |
| [BUG-4](#bug-4) | HIGH | policy | Create | Statements 空数组 + 修复引入 god-mode | FIXING | 未分配 | Story-05 |
| [BUG-5](#bug-5) | MEDIUM | 30+ 资源 (119 处) | Create/Update | Optional+Computed IsNull guard 不完整 | FIXING | 未分配 | Story-03 |
| [BUG-6](#bug-6) | MEDIUM | global_config | Read | 验收测试 inconsistent result | OPEN | 未分配 | — |
| [~~BUG-7~~](#bug-7) | ~~CRITICAL~~ | alarm | Disappears | stateCheckAlarmDisappears 未定义 | FIXED | — | Story-02 |
| [~~BUG-8~~](#bug-8) | ~~MEDIUM~~ | 4 资源 | Update | 错误的 RequiresReplace 导致不必要 destroy | FIXED | — | Story-07 |
| [BUG-9](#bug-9) | **CRITICAL** | instance | Update | UpdateVmInstance SDK 返回空 struct | OPEN | 未分配 | Story-11 |
| [BUG-10](#bug-10) | **HIGH** | instance | Create | stringPtr 传空字符串导致 API 校验失败 | OPEN | 未分配 | Story-11 |
| [BUG-11](#bug-11) | **CRITICAL** | vm_cdrom | Update | UpdateVmCdRom SDK 返回空 struct | OPEN | 未分配 | Story-11 |
| [BUG-12](#bug-12) | **MEDIUM** | instance_scripts | Create/Read | description/timeout/encoding_type 三处缺陷 | OPEN | 未分配 | Story-11 |
| [BUG-13](#bug-13) | **HIGH** | l3network | Create | ipVersion=0 被 API 拒绝 | OPEN | 未分配 | Story-12 |
| [BUG-14](#bug-14) | **CRITICAL** | l3network | Update | UpdateL3Network 返回不完整 | OPEN | 未分配 | Story-12 |
| [BUG-15](#bug-15) | **HIGH** | l3network | Delete | DeleteL3Network URL 缺少 UUID | OPEN | 未分配 | Story-12 |
| [BUG-16](#bug-16) | **MEDIUM** | subnet_ip_range | Read | ip_range_type Read 后丢失 | OPEN | 未分配 | Story-12 |

---

## Bug 详情

<a id="bug-1"></a>
### BUG-1 [CRITICAL] policy — Create 后 description 为 Unknown

| 字段 | 值 |
|------|-----|
| **严重度** | CRITICAL |
| **状态** | OPEN |
| **资源** | `zstack_policy` |
| **资源生命周期阶段** | Create |
| **发现方式** | 验收测试 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | Story-05 |
| **关联分支** | 未创建 |
| **阻塞项** | TestAccPolicyResource 失败 |

**现象**

用户创建 policy 时省略 description，`terraform apply` 成功创建资源但报错：
"provider still indicated an unknown value for zstack_policy.test.description"

**根因**

- **位置**: `resource_zstack_policy.go` Create 方法
- **分析**: description 为 `Optional+Computed` + `UseStateForUnknown()`，用户省略时 plan 阶段值为 Unknown；Create 方法写回 state 时未设置 description，Unknown 值残留
- **对比**: 其他资源（如 alarm）在 Create 中显式设置空字符串

**修复思路**

1. 在 Create 方法写回 state 前，检查 `plan.Description` 是否为 Unknown 或 Null
2. 若是，设置为空字符串 `types.StringValue("")`
3. 确保 API 返回的 description 优先使用（如果 API 有返回）

**验证方法**

```bash
TF_ACC=1 go test ./zstack/provider/ -run TestAccPolicyResource -v -timeout 30m
```

---

<a id="bug-2"></a>
### BUG-2 [CRITICAL] l2vxlan_network — vni 缺少 RequiresReplace

| 字段 | 值 |
|------|-----|
| **严重度** | CRITICAL |
| **状态** | OPEN |
| **资源** | `zstack_l2vxlan_network` |
| **资源生命周期阶段** | Update |
| **发现方式** | 代码审查 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | Story-07 |
| **关联分支** | 未创建 |
| **阻塞项** | 无（静默错误，不阻塞测试） |

**现象**

用户修改 l2vxlan_network 的 vni 值，`terraform apply` 显示成功，但实际 vni 未变更。API 忽略了 vni 变更，Read 回写用旧值覆盖 plan。用户无感知数据丢失。

**根因**

- **位置**: `resource_zstack_l2vxlan_network.go:90-94` vni schema 定义
- **分析**: vni 为 `Optional+Computed` 但无 `RequiresReplace`，Terraform 认为可以 in-place update；但 Update API（UpdateL2Network）只发送 name/description，忽略 vni
- **对比**: 兄弟资源 `l2vlan_network.vlan` 正确标记了 RequiresReplace；同文件 `pool_uuid`、`zone_uuid` 也有

**修复思路**

1. 在 vni 的 schema 定义中添加 `PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()}`
2. 用户修改 vni 时 Terraform 将自动执行 destroy + recreate 而非 in-place update
3. 与 pool_uuid、zone_uuid 的处理方式保持一致

**验证方法**

```bash
# 编译验证
go build ./...

# 行为验证：修改 vni 后 terraform plan 应显示 "must be replaced"
# terraform plan（检查输出包含 forces replacement）
```

---

<a id="bug-3"></a>
### BUG-3 [CRITICAL] iam2_project — Expunge 逻辑缺失

| 字段 | 值 |
|------|-----|
| **严重度** | CRITICAL |
| **状态** | FIXING（待审查） |
| **资源** | `zstack_iam2_project` |
| **资源生命周期阶段** | Delete |
| **发现方式** | 代码审查 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | Story-04 |
| **关联分支** | `refactor/provider-quality-hardening` @ `9c30736`（已推送 myorigin） |
| **阻塞项** | 重复 apply/destroy 测试会因 name 冲突失败 |

**现象**

`terraform destroy` 成功删除 iam2_project，但无法立即用相同 name 创建新项目。ZStack 要求 Expunge 后才释放 name 唯一性约束。

**根因**

- **位置**: `resource_zstack_iam2_project.go` Delete 方法
- **分析**: Delete 仅调用 `DeleteIAM2Project`，未调用 `ExpungeIAM2Project`。资源进入"已删除未彻底清除"状态
- **对比**: Instance/Image 资源通过 `expunge` boolean 让用户选择；IAM2Project 因 name 唯一性约束应无条件 Expunge

**修复思路**

1. 在 Delete 方法中，`DeleteIAM2Project` 成功后追加 `ExpungeIAM2Project` 调用
2. Expunge 失败时使用 `AddWarning` 而非 return error — Delete 已成功，不阻止 Terraform 清除 state
3. Warning 消息中说明 "name may not be immediately reusable"
4. 添加注释解释为什么 IAM2Project 必须无条件 Expunge（name 唯一性约束）

**修复进度**

- [x] 代码修改 — Delete 方法追加 `ExpungeIAM2Project` 调用
- [ ] **代码审查发现缺陷** — Expunge 失败时使用了 `AddError`+`return`，应改为 `AddWarning` 不 return（Delete 已成功，不应阻止 Terraform 清除 state）
- [ ] 验收测试通过
- [ ] 提交审查

**验证方法**

```bash
TF_ACC=1 go test ./zstack/provider/ -run TestAccIAM2ProjectResource -v -timeout 30m
# 验证：连续两次 apply+destroy 同名 project，第二次不应报 name 冲突
```

---

<a id="bug-4"></a>
### BUG-4 [HIGH] policy — Create Statements 处理缺陷

| 字段 | 值 |
|------|-----|
| **严重度** | HIGH |
| **状态** | FIXING |
| **资源** | `zstack_policy` |
| **资源生命周期阶段** | Create |
| **发现方式** | 验收测试 + 代码审查 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | Story-05 |
| **关联分支** | `fix/policy-empty-statements` @ `bfb4b9a`（已推送 myorigin） |
| **阻塞项** | 分支存在审查缺陷，需修正后方可合入 |

**现象**

policy Create 存在两层问题：

1. **原始问题**：master 上 Create 发送空 `Statements: []`，ZStack API 要求非空，导致 Create 完全不可用（测试已用 `t.Skip` 标注）
2. **审查发现**：`fix/policy-empty-statements` 分支修复空数组时，硬编码了 `Effect: "Allow", Actions: ["**"]`（god-mode），用户无法控制 statements 内容，存在安全风险

**根因**

- **位置**: `resource_zstack_policy.go:120-126`
- **原始问题**：Create 方法传空 Statements 数组，API 拒绝
- **审查缺陷**：修复分支用 god-mode 默认值替代空数组，解决了 API 拒绝但引入全权限安全风险。schema 中无 statements 属性供用户自定义

**修复思路**

1. **最低修复**：保留默认 statement，在 policy schema Description 中文档化此默认行为，明确告知用户
2. **完整修复**（后续）：将 statements 作为 schema 属性暴露，支持用户自定义 Effect/Actions
3. 在 `fix/policy-empty-statements` 分支上完成，修正后合入

**修复进度**

- [x] 空数组问题代码修复（`bfb4b9a`）
- [ ] 审查缺陷修正 — god-mode 默认值需文档化或暴露 statements 属性
- [ ] 验收测试通过
- [ ] 提交审查

**验证方法**

```bash
TF_ACC=1 go test ./zstack/provider/ -run TestAccPolicyResource -v -timeout 30m
```

---

<a id="bug-5"></a>
### BUG-5 [MEDIUM] 30+ 资源 (119 处) — Optional+Computed IsNull guard 不完整

| 字段 | 值 |
|------|-----|
| **严重度** | MEDIUM |
| **状态** | FIXING（port_forwarding_rule 待审查，其余待修复） |
| **资源** | 全量扫描发现 30+ 个资源文件共 119 处（见下方清单） |
| **资源生命周期阶段** | Create / Update |
| **发现方式** | 代码审查 + 全量 grep 扫描 |
| **发现日期** | 2026-04-17（初次发现）/ 2026-04-21（全量扫描） |
| **负责人** | 未分配 |
| **关联 Story** | [Story-03](test-status-overview.md#story-03) |
| **关联分支** | `fix/port-forwarding-rule` @ `cfe1d76`（已推送 myorigin） |
| **阻塞项** | 大部分偶然安全（String 零值=""），Int64/Bool 零值可能触发 API 错误（如 BUG-13） |

**现象**

Optional+Computed 字段在 Create/Update 方法中只检查 `IsNull()` 未检查 `IsUnknown()`。当 Unknown 时 `ValueString()` 返回 `""`（偶然安全），但 `ValueInt64()` 返回 `0`、`ValueBool()` 返回 `false`，可能导致 API 校验失败或语义错误。

**根因**

- **分析**: Optional+Computed 字段在用户未设置时 plan 值为 Unknown 而非 Null，只检查 `IsNull()` 会漏掉 Unknown 情况
- **对比**: `port_forwarding_rule` 的 3 个 Int64 字段已在 `fix/port-forwarding-rule` 分支正确修复

**受影响资源全量清单**（`grep -rn 'IsNull()' --include='resource_zstack_*.go'` 对比 `IsUnknown` 扫描）:

| 资源文件 | 受影响处数 | 关键字段 | 风险等级 |
|---------|:--------:|---------|---------|
| `resource_zstack_instance.go` | 28 | RootDisk, DataDisks, Strategy, GPUDevices, NeverStop 等 | 高（Int64/Bool 零值） |
| `resource_zstack_port_forwarding_rule.go` | 6 | VipPortEnd, PrivatePortStart/End (Int64) | 中（已修复） |
| `resource_zstack_primary_storage.go` | 9 | Description, Url, PoolName 等 | 低（String 零值偶然安全） |
| `resource_zstack_networking_secgroup_rule.go` | 6 | Description, DestinationPortRanges | 低 |
| `resource_zstack_l3network.go` | 3 | IpVersion (Int64), System (Bool) | **高**（触发 BUG-13） |
| `resource_zstack_instance_scripts.go` | 5 | ScriptTimeout (Int64), Description | **高**（触发 BUG-12b） |
| `resource_zstack_volume.go` | 6 | DiskSize (Int64), PrimaryStorageUuid | 中 |
| `resource_zstack_vip_qos.go` | 6 | Port, Bandwidth (Int64) | 中 |
| `resource_zstack_policy_route_rule.go` | 5 | DestIp, SourceIp, Port 等 | 中 |
| 其他 17 个资源 | 45 | 多为 Description (String) | 低 |

**修复思路**

1. port_forwarding_rule 随 Story-03 分支合入即可
2. **优先修复 Int64/Bool 类型字段**（零值有实际语义），String 类型可批量处理
3. 统一模式：`if !field.IsNull() && !field.IsUnknown()` 再取值

**修复进度**

- [x] port_forwarding_rule 代码修改（3 处）
- [ ] l3network、instance_scripts 代码修改（8 处）
- [ ] instance.go 高风险字段（28 处）
- [ ] 其余资源批量修复（75 处）

**验证方法**

```bash
# 全量扫描命令（验证修复完整性）
grep -rn 'IsNull()' zstack/provider/resource_zstack_*.go | grep -v '_test.go' | grep -v 'IsUnknown'

# port_forwarding_rule 验收测试
TF_ACC=1 go test ./zstack/provider/ -run TestAccPortForwardingRuleResource -v -timeout 30m
```

---

<a id="bug-6"></a>
### BUG-6 [MEDIUM] global_config — Create 后 state 全空（SDK PutWithSpec responseKey 缺失）

| 字段 | 值 |
|------|-----|
| **严重度** | MEDIUM |
| **状态** | OPEN |
| **资源** | `zstack_global_config` |
| **资源生命周期阶段** | Create |
| **发现方式** | 验收测试 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | — |
| **关联分支** | 未创建 |
| **阻塞项** | TestAccGlobalConfigResource 失败（不阻塞其他测试） |

**现象**

`TestAccGlobalConfigResource` 报错 "Provider produced inconsistent result after apply"。master 上同样失败，非特定分支引入。

**根因**

- **位置**: SDK `global_config_actions.go:17` + Provider `resource_zstack_global_config.go:126`
- **分析**:
  1. `UpdateGlobalConfig` SDK 方法调用 `PutWithSpec("v1/global-configurations", category, name+"/actions", "", params, &resp)`
  2. 第 4 个参数 `responseKey = ""`（空字符串）
  3. API 实际响应格式为 `{"inventory": {"name": "...", "value": "...", ...}}`（见 `UpdateGlobalConfigEventView` 定义）
  4. `PutWithSpec` 在 `responseKey == ""` 时直接反序列化整个 response body 到 `GlobalConfigInventoryView`
  5. 顶层 key 是 `"inventory"`，而 struct 字段期望 `"name"`, `"value"` 等 → 全部字段为零值
  6. Provider Create 方法 L135-139 用零值写回 state → `value = ""` 而 plan 期望 `"Delay"` → "inconsistent result"
- **对比**: 与 BUG-9/11/14 的 `PutWithRespKey` 问题同源 — SDK 响应解析层统一缺少 responseKey

**不一致字段**: `name`（`""` vs `"deletionPolicy"`）、`category`（`""` vs `"vm"`）、`value`（`""` vs `"Delay"`）— 所有字段都因 SDK 返回空 struct 而零值

**修复思路**

> ⚠️ 本 Bug 属于 [系统性模式 1: SDK PutWithRespKey 返回空 struct](#模式-1-sdk-putwithrespkey--putwithspec-返回空-struct)，与 BUG-9/11/14 同根因，建议统一修复而非单独处理。

| 方案 | 说明 |
|------|------|
| A: 修 SDK `responseKey` | `UpdateGlobalConfig` 调用改为 `PutWithSpec(..., "inventory", ...)` |
| B: Provider re-read | Create 后用 `QueryGlobalConfig` 重新读取（与 BUG-9 模式一致） |
| C: Provider 保留 plan 值 | Create 中不覆盖 `name`/`category`/`value`（Required 字段不应从 API 回写） |

**测试绕过**: 无（Create 必然失败）

**最小复现 HCL**:
```hcl
resource "zstack_global_config" "test" {
  category = "vm"
  name     = "deletionPolicy"
  value    = "Delay"
}
```

**验证方法**

```bash
TF_ACC=1 go test ./zstack/provider/ -run TestAccGlobalConfigResource -v -timeout 30m
```

---

<a id="bug-7"></a>
### ~~BUG-7~~ [FIXED] alarm — stateCheckAlarmDisappears 未定义

| 字段 | 值 |
|------|-----|
| **严重度** | ~~CRITICAL~~ |
| **状态** | FIXED |
| **资源** | `zstack_alarm` |
| **资源生命周期阶段** | Disappears |
| **发现方式** | 编译错误 |
| **发现日期** | 2026-04-17 |
| **关联 Story** | Story-02 |
| **修复分支** | `fix/alarm-update-and-test` |
| **合入 commit** | `86dd7b3` |
| **修复日期** | 2026-04-20 |

**修复内容**

- Alarm Update 改用 read-after-write（SDK PutWithRespKey bug 返回空字段）
- 修正测试 metric_name：`CPUAverageUtilization` → `CPUAverageUsedUtilization`
- 新增 `stateCheckAlarmDisappears` 辅助函数和 disappears 验收测试

**已完成测试项**

| 测试 | 类型 | 说明 |
|------|------|------|
| `TestAlarmResource_Schema` | 单元测试 | 已有，修复后通过 |
| `TestAlarmResource_Metadata` | 单元测试 | 已有，修复后通过 |
| `TestAccAlarmResource` | 验收测试 | metric_name 修正后通过 |
| `TestAccAlarmResource_disappears` | 验收测试 | **新增** — 外部删除后 plan 检测 |
| `stateCheckAlarmDisappears` | 测试辅助 | **新增** — disappears 检查函数 |
| `testAccCheckAlarmDestroy` | 测试辅助 | **新增** — destroy 验证函数 |

---

<a id="bug-8"></a>
### ~~BUG-8~~ [FIXED] 4 资源 — 错误的 RequiresReplace 导致不必要 destroy

| 字段 | 值 |
|------|-----|
| **严重度** | ~~MEDIUM~~ |
| **状态** | FIXED |
| **资源** | `monitor_group`、`monitor_template`、`preconfiguration_template`、`price_table` |
| **资源生命周期阶段** | Update |
| **发现方式** | 代码审查 |
| **发现日期** | 2026-04-20 |
| **关联 Story** | Story-07 |
| **修复分支** | 4 个独立分支（见下） |
| **修复日期** | 2026-04-20 |

**修复内容**

name/description 等可 in-place update 的字段被错误标记 `RequiresReplace`，导致修改这些字段时 Terraform 执行不必要的 destroy+recreate。4 个资源均有 Update 方法，确认支持 in-place 更新（代码审查 APPROVE）。

| 分支 | 资源 | 移除的字段 |
|------|------|-----------|
| `fix/monitor-group-requires-replace` @ `4e0e080` | monitor_group | name, description |
| `fix/monitor-template-requires-replace` @ `30887e2` | monitor_template | name, description |
| `fix/preconfiguration-template-requires-replace` @ `8b54d6d` | preconfiguration_template | name, description, distribution, type, content |
| `fix/price-table-requires-replace` @ `0a888c7` | price_table | name, description |

> 以上 4 个分支均已推送 myorigin，待合入 master。

**待补充测试项**

这 4 个资源当前无验收测试，需在 Phase 1 Sprint A 中补充：

- [ ] `TestAccMonitorGroupResource` — Create + Update name/description 不触发 replace
- [ ] `TestAccMonitorTemplateResource` — Create + Update name/description 不触发 replace
- [ ] `TestAccPreconfigurationTemplateResource` — Create + Update 各字段不触发 replace
- [ ] `TestAccPriceTableResource` — Create + Update name/description 不触发 replace

---

<a id="bug-9"></a>
### BUG-9 [CRITICAL] instance — UpdateVmInstance SDK 返回空 struct

| 字段 | 值 |
|------|-----|
| **严重度** | CRITICAL |
| **状态** | OPEN |
| **资源** | `zstack_instance` |
| **资源生命周期阶段** | Update |
| **发现方式** | 验收测试 |
| **发现日期** | 2026-04-21 |
| **负责人** | 未分配 |
| **关联 Story** | Story-11 |
| **关联分支** | 未创建 |
| **阻塞项** | instance Update 测试无法 PASS |

**现象**

`terraform apply` 更新 instance name 后，Terraform 报 "Provider produced inconsistent result after apply"。6 个字段同时为零值：
```
.uuid: was cty.StringVal("c5e416dd..."), but now cty.StringVal("")
.name: was cty.StringVal("acc-test-instance-updated"), but now cty.StringVal("")
.cpu_num: was cty.NumberIntVal(1), but now cty.NumberIntVal(0)
.memory_size: was cty.NumberIntVal(300), but now cty.NumberIntVal(0)
.network_interfaces: was [...], but now null
.vm_nics: was [...], but now null
```

**根因**

- **位置**: `resource_zstack_instance.go` Update 方法 L988
- **分析**: `r.client.UpdateVmInstance(uuid, param)` 调用 SDK 的 `PutWithRespKey` 方法，SDK 解析响应时未能正确提取 `inventory` 字段，返回一个空 struct。Provider 用这个空 struct 的零值覆盖了 state
- **验证**: 通过 `QueryVmInstance` 确认 API 实际更新成功（name="acc-test-instance-updated", memorySize=314572800, cpuNum=1），证明是 SDK 解析问题而非 API 问题

**修复方案评估**

> ⚠️ 本 Bug 属于 [系统性模式 1: SDK PutWithRespKey 返回空 struct](#模式-1-sdk-putwithrespkey--putwithspec-返回空-struct)，与 BUG-6/11/14 同根因，建议统一修复而非单独处理。

| 方案 | 改动层 | 优点 | 缺点 |
|------|--------|------|------|
| A: Provider re-read after update | Provider | 快速、已有 alarm 先例 | 多一次 API 调用；不治本，每个受影响资源都要加 |
| B: 修复 SDK PutWithRespKey 解析 | SDK | 治本，一次修复所有资源 | 需 fork/PR SDK，周期长 |
| C: 混合 — Provider workaround + SDK issue | 两层 | 短期可用，长期可移除 workaround | 需要维护两层代码 |

**影响范围**: 同一个 `PutWithRespKey` 问题至少影响 instance、vm_cdrom、l3network、alarm、global_config 5 个资源的 Update/Create

**测试绕过**: 测试中跳过 Update step，仅测 Create + Import

**最小复现 HCL**:
```hcl
# Step 1: Create
resource "zstack_instance" "test" {
  name                   = "test-vm"
  image_uuid             = "xxx"
  instance_offering_uuid = "xxx"
  expunge                = true
  network_interfaces = [{ l3_network_uuid = "xxx", default_l3 = true }]
}
# Step 2: 修改 name = "test-vm-updated" → terraform apply → "inconsistent result"
```

---

<a id="bug-10"></a>
### BUG-10 [HIGH] instance — stringPtr 传空字符串导致 API 校验失败

| 字段 | 值 |
|------|-----|
| **严重度** | HIGH |
| **状态** | OPEN |
| **资源** | `zstack_instance` |
| **资源生命周期阶段** | Create |
| **发现方式** | 验收测试 |
| **发现日期** | 2026-04-21 |
| **负责人** | 未分配 |
| **关联 Story** | Story-11 |
| **关联分支** | 未创建 |
| **阻塞项** | HCL 不显式设 RootDiskOfferingUuid/Strategy 时 Create 失败 |

**现象**

创建 VM 时 API 报参数校验错误，`RootDiskOfferingUuid` 和 `Strategy` 字段收到空字符串 `""` 而非 null。

**根因**

- **位置**: `resource_zstack_instance.go` L729, L740
- **分析**: `stringPtr("")` 返回 `&""` (指向空字符串的指针)，而 API 期望这些 Optional 字段要么传有效值，要么不传（nil）。`stringPtr` 不区分空字符串和有效值
- **对比**: 同文件其他位置已使用 `stringPtrOrNil()` 正确处理

**修复思路**

`stringPtr()` → `stringPtrOrNil()`，空字符串时返回 nil。需检查是否有其他 Optional 字段存在同样问题。

**测试绕过**: HCL 中不设 `root_disk_offering_uuid` 和 `strategy`（不传这两个字段）

---

<a id="bug-11"></a>
### BUG-11 [CRITICAL] vm_cdrom — UpdateVmCdRom SDK 返回空 struct

| 字段 | 值 |
|------|-----|
| **严重度** | CRITICAL |
| **状态** | OPEN |
| **资源** | `zstack_vm_cdrom` |
| **资源生命周期阶段** | Update |
| **发现方式** | 验收测试 |
| **发现日期** | 2026-04-21 |
| **负责人** | 未分配 |
| **关联 Story** | Story-11 |
| **关联分支** | 未创建 |
| **阻塞项** | vm_cdrom Update 测试无法 PASS |

**现象**

`terraform apply` 更新 vm_cdrom name 后，state 中 uuid/name/vm_instance_uuid 等字段全部变为零值。与 BUG-9 (instance) 相同模式。

**根因**

- **位置**: `resource_zstack_vm_cdrom.go` Update 方法 L200
- **分析**: `r.client.UpdateVmCdRom()` SDK 返回空 struct（同一个 `PutWithRespKey` 解析问题）
- **与 BUG-9 同根因**: SDK `PutWithRespKey` 系统性问题

**错误日志**:
```
TestAccVmCdRomResource: Step 2/3 error: After applying this test step, the plan was not empty.
  # zstack_vm_cdrom.test will be updated in-place
  ~ resource "zstack_vm_cdrom" "test" {
      ~ device_id        = 0 -> (known after apply)
      ~ name             = "" -> "acc-test-cdrom-updated"
      ~ uuid             = "" -> (known after apply)
      ~ vm_instance_uuid = "" -> "..."
    }
```

**修复方案**

> ⚠️ 本 Bug 属于 [系统性模式 1: SDK PutWithRespKey 返回空 struct](#模式-1-sdk-putwithrespkey--putwithspec-返回空-struct)，与 BUG-6/9/14 同根因，建议统一修复而非单独处理。

同 BUG-9 方案评估，待统一决策。

**测试绕过**: 测试中跳过 Update step，仅测 Create + Import

**最小复现 HCL**:
```hcl
# 需先创建 Stopped VM
resource "zstack_vm_cdrom" "test" {
  name             = "test-cdrom"
  vm_instance_uuid = zstack_instance.vm.uuid
}
# 修改 name = "test-cdrom-updated" → terraform apply → plan not empty (state 全零值)
```

---

<a id="bug-12"></a>
### BUG-12 [MEDIUM] instance_scripts — description/timeout/encoding_type 三处缺陷

| 字段 | 值 |
|------|-----|
| **严重度** | MEDIUM |
| **状态** | OPEN |
| **资源** | `zstack_instance_scripts` (Terraform type: `zstack_script`) |
| **资源生命周期阶段** | Create / Read |
| **发现方式** | 验收测试 + 编译 |
| **发现日期** | 2026-04-21 |
| **负责人** | 未分配 |
| **关联 Story** | Story-11 |
| **关联分支** | 未创建 |
| **阻塞项** | script Create 和 Import 测试可能失败 |

**现象**

三个独立缺陷，均在 `resource_zstack_instance_scripts.go` 中：

1. **description "inconsistent result"**: HCL 不设 description 时，Create 返回 `""`，Terraform 期望 null。
   - **位置**: Schema 定义 L84-L88
   - **修复思路**: Schema 加 `Computed: true` + `UseStateForUnknown()`

2. **script_timeout Unknown→0→API 拒绝**: Optional+Computed 字段 plan 值为 Unknown，`ValueInt64()` 返回 0，API 要求 [1, 86400]。
   - **位置**: Create L166, Update L312
   - **修复思路**: `!IsNull()` → `!IsNull() && !IsUnknown()`
   - > ⚠️ 属于 [系统性模式 2: Optional+Computed 缺 IsUnknown 检查](#模式-2-optionalcomputed-缺-isunknown-检查)，与 BUG-5/13 同根因，建议随全量扫描统一修复。

3. **Import 后 encoding_type 丢失**: Read 方法未从 API 响应恢复 `encoding_type`，Import 后该字段为 null，verify 失败。
   - **位置**: Read L277
   - **修复思路**: 加 `state.EncodingType = types.StringValue(scripts.EncodingType)`
   - > ⚠️ 属于 [系统性模式 3: Read 方法字段遗漏](#模式-3-read-方法字段遗漏)，与 BUG-16 同根因，建议对所有资源做 Schema↔Read 字段对齐检查。

---

<a id="bug-13"></a>
### BUG-13 [HIGH] l3network — ipVersion=0 被 API 拒绝

| 字段 | 值 |
|------|-----|
| **严重度** | HIGH |
| **状态** | OPEN |
| **资源** | `zstack_l3network` |
| **资源生命周期阶段** | Create |
| **发现方式** | 验收测试 |
| **发现日期** | 2026-04-21 |
| **负责人** | 未分配 |
| **关联 Story** | Story-12 |
| **关联分支** | 未创建 |
| **阻塞项** | HCL 不设 ip_version 时 Create 失败 |
| **系统性模式** | BUG-5 (Optional+Computed IsNull guard) 的具体实例 |

**现象**

用户创建 L3 network 不设 `ip_version`（Optional+Computed 字段），API 报 "ipVersion=0 is invalid"。

**根因**

- **位置**: `resource_zstack_l3network.go` Create 方法 L184-187
- **传值路径**:
  ```go
  var ipVersion *int              // L184: nil 初始化
  if !plan.IpVersion.IsNull() {   // L185: Unknown ≠ Null → IsNull()=false → 进入分支
      ipVersion = intPtr(int(plan.IpVersion.ValueInt64()))  // L186: ValueInt64() 对 Unknown 返回 0 → intPtr(0)
  }
  ```
  1. `ip_version` 是 `Optional+Computed`，用户不设时 plan 值为 **Unknown**（非 Null）
  2. `IsNull()` 返回 `false`（Unknown ≠ Null）→ **错误进入赋值分支**
  3. `ValueInt64()` 对 Unknown 状态返回零值 `0`
  4. `intPtr(0)` 创建 `*int` 指向 `0` → API 收到 `ipVersion=0`
  5. API 仅接受 `4` 或 `6`，拒绝 `0`
- **正确行为**: Unknown 时应保持 `ipVersion = nil`（不传该参数），让 API 使用默认值

**修复思路**

> ⚠️ 本 Bug 属于 [系统性模式 2: Optional+Computed 缺 IsUnknown 检查](#模式-2-optionalcomputed-缺-isunknown-检查)，与 BUG-5/12b 同根因，建议随全量扫描统一修复。

1. `!plan.IpVersion.IsNull()` → `!plan.IpVersion.IsNull() && !plan.IpVersion.IsUnknown()`（与 BUG-5 统一模式）
2. 或在 Schema 中添加 `Default: int64default.StaticInt64(4)`（更彻底，消除 Unknown 状态）

**测试绕过**: HCL 显式设 `ip_version = 4`

**最小复现 HCL**:
```hcl
resource "zstack_l3network" "test" {
  name             = "test-l3"
  type             = "L3BasicNetwork"
  l2_network_uuid  = "xxx"
  category         = "Public"
  # 不设 ip_version → Unknown → 0 → API 拒绝
}
```

---

<a id="bug-14"></a>
### BUG-14 [CRITICAL] l3network — UpdateL3Network 返回不完整

| 字段 | 值 |
|------|-----|
| **严重度** | CRITICAL |
| **状态** | OPEN |
| **资源** | `zstack_l3network` |
| **资源生命周期阶段** | Update |
| **发现方式** | 验收测试 |
| **发现日期** | 2026-04-21 |
| **负责人** | 未分配 |
| **关联 Story** | Story-12 |
| **关联分支** | 未创建 |
| **阻塞项** | L3 network Update 测试无法 PASS |
| **系统性模式** | SDK PutWithRespKey 返回空 struct（与 BUG-9/11/6 同根因） |

**现象**

Update L3 network name 后 Terraform 报 "Provider produced inconsistent result"。`type`、`category`、`l2_network_uuid`、`zone_uuid`、`ip_version` 等字段从有值变为空/零。

**根因**

- **位置**: `resource_zstack_l3network.go` Update 方法 L319-340
- **SDK 调用**: `r.client.UpdateL3Network(uuid, p)` → SDK `PutWithRespKey("v1/l3-networks", ..., "inventory", ...)` → 返回空 struct
- **Provider 代码**: L328-340 直接用 `result.UUID`、`result.Name` 等零值覆盖 state
- **SDK Update 参数**: `UpdateL3NetworkParamDetail` 只包含 `Name`、`Description`、`DnsDomain`、`Category`、`System`（SDK `l3network_params.go:10-16`），不包含 `IpVersion`、`L2NetworkUuid`、`Type` — 但这与 bug 无关，关键是 SDK 返回的 struct 为空

**API 侧证据**（与 BUG-9 相同模式）:

| 验证方式 | 结果 |
|---------|------|
| SDK `UpdateL3Network` 返回值 | 空 struct（所有字段零值） |
| `QueryL3Network` / `GetL3Network` 查询 | 数据完整（name/type/category/l2_network_uuid 等均有值） |
| **结论** | API 实际更新成功，是 SDK 解析问题 |

**可用的 re-read 方法**:
- `GetL3Network(uuid)` — SDK `l3network_actions.go:19-25`
- `QueryL3Network(params)` — SDK `l3network_actions.go:14-17`
- Provider Read 方法已使用 `findResourceByQuery(r.client.QueryL3Network, ...)`

**修复思路**

> ⚠️ 本 Bug 属于 [系统性模式 1: SDK PutWithRespKey 返回空 struct](#模式-1-sdk-putwithrespkey--putwithspec-返回空-struct)，与 BUG-6/9/11 同根因，建议统一修复而非单独处理。

Update 后用 `findResourceByQuery(r.client.QueryL3Network, uuid)` 重新读取完整 state（与 Read 方法一致）。

**最小复现 HCL**:
```hcl
# Step 1: Create
resource "zstack_l3network" "test" {
  name            = "test-l3"
  type            = "L3BasicNetwork"
  l2_network_uuid = "xxx"
  category        = "Public"
  ip_version      = 4
}
# Step 2: Update name → "test-l3-updated" → "inconsistent result"
```

---

<a id="bug-15"></a>
### BUG-15 [HIGH] l3network — DeleteL3Network URL/参数传递异常

| 字段 | 值 |
|------|-----|
| **严重度** | HIGH |
| **状态** | OPEN |
| **资源** | `zstack_l3network` |
| **资源生命周期阶段** | Delete |
| **发现方式** | 测试残留资源清理 |
| **发现日期** | 2026-04-21 |
| **负责人** | 未分配 |
| **关联 Story** | Story-12 |
| **关联分支** | 未创建 |
| **阻塞项** | 测试清理可能不生效 |
| **SDK 版本** | `github.com/zstackio/zstack-sdk-go-v2 v0.0.4` |
| **SDK issue** | 未提交（SDK 仓库: github.com/zstackio/zstack-sdk-go-v2） |

**现象**

调用 `DeleteL3Network` 时删除行为异常，测试后残留资源未被清理。

**根因**

- **Provider 调用**: `resource_zstack_l3network.go:359`
  ```go
  err := r.client.DeleteL3Network(state.Uuid.ValueString(), param.DeleteModePermissive)
  ```
- **SDK 实现**: `l3network_actions.go:51-54`
  ```go
  func (cli *ZSClient) DeleteL3Network(uuid string, deleteMode param.DeleteMode) error {
      return cli.Delete("v1/l3-networks", uuid, string(deleteMode))
  }
  ```
- **SDK Delete 方法**: `http_client.go:408-410`
  ```go
  func (cli *ZSHttpClient) Delete(ctx context.Context, resource, resourceId, deleteMode string) error {
      return cli.DeleteWithSpec(ctx, resource, resourceId, "", fmt.Sprintf("deleteMode=%s", deleteMode), nil)
  }
  ```
- **URL 构造**: `getDeleteURL` → `getURL("v1/l3-networks", uuid, "")` → `http://host:port/v1/l3-networks/{uuid}?deleteMode=Permissive`
- **问题分析**: URL 路径本身正确（包含 uuid），但 `deleteMode` 通过 query string 传递而非 request body。ZStack API 对部分资源的删除操作可能期望 body 中的参数格式（如 `{"deleteL3Network": {"deleteMode": "Permissive"}}`），具体行为需抓包确认

**Workaround 路径**

1. **Provider 层**: 直接调用 HTTP client 构造正确的 Delete 请求（绕过 SDK）
   ```go
   // 伪代码
   cli.HttpDelete(fmt.Sprintf("v1/l3-networks/%s?deleteMode=Permissive", uuid))
   ```
2. **SDK 层**: 向 `zstackio/zstack-sdk-go-v2` 提 issue 修复 Delete 参数传递方式

**待确认**

- [ ] 抓包确认 API 期望的 Delete 参数格式（query string vs body）
- [ ] 确认是否影响其他资源的 Delete（统一使用 `cli.Delete` 方法）

---

<a id="bug-16"></a>
### BUG-16 [MEDIUM] subnet_ip_range — ip_range_type Read 后丢失

| 字段 | 值 |
|------|-----|
| **严重度** | MEDIUM |
| **状态** | OPEN |
| **资源** | `zstack_subnet_ip_range` |
| **资源生命周期阶段** | Read |
| **发现方式** | 验收测试 |
| **发现日期** | 2026-04-21 |
| **负责人** | 未分配 |
| **关联 Story** | Story-12 |
| **关联分支** | 未创建 |
| **阻塞项** | refresh 后 plan 永远 non-empty |

**现象**

每次 `terraform plan` 都显示 `ip_range_type` 需要变更（从 null → 用户设的值），因为该字段是 ForceNew，每次 refresh 后 plan 都建议 destroy+recreate。

**根因**

- **位置**: `resource_zstack_subnet_ip_range.go` Read 方法
- **分析**: Read 方法未从 API 响应中恢复 `ip_range_type` 字段，每次 refresh 后值变为 null

**修复思路**

在 Read 方法中加 `state.IpRangeType = types.StringValue(ipRange.IpRangeType)`

> ⚠️ 属于 [系统性模式 3: Read 方法字段遗漏](#模式-3-read-方法字段遗漏)，建议对所有资源做 Schema↔Read 字段对齐检查，一次性修复所有遗漏。

**测试绕过**: `ExpectNonEmptyPlan: true` + `ImportStateVerifyIgnore: ["ip_range_type"]`

---

## 系统性模式问题

> 以下 3 个模式是多个 Bug 的共同根因，修复时应统一处理而非逐个打补丁。

### 模式 1: SDK PutWithRespKey / PutWithSpec 返回空 struct

| 涉及 Bug | 资源 | SDK 方法 |
|---------|------|---------|
| BUG-6 | global_config | `PutWithSpec` (responseKey="") |
| BUG-9 | instance | `PutWithRespKey` |
| BUG-11 | vm_cdrom | `PutWithRespKey` |
| BUG-14 | l3network | `PutWithRespKey` |
| (已修) alarm | alarm | `PutWithRespKey` |

**根因**: SDK 的 Put 系列方法在解析 API 响应时未正确提取 `inventory` 字段，返回空 struct。API 响应格式为 `{"inventory": {...}}`，但 SDK 要么缺少 responseKey 要么解析失败。

**一次性修复路径**:
- **短期 (Provider)**: 每个受影响资源的 Update 方法添加 read-after-update（已有 alarm 先例）
- **长期 (SDK)**: 修复 `PutWithRespKey` 的 response 解析逻辑，或统一为调用方传入正确的 responseKey

### 模式 2: Optional+Computed 缺 IsUnknown 检查

| 涉及 Bug | 资源 | 关键字段 |
|---------|------|---------|
| BUG-5 | 30+ 资源 (119 处) | 各种 Optional+Computed 字段 |
| BUG-12b | instance_scripts | ScriptTimeout (Int64) |
| BUG-13 | l3network | IpVersion (Int64) |

**根因**: `Optional+Computed` 字段在用户未设置时 plan 值为 Unknown（非 Null）。代码只检查 `IsNull()` 会漏掉 Unknown，导致零值被发送到 API。

**一次性修复路径**:
```bash
# 全量扫描
grep -rn 'IsNull()' zstack/provider/resource_zstack_*.go | grep -v '_test.go' | grep -v 'IsUnknown'
# 统一改为: !field.IsNull() && !field.IsUnknown()
```
**优先级**: Int64/Bool 类型优先（零值有语义），String 类型可批量处理（零值 `""` 多数偶然安全）

### 模式 3: Read 方法字段遗漏

| 涉及 Bug | 资源 | 遗漏字段 |
|---------|------|---------|
| BUG-12c | instance_scripts | encoding_type |
| BUG-16 | subnet_ip_range | ip_range_type |

**根因**: Read 方法未回写所有 Schema 定义的字段。遗漏字段在 refresh 后变为 null，如果该字段是 ForceNew 则每次 plan 都建议 destroy+recreate。

**一次性修复路径**: 对每个资源做 Schema 字段 ↔ Read 方法赋值的对齐检查：
```bash
# 提取 Schema 字段名
grep -A2 '"[a-z_]*":' resource_zstack_xxx.go | grep schema
# 提取 Read 方法赋值
grep 'state\.\w* = ' resource_zstack_xxx.go
# 对比找遗漏
```

---

## 基础设施问题

> 非代码 bug，属环境或 ZStack API 层面限制。

| 问题 | 影响范围 | 状态 | 备注 |
|------|---------|------|------|
| ZStack API 503 间歇性不可达 | TestAccPortForwardingRuleResource, TestAccIAM2ProjectResource | 持续 | 环境: 172.24.248.129:8080 |
| 部分资源 Update API 返回 404 | 涉及 Update 测试的资源 | 持续 | ZStack 不支持，已在测试中标注跳过 |

---

## 其他已知问题

> 低优先级问题或当前暂不处理问题

| 问题 | 说明 | 状态 |
|------|------|------|
| alarm SDK需要修改 | 当前Terraform provider做workaround处理，修复分支：fix/alarm-update-and-test  | 待处理 |
| SDK PutWithRespKey/PutWithSpec 系统性 bug | Update/Create API 返回空 struct，影响 instance/vm_cdrom/l3network/alarm/global_config 共 5 个资源。详见[系统性模式 1](#模式-1-sdk-putwithrespkey--putwithspec-返回空-struct) | 待 SDK 修复 |


---

## 跨仓同步记录

> 记录主 provider 仓库中已完成、但不直接映射到本 QA tracker 编号体系的修复状态。

| 日期 | 来源 | 状态 | 说明 |
|------|------|------|------|
| 2026-04-21 | 主仓库 `_bmad-output/bug-tracker.md` / BUG-019 | FIXED | `resource_zstack_instance` 已将 `data_disks` 的 data volume UUID 持久化到 state，Delete 不再重新查询 VM，而是直接使用 state 中的磁盘 UUID 执行删除。由于本 QA tracker 当前 BUG-1~BUG-8 无对应编号项，此处以同步记录方式标注。 |


---

## 附录

### 严重度定义

| 等级 | 定义 | 本项目典型场景 |
|------|------|---------------|
| **CRITICAL** | 数据丢失或状态损坏，用户无感知 | state 不一致、RequiresReplace 缺失致静默丢失、Expunge 失败阻塞 state 清理 |
| **HIGH** | 功能不可用或安全风险 | Create/Delete 报错、权限策略默认全放开 |
| **MEDIUM** | 功能降级但有 workaround | 部分字段类型检查不完整、验收测试间歇失败 |
| **LOW** | 代码规范或文档缺失 | 命名不一致、注释缺失、文档未覆盖 |

### 资源生命周期阶段

| 阶段 | Terraform 操作 | 本项目常见 bug 模式 |
|------|---------------|-------------------|
| **Create** | `terraform apply` (新建) | Optional+Computed 字段 Unknown 未处理、默认值不合理 |
| **Read** | `terraform plan/refresh` | API 返回字段与 Schema 不匹配、Query vs Get 行为差异 |
| **Update** | `terraform apply` (变更) | RequiresReplace 缺失、read-after-write 未实现 |
| **Delete** | `terraform destroy` | Expunge 逻辑缺失、DeleteMode 选择不当 |
| **Import** | `terraform import` | composite ID 解析、ImportState 字段映射 |
| **Disappears** | 资源被外部删除后 plan | stateCheck 函数缺失、API 404 处理 |

### Bug 详情模板

新增 BUG 时复制以下模板，填入实际内容：

```markdown
<a id="bug-N"></a>
### BUG-N [严重度] 资源名 — 一句话描述

| 字段 | 值 |
|------|-----|
| **严重度** | CRITICAL / HIGH / MEDIUM / LOW |
| **状态** | OPEN / FIXING / FIXED / WONTFIX |
| **资源** | `zstack_xxx` |
| **资源生命周期阶段** | Create / Read / Update / Delete / Import / Disappears |
| **发现方式** | 验收测试 / 代码审查 / 编译错误 |
| **发现日期** | YYYY-MM-DD |
| **负责人** | 未分配 |
| **关联 Story** | Story-XX |
| **关联分支** | `fix/xxx` @ `commit`（本地未推送 / 已推送） 或 未创建 |
| **阻塞项** | 具体阻塞的测试或功能 |

**现象**

用户执行什么操作，期望什么结果，实际发生了什么。

**根因**

- **位置**: `resource_zstack_xxx.go` 具体方法
- **分析**: 根本原因说明
- **对比**: 其他资源的正确实现（如有）

**修复思路**

1. 具体修复步骤一
2. 具体修复步骤二
3. 具体修复步骤三

**修复进度**（FIXING 状态时填写）

- [x] 已完成项
- [ ] 待完成项

**验证方法**

​```bash
TF_ACC=1 go test ./zstack/provider/ -run TestAccXxxResource -v -timeout 30m
​```
```

> **注意**：新增 BUG 后需同步更新总览表和统计数字。
