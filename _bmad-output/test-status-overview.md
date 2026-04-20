# Terraform Provider ZStack — 测试全貌

> **日期**: 2026-04-20
> **基准**: master @ `86dd7b3`
> **总资源数**: 111 resources + 42 data sources

---

## 一、总体指标

| 维度 | 覆盖数 | 总数 | 覆盖率 |
|------|--------|------|--------|
| Resource 测试文件 | 111 | 111 | **100%** |
| Resource Schema 单元测试 | 107 | 111 | 96.4% |
| Resource Metadata 单元测试 | 106 | 111 | 95.5% |
| Resource 验收测试 (Create) | 41 | 111 | **36.9%** |
| Resource Update 测试 | 5 | 111 | 4.5% |
| Resource Import 测试 | 35 | 111 | 31.5% |
| Resource Disappears 测试 | 33 | 111 | 29.7% |
| Data Source 测试文件 | 42 | 42 | **100%** |
| Data Source 验收测试 | 42 | 42 | **100%** |

> **注**: docs/resources 中 `disk_offer`、`instance_offer`、`virtual_router_offer` 是文档命名，代码中为 `disk_offering`、`instance_offering`、`virtual_router_offering`，测试文件存在。
> 真正无测试文件的仅 `script`、`script_execution`、`guest_tools_attachment` 三个资源（docs 命名映射确认）。

---

## 二、已知 Bug

| # | 严重度 | 资源 | 问题 | 状态 | 修复工作量 |
|---|--------|------|------|------|-----------|
| 1 | **CRITICAL** | policy | Create 后 description 仍为 Unknown，Terraform 报错 | OPEN | 1 行 |
| 2 | **CRITICAL** | l2vxlan_network | vni 缺少 RequiresReplace，修改 vni 静默丢失 | OPEN | 5 行 |
| 3 | **CRITICAL** | iam2_project | Expunge 逻辑缺失，name 不可复用 | OPEN | 5 行 |
| 4 | ~~CRITICAL~~ | alarm | stateCheckAlarmDisappears 未定义导致编译失败 | **FIXED** | — |
| 5 | HIGH | policy | 默认 statement `Allow/**` 为 god-mode 全权限 | OPEN | 文档 1 行 |
| 6 | MEDIUM | l3network, instance_scripts | Optional+Computed 字段仅检查 IsNull() 未检查 IsUnknown() | OPEN | 新分支 |
| 7 | MEDIUM | global_config | "Provider produced inconsistent result after apply" | OPEN (预存) | 待排查 |

---

## 三、验收测试覆盖矩阵（41 个有验收测试的资源）

### 完整级（Create + Update + Import + Disappears）

| 资源 | Create | Update | Import | Disappears | 备注 |
|------|--------|--------|--------|------------|------|
| alarm | Y | Y | Y | Y | metric_name 已修正 |
| iam2_project | Y | Y | Y | Y | Expunge bug 未修 |
| image | Y | Y | Y | Y | |
| networking_secgroup | Y | Y | Y | Y | |
| networking_secgroup_rule | Y | Y | Y | Y | |
| tag | Y | Y | Y | Y | |

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
| instance_offering | Y | Y | Y | |
| l2vlan_network | Y | Y | Y | |
| load_balancer | Y | Y | Y | |
| load_balancer_listener | Y | Y | Y | |
| policy | Y | Y | Y | Create description bug |
| port_forwarding_rule | Y | Y | Y | |
| reserved_ip | Y | Y | Y | |
| role | Y | Y | Y | |
| scheduler_trigger | Y | Y | Y | |
| sns_topic | Y | Y | Y | |
| ssh_key_pair | Y | Y | Y | |
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
| access_control_list | Y | - | - | - | Import, Disappears |
| access_key | Y | - | - | - | Import, Disappears |
| global_config | Y | - | - | - | 按设计无 Destroy/Import |
| networking_secgroup_attachment | Y | - | Y | - | Disappears |
| scheduler_job | Y | - | - | Y | Import |
| sns_email_endpoint | Y | - | Y | - | Disappears |
| sns_http_endpoint | Y | - | Y | - | Disappears |
| tag_attachment | Y | - | - | - | Import, Disappears |
| virtual_router_instance | Y | - | - | - | Import, Disappears |

