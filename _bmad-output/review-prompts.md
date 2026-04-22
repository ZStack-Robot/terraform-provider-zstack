# 分支审查 Prompt 模板

> 按推荐合入顺序排列。每个 prompt 直接粘贴到新 Claude Code 会话中使用。
> 审查完成后在对应条目后标注 ✅ 或 ❌。

---

## 1. `fix/policy-empty-statements` — policy Create 空 statements Bug

```
审查分支 fix/policy-empty-statements (基于 master)。

背景：ZStack API 创建 policy 时要求 statements 非空，当前代码发送空数组导致 API 拒绝，测试被 t.Skip 跳过。

变更范围（2 文件, +7 -4）：
- resource_zstack_policy.go: Create 方法中 Statements 从空数组改为包含默认策略 {Name:"default", Effect:"Allow", Actions:["**"]}
- resource_zstack_policy_test.go: 移除 t.Skip

审查要点：
1. 默认 statement "Allow + **" 是全权限策略（god-mode）。这是否是安全的默认值？是否应该暴露 statements 为 schema 属性让用户自定义？至少需要在 schema description 中文档化这个默认行为
2. 只改了 Create，Update 方法仍返回 "not supported" — 一致性是否 OK？
3. ⚠️ 关键：t.Skip 只从 TestPolicyResource_Metadata（第54行，单元测试）移除了，但 TestAccPolicyResource（第57行，验收测试）仍有 t.Skip！验收测试才是真正验证 Create API 调用的测试，必须同时移除
4. 修复后是否应跑一次验收测试确认 API 接受新的 statements？

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 2. `test/composite-id-import-steps` — scheduler_job / global_config Import

```
审查分支 test/composite-id-import-steps (基于 master)。

背景：scheduler_job 和 global_config 之前缺 Import Step，属于测试计划 2.6 节已知问题。两个资源使用 composite ID 导入。

变更范围（2 文件, +51）：
- resource_zstack_scheduler_job_test.go: 新增 Import Step + importStateIdSchedulerJob 函数 (composite ID: uuid:type)
- resource_zstack_global_config_test.go: 新增 Import Step + importStateIdGlobalConfig 函数 (composite ID: category/name)

审查要点：
1. composite ID 格式是否与对应资源 ImportState 方法的解析逻辑一致？请交叉检查 resource_zstack_scheduler_job.go 和 resource_zstack_global_config.go 的 ImportState 方法
2. ImportStateVerifyIdentifierAttribute 设置是否正确（scheduler_job→"uuid", global_config→"name"）
3. 函数命名 importStateIdSchedulerJob / importStateIdGlobalConfig 是否与项目现有惯例一致（对比 importStateIdFromUUID）
4. 分隔符不一致：global_config 用 `/`，scheduler_job 用 `:`。这是有意的设计还是应该统一？
5. 缺少负面测试：畸形 import ID（错误分隔符、缺少组件）是否会产生清晰的错误信息？

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 3. `fix/port-forwarding-rule` — Optional+Computed Int64 Unknown 状态修复

```
审查分支 fix/port-forwarding-rule (基于 master)。

背景：port_forwarding_rule 的 vip_port_end、private_port_start、private_port_end 是 Optional+Computed Int64 属性。Terraform Plugin Framework 中这类属性在 plan 阶段为 Unknown（非 Null），原代码只检查 IsNull() 就发送 ValueInt64()，导致 Unknown 时发送 0。

变更范围（2 文件, +4 -4）：
- resource_zstack_port_forwarding_rule.go: 3 处 `!plan.X.IsNull()` 改为 `!plan.X.IsNull() && !plan.X.IsUnknown()`
- resource_zstack_port_forwarding_rule_test.go: 断言从 knownvalue.StringExact("8080") 改为 knownvalue.Int64Exact(8080)

审查要点：
1. 3 个字段是否全部覆盖？注意 vip_port_start 是 Required 不需要修复，但确认 vip_port_end / private_port_start / private_port_end 确实都是 Optional+Computed
2. 此 bug 模式是否存在于其他资源的 Optional+Computed Int64 属性？是否需要全局排查？
3. 测试断言类型修正 (String→Int64) 是否说明之前测试本身有问题？
4. VmNicUuid 检查了 `!IsNull() && ValueString() != ""` 但没检查 `IsUnknown()` — Optional String 属性是否有同样的 Unknown 风险？

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 4a. `fix/monitor-group-requires-replace` — monitor_group RequiresReplace 清理

```
审查分支 fix/monitor-group-requires-replace (基于 master)。

