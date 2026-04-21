# Story 02: alarm 新增 Disappears 测试 + SDK Bug 文档

> **分支**: `fix/alarm-update-and-test`  
> **状态**: WIP (2 commits, 待审查合入)  
> **优先级**: P1  
> **预计审查时间**: 10 分钟

---

## Story

作为 Provider 维护者，我需要为 alarm 资源补充 Disappears 测试（使其从"标准"升级到"完整"级），并记录发现的 SDK Bug（UpdateAlarm 返回空响应）。

## 背景

- alarm 当前为"标准"级 (C+U+I+D)，缺 Disappears 维度
- 在实现过程中发现 ZStack SDK `UpdateAlarm` API 返回空响应体，需要记录以便后续 SDK 修复
- 此分支与 `refactor/read-error-handling` 在 alarm 代码上有完全重叠，**alarm.go 的改动应以 story-01 的分支为准**

## 验收标准

- [ ] AC1: `check_disappears_test.go` 新增 `stateCheckAlarmDisappears` 函数
- [ ] AC2: `TestAccAlarmResource_disappears` 测试存在且可运行
- [ ] AC3: `docs/SDK-BUG-UpdateAlarm-Empty-Response.md` 文档完整描述 Bug 现象、多环境验证结果、workaround

## Tasks

- [x] 1. 在 `check_disappears_test.go` 新增 alarm 的 disappears 检查函数
- [x] 2. 在 `resource_zstack_alarm_test.go` 新增 `TestAccAlarmResource_disappears` 测试
- [x] 3. 编写 SDK Bug 报告文档 `docs/SDK-BUG-UpdateAlarm-Empty-Response.md`
- [ ] 4. 跑验收测试确认通过

## 涉及文件

| 文件 | 变更类型 |
|------|---------|
| `zstack/provider/check_disappears_test.go` | 修改 |
| `zstack/provider/resource_zstack_alarm.go` | 修改 (与 story-01 重叠) |
| `zstack/provider/resource_zstack_alarm_test.go` | 修改 (与 story-01 重叠) |
| `docs/SDK-BUG-UpdateAlarm-Empty-Response.md` | 新增 |

## 审查要点

1. disappears 测试模式是否符合项目标准（使用 `stateCheckDisappears` 工厂）
2. SDK Bug 文档是否包含足够的复现信息
3. **与 story-01 的合并冲突处理**：alarm.go 和 alarm_test.go 的改动在两个分支中完全相同，合入 story-01 后此分支仅需保留 `check_disappears_test.go` 和 SDK Bug 文档的增量

## 合并建议

**选项 A（推荐）**: 先合 story-01，然后 rebase 此分支，只保留 disappears 和 SDK 文档的增量  
**选项 B**: 先合此分支（包含 alarm 完整改动），然后 story-01 在 alarm 部分会 no-op