---

## 四、无验收测试的资源（70 个，仅有 Schema/Metadata 单元测试）

### 按功能域分类

**计算与实例 (8)**
- instance, host, baremetal_instance, baremetal_chassis, baremetal_pxe_server
- guest_tool_attachment, vm_cdrom, vm_nic

**存储 (11)**
- backup_storage, ceph_backup_storage, ceph_pool, ceph_primary_storage
- image_store_backup_storage, primary_storage, database_backup
- volume_backup, volume_snapshot, zbox_backup, dataset

**网络 (17)**
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

**其他 (5)**
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
auto_scaling_groups, gpu_devices, l2vlan_networks, load_balancer_listeners, load_balancers, port_forwarding_rules

---

## 六、Sprint 进度

### Phase 0: 代码质量加固（13 个 Story）

| Story | 描述 | 状态 | 分支 |
|-------|------|------|------|
| 01 | Update read-after-write | **已合入 master** | `refactor/read-error-handling` |
| 02 | alarm Disappears + SDK Bug | **已合入 master** | `fix/alarm-update-and-test` |
| 03 | port_forwarding_rule Unknown | 待审查 | `fix/port-forwarding-rule` |
| 04 | iam2_project Expunge | 待审查 | `refactor/provider-quality-hardening` |
| 05 | policy Create Bug | **OPEN** | 未创建 |
| 06 | scheduler_job/global_config Import | 待创建 | — |
| 07 | RequiresReplace + l2vxlan Update | 待创建 | — |
| 08 | reserved_ip/secgroup Import + VR Destroy | 未开始 | — |
| 09 | 批量 Disappears (~26 个) | **已完成** (33个) | 已合入 master |

### Phase 1 Sprint A: 核心资源验收测试

| Story | 描述 | 状态 | 目标资源数 |
|-------|------|------|-----------|
| 10 | 自包含资源 | 未开始 | 8 |
| 11 | 计算+存储 | 未开始 | 7 |
| 12 | 网络+基础设施 | 未开始 | 6 |
| 13 | 监控运维+备份 | 未开始 | 3 |

### 里程碑目标

| 指标 | 当前 | Phase 0 完成后 | Sprint A 完成后 |
|------|------|---------------|----------------|
| 验收覆盖率 | 36.9% (41/111) | ~37% | **56.8%** (63/111) |
| 完整级 (CUID+Dis) | 6 | 6+ | 6+ |
| 标准级 (CID) | 26 | 27+ | ~35 |
| 已知 Bug | 6 | ~3 | ~3 |

---

## 七、基础设施问题

- ZStack API 503 间歇性不可达，影响以下测试：TestAccPortForwardingRuleResource, TestAccIAM2ProjectResource
- TestAccGlobalConfigResource "inconsistent result after apply" 为预存 bug
- 部分资源 Update API 返回 404（ZStack 不支持），已在测试中标注跳过

---

## 八、文件索引

| 文件 | 用途 |
|------|------|
| **本文档** `_bmad-output/test-status-overview.md` | 测试全貌（权威文档） |
| `_bmad-output/review-results.md` | 分支审查详细结果 & Bug 修复方案 |
| `_bmad-output/sprint-overview.md` | Sprint Story 拆解 & 依赖图（已被本文档覆盖） |
| `_bmad-output/story-01~13.md` | 各 Story 详细描述和验收标准 |
| `_bmad-output/workflow-git-push.md` | Git 推送工作流 |
| `.github/ISSUE_TEMPLATE/bug_report.md` | Bug 报告模板 |
