# Story 05: 修复 policy 资源创建 Bug + 解除测试 Skip

> **分支**: 待创建 (`fix/policy-create-statements`)  
> **状态**: WIP (工作目录未提交)  
> **优先级**: P0 (已知 Bug, 测试被 t.Skip 跳过)  
> **预计审查时间**: 5 分钟

---

## Story

作为 Provider 用户，我需要能成功创建 policy 资源。当前创建时发送空 `statements` 数组导致 ZStack API 拒绝请求。

## 根因

`resource_zstack_policy.go` Create 方法中 `Statements: []param.PolicyStatementParam{}` 发送空数组。ZStack API 要求 statements 非空。修复方案: 传入合理的默认 statement。

## 验收标准

- [ ] AC1: Create 时 statements 包含默认策略 `{Name: "default", Effect: "Allow", Actions: ["**"]}`
- [ ] AC2: `resource_zstack_policy_test.go` 中的 `t.Skip` 已移除
- [ ] AC3: `TestAccPolicyResource` 验收测试通过

## Tasks

- [x] 1. 修改 `resource_zstack_policy.go` Create — 替换空 statements 为默认策略
- [x] 2. 修改 `resource_zstack_policy_test.go` — 移除 `t.Skip(...)` 行
- [ ] 3. 提交到新分支 `fix/policy-create-statements`
- [ ] 4. 跑验收测试确认通过

## 涉及文件

| 文件 | 变更类型 |
|------|---------|
| `zstack/provider/resource_zstack_policy.go` | 修改 (+7 -3) |
| `zstack/provider/resource_zstack_policy_test.go` | 修改 (-1 行) |

## 审查要点

1. 默认 statement `Allow + **` 是否安全合理（全权限策略作为默认值）
2. 是否应该让用户在 HCL 中声明 statements 而非硬编码默认值（后续 story 可能需要将 statements 暴露为 schema 属性）
