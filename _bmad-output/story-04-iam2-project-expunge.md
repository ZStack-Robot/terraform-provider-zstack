# Story 04: iam2_project Delete 追加 Expunge 调用

> **分支**: `refactor/provider-quality-hardening`  
> **状态**: WIP (2 commits, 待审查合入)  
> **优先级**: P1  
> **预计审查时间**: 5 分钟

---

## Story

作为 Provider 用户，当我执行 `terraform destroy` 删除 iam2_project 后，项目名称应被彻底释放，以便后续 `terraform apply` 可以使用相同名称创建新项目。

## 根因

ZStack 的 iam2_project 删除后进入回收站，名称未释放。再次创建同名项目会报 duplicate name 错误。需要在 Delete 后追加 Expunge 调用彻底清除。

## 验收标准

- [ ] AC1: `resource_zstack_iam2_project.go` Delete 方法在 DeleteIAM2Project 后调用 ExpungeIAM2Project
- [ ] AC2: Expunge 失败时返回 Diagnostics Error（不吞错误）
- [ ] AC3: 验收测试中同名资源可重复创建销毁无冲突

## Tasks

- [x] 1. 修改 `resource_zstack_iam2_project.go` Delete — 追加 `r.client.ExpungeIAM2Project(uuid)` 调用
- [x] 2. 添加 Expunge 错误处理（AddError + return）
- [ ] 3. 跑 `TestAccIAM2ProjectResource` 验证无回归

## 涉及文件

| 文件 | 变更类型 |
|------|---------|
| `zstack/provider/resource_zstack_iam2_project.go` | 修改 (+9 行) |

## 审查要点

1. Expunge 是否为幂等操作（资源已不存在时是否报错）
2. 是否有其他资源也有回收站行为需要同样处理（参考 E2E-TEST-PLAN.md 1.4 节: instance, volume, image）

## 备注

此分支还包含项目配置文件变更（`.omc/prd.json`, `_bmad/core/config.yaml`）和删除 `generate_tf` 二进制文件，这些是 chore 性质，不影响 Provider 功能。
