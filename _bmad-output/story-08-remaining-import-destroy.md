# Story 08: 补全 reserved_ip / secgroup_attachment Import 和 virtual_router_instance Destroy

> **分支**: 待创建 (`fix/incomplete-test-dimensions`)  
> **状态**: 未开始  
> **优先级**: P1 — Phase 0.3 收尾  
> **预计工作量**: 0.5 天

---

## Story

作为 Provider 维护者，我需要把最后 3 个"不完整"级资源修复到至少"基础"级，完成 Phase 0.3 收尾。

## 验收标准

- [ ] AC1: `reserved_ip` 测试包含 Import Step，验收通过
- [ ] AC2: `secgroup_attachment` 测试包含 Import Step，验收通过
- [ ] AC3: `virtual_router_instance` 测试包含 CheckDestroy，验收通过
- [ ] AC4: 不完整级资源数从 2 降到 0（secgroup_attachment 和 virtual_router_instance 全部补齐）

## Tasks

### reserved_ip Import

- [ ] 1. 分析 `resource_zstack_reserved_ip.go` ImportState 方法的 ID 格式
- [ ] 2. 在 `resource_zstack_reserved_ip_test.go` 添加 Import Step
- [ ] 3. 跑验收测试确认

### secgroup_attachment Import

- [ ] 4. 分析 `resource_zstack_networking_secgroup_attachment.go` ImportState — 该资源无 UUID，需要 composite ID (secgroup_uuid + vm_nic_uuid)
- [ ] 5. 编写 `importStateIdSecgroupAttachment` 函数
- [ ] 6. 在测试中添加 Import Step
- [ ] 7. 跑验收测试确认

### virtual_router_instance Destroy

- [ ] 8. 在 `check_destroy_test.go` 新增 `testAccCheckVirtualRouterInstanceDestroy`
- [ ] 9. 在测试中添加 CheckDestroy 引用
- [ ] 10. 跑验收测试确认（注：该资源使用 mock server，需确认 mock 支持 destroy 验证）

## 审查要点

1. secgroup_attachment 的 composite ID 格式与 ImportState 解析逻辑是否一致
2. virtual_router_instance 使用 mock server，CheckDestroy 逻辑是否需要适配
