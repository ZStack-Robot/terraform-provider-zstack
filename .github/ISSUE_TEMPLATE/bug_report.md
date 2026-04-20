---
name: Bug Report
about: Terraform Provider for ZStack Cloud bug 提交模板
---

## 标题格式

`[资源类型] 简短描述问题现象`

示例：
- [zstack_instance] Create 后 state 中 memory_size 为 null
- [zstack_volume_snapshot] Import 后 plan 显示 drift
- [data_source/zstack_l3networks] 过滤条件 name_pattern 不生效

---

## Epic (可选)

关联 Epic: （如 [ZSTAC-84425](http://jira.zstack.io/browse/ZSTAC-84425)）

## 影响

- Resource/DataSource 名称: `zstack_xxx`
- 对应源文件: `zstack/provider/resource_zstack_xxx.go`

## 环境

### ZStack 环境

| 项目 | 值 |
|------|-----|
| ZStack Cloud 版本 | |
| 管理节点 IP | |
| Web UI 地址 | (如 http://172.20.x.x:5000) |
| MN SSH 用户/密码 | (如 root / password) |

### Terraform 执行环境

| 项目 | 值 |
|------|-----|
| Provider 版本 | (git commit hash 或 tag) |
| Provider 来源 | registry / 本地编译 / dev_overrides |
| Terraform 版本 | |
| Go 版本 | (如果是编译问题) |
| OS / Arch | (如 CentOS 7.9 x86_64) |

<details>
<summary>快速获取环境信息</summary>

```bash
# Terraform 执行环境（本地执行）
terraform version          # Terraform + Provider 版本
uname -a                   # OS / Arch
go version                 # Go 版本（编译问题时需要）

# ZStack 版本（SSH 到管理节点）
zstack-ctl status          # ZStack Cloud 版本
```

</details>

## 操作阶段

- [ ] terraform plan
- [ ] terraform apply (Create)
- [ ] terraform apply (Update)
- [ ] terraform apply (Delete)
- [ ] terraform import
- [ ] terraform refresh / state drift
- [ ] terraform destroy

## 复现步骤

1. 使用以下 HCL 配置:

```hcl
resource "zstack_xxx" "test" {
  # 最小复现配置
}
```

2. 执行 `terraform apply`
3. 观察到...

复现率: ___% （100 = 必现，填写大致比例即可）

## 预期行为

描述正确情况应该是什么。

## 实际行为

描述实际发生了什么。

## 关键日志

<details>
<summary>TF_LOG=DEBUG 输出（点击展开）</summary>

执行 `TF_LOG=DEBUG terraform apply 2>&1 | tail -100` 贴到这里：

```
```

</details>

## State 片段 (state drift / import 类 bug 必填)

`terraform state show zstack_xxx.test` 的输出：

```json
```

## 根因分类 (如果已定位)

- [ ] Schema 定义错误 (类型、Required/Optional/Computed 标记)
- [ ] CRUD 逻辑错误 (Create/Read/Update/Delete handler)
- [ ] State 映射错误 (API 响应 → Terraform state 字段映射)
- [ ] Import 实现错误 (ImportState 函数)
- [ ] SDK 调用参数错误 (传给 zstack-sdk-go-v2 的参数)
- [ ] API 兼容性问题 (ZStack API 行为与 Provider 假设不一致)
- [ ] 过滤/查询逻辑错误 (Data Source 的 filter 实现)
- [ ] 并发/依赖问题 (资源间依赖顺序)

## 严重程度

- [ ] P0 - 阻塞：资源完全不可用或造成数据丢失
- [ ] P1 - 严重：核心 CRUD 某个操作失败
- [x] P2 - 一般：非核心功能或 edge case（默认）
- [ ] P3 - 低：文档/命名/日志等非功能性问题

## 修复建议 (可选)

如已定位根因，简述修复方向：
- 文件: `zstack/provider/resource_zstack_xxx.go`
- 行数: L123
- 修改内容: ...
