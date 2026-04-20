# Story 03: 修复 port_forwarding_rule Optional+Computed Int64 Unknown 状态

> **分支**: `fix/port-forwarding-rule`  
> **状态**: WIP (1 commit, 待审查合入)  
> **优先级**: P0 (影响资源创建)  
> **预计审查时间**: 5 分钟

---

## Story

作为 Provider 用户，当我创建 port_forwarding_rule 时如果不指定 `vip_port_end`、`private_port_start`、`private_port_end` 这些 Optional+Computed 属性，Provider 不应将 Unknown 状态的零值发送给 API，而是应跳过这些字段让 API 使用默认值。

## 根因

Terraform Plugin Framework 中 Optional+Computed 的 Int64 属性在 plan 阶段为 Unknown（不是 Null）。原代码只检查 `!plan.Field.IsNull()` 就发送值，但 Unknown 状态下 `ValueInt64()` 返回 0，导致发送错误的端口值。

## 验收标准

- [ ] AC1: `vip_port_end` 检查同时包含 `!IsNull() && !IsUnknown()`
- [ ] AC2: `private_port_start` 检查同时包含 `!IsNull() && !IsUnknown()`
- [ ] AC3: `private_port_end` 检查同时包含 `!IsNull() && !IsUnknown()`
- [ ] AC4: 测试断言类型正确（`knownvalue.Int64Exact(8080)` 而非 `StringExact("8080")`）
- [ ] AC5: 编译通过，测试通过

## Tasks

- [x] 1. 修改 `resource_zstack_port_forwarding_rule.go` Create 方法 — 三个字段增加 `!IsUnknown()` 检查
- [x] 2. 修改 `resource_zstack_port_forwarding_rule_test.go` — 断言类型从 StringExact 改为 Int64Exact
- [ ] 3. 编译确认无错误

## 涉及文件

| 文件 | 变更类型 |
|------|---------|
| `zstack/provider/resource_zstack_port_forwarding_rule.go` | 修改 (3 行) |
| `zstack/provider/resource_zstack_port_forwarding_rule_test.go` | 修改 (1 行) |

## 审查要点

1. 所有 Optional+Computed Int64 属性是否都加了 Unknown 检查（有没有遗漏的字段）
2. 此模式是否应推广到其他资源的 Optional+Computed Int64 属性（可作为后续 story）
