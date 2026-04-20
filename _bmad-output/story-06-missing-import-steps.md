# Story 06: 补全 scheduler_job 和 global_config 的 Import Step

> **分支**: 待创建 (`fix/missing-import-steps`)  
> **状态**: WIP (工作目录未提交)  
> **优先级**: P1  
> **预计审查时间**: 10 分钟

---

## Story

作为 Provider 用户，我需要能通过 `terraform import` 导入已有的 scheduler_job 和 global_config 资源。这两个资源使用 composite ID（而非简单 UUID），需要自定义 ImportStateIdFunc。

## 背景

- `scheduler_job` 的 import ID 格式: `uuid:type`（需要 type 来区分不同类型的 scheduler job）
- `global_config` 的 import ID 格式: `category/name`（global_config 没有 UUID，以 category+name 唯一标识）
- 两者都在 E2E-TEST-PLAN.md 2.6 节被列为"已知问题"

## 验收标准

- [ ] AC1: `TestAccSchedulerJobResource` 包含 ImportState Step，使用 `uuid:type` composite ID
- [ ] AC2: `importStateIdSchedulerJob` 函数正确从 state 提取 uuid 和 type 拼接
- [ ] AC3: `TestAccGlobalConfigResource` 包含 ImportState Step，使用 `category/name` composite ID
- [ ] AC4: `importStateIdGlobalConfig` 函数正确从 state 提取 category 和 name 拼接
- [ ] AC5: 两个资源的验收测试通过（含 Import Step）
- [ ] AC6: scheduler_job 从"不完整"级升级到"标准"级 (C+U+I+D)
- [ ] AC7: global_config 从"不完整"级升级到"标准"级 (C+U+I+D*)

## Tasks

- [x] 1. 在 `resource_zstack_scheduler_job_test.go` 新增 Import Step + `importStateIdSchedulerJob` 函数
- [x] 2. 在 `resource_zstack_global_config_test.go` 新增 Import Step + `importStateIdGlobalConfig` 函数
- [ ] 3. 提交到新分支 `fix/missing-import-steps`
- [ ] 4. 跑验收测试确认 Import 成功

## 涉及文件

| 文件 | 变更类型 |
|------|---------|
| `zstack/provider/resource_zstack_scheduler_job_test.go` | 修改 (+25 行) |
| `zstack/provider/resource_zstack_global_config_test.go` | 修改 (+26 行) |

## 审查要点

1. composite ID 格式是否与 ImportState 方法中的解析逻辑一致
2. `ImportStateVerifyIdentifierAttribute` 是否设置正确（scheduler_job → "uuid", global_config → "name"）
3. 命名惯例: `importStateIdSchedulerJob` vs `importStateIdFromUUID` — 是否需要统一命名前缀
