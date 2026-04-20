# Sprint 全局概览: Phase 0 收尾 + Phase 1 Sprint A

> **日期**: 2026-04-17  
> **目标**: 完成 Phase 0 全部收尾，启动 Phase 1 Sprint A 核心资源验收测试  
> **预计总工作量**: ~10 天

---

## Story 总览

### Phase 0.4: 代码质量加固 (story 01-07)

| # | Story | 分支 | 状态 | P | 大小 | 合入序 |
|---|-------|------|------|---|------|--------|
| 01 | Update read-after-write | `refactor/read-error-handling` | 待审查 | P0 | M | 1st |
| 02 | alarm Disappears + SDK Bug | `fix/alarm-update-and-test` | 待审查 | P1 | S | 2nd (依赖01) |
| 03 | port_forwarding_rule Unknown | `fix/port-forwarding-rule` | 待审查 | P0 | XS | 独立 |
| 04 | iam2_project Expunge | `refactor/provider-quality-hardening` | 待审查 | P1 | XS | 3rd (依赖01) |
| 05 | policy Create Bug | 待创建 | 未提交 | P0 | XS | 独立 |
| 06 | scheduler_job/global_config Import | 待创建 | 未提交 | P1 | S | 独立 |
| 07 | RequiresReplace + l2vxlan Update | 待创建 | 未提交 | P1 | M | 独立 |

### Phase 0.3: 收尾 (story 08)

| # | Story | 分支 | 状态 | P | 大小 | 合入序 |
|---|-------|------|------|---|------|--------|
| 08 | reserved_ip/secgroup Import + VR Destroy | 待创建 | 未开始 | P1 | S | 独立 |

### Phase 0.2: Disappears 批量补全 (story 09)

| # | Story | 分支 | 状态 | P | 大小 | 合入序 |
|---|-------|------|------|---|------|--------|
| 09 | 批量 Disappears (~26 个) | 待创建 | 未开始 | P2 | L | 依赖08完成 |

### Phase 1 Sprint A: 核心资源验收测试 (story 10-13)

| # | Story | 分支 | 状态 | P | 资源数 | 合入序 |
|---|-------|------|------|---|--------|--------|
| 10 | 自包含资源 | 待创建 | 未开始 | P1 | 8 | 依赖07 |
| 11 | 计算+存储资源 | 待创建 | 未开始 | P1 | 7 | 独立 |
| 12 | 网络+基础设施 | 待创建 | 未开始 | P1 | 6 | 依赖07(l2vxlan) |
| 13 | 监控运维+备份 | 待创建 | 未开始 | P1 | 3 | 依赖10(monitor_template) |

---

## 依赖图

```
Phase 0.4 (代码质量)                    Phase 0.3      Phase 0.2
┌──────────────────────┐              ┌─────────┐    ┌─────────┐
│ story-01 (read-after-write) ──┐     │ story-08 │───→│ story-09 │
│   ├→ story-02 (alarm dis)     │     │ (Import/ │    │ (批量    │
│   └→ story-04 (iam2 expunge)  │     │  Destroy)│    │ Disappears)
├──────────────────────┤        │     └─────────┘    └─────────┘
│ story-03 (port-fwd)  │ 独立   │
│ story-05 (policy)    │ 独立   │
│ story-06 (import)    │ 独立   │
│ story-07 (RR+l2vxlan)│────────┼────────────────────┐
└──────────────────────┘        │                    │
                                │  Phase 1 Sprint A   │
                                │  ┌──────────────────┘
                                │  │
                                ▼  ▼
                          ┌───────────┐
                          │ story-10   │──→ story-13
                          │ (自包含 8个)│   (监控运维)
                          └───────────┘
                          ┌───────────┐
                          │ story-11   │ 独立
                          │ (计算存储)  │
                          └───────────┘
                          ┌───────────┐
                          │ story-12   │ 依赖 story-07
                          │ (网络基础)  │ (l2vxlan Update)
                          └───────────┘
```

