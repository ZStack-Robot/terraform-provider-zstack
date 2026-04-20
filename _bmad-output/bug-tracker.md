# Bug 状态跟踪

> **项目**: Terraform Provider ZStack (zstack-sdk-go-v2)
> **最后更新**: 2026-04-20
> **统计**: OPEN 3 | FIXING 3 | FIXED 2 | WONTFIX 0

---

## 总览

| ID | 严重度 | 资源 | 资源生命周期 | 一句话描述 | 状态 | 负责人 | 关联 |
|----|--------|------|-------------|-----------|------|--------|------|
| [BUG-1](#bug-1) | **CRITICAL** | policy | Create | Create 后 description 为 Unknown | OPEN | 未分配 | Story-05 |
| [BUG-2](#bug-2) | **CRITICAL** | l2vxlan_network | Update | vni 缺少 RequiresReplace，静默数据丢失 | OPEN | 未分配 | Story-07 |
| [BUG-3](#bug-3) | **CRITICAL** | iam2_project | Delete | Expunge 逻辑缺失，name 不可复用 | FIXING | 未分配 | Story-04 |
| [BUG-4](#bug-4) | HIGH | policy | Create | Statements 空数组 + 修复引入 god-mode | FIXING | 未分配 | Story-05 |
| [BUG-5](#bug-5) | MEDIUM | 多资源 | Create | Optional+Computed IsNull guard 不完整 | FIXING | 未分配 | Story-03 |
| [BUG-6](#bug-6) | MEDIUM | global_config | Read | 验收测试 inconsistent result | OPEN | 未分配 | — |
| [~~BUG-7~~](#bug-7) | ~~CRITICAL~~ | alarm | Disappears | stateCheckAlarmDisappears 未定义 | FIXED | — | Story-02 |
| [~~BUG-8~~](#bug-8) | ~~MEDIUM~~ | 4 资源 | Update | 错误的 RequiresReplace 导致不必要 destroy | FIXED | — | Story-07 |

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
| **关联分支** | `refactor/provider-quality-hardening` @ `153cb74`（本地未推送） |
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
| **关联分支** | `fix/policy-empty-statements` @ `bfb4b9a`（本地未推送） |
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
### BUG-5 [MEDIUM] 多资源 — Optional+Computed IsNull guard 不完整

| 字段 | 值 |
|------|-----|
| **严重度** | MEDIUM |
| **状态** | FIXING（port_forwarding_rule 待审查，l3network/instance_scripts 待修复） |
| **资源** | `zstack_port_forwarding_rule`、`zstack_l3network`、`zstack_instance_scripts` |
| **资源生命周期阶段** | Create |
| **发现方式** | 代码审查 + 模式扫描 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | [Story-03](test-status-overview.md#story-03) |
| **关联分支** | `fix/port-forwarding-rule` @ `cfe1d76`（本地未推送） |
| **阻塞项** | 无（偶然安全，Unknown 时 ValueXxx() 返回零值） |

**现象**

Optional+Computed 字段在 Create 方法中只检查 `IsNull()` 未检查 `IsUnknown()`。当前偶然安全（Unknown 时 `ValueString()` 返回空字符串），但与已修复模式不一致。

**根因**

- **位置**:
  - `resource_zstack_l3network.go` — IpVersion (Int64)、System (Bool)
  - `resource_zstack_instance_scripts.go` — ScriptTimeout (Int64)
- **分析**: Optional+Computed 字段在用户未设置时 plan 值为 Unknown 而非 Null，只检查 IsNull() 会漏掉 Unknown 情况
- **对比**: `port_forwarding_rule` 的 3 个 Int64 字段已在 `fix/port-forwarding-rule` 分支正确修复

**修复思路**

1. port_forwarding_rule 随 Story-03 分支合入即可
2. l3network、instance_scripts 后续创建 `fix/optional-computed-isnull-guard` 分支
3. 统一模式：`if !field.IsNull() && !field.IsUnknown()` 再取值

**修复进度**

- [x] port_forwarding_rule 代码修改 — 3 个 Int64 字段 `IsNull()` → `!IsNull() && !IsUnknown()`（代码审查 APPROVE）
- [ ] port_forwarding_rule 验收测试通过
- [ ] l3network、instance_scripts 代码修改
- [ ] 提交审查

**验证方法**

```bash
# 编译验证
go build ./...

# port_forwarding_rule 验收测试
TF_ACC=1 go test ./zstack/provider/ -run TestAccPortForwardingRuleResource -v -timeout 30m
```

---

<a id="bug-6"></a>
### BUG-6 [MEDIUM] global_config — 验收测试 inconsistent result

| 字段 | 值 |
|------|-----|
| **严重度** | MEDIUM |
| **状态** | OPEN |
| **资源** | `zstack_global_config` |
| **资源生命周期阶段** | Read |
| **发现方式** | 验收测试 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | — |
| **关联分支** | 未创建 |
| **阻塞项** | TestAccGlobalConfigResource 失败（不阻塞其他测试） |

**现象**

`TestAccGlobalConfigResource` 报错 "Provider produced inconsistent result after apply"。master 上同样失败，非特定分支引入。

**根因**

- **位置**: `resource_zstack_global_config.go` Read 方法（待排查）
- **分析**: Apply 后 Read 返回的 state 与 plan 不一致，可能是某个 Computed 字段的值与 API 返回不匹配
- **对比**: 待排查具体不一致字段

**修复思路**

1. 运行 `TestAccGlobalConfigResource` 获取完整错误输出，定位不一致的字段名
2. 对比 Read 方法中该字段的赋值逻辑与 Schema 定义
3. 修复 Read 返回值或调整 Schema 的 Optional/Computed 标记

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

> 以上 4 个分支均为本地未推送，待合入 master。

**待补充测试项**

这 4 个资源当前无验收测试，需在 Phase 1 Sprint A 中补充：

- [ ] `TestAccMonitorGroupResource` — Create + Update name/description 不触发 replace
- [ ] `TestAccMonitorTemplateResource` — Create + Update name/description 不触发 replace
- [ ] `TestAccPreconfigurationTemplateResource` — Create + Update 各字段不触发 replace
- [ ] `TestAccPriceTableResource` — Create + Update name/description 不触发 replace

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
