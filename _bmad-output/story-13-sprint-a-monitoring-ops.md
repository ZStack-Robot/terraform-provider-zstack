# Story 13: Sprint A — 监控运维和备份资源验收测试 (3 个)

> **分支**: 待创建 (`test/sprint-a-monitoring-ops`)  
> **状态**: 未开始  
> **优先级**: P1 — Phase 1 Sprint A  
> **前置依赖**: story-07 (monitor_group RequiresReplace 清理); story-10 (monitor_template 测试模式参考)  
> **预计工作量**: 1 天

---

## Story

作为 Provider 维护者，我需要为监控运维和数据保护相关的资源新增验收测试，完成 Sprint A 最后一组。

## 资源清单

| # | Resource | Terraform Type Name | Update 状态 | 目标级别 | SDK Read 方式 | 环境依赖 |
|---|----------|-------------------|-----------|---------|-------------|---------|
| 1 | `monitor_group` | `zstack_monitor_group` | 有 (UpdateMonitorGroup) | 标准 | Query (`QueryMonitorGroup`) | monitor_template_uuid (需在测试中创建) |
| 2 | `cdp_task` | `zstack_cdp_task` | 有 (UpdateCdpTask) | 标准 | Query (`QueryCdpTask`) | policy_uuid + backup_storage_uuid + resource_uuids (依赖链长) |
| 3 | `backup_storage` | `zstack_backup_storage` | 有 (UpdateBackupStorage + zone attach/detach) | 标准 | Get (`GetBackupStorage`) | zone_uuid |

### 已知问题

- **monitor_group RequiresReplace 未清理** (依赖 story-07): `name` 和 `description` 仍有 RequiresReplace。story-07 完成前，Update Step 实际触发 destroy+recreate。
- **cdp_task Required 属性远多于原始 Story 描述**: 实际有 5 个 Required 字段，且 4 个 ForceNew（见下方详细说明）。

### env.json 数据可用性

| env.json 字段 | 数据量 | 说明 |
|--------------|:------:|------|
| `monitor_groups` | NULL | 需在测试中 Create |
| `monitor_templates` | NULL | 需在测试中 Create（monitor_group 的依赖）|
| `backup_storages` | 1 | 可用（cdp_task 和 backup_storage 使用）|
| `zones` | 1 | 可用（backup_storage zone attach 使用）|
| `volumes` | 5 | 可用（cdp_task 的 resource_uuids 使用）|

## 关键风险

- **monitor_group**: 依赖 monitor_template，而 monitor_template API 在测试环境返回 404。**需用户确认环境是否支持 monitor_template API**。
- **cdp_task**: 依赖链长（cdp_policy + backup_storage + volume），且 cdp_policy Create 在测试环境返回 503。**需用户确认 CDP 服务状态**。
- **backup_storage**: 创建需要指定 `type`（SftpBackupStorage / ImageStoreBackupStorage），不同类型需要不同额外参数。需在 Task 0 中确认可创建的类型。

## 测试执行要求

> **所有验收测试必须在真实 ZStack 环境上运行通过，不接受仅编译通过或 t.Skip。**

- 测试环境：`.env.test` 中配置的 ZStack 实例 (当前: 172.24.248.129:8080)
- 运行方式：`source .env.test && go test -v -run 'TestAcc<X>Resource' -count=1 -timeout 300s ./zstack/provider/`
- **每个资源的验收测试必须实际执行 Create → (Update) → Import → Destroy 全流程并 PASS**
- 如果某个 API 返回 404/503（服务未启用）或依赖链中的前置资源创建失败，记录具体错误信息包括tf文件及所有测试参数并**立即上报用户**，等待用户解决环境问题后重跑
- **不允许**用 `t.Skip` 跳过环境问题来标记 Story 完成
- Story 完成的标志：本 Story 的 3 个验收测试全部 PASS

### 已知环境风险 (来自 Story-10 环境探测)

| 依赖 API | 状态 | 影响 |
|----------|:----:|------|
| monitor_template (Create) | **404** | monitor_group 测试无法创建依赖 |
| cdp_policy (Create) | **503** | cdp_task 测试无法创建依赖 |
| backup_storage (Query) | 200 | 可用 |

**注意**：monitor_group 和 cdp_task 的测试依赖链中涉及的 API 在 Story-10 环境探测中有问题。用户需确认这些 API 的可用性后才能开始相关测试。

