# SDK Bug Report: UpdateAlarm 返回空值导致 State 丢失

> 发现日期: 2026-04-15  
> 验证日期: 2026-04-17（多环境复现 + 修复验证）  
> 影响范围: `zstack_alarm` 资源的 Update 操作  
> SDK 版本: `zstackio/zstack-sdk-go-v2 v0.0.4`  
> 严重度: **高** — Update 后 Terraform state 全部归零，apply 必定失败  
> 修复分支: `fix/alarm-update-and-test`

---

## 1. 问题现象

对 `zstack_alarm` 资源执行 `terraform apply`（修改 name/description）时，Terraform 报错：

```
Error: Provider produced inconsistent result after apply

When applying changes to zstack_alarm.test, provider produced an unexpected new value:
  .uuid:                was "688a692a452d4dbd95a40fb3e36e5702", but now ""
  .name:                was "acc-test-alarm-updated",             but now ""
  .namespace:           was "ZStack/VM",                          but now ""
  .metric_name:         was "CPUAverageUsedUtilization",          but now ""
  .comparison_operator: was "GreaterThanOrEqualTo",               but now ""
  .threshold:           was 90,                                   but now 0
  .period:              was 60,                                   but now 0
  .repeat_interval:     was 600,                                  but now 0
  .description:         was "Updated acceptance test alarm",      but now null
```

**所有字段都变成了零值/空字符串**，而非 Update 后的正确值。

---

## 2. 暴露过程

在为 `zstack_alarm` 补充 Update Step 验收测试时发现此问题。测试流程：

```
Step 1: Create alarm (name="acc-test-alarm")              → PASS
Step 2: Update alarm (name="acc-test-alarm-updated")      → FAIL (inconsistent state)
Step 3: Import                                             → 未执行
```

Step 1 证明 Create 正常，Step 2 证明 Update API 调用本身成功（ZStack 返回了 200），但 Provider 写回 Terraform state 的数据全部为零值。

---

## 3. 排查过程

### 3.1 确认 API 调用成功

从测试日志中提取 Update 的 HTTP 交互：

**请求** (Provider → ZStack):
```
PUT /zstack/v1/zwatch/alarms/688a692a.../actions
{
  "updateAlarm": {
    "name": "acc-test-alarm-updated",
    "description": "Updated acceptance test alarm",
    "comparisonOperator": "GreaterThanOrEqualTo",
    "threshold": 90,
    "period": 60,
    "repeatInterval": 600
  }
}
```

**响应** (ZStack → Provider):
```json
{
  "inventory": {
    "uuid": "688a692a452d4dbd95a40fb3e36e5702",
    "name": "acc-test-alarm-updated",
    "description": "Updated acceptance test alarm",
    "comparisonOperator": "GreaterThanOrEqualTo",
    "period": 60,
    "namespace": "ZStack/VM",
    "metricName": "CPUAverageUsedUtilization",
    "threshold": 90.0,
    "repeatInterval": 600,
    "status": "OK",
    "state": "Enabled"
  }
}
```

**结论**：API 返回了完整正确的数据，问题在 SDK 反序列化环节。

### 3.2 追踪 SDK 源码

Provider 的 Update 方法调用链：

```go
// resource_zstack_alarm.go:284
result, err := r.client.UpdateAlarm(state.Uuid.ValueString(), p)
// result.UUID == ""  ← 所有字段都是零值！
```

进入 SDK 实现：

```go
// zstack-sdk-go-v2@v0.0.4/pkg/client/alarm_actions.go
func (cli *ZSClient) UpdateAlarm(uuid string, params param.UpdateAlarmParam) (*view.AlarmInventoryView, error) {
    resp := view.AlarmInventoryView{}
    if err := cli.PutWithRespKey("v1/zwatch/alarms", uuid, "", /*←responseKey为空*/ map[string]interface{}{
        "updateAlarm": params.Params,
    }, &resp); err != nil {
        return nil, err
    }
    return &resp, nil  // resp 全部零值
}
```

关键参数：**第三个参数 `responseKey` 传了空字符串 `""`**。

继续追踪到 HTTP 层：

```go
// zstack-sdk-go-v2@v0.0.4/pkg/client/http_client.go - PutWithAsync()
if len(responseKey) == 0 {
    return location, resp.Unmarshal(retVal)           // ← 走这条路径
}
return location, resp.Unmarshal(retVal, responseKey)  // ← 应该走这条
```

### 3.3 根因定位

API 返回的 JSON 结构：

```json
{
  "inventory": {        ← 外层 envelope key
    "uuid": "...",      ← 实际数据在这一层
    "name": "...",
    ...
  }
}
```

SDK 反序列化行为对比：

| responseKey 值 | 调用方式 | 行为 | 结果 |
|:---:|---|---|:---:|
| `""` (空) | `resp.Unmarshal(retVal)` | 把 `{"inventory":{...}}` 直接映射到 `AlarmInventoryView` | **全部零值** — struct 中无 `inventory` 字段 |
| `"inventory"` | `resp.Unmarshal(retVal, "inventory")` | 先提取 `inventory` key 下的对象，再映射 | **正确** |

**根因**：`UpdateAlarm` 的 `responseKey` 应该传 `"inventory"` 而不是 `""`。

### 3.4 交叉验证

对比同文件中 Create 的实现——Create 使用 `Post` 方法，内部自动处理了 `inventory` envelope：

