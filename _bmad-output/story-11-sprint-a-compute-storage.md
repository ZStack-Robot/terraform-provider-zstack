# Story 11: Sprint A — 计算和存储资源验收测试 (7 个)

> **分支**: 待创建 (`test/sprint-a-compute-storage`)  
> **状态**: 未开始  
> **优先级**: P1 — Phase 1 Sprint A  
> **前置依赖**: env.json 中有 vm_instances, volumes, images, instance_offerings, l3_networks  
> **预计工作量**: 2 天

---

## Story

作为 Provider 维护者，我需要为计算和存储相关的核心资源新增验收测试。这组资源依赖 env.json 中的基础设施 UUID。

## 资源清单

| # | Resource | Update 状态 | 目标级别 | 环境依赖 |
|---|----------|-----------|---------|---------|
| 1 | `instance` | 有 (UpdateVmInstance) | 标准 | image, l3network, offering |
| 2 | `vm_cdrom` | 有 (UpdateVmCdRom) | 标准 | instance_uuid, image_uuid |
| 3 | `vm_nic` | N/A (not supported) | 基础 | instance_uuid, l3_uuid |
| 4 | `guest_tool_attachment` | N/A (not supported) | 基础 | instance_uuid |
| 5 | `instance_scripts_execution` | N/A (not supported) | 基础 | instance_uuid, script |
| 6 | `volume_snapshot` | 有 (UpdateVolumeSnapshot) | 标准 | volume_uuid |
| 7 | `volume_backup` | N/A (not supported) | 基础 | volume_uuid, bs_uuid |

## 关键风险

- **instance 是最复杂的资源**：涉及 VM 全生命周期、回收站、多属性 Update (name, CPU, memory)
- **instance 的 CheckDestroy 需要处理回收站**：使用策略 2 或 3 (参见 E2E-TEST-PLAN.md 1.4 节)
- **vm_cdrom / vm_nic / guest_tool_attachment** 依赖运行中的 VM 实例
- **volume_snapshot** 需要已有卷

## 验收标准

- [ ] AC1: `instance` 验收测试含 Create + Update (name, 或 CPU/memory) + Import + Destroy (含回收站处理)
- [ ] AC2: `vm_cdrom` / `vm_nic` / `guest_tool_attachment` 各有基础级验收测试
- [ ] AC3: `volume_snapshot` 验收测试含 Create + Update (name) + Import + Destroy
- [ ] AC4: `volume_backup` 验收测试含 Create + Import + Destroy
- [ ] AC5: 所有 Destroy 检查函数正确处理各资源的删除语义
- [ ] AC6: 编译通过，抽样运行验证

## Tasks

### 基础设施

- [ ] 1. 确认 env.json 中 vm_instances、volumes 数据可用
- [ ] 2. 在 `check_destroy_test.go` 新增 7 个 Destroy 检查函数（instance 需特殊处理回收站）

### 计算资源

- [ ] 3. `resource_zstack_instance_test.go` — 最复杂：Create (需 image+l3+offering) + Update (name) + Import + Destroy (expunge)
- [ ] 4. `resource_zstack_vm_cdrom_test.go` — Create + Update (image_uuid) + Import + Destroy
- [ ] 5. `resource_zstack_vm_nic_test.go` — Create + Import + Destroy (U=N/A)
- [ ] 6. `resource_zstack_guest_tool_attachment_test.go` — Create + Import + Destroy (U=N/A)
- [ ] 7. `resource_zstack_instance_scripts_execution_test.go` — Create + Import + Destroy (U=N/A)

### 存储资源

- [ ] 8. `resource_zstack_volume_snapshot_test.go` — Create + Update (name) + Import + Destroy
- [ ] 9. `resource_zstack_volume_backup_test.go` — Create + Import + Destroy (U=N/A)

### 验证

- [ ] 10. 编译确认
- [ ] 11. 抽样运行 instance + volume_snapshot 验收测试

## 审查要点

1. instance 的 CheckDestroy 是否正确处理回收站（检查 state == "Destroyed" 或 expunge 后 404）
2. vm_cdrom / vm_nic 是否在测试中创建自己的 VM 还是依赖 env.json 中已有的 VM
3. volume_snapshot 创建后是否影响原 volume 的状态