## 验收标准

- [ ] AC1: `monitor_group` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC2: `cdp_task` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC3: `backup_storage` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC4: **全部 3 个验收测试在真实环境运行 PASS**（不接受 Skip 或仅编译通过）

## Tasks

### Task 0: 为每个资源生成最小 HCL Config

对 3 个资源逐一执行：

- [ ] 0.1 读取 `resource_zstack_<x>.go` 的 `Schema()` 方法，列出所有 `Required: true` 属性
- [ ] 0.2 检查每个 Required 属性的 Validators 确定合法值范围
- [ ] 0.3 对需要 env.json UUID 的属性，确认 `EnvData` 结构体中对应字段名和取值方式
- [ ] 0.4 输出 Create Config + Update Config（Update Config 只改 name）

#### 已知 Required 属性（需通过 Task 0 验证）

| Resource | Required 属性 | ForceNew 属性 |
|----------|-------------|--------------|
| `monitor_group` | `name` | `name`, `description` (依赖 story-07 清理) |
| `cdp_task` | `name`, `task_type`, `policy_uuid`, `backup_storage_uuid`, `resource_uuids` | `task_type`, `policy_uuid`, `backup_storage_uuid`, `resource_uuids` |
| `backup_storage` | `name`, `type` | `type` |

#### cdp_task 详细说明

`cdp_task` 的 5 个 Required 属性中有 4 个是 ForceNew，只有 `name` 可做 in-place update。Update 方法还可更新 `description`, `backup_bandwidth`, `max_capacity`, `max_latency`（均为 Optional）。

测试中 cdp_task 的依赖链：
1. 需要 `cdp_policy`（可在同一测试中 Create，或先运行 Story-10 cdp_policy 测试留下资源）
2. 需要 `backup_storage` UUID（可从 env.json 获取）
3. 需要 `resource_uuids`（volume UUID 列表，可从 env.json volumes 获取）

#### backup_storage 详细说明

`backup_storage` 的 Required 只有 `name` 和 `type`（不是 zone_uuid + SFTP 配置）。`type` 是 ForceNew。不同 type 的创建行为不同：
- **SftpBackupStorage**: 需要 SFTP 服务器信息（hostname, username, password, path）
- **ImageStoreBackupStorage**: 需要 hostname, username, password
- Zone 通过 `attached_zone_uuids` (Optional) 关联，Update 时可 attach/detach

需在 Task 0 中确认测试环境可创建的 backup_storage 类型。

### Task 1: 基础设施 — Destroy 检查函数

- [ ] 1.1 如 Story-10/11 尚未新增 `testAccCheckResourceDestroyByQuery` 泛型 helper，在此 Story 中新增
- [ ] 1.2 在 `check_destroy_test.go` 新增 Destroy 检查：

| 变量/函数名 | 模式 | SDK 方法 | 资源类型 |
|------------|------|---------|---------|
| `testAccCheckMonitorGroupDestroy` | Query 泛型 | `QueryMonitorGroup` | `zstack_monitor_group` |
| `testAccCheckCdpTaskDestroy` | Query 泛型 | `QueryCdpTask` | `zstack_cdp_task` |
| `testAccCheckBackupStorageDestroy` | Get | `GetBackupStorage` | `zstack_backup_storage` |

### Task 2: 处理 monitor_group 测试依赖

- [ ] 2.1 monitor_group 需要 monitor_template_uuid — 在测试 HCL 中先 Create 一个 `zstack_monitor_template` 资源，然后引用其 uuid
- [ ] 2.2 **不要依赖 env.json 中的 monitor_templates**（数据为 NULL），也不要假设 Story-10 的测试会留下资源

### Task 3: 处理 cdp_task 测试依赖

- [ ] 3.1 在测试 HCL 中先 Create 一个 `zstack_cdp_policy` 资源（name + recovery_point_per_second）
- [ ] 3.2 从 env.json 获取 backup_storage UUID 和 volume UUID
- [ ] 3.3 构造完整的 cdp_task Create Config，包含全部 5 个 Required 属性

### 验收测试

