# Story 09: 批量补全 Disappears 测试 (Phase 0.2)

> **分支**: 待创建 (`test/batch-disappears`)  
> **状态**: 未开始  
> **优先级**: P2 — Phase 0.2  
> **预计工作量**: 1-2 天

---

## Story

作为 Provider 维护者，我需要为现有 39 个有验收测试的资源中缺少 Disappears 测试的资源批量补全，提升测试完整度。

## 背景

- 当前已有 Disappears: zone, account, ssh_key_pair, affinity_group, cluster, iam2_project (6 个) + alarm (WIP, story-02)
- 需新增: 约 26 个（排除已有 6 个、WIP 1 个、不适用 4 个: sns_email_endpoint, sns_http_endpoint, global_config, tag — 另 secgroup_attachment 也不适用）
- 模式固定: 在 `check_disappears_test.go` 加工厂函数 + 在各 `_test.go` 加 `_disappears` 测试

## 不适用的资源（跳过）

- `sns_email_endpoint` — API 不支持删除
- `sns_http_endpoint` — API 不支持删除
- `global_config` — 语义为"重置"而非"删除"
- `networking_secgroup_attachment` — 无独立删除语义

## 验收标准

- [ ] AC1: `check_disappears_test.go` 新增 ~26 个 `stateCheck<X>Disappears` 工厂函数
- [ ] AC2: 每个资源有对应的 `TestAcc<X>Resource_disappears` 测试函数
- [ ] AC3: 所有 disappears 测试编译通过
- [ ] AC4: 抽样运行 5 个验证通过（全量依赖环境时间）

## Tasks

建议按批次拆分提交，每批 ~8 个资源，每批一个 commit：

### Batch 1: 标准级资源（有 Update）

- [ ] 1. volume, l2vlan_network, networking_secgroup_rule, load_balancer
- [ ] 2. role, scheduler_trigger, certificate, iam2_organization

### Batch 2: 标准级资源（续）+ scheduler_job / global_config（跳过后者）

- [ ] 3. sns_topic, scheduler_job
- [ ] 4. alarm — 由 story-02 覆盖，此处跳过

### Batch 3: 基础级资源 (U=N/A)

- [ ] 5. image, vip, disk_offering, instance_offering
- [ ] 6. networking_secgroup, tag, virtual_router_image, virtual_router_offering
- [ ] 7. policy, reserved_ip, webhook, user, iam2_virtual_id

### 验证

- [ ] 8. 编译通过 (`go build ./...`)
- [ ] 9. 抽样运行 5 个 disappears 测试验证

## 涉及文件

| 文件 | 变更类型 |
|------|---------|
| `zstack/provider/check_disappears_test.go` | 修改 (~26 个新函数) |
| `zstack/provider/resource_zstack_*_test.go` (26 个) | 修改 (各加 1 个测试函数) |

## 审查要点

1. 每个工厂函数的 delete 调用是否使用正确的 SDK 方法和删除模式
2. 有回收站的资源 (instance, volume) 的 disappears 函数是否用 Permissive 删除
3. 模式一致性: 全部使用 `stateCheckDisappears` 通用工厂
