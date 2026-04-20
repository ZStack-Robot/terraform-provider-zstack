# Story 13: Sprint A — 监控运维和备份资源验收测试 (3 个)

> **分支**: 待创建 (`test/sprint-a-monitoring-ops`)  
> **状态**: 未开始  
> **优先级**: P1 — Phase 1 Sprint A  
> **前置依赖**: story-10 (monitor_template 测试，monitor_group 依赖 monitor_template_uuid)  
> **预计工作量**: 1 天

---

## Story

作为 Provider 维护者，我需要为监控运维和数据保护相关的资源新增验收测试，完成 Sprint A 最后一组。

## 资源清单

| # | Resource | Update 状态 | 目标级别 | 环境依赖 |
|---|----------|-----------|---------|---------|
| 1 | `monitor_group` | 有 (UpdateMonitorGroup) | 标准 | monitor_template_uuid |
| 2 | `cdp_task` | 有 (UpdateCdpTask) | 标准 | cdp_policy_uuid, volume_uuid |
| 3 | `backup_storage` | 有 (UpdateBackupStorage + zone attach/detach) | 标准 | zone_uuid |

## 验收标准

- [ ] AC1: `monitor_group` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC2: `cdp_task` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC3: `backup_storage` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC4: 编译通过，抽样运行验证

## Tasks

- [ ] 1. `resource_zstack_monitor_group_test.go` — Create (需 monitor_template_uuid) + Update (name) + Import + Destroy
- [ ] 2. `resource_zstack_cdp_task_test.go` — Create (需 cdp_policy_uuid + volume_uuid) + Update + Import + Destroy
- [ ] 3. `resource_zstack_backup_storage_test.go` — Create (需 zone_uuid + SFTP 配置) + Update (name) + Import + Destroy
- [ ] 4. 在 `check_destroy_test.go` 新增 3 个 Destroy 检查函数
- [ ] 5. 编译确认 + 抽样运行

## 审查要点

1. monitor_group 是否需要先创建 monitor_template（测试内创建 vs 引用 env.json）
2. cdp_task 的 cdp_policy 依赖如何处理（同一测试内创建 vs 引用 env.json）
3. backup_storage 创建什么类型（SftpBackupStorage）？需要 SFTP 服务器信息
