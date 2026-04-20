# Story 12: Sprint A — 网络和基础设施资源验收测试 (6 个)

> **分支**: 待创建 (`test/sprint-a-network-infra`)  
> **状态**: 未开始  
> **优先级**: P1 — Phase 1 Sprint A  
> **前置依赖**: env.json 中有 l2_networks, l3_networks, zones, clusters, vips  
> **预计工作量**: 1.5 天

---

## Story

作为 Provider 维护者，我需要为网络和基础设施管理相关的资源新增验收测试。

## 资源清单

| # | Resource | Update 状态 | 目标级别 | 环境依赖 |
|---|----------|-----------|---------|---------|
| 1 | `l3network` | 有 (UpdateL3Network) | 标准 | l2_network_uuid |
| 2 | `subnet_ip_range` | N/A (not supported) | 基础 | l3_network_uuid |
| 3 | `l2vxlan_network` | 有 (UpdateL2Network, WIP story-07) | 标准 | zone_uuid, cluster |
| 4 | `eip` | N/A (not supported) | 基础 | vip_uuid, vm_nic_uuid |
| 5 | `host` | 有 (UpdateHost + UpdateKvmHost) | 标准 | cluster_uuid |
| 6 | `primary_storage` | 有 (UpdatePrimaryStorage + cluster attach/detach) | 标准 | zone_uuid |

## 关键风险

- **host**: 添加主机会改变集群状态，测试需确保不影响其他测试
- **primary_storage**: 涉及 cluster attach/detach，需小心不要误操作共享环境的存储
- **l2vxlan_network**: 依赖 SDN controller 或 VXLAN pool，环境可能不具备
- **eip**: 需要 VIP + VM NIC，依赖链较长

## 验收标准

- [ ] AC1: `l3network` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC2: `subnet_ip_range` 验收测试含 Create + Import + Destroy
- [ ] AC3: `l2vxlan_network` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC4: `eip` 验收测试含 Create + Import + Destroy (U=N/A)
- [ ] AC5: `host` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC6: `primary_storage` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC7: 编译通过，抽样运行验证

## Tasks

### 网络资源

- [ ] 1. `resource_zstack_l3network_test.go` — Create (需 l2_network_uuid) + Update (name) + Import + Destroy
- [ ] 2. `resource_zstack_subnet_ip_range_test.go` — Create (需 l3_network_uuid + IP range) + Import + Destroy
- [ ] 3. `resource_zstack_l2vxlan_network_test.go` — Create (需 zone_uuid + pool_uuid) + Update (name) + Import + Destroy
- [ ] 4. `resource_zstack_eip_test.go` — Create (需 vip_uuid) + Import + Destroy

### 基础设施管理

- [ ] 5. `resource_zstack_host_test.go` — Create (需 cluster_uuid + host IP/credentials) + Update (name) + Import + Destroy
- [ ] 6. `resource_zstack_primary_storage_test.go` — Create (需 zone_uuid) + Update (name) + Import + Destroy

### 基础设施

- [ ] 7. 在 `check_destroy_test.go` 新增 6 个 Destroy 检查函数
- [ ] 8. 确认 env.json 中相关数据可用

### 验证

- [ ] 9. 编译确认
- [ ] 10. 抽样运行 l3network + l2vxlan_network 验收测试

## 审查要点

1. host 测试需要真实主机 IP 和 SSH 凭据，是否需要从 env.json 获取？还是 t.Skip?
2. primary_storage 创建什么类型（LocalStorage / NFS / Ceph）？是否需要特定后端？
3. l2vxlan_network 是否有 VXLAN pool 可用？不可用时 t.Skip
4. eip 的 VIP 和 VM NIC 引用方式