- [ ] 4. `resource_zstack_monitor_group_test.go` — 在 HCL 中先创建 monitor_template 依赖; Create + Update (name) + Import + Destroy
- [ ] 5. `resource_zstack_cdp_task_test.go` — 在 HCL 中先创建 cdp_policy 依赖; Create + Update (name) + Import + Destroy
- [ ] 6. `resource_zstack_backup_storage_test.go` — Create (需 type + 类型特定参数) + Update (name) + Import + Destroy

### 验证

- [ ] 7. 编译确认 (`go build ./...` + `go test -short ./zstack/provider/`)
- [ ] 8. **全量运行 3 个验收测试**：`source .env.test && go test -v -run 'TestAcc(MonitorGroup|CdpTask|BackupStorage)Resource' -count=1 -timeout 600s ./zstack/provider/`
- [ ] 9. 全部 PASS 后记录测试输出到 Dev Agent Record；如有 FAIL，记录错误信息并上报用户等待环境修复后重跑

## 审查要点

1. 每个资源的 HCL Config 是否包含所有 Required 属性（通过 Task 0 动态获取）
2. monitor_group 的 Update Step：在 story-07 完成前改 name 会触发 destroy+recreate，需知晓
3. cdp_task 测试是否在同一 TestCase 中创建了完整依赖链（cdp_policy + 引用 env.json 的 bs/volume）
4. Import Step 统一使用 `importStateIdFromUUID("zstack_xxx.test")`
5. Create Step 必须包含 `ConfigStateChecks`：至少验证 `uuid` NotNull + `name` StringExact
6. backup_storage 用 `testAccCheckResourceDestroyByGet`，其余 2 个用 Query 泛型
7. Sprint A 先只测 name 变更；disappears 测试由 Story-09 覆盖

---

## Dev Agent Record

> **重要**：执行过程中发现的所有注意事项（环境问题、API 不可用、参数修正、依赖变更、调试结论等）必须**立即写入**下方对应部分。口头汇总不算完成记录，信息必须持久化到本文档。

### Implementation Plan

1. 在 check_destroy_test.go 添加 3 个 destroy check 函数
2. 为 monitor_group 添加验收测试（Create + Update + Import）
3. 为 cdp_task 添加验收测试（内联 cdp_policy 依赖 + env.json bs/volume）
4. 为 backup_storage 添加验收测试（ImageStoreBackupStorage 类型）
5. 编译 + 运行全部 3 个验收测试

### Debug Log

**2026-04-20 环境验收测试执行**

| 资源 | API | 返回码 | 结果 |
|------|-----|:------:|:----:|
| `monitor_group` | CreateMonitorGroup | 200 | PASS |
| `cdp_task` | CreateCdpPolicy (依赖) | **503** | FAIL |
| `backup_storage` | AddImageStoreBackupStorage | **503** | FAIL |

- `TestAccMonitorGroupResource`: **全流程通过** (Create → Update/Replace → Import → Destroy)
- `TestAccCdpTaskResource`: 失败于 Step 1 — `zstack_cdp_policy.dep` 创建时返回 `StatusCode: 503`，CDP 服务未在测试环境 (172.24.248.129) 启用
- `TestAccBackupStorageResource`: 失败于 Step 1 — `AddImageStoreBackupStorage` 返回 `StatusCode: 503`，ImageStore 创建服务不可用

**结论**: 测试代码本身无问题（monitor_group 已验证全流程）。cdp_task 和 backup_storage 的失败是环境限制，需用户确认：
1. CDP 服务是否可在测试环境启用？
2. 是否允许在测试环境创建新的 ImageStoreBackupStorage？

### Completion Notes

- monitor_group 验收测试已通过，代码可提交
- cdp_task 和 backup_storage 测试代码已就绪，等待环境修复后重跑

## File List

- `zstack/provider/check_destroy_test.go` — 新增 3 个 destroy check 函数
- `zstack/provider/resource_zstack_monitor_group_test.go` — 新增 TestAccMonitorGroupResource
- `zstack/provider/resource_zstack_cdp_task_test.go` — 新增 TestAccCdpTaskResource
- `zstack/provider/resource_zstack_backup_storage_test.go` — 新增 TestAccBackupStorageResource

## Change Log

| 日期 | 变更 |
|------|------|
| 2026-04-20 | 初始实现：3 个资源验收测试 + destroy check。monitor_group PASS，其余 2 个环境 503 |

## Status

部分完成（1/3 PASS，2/3 等待环境修复）