背景：commit eeaf3e8f (compliance audit) 批量给资源属性添加了 RequiresReplace，其中 monitor_group 的 name 和 description 不合理——该资源有可用的 UpdateMonitorGroup API。

变更范围（1 文件, -4 行，纯删除）：
- resource_zstack_monitor_group.go: 移除 name 的整个 PlanModifiers 块（含 RequiresReplace），移除 description 的 RequiresReplace（保留 UseStateForUnknown）

审查要点：
1. 交叉检查 Update 方法（第167-200行）：UpdateMonitorGroupParamDetail 是否包含 Name 和 Description 字段？
2. Update 方法是否正确读取 plan.Name 和 plan.Description 并传入 API？
3. 移除后 name 属性不再有任何 PlanModifiers — Required+无修饰符是否为正确配置？

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 4b. `fix/monitor-template-requires-replace` — monitor_template RequiresReplace 清理

```
审查分支 fix/monitor-template-requires-replace (基于 master)。

背景：同 4a，monitor_template 的 name 和 description 的 RequiresReplace 不合理——该资源有可用的 UpdateMonitorTemplate API。

变更范围（1 文件, -4 行，纯删除）：
- resource_zstack_monitor_template.go: 移除 name 的整个 PlanModifiers 块（含 RequiresReplace），移除 description 的 RequiresReplace（保留 UseStateForUnknown）

审查要点：
1. 交叉检查 Update 方法（第161-193行）：UpdateMonitorTemplateParamDetail 是否包含 Name 和 Description 字段？
2. Update 方法是否正确读取 plan.Name 和 plan.Description 并传入 API？
3. 移除后 name 属性不再有任何 PlanModifiers — 与 monitor_group 保持一致

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 4c. `fix/preconfiguration-template-requires-replace` — preconfiguration_template RequiresReplace 清理

```
审查分支 fix/preconfiguration-template-requires-replace (基于 master)。

背景：同 4a，preconfiguration_template 有 5 个属性的 RequiresReplace 不合理——该资源有可用的 UpdatePreconfigurationTemplate API。

变更范围（1 文件, -13 行，纯删除）：
- resource_zstack_preconfiguration_template.go: 移除 name, description, distribution, type, content 共 5 个属性的 RequiresReplace

审查要点：
1. ⚠️ 关键：交叉检查 Update 方法（第235-282行）的 UpdatePreconfigurationTemplateParamDetail — 是否**具体包含** Distribution、Type、Content 字段？（有 Update 方法不等于处理了所有字段）
2. distribution/type/content 是 Required 属性，移除 RR 后 Update 会尝试发送新值 — API 是否接受修改这些字段？（某些 API 可能只接受 name/description 的 Update）
3. description 移除 RR 后保留 UseStateForUnknown — 正确
4. name/distribution/type/content 移除后不再有任何 PlanModifiers — Required 属性无修饰符是否为正确配置？

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 4d. `fix/price-table-requires-replace` — price_table RequiresReplace 清理

