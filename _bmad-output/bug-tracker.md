# Bug 状态跟踪

> **项目**: Terraform Provider ZStack (zstack-sdk-go-v2)
> **最后更新**: 2026-04-20
> **统计**: OPEN 6 | FIXING 0 | FIXED 1 | WONTFIX 0

---

## 总览

| ID | 严重度 | 资源 | 生命周期 | 一句话描述 | 状态 | 负责人 | 关联 |
|----|--------|------|----------|-----------|------|--------|------|
| BUG-1 | **CRITICAL** | policy | Create | Create 后 description 为 Unknown | OPEN | 未分配 | Story-05 |
| BUG-2 | **CRITICAL** | l2vxlan_network | Update | vni 缺少 RequiresReplace，静默数据丢失 | OPEN | 未分配 | Story-07 |
| BUG-3 | **CRITICAL** | iam2_project | Delete | Expunge 逻辑缺失，name 不可复用 | OPEN | 未分配 | Story-04 |
| BUG-4 | HIGH | policy | Create | 默认 statement 为全权限 god-mode | OPEN | 未分配 | Story-05 |
| BUG-5 | MEDIUM | 多资源 | Create | Optional+Computed IsNull guard 不完整 | OPEN | 未分配 | Story-03 |
| BUG-6 | MEDIUM | global_config | Read | 验收测试 inconsistent result | OPEN | 未分配 | — |
| ~~BUG-7~~ | ~~CRITICAL~~ | alarm | Disappears | stateCheckAlarmDisappears 未定义 | FIXED | — | Story-02 |

---

## OPEN — 待修复

<a id="bug-1"></a>
### BUG-1 [CRITICAL] policy — Create 后 description 为 Unknown

| 字段 | 值 |
|------|-----|
| **严重度** | CRITICAL |
| **资源** | `zstack_policy` |
| **生命周期阶段** | Create |
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
| **资源** | `zstack_l2vxlan_network` |
| **生命周期阶段** | Update |
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
| **资源** | `zstack_iam2_project` |
| **生命周期阶段** | Delete |
| **发现方式** | 代码审查 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | Story-04 |
| **关联分支** | `refactor/provider-quality-hardening`（待审查） |
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

**验证方法**

```bash
TF_ACC=1 go test ./zstack/provider/ -run TestAccIAM2ProjectResource -v -timeout 30m
# 验证：连续两次 apply+destroy 同名 project，第二次不应报 name 冲突
```

---

<a id="bug-4"></a>
### BUG-4 [HIGH] policy — 默认 statement 为全权限 god-mode

| 字段 | 值 |
|------|-----|
| **严重度** | HIGH |
| **资源** | `zstack_policy` |
| **生命周期阶段** | Create |
| **发现方式** | 代码审查 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | Story-05 |
| **关联分支** | 未创建 |
| **阻塞项** | 无（功能可用，但有安全风险） |

**现象**

通过 `zstack_policy` 资源创建的每个 policy 都自动获得 `Effect: "Allow", Actions: ["**"]`（全权限），用户无法在 Terraform 配置中控制 statements 内容。

**根因**

- **位置**: `resource_zstack_policy.go:120-126`
- **分析**: Create 方法硬编码了默认 statement，schema 中无 statements 属性供用户自定义

**修复思路**

1. **最低修复**：在 policy schema 的 Description 中文档化此默认行为，明确告知用户创建的 policy 为全权限
2. **完整修复**（后续）：将 statements 作为 schema 属性暴露，支持用户自定义 Effect/Actions

**验证方法**

```bash
# 最低修复：检查文档更新
go generate ./...
# 查看 docs/resources/policy.md 中是否包含默认行为说明
```

---

<a id="bug-5"></a>
### BUG-5 [MEDIUM] 多资源 — Optional+Computed IsNull guard 不完整

| 字段 | 值 |
|------|-----|
| **严重度** | MEDIUM |
| **资源** | `zstack_port_forwarding_rule`、`zstack_l3network`、`zstack_instance_scripts` |
| **生命周期阶段** | Create |
| **发现方式** | 代码审查 + 模式扫描 |
| **发现日期** | 2026-04-17 |
| **负责人** | 未分配 |
| **关联 Story** | [Story-03](test-status-overview.md#story-03) |
| **关联分支** | `fix/port-forwarding-rule`（port_forwarding_rule 已修复，待合入） |
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
| **资源** | `zstack_global_config` |
| **生命周期阶段** | Read |
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

## FIXING — 修复中

<!-- 当 bug 进入修复状态时，从 OPEN 移到此区域，格式同 OPEN，额外增加修复进度 -->

<!-- **修复进度** -->
<!-- - [ ] 代码修改 -->
<!-- - [ ] 单元测试通过 -->
<!-- - [ ] 验收测试通过 -->
<!-- - [ ] 提交审查 -->

---

## FIXED — 已修复

| ID | 资源 | 问题 | 修复分支 | 合入 commit | 日期 |
|----|------|------|---------|-------------|------|
| ~~BUG-7~~ | alarm | stateCheckAlarmDisappears 未定义致编译失败 | `fix/alarm-update-and-test` | `86dd7b3` | 2026-04-20 |

---

## WONTFIX — 不修复

| ID | 资源 | 问题 | 原因 |
|----|------|------|------|
| — | — | — | — |

---

## 基础设施问题（非代码 bug）

| 问题 | 影响范围 | 状态 | 备注 |
|------|---------|------|------|
| ZStack API 503 间歇性不可达 | TestAccPortForwardingRuleResource, TestAccIAM2ProjectResource | 持续 | 环境: 172.24.248.129:8080 |
| 部分资源 Update API 返回 404 | 涉及 Update 测试的资源 | 持续 | ZStack 不支持，已在测试中标注跳过 |

---

## 附录：分类说明

### 严重度定义

| 等级 | 定义 | 本项目典型场景 |
|------|------|---------------|
| **CRITICAL** | 数据丢失或状态损坏，用户无感知 | state 不一致、RequiresReplace 缺失致静默丢失、Expunge 失败阻塞 state 清理 |
| **HIGH** | 功能不可用或安全风险 | Create/Delete 报错、权限策略默认全放开 |
| **MEDIUM** | 功能降级但有 workaround | 部分字段类型检查不完整、验收测试间歇失败 |
| **LOW** | 代码规范或文档缺失 | 命名不一致、注释缺失、文档未覆盖 |

### 生命周期阶段

| 阶段 | Terraform 操作 | 本项目常见 bug 模式 |
|------|---------------|-------------------|
| **Create** | `terraform apply` (新建) | Optional+Computed 字段 Unknown 未处理、默认值不合理 |
| **Read** | `terraform plan/refresh` | API 返回字段与 Schema 不匹配、Query vs Get 行为差异 |
| **Update** | `terraform apply` (变更) | RequiresReplace 缺失、read-after-write 未实现 |
| **Delete** | `terraform destroy` | Expunge 逻辑缺失、DeleteMode 选择不当 |
| **Import** | `terraform import` | composite ID 解析、ImportState 字段映射 |
| **Disappears** | 资源被外部删除后 plan | stateCheck 函数缺失、API 404 处理 |
