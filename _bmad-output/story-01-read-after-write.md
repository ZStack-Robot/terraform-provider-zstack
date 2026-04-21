# Story 01: Update 方法改用 read-after-write 模式

> **分支**: `refactor/read-error-handling`  
> **状态**: WIP (2 commits, 待审查合入)  
> **优先级**: P0  
> **预计审查时间**: 15 分钟

---

## Story

作为 Provider 维护者，我需要确保 Update 方法在写入 state 前始终从 API 回读最新数据，而不是信赖 Update API 的返回值，因为 ZStack SDK 多个 Update API 返回值不完整或为空（已确认 alarm 的 UpdateAlarm 返回空响应）。

## 背景

- ZStack SDK 的 Update API 返回值不可靠（已记录: `docs/SDK-BUG-UpdateAlarm-Empty-Response.md`）
- 当前模式: `result, err := cli.Update<X>(...); plan.Field = result.Field` — 如果 result 字段为空，state 会被写入零值
- 目标模式: `cli.Update<X>(...); obj, err := cli.Get<X>(id); plan.Field = obj.Field` — 始终从 Get 回读

## 验收标准

- [ ] AC1: account Update 使用 read-after-write（Update → GetAccount → 写 state）
- [ ] AC2: affinity_group Update 使用 read-after-write（Update → GetAffinityGroup → 写 state）
- [ ] AC3: ssh_key_pair Update 使用 read-after-write（Update → GetSshKeyPair → 写 state）
- [ ] AC4: iam2_project Update 使用 read-after-write（Update → GetIAM2Project → 写 state）
- [ ] AC5: alarm Update 使用 read-after-write（Update → QueryAlarm → 写 state）
- [ ] AC6: 所有修改的资源现有验收测试通过（`TestAccAccountResource`, `TestAccAffinityGroupResource`, `TestAccSshKeyPairResource`, `TestAccIAM2ProjectResource`, `TestAccAlarmResource`）
- [ ] AC7: alarm 测试中 `metric_name` 使用正确值 `CPUAverageUsedUtilization`

## Tasks

- [x] 1. 修改 `resource_zstack_account.go` — Update 方法拆分为 Update + Get 回读
- [x] 2. 修改 `resource_zstack_affinity_group.go` — 同上
- [x] 3. 修改 `resource_zstack_ssh_key_pair.go` — 同上
- [x] 4. 修改 `resource_zstack_iam2_project.go` — 同上
- [x] 5. 修改 `resource_zstack_alarm.go` — Update 方法改用 QueryAlarm 回读（alarm 无 GetAlarm，只有 QueryAlarm）
- [x] 6. 修复 `resource_zstack_alarm_test.go` — metric_name 从 `CPUAverageUtilization` 改为 `CPUAverageUsedUtilization`
- [ ] 7. 跑验收测试确认全部通过

## 涉及文件

| 文件 | 变更类型 |
|------|---------|
| `zstack/provider/resource_zstack_account.go` | 修改 |
| `zstack/provider/resource_zstack_affinity_group.go` | 修改 |
| `zstack/provider/resource_zstack_ssh_key_pair.go` | 修改 |
| `zstack/provider/resource_zstack_iam2_project.go` | 修改 |
| `zstack/provider/resource_zstack_alarm.go` | 修改 |
| `zstack/provider/resource_zstack_alarm_test.go` | 修改 |

## 审查要点

1. 每个资源的 Get/Query 调用是否使用正确的方法（Get vs Query）
2. 错误处理是否完整（Update 失败 return、Get 失败 return）
3. state 字段赋值是否完整（与原有字段列表一致）

## 合并注意

- **此分支覆盖了 `fix/alarm-update-and-test` 分支的 alarm 部分**，合入此分支后该分支的 alarm 改动可跳过
- 与 `refactor/provider-quality-hardening` 有交叉（iam2_project），需确认合并顺序