```
审查分支 fix/price-table-requires-replace (基于 master)。

背景：同 4a，price_table 的 name 和 description 的 RequiresReplace 不合理——该资源有可用的 UpdatePriceTable API。

变更范围（1 文件, -4 行，纯删除）：
- resource_zstack_price_table.go: 移除 name 的整个 PlanModifiers 块（含 RequiresReplace），移除 description 的 RequiresReplace（保留 UseStateForUnknown）

审查要点：
1. 交叉检查 Update 方法（第172-209行）：UpdatePriceTableParamDetail 是否包含 Name 和 Description 字段？
2. Update 方法是否正确读取 plan.Name 和 plan.Description 并传入 API？
3. 移除后 name 属性不再有任何 PlanModifiers — 与 monitor_group/monitor_template 保持一致

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 5. `feat/l2vxlan-network-update` — l2vxlan_network Update 实现

```
审查分支 feat/l2vxlan-network-update (基于 master)。

背景：l2vxlan_network 的 Update 方法原来返回 "Update Not Supported" 错误，但 ZStack 有通用的 UpdateL2Network API 可用于更新 name/description。

变更范围（1 文件, +49 -4）：
- resource_zstack_l2vxlan_network.go: 替换 "not supported" 错误为实际 Update 实现，使用 read-after-write 模式

审查要点：
1. 使用通用 UpdateL2Network（而非 L2VxlanNetwork 专用接口）是否正确？L2VxlanNetwork 特有字段（vni, pool_uuid）是否可能被覆盖？
2. read-after-write 中使用 QueryL2VxlanNetwork 回读 — 与 Update 用的 UpdateL2Network 是否一致？
3. Update 后回读的字段列表是否完整？对比 Read 方法中的字段赋值
4. ⚠️ vni 是 Optional+Computed 且无 RequiresReplace，但 UpdateL2NetworkParam 只发送 name/description。如果用户修改 vni，Update 会静默忽略，read-back 用旧值覆盖 plan — 需要对 vni/pool_uuid/zone_uuid 等不可变属性加 ForceNew 或在 Update 中处理
5. 无验收测试覆盖新的 Update 路径

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 6. `refactor/read-error-handling` — Update read-after-write 模式 (核心改动)

```
审查分支 refactor/read-error-handling (基于 master, 2 commits)。

背景：ZStack SDK 多个 Update API 返回值不完整或为空（已确认 alarm 的 UpdateAlarm 返回空响应体）。此分支将 5 个资源的 Update 方法从"信赖 Update 返回值"改为"Update 后 Get/Query 回读"。

变更范围（6 文件, +92 -22）：
- resource_zstack_account.go: Update → 忽略返回值 → GetAccount 回读
- resource_zstack_affinity_group.go: Update → 忽略返回值 → GetAffinityGroup 回读
- resource_zstack_ssh_key_pair.go: Update → 忽略返回值 → GetSshKeyPair 回读
- resource_zstack_iam2_project.go: Update → 忽略返回值 → GetIAM2Project 回读
- resource_zstack_alarm.go: Update → 忽略返回值 → QueryAlarm 回读（alarm 无 Get 方法）
- resource_zstack_alarm_test.go: metric_name 从 CPUAverageUtilization 改为 CPUAverageUsedUtilization + 新增 disappears 测试

审查要点：
1. 模式一致性：5 个资源的改动结构是否一致（Update → 忽略返回 → Get/Query → 写 state）？
2. alarm 使用 findResourceByQuery 而非直接 Get — 是否因为 SDK 确实没有 GetAlarm？代码中缺少注释解释"为什么用 Query 不用 Get"，建议补充
3. 每个资源的 Get/Query 后字段赋值列表是否与改动前完全一致？（不应遗漏或新增字段）
4. alarm_test.go 的 metric_name 修正是否合理？CPUAverageUsedUtilization 是否为 ZStack 的正确指标名？旧值 CPUAverageUtilization 是否说明 master 上的测试本身就会失败？
5. 此模式是否应推广到更多资源？（目前只改了有验收测试的 5 个）
6. 新增的 Get/Query 失败路径（Update 成功但 re-read 失败）没有测试覆盖 — 是否可接受？

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 7. `fix/alarm-update-and-test` — alarm Disappears + SDK Bug 文档

> ⚠️ 此分支与 #6 在 alarm.go / alarm_test.go 上有完全重叠。合入 #6 后需 rebase 此分支。

```
审查分支 fix/alarm-update-and-test (基于 master, 2 commits)。

