# Sprint Plan: terraform-provider-zstack — ZCF-1317 自动化与IaC

**Date:** 2026-04-03 (Updated)
**Project Level:** 2 (Medium)
**Epic:** ZCF-1317 — 自动化与IaC
**Source PRD:** `zcf_prd/stories/STORY-TF-001_terraform_provider_zstack_cloud.md`
**Source ROADMAP:** `docs/ROADMAP.md`
**Total Stories:** 20 (11 original + 9 new)
**Total Points:** 127 (Completed: 117 + Deferred: 5 + New: 5)
**Planned Sprints:** 7 (2 weeks each, 14 weeks total)
**Team:** 1 Senior Developer
**Capacity:** ~30 pts/sprint (target 80% = 24 committed)
**Rolling Velocity:** ~20 pts/sprint (based on Sprint 1–4)

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

**Priority:** P2 | **Points:** 8 | **Status:** Done | **Sprint:** 4 | **Jira:** ZCF-1542

使用 Terraform 管理区域、集群和物理机，实现平台初始化 IaC。

---

### STORY-TF-001-10: 基础设施域 — Primary Storage + Backup Storage

**Priority:** P2 | **Points:** 8 | **Status:** Done | **Sprint:** 4 | **Jira:** ZCF-1543

使用 Terraform 管理主存储和镜像存储，支持多种存储类型。

---

### STORY-TF-001-11: 文档完善与 Terraform Registry 发布

**Priority:** P1 | **Points:** 5 | **Status:** Done | **Sprint:** 4 | **Jira:** ZCF-1544

完成全部 Resource/Data Source 文档，发布到 Terraform Registry。

---

### STORY-TF-001-12: SDK URL Workaround 修复 + 已知问题修复

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 5 | **Jira:** TBD

修复 8 处 SDK URL 占位符运行时 bug（6 个文件），修复 VPC error handling 和 Tag Value 字段问题。

---

### STORY-TF-001-13: 补充 Data Source — EIP + Reserved IP + Subnet IP Range

**Priority:** P0 | **Points:** 3 | **Status:** Done | **Sprint:** 5 | **Jira:** TBD

为已有 Resource 补充缺失的 Data Source（eip, reserved_ip, subnet_ip_range），实现完整查询能力。

---

### STORY-TF-001-14: 补充 Resource — SDN Controller + L3 Network

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 5 | **Jira:** TBD

为已有 Data Source 补充缺失的 Resource（sdn_controller, l3network），实现完整生命周期管理。

---

### STORY-TF-001-15: 网络域 — VM NIC 管理

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 5 | **Jira:** TBD

独立管理云主机网卡，支持多网卡场景和精细化网络配置。

---

### STORY-TF-001-16: 网络域 — IPSec VPN 连接 + VPC VPN 网关

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 6 | **Jira:** TBD

管理 IPSec VPN 连接和 VPC VPN 网关，实现站点到站点 VPN 隧道自动化。

---

### STORY-TF-001-17: 网络域 — VPC Firewall + 路由表/路由条目

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 6 | **Jira:** TBD

管理 VPC 防火墙和虚拟路由表/条目，实现声明式网络安全策略和路由规则。

---

### STORY-TF-001-18: 存储域 — Ceph 存储 + ImageStore 备份存储

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 6 | **Jira:** TBD

管理 Ceph 主存储、Ceph 备份存储和镜像仓库备份存储。

---

### STORY-TF-001-19: 系统域 — IAM2 Virtual ID + Access Key + Global Config

**Priority:** P0 | **Points:** 5 | **Status:** Done | **Sprint:** 6 | **Jira:** TBD

管理 IAM2 虚拟身份、Access Key 和全局配置，扩展 IAM 和平台配置能力。

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

### Sprint 5 — SDK 修复 + 覆盖补全 + VM NIC (Apr 28 – May 9) ✅ DONE

| Story | Points | Status |
|-------|--------|--------|
| STORY-TF-001-12 SDK Workaround 修复 + 已知问题 | 5 | Done |
| STORY-TF-001-13 补充 Data Source (EIP/Reserved IP/Subnet IP Range) | 3 | Done |
| STORY-TF-001-14 补充 Resource (SDN Controller + L3 Network) | 5 | Done |
| STORY-TF-001-15 网络域 — VM NIC 管理 | 5 | Done |
| **Total** | **18/30** | |

> Completed ahead of schedule — all work done during Sprint 3–4 timeframe. SDK v0.0.4 resolved all URL placeholder issues.

---

### Sprint 6 — P0 VPN/防火墙/存储/IAM 扩展 (May 12–23) ✅ DONE

| Story | Points | Status |
|-------|--------|--------|
| STORY-TF-001-16 网络域 — IPSec VPN + VPC VPN Gateway | 5 | Done |
| STORY-TF-001-17 网络域 — VPC Firewall + Route Table/Entry | 5 | Done |
| STORY-TF-001-18 存储域 — Ceph + ImageStore Backup Storage | 5 | Done |
| STORY-TF-001-19 系统域 — IAM2 Virtual ID + Access Key + Global Config | 5 | Done |
| **Total** | **20/30** | |

> Completed ahead of schedule. vpc_vpn_gateway skipped (SDK has no Create method). All other resources fully implemented.

---

### Sprint 7 — License 许可证管理 + 测试修复 (TBD)

| Story | Points | Status |
|-------|--------|--------|
| STORY-TF-001-20 系统域 — License 许可证管理 | 5 | Done |
| **Total** | **5/30** | |

**Sprint Goal:** 实现 License 许可证资源管理，补充 Day-0 平台初始化 IaC 能力；修复已有测试问题

**Notes:**
- SDK 有完整 License API（Update/Delete License + Query Authorized Node + Get Authorized Capacity）
- License 是特殊资源：Create = UpdateLicense（上传），无独立 Get API
- 可结合其他 backlog 项（如测试修复、P2/P3 补充资源）充实此 Sprint

---

## Story Inventory (continued)

### STORY-TF-001-20: 系统域 — License 许可证管理

**Priority:** P1 | **Points:** 5 | **Status:** Done | **Sprint:** 7 | **Jira:** TBD

使用 Terraform 管理 ZStack Cloud 许可证，支持许可证上传/更新/删除，查询授权节点和容量。实现 Day-0 平台初始化 IaC 场景。

**Terraform 资源:**
- `resource.zstack_license` — 许可证上传/更新/删除（SDK: UpdateLicense / DeleteLicense）
- `data.zstack_license_authorized_nodes` — 授权节点列表查询（SDK: QueryLicenseAuthorizedNode）
- `data.zstack_license_authorized_capacity` — 授权容量查询（SDK: GetLicenseAuthorizedCapacity）

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
