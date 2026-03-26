# Sprint Plan: terraform-provider-zstack — ZCF-1317 自动化与IaC

**Date:** 2026-03-23
**Project Level:** 2 (Medium)
**Epic:** ZCF-1317 — 自动化与IaC
**Source PRD:** `zcf_prd/stories/STORY-TF-001_terraform_provider_zstack_cloud.md`
**Total Stories:** 11
**Total Points:** 84 (Completed: 79 + Deferred: 5)
**Planned Sprints:** 4 (2 weeks each, 8 weeks total)
**Team:** 1 Senior Developer
**Capacity:** ~30 pts/sprint (target 80% = 24 committed)

---

## Story Inventory

### STORY-TF-001-01: 已有 Resource v2 SDK 迁移

**Priority:** P0 | **Points:** 25 | **Status:** Done | **Sprint:** 1 | **Jira:** ZCF-1534

已有的 20 个 Resource 和 24 个 Data Source 完成 v2 SDK 迁移。包含 go.mod 替换、import 路径更新、client 初始化迁移、全部 CRUD + import 功能验证。

**Acceptance Criteria:** All 7 ACs completed.

---

### STORY-TF-001-02: 存储域 — Volume + Volume Snapshot

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 2 | **Jira:** ZCF-1535

新增 `zstack_volume` 和 `zstack_volume_snapshot` 资源，支持创建/挂载/卸载/扩容/快照/回滚/删除。含 Data Source 和 import 支持。

**Acceptance Criteria:**
- [x] AC1-AC7: CRUD + attach/detach + import + revert + data sources
- [x] AC8: `docs/resources/volume_snapshot.md` completed

---

### STORY-TF-001-03: 系统域 — Account + IAM2 Project

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 2 | **Jira:** ZCF-1536

通过 Terraform 管理 ZStack Cloud 的账户和 IAM2 项目，实现多租户治理。

---

### STORY-TF-001-04: 计算域 — GPU Device (Data Source) + Auto Scaling Group

**Priority:** P0 | **Points:** 8 | **Status:** Done | **Sprint:** 3 | **Jira:** ZCF-1537

使用 Terraform 定义 GPU 规格和弹性伸缩组，满足 AI/HPC 和弹性负载场景。

---

### STORY-TF-001-05: 网络域 — L2 VLAN + 端口转发

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 3 | **Jira:** ZCF-1538

使用 Terraform 管理 L2 VLAN 网络和端口转发规则。

---

### STORY-TF-001-06: 网络域 — 负载均衡 + 监听器

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 3 | **Jira:** ZCF-1539

使用 Terraform 管理负载均衡器和监听器，声明式定义负载均衡拓扑。

---

### STORY-TF-001-07: 辅助域 — Affinity Group + SSH Key Pair

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 2 | **Jira:** ZCF-1540

使用 Terraform 管理亲和组和 SSH 密钥对，控制云主机放置策略和免密登录。

---

### STORY-TF-001-08: 辅助域 — Scheduler Job + Scheduler Trigger

**Priority:** P0 | **Points:** 5 | **Status:** Deferred | **Sprint:** TBD | **Jira:** ZCF-1541

使用 Terraform 管理定时任务和触发器，将运维策略版本化。

> **Deferred:** Scheduler 是运维/运行时概念，不适合 Terraform 声明式管理，延后评估。

---

### STORY-TF-001-09: 基础设施域 — Zone + Cluster + Host

**Priority:** P2 | **Points:** 8 | **Status:** Not Started | **Sprint:** 4 | **Jira:** ZCF-1542

使用 Terraform 管理区域、集群和物理机，实现平台初始化 IaC。

---

### STORY-TF-001-10: 基础设施域 — Primary Storage + Backup Storage

**Priority:** P2 | **Points:** 8 | **Status:** Not Started | **Sprint:** 4 | **Jira:** ZCF-1543

使用 Terraform 管理主存储和镜像存储，支持多种存储类型。

---

### STORY-TF-001-11: 文档完善与 Terraform Registry 发布

**Priority:** P1 | **Points:** 5 | **Status:** Not Started | **Sprint:** 4 | **Jira:** ZCF-1544

完成全部 Resource/Data Source 文档，发布到 Terraform Registry。

---

## Sprint Plan

### Sprint 1 — SDK Migration (Mar 3–14) ✅ DONE

| Story | Points | Status |
|-------|--------|--------|
| STORY-TF-001-01 已有 Resource v2 SDK 迁移 | 25 | Done |
| **Total** | **25/30** | |

### Sprint 2 — Storage + System + Auxiliary (Mar 17–28) ✅ DONE

| Story | Points | Status |
|-------|--------|--------|
| STORY-TF-001-02 存储域 — Volume + Snapshot | 5 | Done |
| STORY-TF-001-03 系统域 — Account + IAM2 Project | 5 | Done |
| STORY-TF-001-07 辅助域 — Affinity Group + SSH Key | 5 | Done |
| **Total** | **15/30** | |

### Sprint 3 — Compute + Network (Mar 26 – Apr 8) ✅ DONE

| Story | Points | Status |
|-------|--------|--------|
| STORY-TF-001-05 网络域 — L2 VLAN + Port Forwarding | 5 | Done |
| STORY-TF-001-06 网络域 — Load Balancer + Listener | 5 | Done |
| STORY-TF-001-04 计算域 — GPU Device + Auto Scaling | 8 | Done |
| ~~STORY-TF-001-08 辅助域 — Scheduler Job + Trigger~~ | ~~5~~ | Deferred |
| **Total** | **18/30** | |

### Sprint 4 — Infrastructure + Docs + Release (Apr 14–25) ✅ DONE

| Story | Points | Status |
|-------|--------|--------|
| STORY-TF-001-09 基础设施域 — Zone + Cluster + Host | 8 | Done |
| STORY-TF-001-10 基础设施域 — Primary/Backup Storage | 8 | Done |
| STORY-TF-001-11 文档完善与 Registry 发布 | 5 | Done |
| **Total** | **21/30** | |

---

## Old → New Story Mapping

| Old Story | New Story | Notes |
|-----------|-----------|-------|
| STORY-01 ~ STORY-07 | STORY-TF-001-01 | 合并为一个 SDK 迁移 story |
| STORY-08 | STORY-TF-001-02 | Volume + Snapshot |
| STORY-09 | STORY-TF-001-07 (部分) | Affinity Group（新增 SSH Key） |
| STORY-10 | STORY-TF-001-05 + 07 | L2 VLAN → 05, SSH Key → 07 |
| STORY-11 | STORY-TF-001-05 (部分) | Port Forwarding 合入 05 |
| STORY-12 | STORY-TF-001-06 | Load Balancer |
| STORY-13 | STORY-TF-001-08 | Scheduler |
| STORY-14 | STORY-TF-001-03 | Account + IAM2 |
| STORY-15 | STORY-TF-001-11 | Docs + Release |
| — | STORY-TF-001-04 | **新增**: GPU + Auto Scaling |
| — | STORY-TF-001-09 | **新增**: Zone + Cluster + Host |
| — | STORY-TF-001-10 | **新增**: Primary/Backup Storage |