```go
func (cli *ZSClient) CreateAlarm(params param.CreateAlarmParam) (*view.AlarmInventoryView, error) {
    resp := view.AlarmInventoryView{}
    if err := cli.Post("v1/zwatch/alarms", params, &resp); err != nil {  // Post 内部处理了 envelope
        return nil, err
    }
    return &resp, nil  // resp 正确填充
}
```

这解释了为什么 Create 正常而 Update 异常——两者使用了不同的 HTTP 方法封装，envelope 处理逻辑不同。

---

## 4. 影响范围评估

### 4.1 直接影响

任何对 `zstack_alarm` 资源修改以下属性的 `terraform apply` 都会失败：

- `name`
- `description`
- `comparison_operator`
- `threshold`
- `period`
- `repeat_interval`

用户会看到 "Provider produced inconsistent result after apply" 错误。

### 4.2 潜在影响

其他使用 `PutWithRespKey(path, uuid, "", ...)` 模式的 SDK 方法**可能存在同样问题**。快速扫描：

```bash
grep -n 'PutWithRespKey.*""' ~/go/pkg/mod/github.com/.../pkg/client/*_actions.go
```

需要逐个检查是否存在相同的 responseKey 为空的调用模式。

---

## 5. 修复方案

### 方案 A: 修复 SDK（正确修复）

```go
// alarm_actions.go — 修改 responseKey 为 "inventory"
func (cli *ZSClient) UpdateAlarm(uuid string, params param.UpdateAlarmParam) (*view.AlarmInventoryView, error) {
    resp := view.AlarmInventoryView{}
-   if err := cli.PutWithRespKey("v1/zwatch/alarms", uuid, "", map[string]interface{}{
+   if err := cli.PutWithRespKey("v1/zwatch/alarms", uuid, "inventory", map[string]interface{}{
        "updateAlarm": params.Params,
    }, &resp); err != nil {
        return nil, err
    }
    return &resp, nil
}
```

**优点**：根因修复，一行改动  
**缺点**：需要发 SDK 新版本，Provider 需要更新依赖

### 方案 B: Provider 端 workaround（当前采用）

Update 后不信任 `UpdateAlarm` 的返回值，改为通过 `QueryAlarm` 重新读取：

```go
// resource_zstack_alarm.go — Update 方法
-   result, err := r.client.UpdateAlarm(state.Uuid.ValueString(), p)
+   _, err := r.client.UpdateAlarm(state.Uuid.ValueString(), p)
    if err != nil { ... }

+   // Re-read: UpdateAlarm SDK method uses PutWithRespKey with empty
+   // responseKey, which fails to unmarshal the nested "inventory" envelope.
+   result, err := findResourceByQuery(r.client.QueryAlarm, state.Uuid.ValueString())
+   if err != nil { ... }

    plan.Uuid = types.StringValue(result.UUID)
    // ... (其余不变)
```

**优点**：不依赖 SDK 发版，立即生效  
**缺点**：多一次 API 调用（Query），不是根因修复

### 建议

- **短期**：方案 B（已实施），解除测试阻塞
- **中期**：提交 SDK PR 修复 responseKey，合并后切回方案 A 并删除 workaround

---

## 6. 验证结果

修复后完整测试通过：

```
TestAccAlarmResource
  Step 1: Create  (name="acc-test-alarm")          → PASS
  Step 2: Update  (name="acc-test-alarm-updated")   → PASS  ← 之前 FAIL
  Step 3: Import  (state 一致性校验)                 → PASS
  Destroy         (CheckDestroy 确认删除)            → PASS
```

---

## 7. 附录：测试中发现的另一个问题

### metric_name 写错导致误导性 SYS.1000 错误

| | 值 |
|---|---|
| 测试原始值 | `CPUAverageUtilization` ❌ |
| ZStack 实际 metric | `CPUAverageUsedUtilization` ✅ |

ZStack API 在收到不存在的 metric_name 时，返回 `SYS.1000: An internal error happened in system`（503），而非参数校验错误 `SYS.1007`。这导致最初误以为是 ZWatch 服务不可用。

**建议**：向 ZStack 平台反馈，CreateAlarm API 对无效 metric_name 应返回 SYS.1007 而非 SYS.1000。

---

## 8. 多环境验证记录（2026-04-17）

| 环境 | metricName 错误时 | metricName 修正后 | read-after-write 修复后 |
|------|:---:|:---:|:---:|
| 172.24.248.129 | 503 SYS.1000 | 未测试（环境 zwatch 有独立问题） | — |
| 172.24.189.211 | 503 SYS.1000 | disappears PASS, 主测试 Step 2 inconsistent result | 全部 PASS |

### 排查方法论

1. **两个环境报同样错误 → 排除环境问题，转向审查参数**
2. **用 QueryAlarm 查真实数据对比测试参数 → 定位 metricName 错误**
3. **从测试日志提取 Response body → 确认 API 正确但 SDK 反序列化失败**
4. **追踪 SDK 源码 `PutWithRespKey` → 定位空 responseKey 根因**

---

## 9. 最终修复变更

分支 `fix/alarm-update-and-test`，基于 master 的独立 MR：

| 文件 | 改动 |
|------|------|
| `resource_zstack_alarm.go` | Update 改用 read-after-write |
| `resource_zstack_alarm_test.go` | 修正 metricName + 新增 disappears 测试 |
| `check_disappears_test.go` | 新增 `stateCheckAlarmDisappears` 辅助函数 |
