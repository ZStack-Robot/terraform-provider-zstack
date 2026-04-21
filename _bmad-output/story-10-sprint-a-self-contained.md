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

| # | Resource | Terraform Type Name | Update 状态 | 目标级别 | SDK Read 方式 | 环境依赖 |
|---|----------|-------------------|-----------|---------|-------------|---------|
| 1 | `access_key` | `zstack_access_key` | N/A (no-op) | 基础 | Query (`QueryAccessKey`) | account_uuid, user_uuid (env.json) |
| 2 | `monitor_template` | `zstack_monitor_template` | 有 (UpdateMonitorTemplate) | 标准 | Query (`QueryMonitorTemplate`) | 无 |
| 3 | `log_server` | `zstack_log_server` | 有 (UpdateLogServer) | 标准 | Query (`QueryLogServer`) | 无 |
| 4 | `cdp_policy` | `zstack_cdp_policy` | 有 (UpdateCdpPolicy) | 标准 | Query (`QueryCdpPolicy`) | 无 |
| 5 | `price_table` | `zstack_price_table` | 有 (UpdatePriceTable) | 标准 | Query (`QueryPriceTable`) | 无 |
| 6 | `stack_template` | `zstack_stack_template` | 有 (UpdateStackTemplate) | 标准 | Query (`QueryStackTemplate`) | 无 |
| 7 | `preconfiguration_template` | `zstack_preconfiguration_template` | 有 (UpdatePreconfigurationTemplate) | 标准 | Query (`QueryPreconfigurationTemplate`) | 无 |
| 8 | `instance_scripts` | **`zstack_script`** | 有 (UpdateGuestVmScript) | 标准 | Get (`GetGuestVmScript`) | 无 |

### 已知问题

- **RequiresReplace 未清理** (依赖 story-07): `monitor_template` (name, description)、`price_table` (name, description)、`preconfiguration_template` (name, description, distribution, type, content) 仍有 RequiresReplace。story-07 完成前，这 3 个资源的 Update Step 实际触发 destroy+recreate，不是 in-place update。
- **`instance_scripts` 类型名注意**: 文件名是 `resource_zstack_instance_scripts.go`，但 Terraform 类型名是 `zstack_script`。HCL 中必须写 `resource "zstack_script"`。

## 测试执行要求

> **所有验收测试必须在真实 ZStack 环境上运行通过，不接受仅编译通过或 t.Skip。**

- 测试环境：`.env.test` 中配置的 ZStack 实例 (当前: 172.24.248.129:8080)
- 运行方式：`source .env.test && go test -v -run 'TestAcc<X>Resource' -count=1 -timeout 300s ./zstack/provider/`
- **每个资源的验收测试必须实际执行 Create → (Update) → Import → Destroy 全流程并 PASS**
- 如果某个 API 返回 404/503（服务未启用），记录具体错误信息并**立即上报用户**，等待用户解决环境问题后重跑
- **不允许**用 `t.Skip` 跳过环境问题来标记 Story 完成
- Story 完成的标志：`go test -v -run 'TestAcc' -count=1 ./zstack/provider/` 中本 Story 的 8 个测试全部 PASS

### 环境探测结果 (2026-04-20, 验收测试实测更新)

| API | Query | Create | 测试结果 | 用户确认 |
|-----|:-----:|:------:|:--------:|:-------:|
| access_key | 200 ✓ | 200 ✓ | **PASS** | - |
| monitor_template | 200 ✓ | 200 ✓ | **PASS** | - |
| log_server | 200 ✓ | **503** | FAIL | 待用户确认 |
| cdp_policy | 200 ✓ | **503** | FAIL | 待用户确认 |
| price_table | 200 ✓ | **503** | FAIL | 待用户确认 |
| stack_template | 200 ✓ | **503** | FAIL | 待用户确认 |
| preconfiguration_template | 200 ✓ | **503** | FAIL | 待用户确认 |
| script (instance_scripts) | 200 ✓ | **503** | FAIL | 待用户确认 |

## 验收标准

- [ ] AC1: 8 个资源各有 `TestAcc<X>Resource` 函数，包含 Create + Import + Destroy (基础)
- [ ] AC2: 7 个有 Update 的资源包含 Update Step (标准)
- [ ] AC3: `check_destroy_test.go` 新增对应的 Destroy 检查变量/函数
- [ ] AC4: 编译通过，所有新测试编译通过
- [ ] AC5: **全部 8 个验收测试在真实环境运行 PASS**（不接受 Skip 或仅编译通过）

## Tasks

### Task 0: 为每个资源生成最小 HCL Config

对 8 个资源逐一执行：