背景：在实现 alarm Update 过程中发现 ZStack SDK UpdateAlarm API 返回空响应。此分支记录了该 Bug 并为 alarm 补充了 Disappears 测试。

⚠️ 注意：alarm.go 和 alarm_test.go 的改动与分支 refactor/read-error-handling 完全重叠。审查重点放在该分支独有的增量上。

独有增量（2 文件）：
- docs/SDK-BUG-UpdateAlarm-Empty-Response.md (+299): SDK Bug 文档，含多环境验证结果
- check_disappears_test.go (+6): 新增 stateCheckAlarmDisappears 工厂函数

重叠部分（已在 #6 审查，此处跳过）：
- resource_zstack_alarm.go: read-after-write（同 #6）
- resource_zstack_alarm_test.go: metric_name 修正 + disappears 测试（同 #6）

审查要点：
1. SDK Bug 文档是否描述清晰？复现步骤、多环境验证、workaround 是否完整？
2. stateCheckAlarmDisappears 是否使用正确的删除方法和模式？
3. 合入策略：alarm.go 和 alarm_test.go 与 #6 **逐字节相同**，这意味着无论谁先合，另一个都会产生合并冲突（即使内容相同）。建议先合 #6，然后 rebase 此分支，冲突文件应全部取 master 版本，只保留 SDK 文档 + disappears 函数的增量
4. metric_name 修正 (CPUAverageUtilization → CPUAverageUsedUtilization) 同时出现在此分支和 #6 — 旧值是否在 master 上就已经是错的？
5. Bug 文档引用 SDK v0.0.4，但没有 go.mod 变更来 pin 版本 — SDK 升级后 workaround 可能过时或失效

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```

---

## 8. `refactor/provider-quality-hardening` — iam2_project Delete + Expunge

> ⚠️ 与 #6 在 iam2_project.go 上有交叉（#6 改 Update，此分支改 Delete）。合入 #6 后需 rebase。

```
审查分支 refactor/provider-quality-hardening (基于 master, 2 commits)。

背景：iam2_project 删除后进入回收站，名称未释放，导致再次创建同名项目报 duplicate name 错误。

功能变更（1 文件, +9）：
- resource_zstack_iam2_project.go: Delete 方法在 DeleteIAM2Project 后追加 ExpungeIAM2Project 调用

Chore 变更（忽略）：
- .omc/prd.json: 项目配置（非功能代码）
- _bmad/core/config.yaml: BMad 配置（非功能代码）
- generate_tf (binary deleted): 清理生成的二进制文件

审查要点：
1. ExpungeIAM2Project 是否幂等？（资源已不存在时调用是否报错）
2. ⚠️ 关键问题：如果 Delete 成功但 Expunge 失败，`return` 会阻止 Terraform 清除 state。后果是资源在 ZStack 已删除但 Terraform state 中仍存在 — 下次 apply 时会尝试删除不存在的资源。是否应该：(a) Expunge 失败仅 warning 不 return？(b) 或先 Expunge 再 Delete？
3. 是否有其他资源也需要 Delete + Expunge？（参考 E2E-TEST-PLAN.md 1.4 节：instance, volume, image 有回收站）
4. chore 文件问题：(a) .omc/prd.json 是项目管理元数据，不属于 provider 代码仓库 (b) generate_tf 二进制删除不会减小 repo 体积（仍在 git history 中）(c) 建议将 chore 变更剥离到单独 commit 或分支

请 diff 审查后给出：通过 / 需修改 / 拒绝，附具体意见。
```
