# SDK v2 URL Template Bug — `ZSClient.Post()` 不解析路径占位符

## 问题描述

`zstack-sdk-go-v2` 的 `other_actions.go` 中大量 `ZSClient` 方法使用了 URL 路径模板占位符（如 `{l2NetworkUuid}`、`{clusterUuid}`），但 `ZSClient.Post()` **不会解析这些占位符**，导致请求发送到带有字面 `{xxx}` 的 URL，API 调用必然失败。

### 根因分析

SDK 存在两套 HTTP client：

| Client | Post 方法 | URL 构建方式 |
|--------|----------|-------------|
| `ZSHttpClient`（旧） | `ZSHttpClient.Post(resource, params, retVal)` | 通过 `getPostURL()` → `getURL()` 正确拼接 |
| `ZSClient`（新） | `ZSClient.Post(path, params, result)` | 直接 `fmt.Sprintf("%s/%s", baseURL, path)` **不做模板替换** |

`ZSClient` 嵌入了 `*ZSHttpClient`，但自身重新定义了 `Post()` 方法（签名相同），因此**遮蔽了** `ZSHttpClient.Post()`。

`other_actions.go` 中的方法都定义在 `*ZSClient` 上，调用 `cli.Post(...)` 时走的是 `ZSClient.Post()` 而非 `ZSHttpClient.Post()`，URL 模板永远不会被解析。

### 影响范围

```
$ grep -c 'cli\.Post("v1/.*{' other_actions.go
321   ← 321 个方法含未解析的 URL 模板
```

### 绕过方案

直接调用嵌入的 `ZSHttpClient.Post()` 方法，手动拼接 URL：

```go
// ❌ 不工作 — URL 模板不会被解析
r.client.AttachL2NetworkToCluster(params)

// ✅ 绕过 — 使用 ZSHttpClient.Post() 直接拼接 URL
var resp view.L2NetworkInventoryView
err := r.client.ZSHttpClient.Post(
    fmt.Sprintf("v1/l2-networks/%s/clusters/%s", l2NetworkUuid, clusterUuid),
    params,
    &resp,
)
```

---

## 按资源 Workaround 记录

> **SDK 修复后**，需根据此表逐个回退为直接调用 SDK 方法。
> 搜索关键字：`ZSHttpClient.Post(` 可快速定位所有 workaround 代码。

### 已绕过的资源

| 资源 | 文件 | 绕过方法 | 绕过的 SDK 方法 | SDK 原始 URL 模板 | 回退方式 |
|------|------|---------|----------------|------------------|---------|
| `zstack_l2vlan_network` | `resource_zstack_l2vlan_network.go` | `attachCluster()` | `AttachL2NetworkToCluster()` | `v1/l2-networks/{l2NetworkUuid}/clusters/{clusterUuid}` | 改回 `r.client.AttachL2NetworkToCluster(params)` |
| `zstack_port_forwarding_rule` | `resource_zstack_port_forwarding_rule.go` | `attachToVmNic()` | `AttachPortForwardingRule()` | `v1/port-forwarding/{ruleUuid}/vm-instances/nics/{vmNicUuid}` | 改回 `r.client.AttachPortForwardingRule(params)` |
| `zstack_load_balancer_listener` | `resource_zstack_load_balancer_listener.go` | `Create()` | `CreateLoadBalancerListener()` | `v1/load-balancers/{loadBalancerUuid}/listeners` | 改回 `r.client.CreateLoadBalancerListener(params)` |

### 未绕过但存在同样问题的资源

| 资源 | 文件 | 受影响方法 | 对应 SDK 方法 | SDK 原始 URL 模板 | 说明 |
|------|------|----------|-------------|------------------|------|
| `zstack_volume` | `resource_zstack_volume.go` | `Create()` L208, `reconcileAttachment()` L356 | `AttachDataVolumeToVm()` | `v1/volumes/{volumeUuid}/vm-instances/{vmInstanceUuid}` | volume attach 功能当前不可用，需同样绕过 |

### 后续新增资源注意

Sprint 3/4 中以下资源如果涉及 attach 操作，大概率需要同样绕过：

| 计划资源 | 可能涉及的 SDK attach 方法 |
|---------|-------------------------|
| `zstack_load_balancer_listener` | ~~`CreateLoadBalancerListener`~~ (已绕过), `AddBackendServerToServerGroup` 等 |
| `zstack_scheduler_trigger` | `AddSchedulerJobToTrigger` |
| `zstack_primary_storage` | `AttachPrimaryStorageToCluster` |
| `zstack_backup_storage` | `AttachBackupStorageToZone` |

---

**发现日期**: 2026-03-26
**SDK 版本**: `github.com/zstackio/zstack-sdk-go-v2 v0.0.2-0.20260303063956-a91754e8617b`
**严重程度**: High — 影响所有使用 URL 模板的 attach/操作类 API
**跟踪方式**: SDK 修复后搜索 `ZSHttpClient.Post(` 定位所有 workaround 点