- [ ] 0.1 读取 `resource_zstack_<x>.go` 的 `Schema()` 方法，列出所有 `Required: true` 属性
- [ ] 0.2 检查每个 Required 属性的 Validators（`stringvalidator.OneOf(...)`, `LengthAtLeast` 等）确定合法值范围
- [ ] 0.3 对枚举类属性，从 Validator 或同文件常量/类型定义中取一个合法值
- [ ] 0.4 对自由文本属性（name, description），使用 `"acc-test-<resource>"` 前缀
- [ ] 0.5 输出 Create Config + Update Config（Update Config 只改 name）

### Task 1: 基础设施 — Destroy 检查函数

- [ ] 1.1 在 `provider_test.go` 新增 `testAccCheckResourceDestroyByQuery` 泛型 helper（参照 `testAccCheckResourceDestroyByGet` 模式，改为接受 `func(params *param.QueryParam) ([]T, error)` 签名）
- [ ] 1.2 在 `check_destroy_test.go` 用泛型 helper 新增 7 个 Query-based destroy 变量：

| 变量名 | SDK 方法 | 资源类型 |
|--------|---------|---------|
| `testAccCheckAccessKeyDestroy` | `QueryAccessKey` | `zstack_access_key` |
| `testAccCheckMonitorTemplateDestroy` | `QueryMonitorTemplate` | `zstack_monitor_template` |
| `testAccCheckLogServerDestroy` | `QueryLogServer` | `zstack_log_server` |
| `testAccCheckCdpPolicyDestroy` | `QueryCdpPolicy` | `zstack_cdp_policy` |
| `testAccCheckPriceTableDestroy` | `QueryPriceTable` | `zstack_price_table` |
| `testAccCheckStackTemplateDestroy` | `QueryStackTemplate` | `zstack_stack_template` |
| `testAccCheckPreconfigurationTemplateDestroy` | `QueryPreconfigurationTemplate` | `zstack_preconfiguration_template` |

- [ ] 1.3 在 `check_destroy_test.go` 用现有 `testAccCheckResourceDestroyByGet` 新增 1 个 Get-based destroy 变量：

| 变量名 | SDK 方法 | 资源类型 |
|--------|---------|---------|
| `testAccCheckScriptDestroy` | `GetGuestVmScript` | `zstack_script` |

### Task 2: 基础设施 — 确认 env.json 数据

- [ ] 2.1 确认 `env.json` 中 `accounts` 包含可用的 account_uuid
- [ ] 2.2 确认 `env.json` 中 `users` 包含可用的 user_uuid（`access_key` 的 Required 属性需要 `account_uuid` + `user_uuid` 两个字段）

### 验收测试（每个资源一个 task）

- [ ] 3. `resource_zstack_access_key_test.go` — Create (需 account_uuid + user_uuid) + Import + Destroy (U=N/A，Update 是 no-op)
- [ ] 4. `resource_zstack_monitor_template_test.go` — Create + Update (name) + Import + Destroy
- [ ] 5. `resource_zstack_log_server_test.go` — Create + Update (name) + Import + Destroy
- [ ] 6. `resource_zstack_cdp_policy_test.go` — Create + Update (name) + Import + Destroy
- [ ] 7. `resource_zstack_price_table_test.go` — Create + Update (name) + Import + Destroy
- [ ] 8. `resource_zstack_stack_template_test.go` — Create + Update (name) + Import + Destroy
- [ ] 9. `resource_zstack_preconfiguration_template_test.go` — Create + Update (name) + Import + Destroy
- [ ] 10. `resource_zstack_instance_scripts_test.go` — HCL 中使用 `resource "zstack_script"`; Create + Update (name) + Import + Destroy

### 验证

- [ ] 11. 编译确认 (`go build ./...`)
- [ ] 12. Schema 单元测试通过 (`go test -short ./zstack/provider/`)
- [ ] 13. **全量运行 8 个验收测试**：`source .env.test && go test -v -run 'TestAcc(AccessKey|MonitorTemplate|LogServer|CdpPolicy|PriceTable|StackTemplate|PreconfigurationTemplate|Script)Resource' -count=1 -timeout 600s ./zstack/provider/`
- [ ] 14. 全部 PASS 后记录测试输出到 Dev Agent Record；如有 FAIL，记录错误信息并上报用户等待环境修复后重跑

## 审查要点

1. 每个资源的 HCL Config 是否包含所有 Required 属性（通过 Task 0 动态获取）
2. Update Step 修改的属性是否为实际可变属性（不是 ForceNew 的）
3. Import Step 统一使用 `importStateIdFromUUID("zstack_xxx.test")`（所有资源的主键都是 `uuid`）
4. Create Step 必须包含 `ConfigStateChecks`：至少验证 `uuid` NotNull + `name` StringExact
5. Destroy 检查函数：7 个 Query-based 用泛型 helper，1 个 Get-based 用现有 helper
6. Sprint A 先只测 name 变更；disappears 测试由 Story-09 覆盖

---

## Dev Agent Record

> **重要**：执行过程中发现的所有注意事项（环境问题、API 不可用、参数修正、依赖变更、调试结论等）必须**立即写入**下方对应部分。口头汇总不算完成记录，信息必须持久化到本文档。