---

## 推荐执行顺序

### 第一批: 并行审查+合入 (Day 1-2)

独立且小的 story，互不冲突：

| 序 | Story | 操作 |
|----|-------|------|
| 1a | story-03 | 审查 → 合入 `fix/port-forwarding-rule` |
| 1b | story-05 | 从工作目录创建分支 → 审查 → 合入 |
| 1c | story-06 | 从工作目录创建分支 → 审查 → 合入 |

### 第二批: 核心改动 (Day 2-3)

| 序 | Story | 操作 |
|----|-------|------|
| 2a | story-01 | 审查 → 合入 `refactor/read-error-handling` |
| 2b | story-02 | rebase on master → 审查 → 合入 |
| 2c | story-04 | rebase on master → 审查 → 合入 |
| 2d | story-07 | 从工作目录创建分支 → 审查 → 合入 |

### 第三批: Phase 0 收尾 (Day 3-4)

| 序 | Story | 操作 |
|----|-------|------|
| 3a | story-08 | 开发 → 审查 → 合入 |
| 3b | story-09 | 开发 → 审查 → 合入 |

### 第四批: Phase 1 Sprint A (Day 4-10)

| 序 | Story | 操作 |
|----|-------|------|
| 4a | story-10 | 开发自包含资源测试 → 审查 → 合入 |
| 4b | story-11 | 开发计算+存储测试（可与 10 并行）→ 审查 → 合入 |
| 4c | story-12 | 开发网络+基础设施测试 → 审查 → 合入 |
| 4d | story-13 | 开发监控运维测试（依赖 10）→ 审查 → 合入 |

---

## 分支操作命令

### 从工作目录创建 3 个新分支

```bash
# 先 stash 当前所有改动
git stash

# Story 05: policy fix
git checkout -b fix/policy-create-statements master
git stash pop
git add zstack/provider/resource_zstack_policy.go zstack/provider/resource_zstack_policy_test.go
git commit -m "fix(policy): use default statements in Create to fix empty array bug"
git checkout master
git stash

# Story 06: import steps
git checkout -b fix/missing-import-steps master
git stash pop
git add zstack/provider/resource_zstack_scheduler_job_test.go zstack/provider/resource_zstack_global_config_test.go
git commit -m "test: add Import Steps for scheduler_job and global_config"
git checkout master
git stash

# Story 07: RR cleanup + l2vxlan
git checkout -b fix/remove-requires-replace master
git stash pop
git add zstack/provider/resource_zstack_monitor_group.go \
        zstack/provider/resource_zstack_monitor_template.go \
        zstack/provider/resource_zstack_preconfiguration_template.go \
        zstack/provider/resource_zstack_price_table.go \
        zstack/provider/resource_zstack_l2vxlan_network.go
git commit -m "fix: remove incorrect RequiresReplace and implement l2vxlan_network Update"
git checkout master
git stash pop  # 恢复残余（应该为空）
```

### 分支清理

```bash
git branch -d fix/storage-acceptance-tests  # 已合入 master
# fix/alarm-update-and-test: 合入 story-01+02 后删除
```

---

## 完成后的项目状态预期

| 指标 | 当前 | Phase 0 完成后 | Sprint A 完成后 |
|------|------|---------------|----------------|
| Resource 验收覆盖率 | 35.1% (39/111) | 35.1% (39/111) | **56.8%** (63/111) |
| 完整级 (C+U+I+D+Dis) | 5 | **6** | 6+ |
| 标准级 (C+U+I+D) | 14 | **15** | **~35** |
| 不完整 | 4 | **0** | 0 |
| 已知问题 | 8 | **~4** | ~4 |
| Disappears 覆盖 | 6/39 | **~32/39** | ~32/63 |
| Data Source 覆盖率 | 100% | 100% | 100% |
