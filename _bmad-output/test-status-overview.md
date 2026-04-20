# Terraform Provider ZStack — 测试全貌

> **日期**: 2026-04-20
> **基准**: master @ `86dd7b3`
> **总资源数**: 111 resources + 42 data sources
> **说明**：针对Terraform Provider 对接zstack sdkv2 的全量测试

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

## 二、Bug 跟踪

### 待修复 — 阻塞或有数据风险

<a id="bug-1"></a>
**BUG-1 [CRITICAL] policy — Create 后 description 为 Unknown**

- **关联**: Story-05 | **位置**: `resource_zstack_policy.go` Create 方法
- **现象**: 用户省略 description 时，Create 成功但 Terraform 报错 "provider still indicated an unknown value"
- **根因**: description 为 `Optional+Computed` + `UseStateForUnknown()`，API 不返回 description，Create 未设置具体值
- **修复** (1 行):
  ```go
  if plan.Description.IsUnknown() || plan.Description.IsNull() {
      plan.Description = types.StringValue("")
  }
  ```

<a id="bug-2"></a>
**BUG-2 [CRITICAL] l2vxlan_network — vni 缺少 RequiresReplace**

- **关联**: Story-07 | **位置**: `resource_zstack_l2vxlan_network.go:90-94`
- **现象**: 用户修改 vni 后 Terraform 显示成功，但实际 vni 未变更（静默数据丢失）
- **根因**: vni 为 `Optional+Computed` 但无 RequiresReplace，Update 只发送 name/description。对比 `l2vlan_network.vlan` 已正确标记
- **修复** (5 行):
  ```go
  "vni": schema.Int64Attribute{
      Optional: true, Computed: true,
      Description: "The VXLAN Network Identifier (VNI).",
      PlanModifiers: []planmodifier.Int64{int64planmodifier.RequiresReplace()},
  },
  ```

<a id="bug-3"></a>
**BUG-3 [CRITICAL] iam2_project — Expunge 逻辑缺失**

- **关联**: Story-04 | **位置**: `resource_zstack_iam2_project.go` Delete 方法
- **现象**: Delete 成功后 name 无法复用（ZStack 要求 Expunge 后才释放 name）
- **根因**: Delete 仅调用 DeleteIAM2Project，未调用 ExpungeIAM2Project
- **修复** (5 行): Delete 成功后调用 Expunge，失败时 AddWarning 而非 return（不阻止 state 清理）
  ```go
  if err := r.client.ExpungeIAM2Project(state.Uuid.ValueString()); err != nil {
      resp.Diagnostics.AddWarning("Warning expunging IAM2 Project",
          "Delete succeeded but expunge failed: "+err.Error())
  }
  ```

### 待修复 — 低优先级

<a id="bug-4"></a>
**BUG-4 [HIGH] policy — 默认 statement 为全权限 god-mode**

- **位置**: `resource_zstack_policy.go:120-126`
- **现象**: 每个 policy 都获得 `Effect: "Allow", Actions: ["**"]`，用户无法控制
- **处置**: 最低要求在 schema Description 中文档化此行为（1 行），完整修复需支持 statements 作为属性

<a id="bug-5"></a>
**BUG-5 [MEDIUM] Optional+Computed IsNull guard 不完整**

- **关联**: [Story-03](#story-03) | **同类 bug，跨多个资源**
- **已修复**: `port_forwarding_rule`（3 个 Int64 字段）— 分支 `fix/port-forwarding-rule` 待合入，详见 [Story-03](#story-03)
- **未修复**: `l3network`（IpVersion, System）、`instance_scripts`（ScriptTimeout）
- **现象**: Optional+Computed 字段只检查 `IsNull()` 未检查 `IsUnknown()`，偶然安全但与已修复模式不一致
- **处置**: port_forwarding_rule 随 Story-03 合入；其余后续创建 `fix/optional-computed-isnull-guard` 分支统一修复

<a id="bug-6"></a>
**BUG-6 [MEDIUM] global_config — 验收测试 inconsistent result**

- **现象**: TestAccGlobalConfigResource "Provider produced inconsistent result after apply"
- **状态**: 预存 bug，master 上同样失败，需排查 Read 返回值与 plan 不一致点

### 已修复

| ID | 资源 | 问题 | 修复于 |
|----|------|------|--------|
| ~~BUG-7~~ | alarm | stateCheckAlarmDisappears 未定义致编译失败 | master @ `86dd7b3` (Story-02) |

### 基础设施问题（非代码 bug）

- **ZStack API 503 间歇性不可达** — 影响 TestAccPortForwardingRuleResource, TestAccIAM2ProjectResource。环境: 172.24.248.129:8080
- **部分资源 Update API 返回 404** — ZStack 不支持，已在测试中标注跳过（非 provider bug）

---

## 三、验收测试覆盖矩阵（41 个有验收测试的资源）

### 完整级（Create + Update + Import + Disappears）

| 资源 | Create | Update | Import | Disappears | 备注 |
|------|--------|--------|--------|------------|------|
| alarm | Y | Y | Y | Y | |
| iam2_project | Y | Y | Y | Y | [BUG-3](#bug-3) Expunge 缺失 |
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
| policy | Y | Y | Y | [BUG-1](#bug-1) description Unknown |
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
| global_config | Y | - | - | - | 按设计无 Destroy/Import; [BUG-6](#bug-6) |
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

| Story | 描述 | 状态 | P | 分支 | 合入序 |
|-------|------|------|---|------|--------|
| 01 | Update read-after-write | **已合入** | P0 | `refactor/read-error-handling` | — |
| 02 | alarm Disappears + SDK Bug | **已合入** | P1 | `fix/alarm-update-and-test` | — |
| <a id="story-03"></a>03 | port_forwarding_rule Unknown | 待审查 | P0 | `fix/port-forwarding-rule` | 独立; [BUG-5](#bug-5) |
| 04 | iam2_project Expunge | 待审查 | P1 | `refactor/provider-quality-hardening` | 依赖01✅ |
| 05 | policy Create Bug | OPEN | P0 | 未创建 | 独立 |
| 06 | scheduler_job/global_config Import | OPEN | P1 | 未创建 | 独立 |
| 07 | RequiresReplace + l2vxlan Update | OPEN | P1 | 未创建 | 独立 |
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
| 验收覆盖率 | 36.9% (41/111) | ~37% | **56.8%** (63/111) |
| 完整级 (CUID+Dis) | 6 | 6+ | 6+ |
| 标准级 (CID) | 26 | 27+ | ~35 |
| 已知 Bug | 6 | ~3 | ~3 |

---

## 七、文件索引

| 文件 | 用途 |
|------|------|
| **本文档** `_bmad-output/test-status-overview.md` | 测试全貌（唯一权威文档） |
| `_bmad-output/review-results.md` | 分支审查详细结果 & Bug 修复方案 |
| `_bmad-output/story-01~13.md` | 各 Story 详细描述和验收标准 |
| `_bmad-output/workflow-git-push.md` | Git 推送工作流 |
| `.github/ISSUE_TEMPLATE/bug_report.md` | Bug 报告模板 |
