# Terraform Provider ZStack — 测试全貌

> **日期**: 2026-04-20
> **基准**: master @ `86dd7b3`
> **总资源数**: 111 resources + 42 data sources
> **说明**：针对Terraform Provider 对接zstack sdkv2 的全量测试

---

## 一、总体指标

| 维度 | 覆盖数 | 总数 | 覆盖率 | 详见 |
|------|--------|------|--------|------|
| Resource 测试文件 | 111 | 111 | **100%** | — |
| Resource Schema 单元测试 | 107 | 111 | 96.4% | — |
| Resource Metadata 单元测试 | 106 | 111 | 95.5% | — |
| Resource 验收测试 (Create) | 39 <sup>*1</sup> | 111 | **35.1%** | [§三、验收测试覆盖矩阵](#三验收测试覆盖矩阵39-个有验收测试的资源) |
| Resource Update 测试 | 3 <sup>*2</sup> | 79 <sup>*3</sup> | 3.8% | [§三、验收测试覆盖矩阵](#完整级create--update--import--disappears) |
| Resource Import 测试 | 35 | 107 <sup>*4</sup> | 32.7% | [§三、验收测试覆盖矩阵](#三验收测试覆盖矩阵39-个有验收测试的资源) |
| Resource Disappears 测试 | 33 | 104 <sup>*5</sup> | 31.7% | [§三、验收测试覆盖矩阵](#三验收测试覆盖矩阵39-个有验收测试的资源) |
| Data Source 测试文件 | 42 | 42 | **100%** | [§五](#五data-source-测试覆盖) |
| Data Source 验收测试 | 42 | 42 | **100%** | [§五](#五data-source-测试覆盖) |

> ***1 验收测试 (Create) 数量说明**：原始审计计为 41，实际代码中含 `resource.Test` / `resource.ParallelTest`
> 调用的资源为 **39 个**。差异来自 `access_control_list` 和 `access_key`——二者测试文件存在且有
> Schema/Metadata 单元测试，但未编写验收测试（无 `resource.Test` 调用）。已将其从第三节"部分级"
> 移至第四节"无验收测试"。需后续确认是遗漏还是因环境依赖暂未实现。
>
> ***2 Update 测试数量说明**：原始审计计为 5，实际具备多 Step 配置变更的验收测试仅 **3 个**
> （alarm、iam2_project、networking_secgroup_rule）。image 的测试仅含 Create+Import 步骤；
> networking_secgroup 和 tag 在代码注释中明确标注 "Update — N/A"（前者 Provider 未实现
> Update method，后者 ZStack API 拒绝 UpdateTag），均无实际 Update 配置变更。
> 已将这三个资源从第三节"完整级"调整至"标准级"。
>
> ***3 Update 总数说明**：111 个资源中 **32 个** 的 Update 方法为空操作，不触发 API 调用。
> 其中 30 个 Provider 显式返回 `"Update not supported"` 错误（**Provider 层设计决策**——
> 所有可变字段标记 `RequiresReplace`，变更走 destroy+recreate；但部分资源的 ZStack SDK
> 实际提供了 Update API，如 `networking_secgroup` 对应 `UpdateSecurityGroup`）；
> 另外 2 个（`ceph_backup_storage`、`ceph_primary_storage`）为 Provider 层 no-op passthrough。
> 有效可测 Update 的资源为 **79 个**。
>
> ***4 Import 总数说明**：4 个资源未实现 `ImportState` 方法：
> - **资源模型限制**（附着/执行类，非独立实体）：`guest_tool_attachment`、`instance_scripts_execution`、`tag_attachment`
> - **Provider 未实现**：`aliyun_nas_access_group`（SDK 有 CRUD 能力，Provider 未补 Import）
>
> 有效可测 Import 的资源为 **107 个**。
>
> ***5 Disappears 总数说明**：7 个资源的 Delete 为空操作或无法真正销毁：
> - **ZStack API 不支持删除**（代码注释明确说明）：`snmp_agent`、`sns_email_endpoint`、`sns_http_endpoint`、`vpc_firewall`
> - **资源语义限制**（非可销毁实体）：`guest_tool_attachment`（ISO 挂载无 detach API）、`instance_scripts_execution`（执行记录）
> - **资源语义限制**（全局配置无法销毁）：`global_config`（Delete 重置为默认值）
>
> 有效可测 Disappears 的资源为 **104 个**。

> **注**: docs/resources 中 `disk_offer`、`instance_offer`、`virtual_router_offer` 是文档命名，代码中为 `disk_offering`、`instance_offering`、`virtual_router_offering`，测试文件存在。
> 真正无测试文件的仅 `script`、`script_execution`、`guest_tools_attachment` 三个资源（docs 命名映射确认）。

---

## 二、Bug 概况

> 详细跟踪见 [bug-tracker.md](bug-tracker.md)

| ID | 严重度 | 资源 | 一句话描述 | 状态 |
|----|--------|------|-----------|------|
| [BUG-1](bug-tracker.md#bug-1) | **CRITICAL** | policy | Create 后 description 为 Unknown | OPEN |
| [BUG-2](bug-tracker.md#bug-2) | **CRITICAL** | l2vxlan_network | vni 缺少 RequiresReplace，静默数据丢失 | OPEN |
| [BUG-3](bug-tracker.md#bug-3) | **CRITICAL** | iam2_project | Expunge 逻辑缺失，name 不可复用 | FIXING |
| [BUG-4](bug-tracker.md#bug-4) | HIGH | policy | Statements 空数组 + 修复引入 god-mode | FIXING |
| [BUG-5](bug-tracker.md#bug-5) | MEDIUM | 多资源 | Optional+Computed IsNull guard 不完整 | FIXING |
| [BUG-6](bug-tracker.md#bug-6) | MEDIUM | global_config | 验收测试 inconsistent result | OPEN |
| ~~BUG-7~~ | ~~CRITICAL~~ | alarm | stateCheckAlarmDisappears 未定义 | FIXED |
| ~~[BUG-8](bug-tracker.md#bug-8)~~ | ~~MEDIUM~~ | 4 资源 | 错误的 RequiresReplace 导致不必要 destroy | FIXED |

---

## 三、验收测试覆盖矩阵（39 个有验收测试的资源）

### 完整级（Create + Update + Import + Disappears）

| 资源 | Create | Update | Import | Disappears | 备注 |
|------|--------|--------|--------|------------|------|
| alarm | Y | Y | Y | Y | |
| iam2_project | Y | Y | Y | Y | [BUG-3](bug-tracker.md#bug-3) Expunge 缺失 |
| networking_secgroup_rule | Y | Y | Y | Y | |

### 标准级（Create + Import + Disappears，无 Update）

| 资源 | Create | Import | Disappears | 备注 |
|------|--------|--------|------------|------|
| account | Y | Y | Y | |
| affinity_group | Y | Y | Y | |
| auto_scaling_group | Y | Y | Y | |
| certificate | Y | Y | Y | |
| cluster | Y | Y | Y | |
| disk_offering | Y | Y | Y | |
| iam2_organization | Y | Y | Y | |
| iam2_virtual_id | Y | Y | Y | |
| image | Y | Y | Y | |
| instance_offering | Y | Y | Y | |
| l2vlan_network | Y | Y | Y | |
| load_balancer | Y | Y | Y | |
| load_balancer_listener | Y | Y | Y | |
| networking_secgroup | Y | Y | Y | Update N/A: Provider 未实现 Update method |
| policy | Y | Y | Y | [BUG-1](bug-tracker.md#bug-1) description Unknown |
| port_forwarding_rule | Y | Y | Y | |
| reserved_ip | Y | Y | Y | |
| role | Y | Y | Y | |
| scheduler_trigger | Y | Y | Y | |
| sns_topic | Y | Y | Y | |
| ssh_key_pair | Y | Y | Y | |
| tag | Y | Y | Y | Update N/A: ZStack API 拒绝 UpdateTag |
| user | Y | Y | Y | |
| vip | Y | Y | Y | |
| virtual_router_image | Y | Y | Y | |
| virtual_router_offering | Y | Y | Y | |
| volume | Y | Y | Y | |
| webhook | Y | Y | Y | |
| zone | Y | Y | Y | |

### 部分级（Create，缺少 Import 或 Disappears）

| 资源 | Create | Update | Import | Disappears | 缺少项 |
|------|--------|--------|--------|------------|--------|
| global_config | Y | - | - | - | 按设计无 Destroy/Import; [BUG-6](bug-tracker.md#bug-6) |
| networking_secgroup_attachment | Y | - | Y | - | Disappears |
| scheduler_job | Y | - | - | Y | Import |
| sns_email_endpoint | Y | - | Y | - | Disappears |
| sns_http_endpoint | Y | - | Y | - | Disappears |
| tag_attachment | Y | - | - | - | Import, Disappears |
| virtual_router_instance | Y | - | - | - | Import, Disappears |

---

## 四、无验收测试的资源（72 个，仅有 Schema/Metadata 单元测试）

### 按功能域分类

**计算与实例 (8)**
- instance, host, baremetal_instance, baremetal_chassis, baremetal_pxe_server
- guest_tool_attachment, vm_cdrom, vm_nic

**认证与访问控制 (2)** †1
- access_control_list, access_key

**存储 (11)**
- backup_storage, ceph_backup_storage, ceph_pool, ceph_primary_storage
- image_store_backup_storage, primary_storage, database_backup
- volume_backup, volume_snapshot, zbox_backup, dataset

**网络 (19)**
- l2vxlan_network, l3network, eip, subnet_ip_range, ipsec_connection
- multicast_router, port_mirror, port_mirror_session
- policy_route_rule, policy_route_rule_set, lb_server_group
- sdn_controller, vpc, vpc_firewall, vpc_ha_group, vpc_shared_qos
- vrouter_route_entry, vrouter_route_table, vip_qos

**安全设备 (5)**
- fi_sec_security_machine, flk_sec_security_machine, info_sec_security_machine
- jit_security_machine, san_sec_security_machine

**监控与运维 (7)**
- monitor_group, monitor_template, preconfiguration_template, price_table
- email_media, log_server, snmp_agent

**混合云 & 第三方 (6)**
- aliyun_nas_access_group, aliyun_proxy_vpc, aliyun_proxy_vswitch
- v2v_conversion_host, vcenter, container_management_endpoint

**脚本与自动化 (4)**
- instance_scripts, instance_scripts_execution, resource_stack, stack_template

**其他 (10)**
- cdp_policy, cdp_task, directory, flow_collector, flow_meter
- iscsi_server, ldap_server, license, nvme_server, pci_device_offering

---

## 五、Data Source 测试覆盖

| 维度 | 覆盖数 | 总数 | 覆盖率 |
|------|--------|------|--------|
| 测试文件 | 42 | 42 | 100% |
| 验收测试 | 42 | 42 | 100% |
| Schema 单元测试 | 6 | 42 | 14.3% |
| Metadata 单元测试 | 6 | 42 | 14.3% |

有 Schema+Metadata 单元测试的 Data Source（6 个）:
auto_scaling_groups, gpu_devices, l2vlan_networks, load_balancer_listeners, load_balancers, port_forwarding_rules。这些函数只验证字段声明（名字、Required/Computed 标记），验收测试已经端到端覆盖了这些。


---

## 六、Sprint 进度

### Phase 0: 代码质量加固（13 个 Story）

| Story | 描述 | 状态 | P | 分支 | 合入序 |
|-------|------|------|---|------|--------|
| 01 | Update read-after-write | **已合入** | P0 | `refactor/read-error-handling` | — |
| 02 | alarm Disappears + SDK Bug | **已合入** | P1 | `fix/alarm-update-and-test` | — |
| <a id="story-03"></a>03 | port_forwarding_rule Unknown | 待审查 | P0 | `fix/port-forwarding-rule` | 独立; [BUG-5](bug-tracker.md#bug-5) |
| 04 | iam2_project Expunge | 待审查（有审查缺陷） | P1 | `refactor/provider-quality-hardening` | 依赖01✅; [BUG-3](bug-tracker.md#bug-3) |
| 05 | policy Create Bug | 待审查（有审查缺陷） | P0 | `fix/policy-empty-statements` | 独立; [BUG-1](bug-tracker.md#bug-1) [BUG-4](bug-tracker.md#bug-4) |
| 06 | scheduler_job/global_config Import | 待审查 | P1 | `test/composite-id-import-steps` | 独立 |
| 07 | RequiresReplace + l2vxlan Update | 部分完成 | P1 | `feat/l2vxlan-network-update` + 4个RR分支([BUG-8](bug-tracker.md#bug-8)) | 独立 |
| 08 | reserved_ip/secgroup Import + VR Destroy | 未开始 | P1 | — | 独立 |
| 09 | 批量 Disappears (~26 个) | **已完成** (33个) | P2 | 已合入 master | — |

### Phase 1 Sprint A: 核心资源验收测试

| Story | 描述 | 状态 | P | 资源数 | 合入序 |
|-------|------|------|---|--------|--------|
| 10 | 自包含资源 | 未开始 | P1 | 8 | 依赖07 |
| 11 | 计算+存储 | 未开始 | P1 | 7 | 独立 |
| 12 | 网络+基础设施 | 未开始 | P1 | 6 | 依赖07(l2vxlan) |
| 13 | 监控运维+备份 | 未开始 | P1 | 3 | 依赖10(monitor_template) |

### 依赖图

```
Phase 0.4 (代码质量)                    Phase 0.3      Phase 0.2
┌──────────────────────┐              ┌─────────┐    ┌─────────┐
│ story-01 ✅ (read-after-write)──┐   │ story-08 │───→│ story-09 ✅│
│   ├→ story-02 ✅ (alarm dis)   │   │ (Import/ │    │ (批量      │
│   └→ story-04 (iam2 expunge)   │   │  Destroy)│    │ Disappears) │
├──────────────────────┤         │   └─────────┘    └───────────┘
│ story-03 (port-fwd)  │ 独立    │
│ story-05 (policy)    │ 独立    │
│ story-06 (import)    │ 独立    │
│ story-07 (RR+l2vxlan)│─────────┼────────────────────┐
└──────────────────────┘         │                    │
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

### 推荐执行顺序

**第一批（并行）**: story-03, 05, 06 — 独立且小
**第二批（顺序）**: story-04, 07 — 依赖 01 已完成，可直接开始
**第三批**: story-08 — Phase 0 收尾
**第四批（并行）**: story-10, 11, 12, 13 — Phase 1 Sprint A

### 里程碑目标

| 指标 | 当前 | Phase 0 完成后 | Sprint A 完成后 |
|------|------|---------------|----------------|
| 验收覆盖率 | 35.1% (39/111) | ~36% | **55.0%** (61/111) |
| 完整级 (CUID+Dis) | 3 | 3+ | 3+ |
| 标准级 (CID) | 29 | 30+ | ~38 |
| 已知 Bug (OPEN+FIXING) | 6 | ~3 | ~3 |

---

## 七、文件索引

| 文件 | 用途 |
|------|------|
| **本文档** `test-status-overview.md` | 测试全貌 — 覆盖率、矩阵、Sprint 进度 |
| [`bug-tracker.md`](bug-tracker.md) | **Bug 唯一跟踪源** — 详情、根因、修复思路、进度 |
| `review-results.md` | 分支审查历史存档（2026-04-17） |
| `story-01~13.md` | 各 Story 详细描述和验收标准 |
| `workflow-git-push.md` | Git 推送工作流 |
| `.github/ISSUE_TEMPLATE/bug_report.md` | Bug 报告模板 |