### Implementation Plan

1. 新增 `testAccCheckResourceDestroyByQuery[T any]()` 泛型 helper 到 `provider_test.go`
2. 新增 7 个 Query-based + 1 个 Get-based destroy 变量到 `check_destroy_test.go`
3. 为 8 个资源在已有的 `*_test.go` 中追加 `TestAcc*Resource` 验收测试函数
4. 通过 MCP API 查询确认 `log_server` 合法字段值（category/type）
5. 在真实环境运行全部 8 个验收测试

### Debug Log

**2026-04-20 第一轮执行：**

1. **log_server HCL 参数修正**：
   - 初始使用 `category=CloudBus`, `type=ElasticSearch` — 错误
   - 通过 `mcp__zstack__describe_api("AddLogServer")` 查到合法值：
     - `category`: `ManagementNodeLog` | `PlatformOperationLog`
     - `type`: `Log4j2` | `FluentBit`
   - 已修正为 `category=ManagementNodeLog`, `type=Log4j2`

2. **script (instance_scripts) platform 字段**：
   - Schema 中 `platform` 标记为 Optional，但 API 实际返回 "the 'platform' field must be provided"
   - 已在 HCL 中增加 `platform = "Linux"`
   - 修复后仍然 503（服务未启用，与 platform 无关）

3. **环境 API 可用性验证结果（2026-04-20）**：

| 资源 | Query API | Create API | 测试结果 |
|------|:---------:|:----------:|:--------:|
| access_key | 200 ✓ | 200 ✓ | **PASS** (6.32s) |
| monitor_template | 200 ✓ | 200 ✓ | **PASS** (11.51s) |
| log_server | 200 ✓ | **503** | FAIL |
| cdp_policy | 200 ✓ | **503** | FAIL |
| price_table | 200 ✓ | **503** | FAIL |
| stack_template | 200 ✓ | **503** | FAIL |
| preconfiguration_template | 200 ✓ | **503** | FAIL |
| script | 200 ✓ | **503** | FAIL |

4. **503 原因**：ZStack 环境中对应的服务模块未启用（如 CDP 备份服务、资源编排服务、预配置模板服务等）。这是环境配置问题，不是代码问题。

5. **注意**：文档原始环境探测表（第 44-55 行）中 monitor_template 标记为 404 不可用，但实际测试 PASS。可能是环境在 4/20 之间有变更，或 Query 与文档探测用了不同 endpoint。

### Completion Notes

- **代码完成度**：8/8 资源验收测试代码已全部就绪，编译通过
- **测试通过率**：2/8 PASS（access_key, monitor_template）
- **阻塞项**：6 个资源的 Create API 返回 503（服务未启用），需用户确认环境是否可以启用这些服务模块
- **已验证的代码行为**：
  - `access_key`: Create + Import + Destroy 全流程 PASS
  - `monitor_template`: Create + Update(name, 走 destroy+recreate 因 RequiresReplace) + Import + Destroy 全流程 PASS

## File List

| 文件 | 变更 |
|------|------|
| `zstack/provider/provider_test.go` | 新增 `testAccCheckResourceDestroyByQuery` 泛型 helper + `param` import |
| `zstack/provider/check_destroy_test.go` | 新增 7 个 Query-based + 1 个 Get-based destroy 变量 |
| `zstack/provider/resource_zstack_access_key_test.go` | 新增 `TestAccAccessKeyResource` |
| `zstack/provider/resource_zstack_monitor_template_test.go` | 新增 `TestAccMonitorTemplateResource` |
| `zstack/provider/resource_zstack_log_server_test.go` | 新增 `TestAccLogServerResource` |
| `zstack/provider/resource_zstack_cdp_policy_test.go` | 新增 `TestAccCdpPolicyResource` |
| `zstack/provider/resource_zstack_price_table_test.go` | 新增 `TestAccPriceTableResource` |
| `zstack/provider/resource_zstack_stack_template_test.go` | 新增 `TestAccStackTemplateResource` |
| `zstack/provider/resource_zstack_preconfiguration_template_test.go` | 新增 `TestAccPreconfigurationTemplateResource` |
| `zstack/provider/resource_zstack_instance_scripts_test.go` | 新增 `TestAccScriptResource` |

## Change Log

| 日期 | 变更 |
|------|------|
| 2026-04-20 | 初始实现：8 个验收测试 + 泛型 destroy helper + destroy 变量 |
| 2026-04-20 | 修正 log_server 字段值 (CloudBus→ManagementNodeLog, ElasticSearch→Log4j2) |
| 2026-04-20 | 修正 script 增加 platform="Linux" 必填字段 |
| 2026-04-20 | 真实环境测试：2/8 PASS，6/8 因 503 被阻塞 |

## Status

**部分完成** — 代码就绪，2/8 测试 PASS，6/8 因环境 503 阻塞，等待用户确认环境
