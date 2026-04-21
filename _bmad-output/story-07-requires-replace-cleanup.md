# Story 07: 移除不合理的 RequiresReplace + l2vxlan_network Update 实现

> **分支**: 待创建 (`fix/remove-requires-replace`)  
> **状态**: WIP (工作目录未提交)  
> **优先级**: P1  
> **预计审查时间**: 15 分钟

---

## Story

作为 Provider 用户，当我修改 monitor_group、monitor_template、preconfiguration_template、price_table 的 name/description 时，Terraform 应执行 in-place update 而非 destroy+recreate。此外，l2vxlan_network 之前 Update 返回 "not supported" 错误，需要实现真正的 Update 方法。

## 背景

- commit `eeaf3e8f` (2026-04-09 compliance audit) 批量添加了 RequiresReplace，其中部分不合理
- 详见 `docs/REQUIRES-REPLACE-AUDIT.md` 审计报告
- l2vxlan_network Update 原来返回 "not supported" 错误，但 ZStack 有 `UpdateL2Network` API 可用

## 验收标准

- [ ] AC1: `monitor_group` name 属性无 `RequiresReplace` PlanModifier
- [ ] AC2: `monitor_template` name 属性无 `RequiresReplace` PlanModifier
- [ ] AC3: `preconfiguration_template` name/description/distribution/type/content 属性无 `RequiresReplace` PlanModifier
- [ ] AC4: `price_table` name 属性无 `RequiresReplace` PlanModifier
- [ ] AC5: `l2vxlan_network` Update 方法调用 `UpdateL2Network` API 并使用 read-after-write 模式回读
- [ ] AC6: 以上 4 个资源的 description 属性也无不合理的 `RequiresReplace`（检查并清理）
- [ ] AC7: 编译通过，Schema 单元测试通过

## Tasks

### RequiresReplace 清理

- [x] 1. `resource_zstack_monitor_group.go` — 移除 name 和 description 的 `stringplanmodifier.RequiresReplace()`
- [x] 2. `resource_zstack_monitor_template.go` — 同上
- [x] 3. `resource_zstack_preconfiguration_template.go` — 移除 name/description/distribution/type/content 的 RequiresReplace
- [x] 4. `resource_zstack_price_table.go` — 移除 name 和 description 的 RequiresReplace

### l2vxlan_network Update 实现

- [x] 5. `resource_zstack_l2vxlan_network.go` — 将 Update 方法从返回 "not supported" 错误改为实际调用 `UpdateL2Network` API
- [x] 6. 确保 Update 使用 read-after-write 模式（`UpdateL2Network` → `QueryL2VxlanNetwork` → 写 state）
- [ ] 7. 提交到新分支 `fix/remove-requires-replace`
- [ ] 8. 编译确认 + Schema 测试通过

## 涉及文件

| 文件 | 变更类型 |
|------|---------|
| `zstack/provider/resource_zstack_monitor_group.go` | 修改 (-4 行) |
| `zstack/provider/resource_zstack_monitor_template.go` | 修改 (-4 行) |
| `zstack/provider/resource_zstack_preconfiguration_template.go` | 修改 (-13 行) |
| `zstack/provider/resource_zstack_price_table.go` | 修改 (-4 行) |
| `zstack/provider/resource_zstack_l2vxlan_network.go` | 修改 (+49 -4 行) |

## 审查要点

1. **RequiresReplace 移除**: 移除后这些资源的 Update 方法是否已实现？如果 Update 未实现只移除 RR，修改属性会导致 state 不一致
2. **l2vxlan_network Update**: 使用 `UpdateL2Network`（L2 Network 通用接口）而非 L2VxlanNetwork 专用接口是否正确？
3. **l2vxlan_network 字段覆盖**: Update 后回读时是否覆盖了所有需要的字段（Vni, PoolUuid, ZoneUuid, PhysicalInterface, Type）
4. 本 story 中的 4 个 RR 清理资源目前无验收测试（仅 Schema 单元测试），后续 Phase 1 Sprint A 新增验收测试时需验证 Update 行为
