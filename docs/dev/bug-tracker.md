# Bug Tracker — terraform-provider-zstack

> Generated: 2026-04-20  
> Updated: 2026-05-06（real-env bugcheck against alternate cluster `172.24.189.211`；入账 BUG-086..091；BUG-086 标记 Won't Fix 并先从 provider 注册器移除 `zstack_resource_stack` / `zstack_stack_template`；provider 升级 `zstack-sdk-go-v2` 到 v0.0.8，BUG-087/088 SDK event response unwrap 已修；BUG-087 acceptance 已改为创建临时 SecurityGroup-capable NIC 后测试；BUG-088 provider Read 已改为按 `tagPatternUuid` 查询 user-tags；BUG-089/090/091 已修；所有成功/部分成功创建的 test resources 已 destroy 清理）
> Branch: `test/progress`  
> Tools used: `golangci-lint run`, `go vet`, `go test -short`, manual code review, automated codebase scanning, **real-env Terraform apply→destroy sweep（13 categories × user/mixed/admin plane, RUN_ID `r144465c0a`）**, targeted real-env Terraform bugcheck (2026-05-06)

---

## Fix Status

**Commit**: `5a52557` — `fix: resolve 17 bugs from QA audit (BUG-001 through BUG-021)`

| Verification | Before | After |
|---|---|---|
| `go build ./...` | ✅ clean | ✅ clean |
| `go vet ./...` | ✅ clean | ✅ clean |
| `golangci-lint run ./...` | 13 issues | ✅ 0 issues |
| `go test ./... -short` | 53 pass, 1 fail | ✅ 53 pass, 0 fail |

### Fixed (22 bugs)

| Bug | Priority | Status | Description |
|-----|----------|--------|-------------|
| BUG-001 | P0 | ✅ Fixed | Remove dead code loop in tag_attachment Delete |
| BUG-002 | P1 | ✅ Fixed | Remove empty if-branch in instance resource |
| BUG-003 | P0 | ✅ Fixed | Replace deprecated strings.Title with cases.Title |
| BUG-004 | P0 | ✅ Fixed | Fix nil pointer panic in testAccClientLoggedIn |
| BUG-005 | P0 | ✅ Fixed | Fix TestQueryEnvironment failure (t.Fatal → t.Skip) |
| BUG-006 | P1 | ✅ Fixed | Remove hardcoded default credentials from tests |
| BUG-007 | P1 | ✅ Fixed | Fix 6 unchecked error return values |
| BUG-008 | P1 | ✅ Fixed | Simplify embedded field selectors |
| BUG-009 | P2 | ✅ Fixed | Refactor if-chain to switch in filter.go |
| BUG-010 | P2 | ✅ Fixed | Remove unused copyStringValues function |
| BUG-011 | P2 | ✅ Fixed | Remove unused envInt function |
| BUG-012 | P2 | ✅ Fixed | Fix typo virtural → virtual in filename |
| BUG-013 | P2 | ✅ Fixed | Fix typo BackupStorges → BackupStorages |
| BUG-014 | P2 | ✅ Fixed | Fix typo gpuDeviceTyp → gpuDeviceType |
| BUG-015 | P2 | ✅ Fixed | Fix truncated description in clusters data source |
| BUG-016 | P2 | ✅ Fixed | Standardize factory function naming in provider.go |
| BUG-017 | P3 | ✅ Fixed | Align vmResource/InstanceResource naming |
| BUG-020 | P2 | ✅ Fixed | Convert vague TODO to actionable comment |
| BUG-021 | P1 | ✅ Fixed | Fix antipattern test glob to use absolute path |
| BUG-022 | P2 | ✅ Fixed | Replace string scanning with AST in antipattern tests |
| BUG-023 | P2 | ✅ Fixed | Randomize test resource names |
| BUG-019 | P2 | ✅ Fixed | Move disk state logic to Read() per TODO |

### Remaining / Recent Real-Env Findings

| Bug | Priority | Status | Description |
|-----|----------|--------|-------------|
| BUG-018 | P3 | ⏸ Deferred / Won't Fix for now (2026-05-03) | Standardize acronym casing (UUID/Uuid/IP/Ip)；仅属内部 Go 命名风格问题，不影响 HCL schema 或 provider 行为。当前修复会产生大范围低价值 diff，留到未来生成器/大规模命名重构时再统一处理。 |
| BUG-024 | P3 | ✅ Fixed (2026-05-03) | Add update steps to acceptance tests；已覆盖所有有明确 in-place Update 行为的 acceptance resources，剩余仅显式不支持 Update/no-op/环境或 API 限制项 |
| BUG-025 | P2 | ✅ Fixed (2026-04-28) | Clean up commented-out code blocks in 19+ files |
| BUG-040 | P1 | ✅ Fixed (2026-04-24) | TypeName 改为 `zstack_virtual_routers` (与文件名/SDK 一致) |
| BUG-041–046 | P3 | ✅ Fixed (2026-04-24) | DataSource TypeName 已对齐 SDK 命名 (6 处) |
| BUG-047–052 | P3 | ✅ Fixed (2026-04-24) | Resource TypeName 已对齐 SDK 命名 (6 处) |
| BUG-053 | P0 | ✅ Fixed (2026-04-24) | `iam2_project` Delete 增加 `ExpungeIAM2Project` 调用 |
| BUG-054 | P0 | ✅ Fixed (2026-04-25) | `policy` 加 `statements` 嵌套 Schema (effect/actions/resources/principals/name)，Create 真发送，Read 真映射；RequiresReplace 因 ZStack policy immutable |
| BUG-055 | P1 | ✅ Fixed (2026-04-24) | 加 IsUnknown guard：virtual_router_offering.description/is_default + vpc_shared_qos.description/bandwidth |
| BUG-056 | P0 | ✅ Fixed (2026-04-25) | Provider 侧 re-query workaround: vm_cdrom/l3network/global_config Update 后用 findResourceByQuery；instance Update 已有该模式 |
| BUG-057 | P0 | ✅ Fixed (2026-04-24) | `l2vxlan_network.vni` 加 `int64planmodifier.RequiresReplace` |
| BUG-058 | P0 | ✅ Fixed (已在代码中) | `l3network` Create ipVersion 已有 IsUnknown guard (line 185-187)，核对后确认修过 |
| BUG-059 | P0 | ✅ Closed as Not Reproducible (2026-05-03) | `l3network` Delete URL 缺 UUID 参数的判断按当前 SDK 代码不成立：`DeleteL3Network` 调用 `Delete("v1/l3-networks", uuid, ...)`，底层 `getURL` 会拼 `/v1/l3-networks/{uuid}`。 |
| BUG-060 | P1 | ✅ Fixed (2026-04-24) | `instance` Create Description 改用 `stringPtrOrNil` |
| BUG-061 | P1 | ✅ Fixed (2026-04-25) | Provider 侧 re-query workaround：global_config Create+Update 后用 QueryGlobalConfig 重读 |
| BUG-062 | P2 | ✅ Fixed (2026-04-24) | 子缺陷 a/b 核实已修；c (encoding_type Read 映射) 新修 |
| BUG-063 | P2 | ✅ Fixed (2026-04-24) | `subnet_ip_range` Read 加上 `IpRangeType` 字段映射 |
| BUG-064 | P1 | ✅ Fixed (2026-04-26, commit `1275c54`) | `reserved_ips` ZQL envelope decode 错位（SDK Zql 把 varargs 当 unmarshal-key 钻进 JSON envelope，但 ZQL 响应是 `{"results":[{"inventories":[...]}]}`），导致 Read 失败；同 commit 给 `testdata/generate_tf` 加 `try(length)` guard |
| BUG-065 | P1 | ✅ Fixed (2026-04-26, commit `d3aafc6`) | 24 处 Optional+Computed 空字符串清洗：引入 `stringValueOrNull()` 让空 API 响应 → state null（消除 plan-time `""` vs null drift）；instance Update 后 `findResourceByGet` re-read 重建 state（BUG-019 后续 + SDK-WA-001 模式扩展） |
| **F1/F4 followup** | TBD | 🔲 待 QA | MR #31 最终验证 F1（Plan Compliance）/ F4（Scope Fidelity）reject；具体要求待 QA 下周给报告，到时再编号入账 |
| BUG-066 | P1 | ✅ Fixed (2026-04-28) | `zstack_access_key.user_uuid` 漂移：schema 改 Optional+Computed+UseStateForUnknown，API 返回的 owner UUID 写回 state，消除 "inconsistent result"。原 BUG-NEW-100。 |
| BUG-067 | P3 | ✅ Fixed (2026-05-03) | `zstack_role.identity` 不硬编码 OneOf；schema/docs 明确说明该字段是 ZStack 角色身份/类型标识且合法值依赖版本，并在 CreateRole 因 identity 被服务端拒绝时返回清晰 diagnostic。原 BUG-NEW-101。 |
| BUG-068 | P1 | ✅ Fixed in SDK v0.0.6 (2026-04-27) | `CreateVipQos` 响应 `lastOpDate` 是 `Apr 26, 2026 11:05:50 PM`，SDK 期望 RFC3339 → decode 崩。根因 `pkg/client/http_client.go` PostWithAsync/PutWithAsync 走 stdlib json.Unmarshal 解 inventory 子树。SDK v0.0.6 改成 `resp.Unmarshal(retVal, responseKeyInventory)` 走 jsonutils。原 BUG-NEW-102 / SDK-BUG-005。 |
| BUG-069 | P0 | ✅ Fixed (2026-04-27) | `zstack_instance` Create 始终发 `instanceOfferingUuid: ""`（unset 时也发空字符串），API SYS.1003 拒绝；同模式还有 `rootDiskOfferingUuid` / `defaultL3NetworkUuid` / `strategy`。修复：把 4 处 `stringPtr(...)` 改为 `stringPtrOrNil(...)`（resource_zstack_instance.go 的 Create 入参组装段）。real-env RUN_ID `r144465c0a` 全跑通：vm_offering_path（带 instance_offering_uuid）+ vm_cpu_mem_path（cpu_num+memory_size）Create+Update+Destroy 均成功。原 BUG-NEW-103。 |
| BUG-070 | P2 | ✅ Fixed (2026-05-02) | `zstack_volume_backup` Create 前查询 `backup_storage_uuid` 类型，仅允许 `ImageStoreBackupStorage`；`SftpBackupStorage`/`CephBackupStorage` 在 provider 侧返回 attribute diagnostic。acceptance fixture 不再 fallback 到普通 BackupStorages，docs/example 同步说明。原 BUG-NEW-104。 |
| BUG-071 | P0 | ✅ Fixed (2026-04-27) | `zstack_instance.root_disk` Schema 与 `diskModel` 结构体不一致：`tfsdk:"volume_uuid"` 在 struct 但不在 schema → "Value Conversion Error"。设置 `root_disk` 必崩。修复：root_disk 与 data_disks 两个 SingleNested/ListNested 都加上 `volume_uuid`（Computed+UseStateForUnknown）；root_disk 的 `primary_storage_uuid` 升级为 Optional+Computed（API 自动填充时不再 drift）；Create+Read 把服务端 `AllVolumes` 按 `Type=="Root"` 切给 root_disk，其余切给 data_disks（plan 里没写 data_disks 时不写 state，避免引入空列表 drift）。real-env RUN_ID `r144465c0a` 验证通过。原 BUG-NEW-105。 |
| BUG-072 | P2 | ✅ Fixed (2026-04-28) | `zstack_instance_scripts.script_content`：API 去尾 `\n`，provider 不重读 → drift。修复：在 `state_helpers.go` 新增 `preserveIfEquivAfterTrim()`，Create/Read 后用 helper 做 TrimRight 等价比较，相等则保留用户原值，不等则用服务端值。原 BUG-NEW-106。 |
| BUG-073 | P2 | ✅ Fixed (2026-05-03) | `zstack_log_server` 增加 `category`/`type`/`level` 枚举校验；`configuration` 改为兼容 raw nested JSON，并新增 `appender_type` + `appender_configuration` 结构化输入，由 provider 生成 `{"appenderType": "...", "configuration": {...}}`。flat host/port JSON 在 provider 侧报清晰 diagnostic。同步更新 acceptance fixture、docs、example。原 BUG-NEW-107。 |
| BUG-074 | P1 | ✅ Fixed (2026-04-28) | `zstack_vpc.enable_ipam` Optional 但缺 Computed，API 总是回 `false` → drift。修复：加 `Computed: true` + `boolplanmodifier.UseStateForUnknown()`，保留现有 RequiresReplace。原 BUG-NEW-108。 |
| BUG-075 | P0 | ✅ Fixed (2026-04-30) | `zstack_iam2_organization.type` API 返回空字符串，provider Create/Read/Update 遇到空 Type 时保留配置/既有 state，避免 Required 字段被空响应覆盖导致永久 drift。原 BUG-NEW-109。 |
| BUG-076 | P0 | ✅ Fixed in SDK v0.0.7 (2026-05-03) | `DeleteDirectory` 是动作式 DELETE：`v0.0.7` 保持 `DeleteDirectory(uuid, deleteMode)` 签名，内部改为 `DeleteWithBody("v1/delete/directory", DeleteDirectoryParam{...})`，会发 `DELETE /zstack/v1/delete/directory` 且 body 为 `{"deleteDirectory":{"uuid":"...","deleteMode":"Permissive"}}`。provider 已升级 SDK 至 `v0.0.7`。原 BUG-NEW-110。 |
| BUG-077 | P2 | ✅ Fixed (2026-04-28) | `zstack_certificate.certificate` 服务端去尾空白，provider 不重读 → drift。修复：复用 BUG-072 引入的 `preserveIfEquivAfterTrim()` helper，Create/Read/Update 三处统一处理（虽然 `certificate` 字段是 RequiresReplace，Update 路径也加上以防御）。原 BUG-NEW-111。 |
| BUG-078 | P3 | ✅ Fixed (2026-05-03) | 新增 `data zstack_global_configs` 查询合法 global config `(category, name)` 及当前/default/description；`zstack_global_config` 服务端失败 diagnostic 增加用 data source 发现合法键的提示。原 BUG-NEW-112。 |
| BUG-079 | P3 | ⏸ Superseded / Removed from provider registry (2026-05-06) | `zstack_stack_template.template_content` 的校验曾在 2026-05-03 修复；后续确认该资源属于 ZStack Resource Stack / CloudFormation-style 编排模板接口，与 BUG-086 同类，先从 provider 注册器移除。原 BUG-NEW-113。 |
| BUG-080 | P1 | ✅ Fixed (2026-04-28) | `zstack_pci_device_offering.name` 改为 Required（API 非指针 string + DB NOT NULL）。原 BUG-NEW-114。 |
| BUG-081 | P0 | ✅ Fixed (2026-04-30) | `zstack_price_table` 补齐必填 `prices` 嵌套 schema，并在 Create 中真实发送到 SDK；因 QueryPriceTable 不返回 price 明细，`prices` 标记 RequiresReplace 并在 Read 中保留配置 state。原 BUG-NEW-115。 |
| BUG-082 | P3 | ✅ Fixed (2026-05-03) | `zstack_resource_stack` Create/Update 前校验必须提供 `template_content` 或 `template_uuid`；直接使用 `template_content` 时必须含 `ZStackTemplateFormatVersion`，否则 provider 返回清晰 attribute diagnostic，避免服务端 `invalid decoder: %s!`。docs/example 同步。原 BUG-NEW-116。 |
| BUG-083 | P1 | ✅ Fixed (2026-05-02) | `zstack_preconfiguration_template` 增加 plan-time 校验：`type` 仅允许小写 `[kickstart, preseed, autoyast, autoinstall]`；`content` 必须包含基础系统变量 markers（REPO_URL/USERNAME/PASSWORD/NETWORK_CFGS/FORCE_INSTALL/PRE_SCRIPTS/POST_SCRIPTS）。同步更新 acceptance fixture、example 和 docs。原 BUG-NEW-117。 |
| BUG-084 | P1 | ✅ Fixed in SDK v0.0.6 + Provider workaround removed (2026-04-27) | `UpdateVmInstance` 同 BUG-068 的根因（PostWithAsync/PutWithAsync 用 stdlib `json.Unmarshal` 解 inventory 子树）。SDK v0.0.6 修好后 provider 端 `isSDKTimeParseError` helper + Update 里的 swallow 分支已删除（resource_zstack_instance.go），保留 `findResourceByGet(GetVmInstance)` 兜底（成本忽略，作为 state 一致性保险）。real-env RUN_ID `r144465c0a` v0.0.6 + 清理后 Create(4) + Update v1↔v2 双向各 2 changed + Destroy(4) 全过。原 BUG-NEW-118 / SDK-WA-006。 |
| BUG-085 | P1 | ✅ Fixed (2026-04-27, SDK v0.0.6 后 workaround removed) | `zstack_image` Update 之前被硬拒（`Update not supported`），但 SDK 早就有 `UpdateImage` 支持 Name/Description/GuestOsType/MediaType/Format/System/Platform/Architecture/Virtio。修复：(1) 把 `name` / `description` / `guest_os_type` / `platform` 从 RequiresReplace 改为 in-place updatable；(2) 实装真实 Update（diff plan vs state，仅传变更字段）；(3) Read 改为无条件 refresh + `stringValueOrNull` 防 null↔"" drift；(4) `last_updated` 去掉 UseStateForUnknown 让 UpdateImage bump 之后成为 known-after-apply。SDK v0.0.6 落地后 image.go Update 里 `isSDKTimeParseError` swallow 分支已删除，只留 `findResourceByGet(GetImage)` 兜底。real-env RUN_ID `r144465c0a` 验证：Create (2) + Update v2↔v1 双向各 1 changed + Destroy (2)，cluster cross-check 0 残留。 |
| BUG-086 | P0 | ⏸ Won't Fix / Removed from provider registry (2026-05-06) | `zstack_resource_stack` 是 ZStack Resource Stack / CloudFormation-style 编排入口，与 Terraform 原生资源编排模型重叠；`zstack_stack_template` 同属其模板接口；二者先从 provider 注册器移除，不再投入修复/推广。 |
| BUG-087 | P1 | ✅ Fixed in SDK v0.0.8 (2026-05-06) | `zstack_networking_secgroup_attachment` 在有效私网 NIC 场景下 AddVmNicToSecurityGroup 失败：`Get: key not found`；SDK v0.0.8 改为 `PostWithAsync(..., responseKey="", ...)` 解 event/empty response，provider 已升级依赖。 |
| BUG-088 | P1 | ✅ Fixed in SDK v0.0.8 + Provider Read hardened (2026-05-06) | `zstack_tag_attachment` Create 调 SDK `AttachTagToResources` 失败：`Get: key not found`；SDK v0.0.8 改为 `PostWithAsync(..., responseKey="", ...)` 解 event response，provider 已升级依赖。真实 acceptance 随后暴露 provider Read 不能用 tag pattern UUID 直接 `GetUserTag`，已改为按 `tagPatternUuid` 查询 user-tags 并匹配 `resource_uuids`。 |
| BUG-089 | P1 | ✅ Fixed (2026-05-06) | `zstack_access_key.user_uuid` 改为 Required + RequiresReplace，docs/schema tests 同步；避免不传 `user_uuid` 时 API 拒绝 `field[userUuid] ... is mandatory, can not be null`。 |
| BUG-090 | P1 | ✅ Fixed (2026-05-06) | `data.zstack_license_authorized_nodes` 移除不受 API 支持的 `name` / `name_pattern` 查询字段，仅保留 `uuid` 与本地 `filter`；Read 初始化 `nodes` 为空 list，避免无结果时为 null。 |
| BUG-091 | P1 | ✅ Fixed (2026-05-06) | `zstack_sns_email_endpoint.platform_uuid` 改为 Required + RequiresReplace，docs/schema tests/test generators 同步；避免省略 platform 时 API 返回 `id to load is required for loading`。 |
| **NEW-SCHEMA-NOTE** | — | ✅ 落地 (2026-04-27) | `zstack_instance.network_interfaces` 重设计：`default_l3` 改 Optional+Computed（多 NIC 仅允许一个真值；全省略时 provider 在 Create 自动选第一个 NIC；服务端解析后 echo 回 state，UseStateForUnknown 防 plan drift），`static_ip` 改 Optional+Computed。同步新增顶层 `platform` / `guest_os_type` / `architecture`（前两个 Optional+Computed 走 Update，`architecture` 因 SDK UpdateVmInstanceParamDetail 没有该字段、必须 RequiresReplace）。real-env 用例：`vm_offering_default_l3=[false,true]`、`vm_cpu_default_l3=[true]`、`platform=Linux`、`guest_os_type` 在 Update 里 `CentOS 7.6 ↔ CentOS 7.9` 成功 flip。 |

---

## 2026-05-06 real-env bugcheck addendum

Environment: alternate real-env cluster `172.24.189.211` using AccessKey auth. Credentials are intentionally not recorded here. `data.zstack_zone` smoke test passed (`zone_count = 1`), so Terraform CLI, local dev override provider, and ZStack API connectivity were valid for this run. Temporary resources that reached state (`alarm`, `tag`, `security group`, `resource_stack`) were destroyed successfully.

### BUG-086 — resource stack orchestration removed from provider registry by design

- **Severity**: P0
- **Status**: ⏸ Won't Fix / Removed from provider registry (2026-05-06)
- **Files**: `zstack/provider/resource_zstack_resource_stack.go`, `zstack/provider/resource_zstack_stack_template.go`
- **Confirmed by**: `terraform apply` against minimal `zstack_resource_stack` with inline `template_content`
- **Decision**: remove `ResourceStackResource` and `StackTemplateResource` from `ZStackProvider.Resources()`

**Observed**: Terraform reported `Provider returned invalid result object after apply` for:
- `zstack_resource_stack.stack.parameters`
- `zstack_resource_stack.stack.rollback`
- `zstack_resource_stack.stack.template_uuid`

**Decision rationale**: `zstack_resource_stack` wraps ZStack Resource Stack / CloudFormation-style template orchestration (`cloudformation/stack`), and `zstack_stack_template` manages the companion template API (`cloudformation/template`). That is a second orchestration engine: Terraform can only see the stack/template objects, not the internal resources created by the stack template. This overlaps with Terraform's native resource graph, state ownership, diff, drift detection, and destroy ordering. For this provider, users should model ZStack resources directly with Terraform resources instead of nesting another orchestration system behind opaque orchestration objects.

**Disposition**: do not fix the post-apply unknown-state bug or continue exposing the template resource. Both resources are first removed from provider registration so new configurations cannot use them through this provider. If backward compatibility becomes required later, prefer an explicit deprecation path rather than repairing and promoting this resource family.

### BUG-087 — `zstack_networking_secgroup_attachment` fails with `Get: key not found`

- **Severity**: P1
- **Status**: ✅ Fixed in SDK v0.0.8 (2026-05-06)
- **File**: `zstack/provider/resource_zstack_networking_secgroup_attachment.go`
- **SDK file**: `github.com/zstackio/zstack-sdk-go-v2/pkg/client/other_actions.go`
- **Confirmed by**: creating a new security group and attaching it to an existing private NIC on L3 `l3-private-1`

**Observed**: Create failed during add:
`Could not add VM NIC to security group: Get: key not found`

**Root cause / fix**: provider Create correctly sends `AddVmNicToSecurityGroupParam.Params.VmNicUuids`. SDK v0.0.7 called `cli.Post("v1/security-groups/{uuid}/vm-instances/nics", params, &resp)`, which forced `inventory` response-key unmarshal even though this API returns an event/empty response. SDK v0.0.8 changes the call to `PostWithAsync(..., responseKey="", ...)`, so top-level/empty event responses no longer fail with `Get: key not found`. Provider has been upgraded to `zstack-sdk-go-v2 v0.0.8`. Acceptance coverage now creates a temporary NIC by attaching a SecurityGroup-enabled L3 to an existing running UserVm, then verifies `zstack_networking_secgroup_attachment` against that fresh NIC.

### BUG-088 — SDK `AttachTagToResources` unwraps event response with wrong key

- **Severity**: P1
- **Status**: ✅ Fixed in SDK v0.0.8 (2026-05-06)
- **Provider file**: `zstack/provider/resource_zstack_tag_attachment.go`
- **SDK file**: `github.com/zstackio/zstack-sdk-go-v2/pkg/client/other_actions.go`
- **Confirmed by**: creating `zstack_tag` and attaching it to the real zone UUID

**Observed**: `zstack_tag` Create succeeded, then `zstack_tag_attachment` failed:
`Error attaching tag: Get: key not found`

**Root cause / fix**: provider Create correctly builds `AttachTagToResourcesParam.Params.ResourceUuids` and calls SDK `AttachTagToResources`. SDK v0.0.7 implemented that method with `cli.Post("v1/tags/{tagUuid}/resources", params, &resp)`. `ZSHttpClient.Post()` delegates to `PostWithRespKey(..., responseKeyInventory, ...)`, so it expected an `inventory` key. `AttachTagToResourcesEventView` is an event view with top-level `success` / `results`, so responses without `inventory` failed inside SDK JSON unwrap with `Get: key not found`. SDK v0.0.8 changes this method to `PostWithAsync(..., responseKey="", ...)`; provider has been upgraded to `zstack-sdk-go-v2 v0.0.8`.

**Provider follow-up**: real acceptance after the SDK upgrade showed Create succeeds, but post-apply refresh failed because provider Read called `GetUserTag(tagUuid)` using the tag pattern UUID. The attach event returns a separate user tag UUID, while `tag_uuid` in Terraform is the tag pattern UUID. Provider Read now queries `QueryUserTag` by `tagPatternUuid` and matches the configured `resource_uuids`.

### BUG-089 — `zstack_access_key.user_uuid` is effectively required

- **Severity**: P1
- **Status**: ✅ Fixed (2026-05-06)
- **File**: `zstack/provider/resource_zstack_access_key.go`
- **Confirmed by**: creating `zstack_access_key` with `account_uuid` only

**Observed**: API rejected Create:
`field[userUuid] ... is mandatory, can not be null`

**Fix**: `user_uuid` is now Required and RequiresReplace. The provider no longer advertises server-side defaulting for this field; callers must explicitly choose the user that owns the access key.

### BUG-090 — `zstack_license_authorized_nodes` name filters do not match API inventory

- **Severity**: P1
- **Status**: ✅ Fixed (2026-05-06)
- **File**: `zstack/provider/data_source_zstack_license_authorized_nodes.go`
- **Confirmed by**: applying with `name_pattern = "%"`

**Observed**: API rejected the query:
`LicenseAuthorizedNodeInventory not having field[name]`

**Fix**: remove unsupported `name` / `name_pattern` schema attributes and stop sending name-based API query conditions. Keep exact `uuid` query and local `filter` blocks for fields returned by the inventory. Initialize `state.Nodes` to an empty slice at Read start so empty results become `[]` instead of null.

### BUG-091 — `zstack_sns_email_endpoint.platform_uuid` is not safely optional

- **Severity**: P1
- **Status**: ✅ Fixed (2026-05-06)
- **File**: `zstack/provider/resource_zstack_sns_email_endpoint.go`
- **Confirmed by**: creating `zstack_sns_email_endpoint` with only `name` and `email`

**Observed**: API rejected Create:
`id to load is required for loading`

**Fix**: `platform_uuid` is now Required and RequiresReplace. The provider no longer advertises implicit platform selection; callers must explicitly provide the SNS email platform UUID. Acceptance/generator fixtures now read `sns_email_platforms` from env data and skip endpoint creation when no platform exists.

---

## Summary

| Category | Critical | High | Medium | Low | Total |
|----------|----------|------|--------|-----|-------|
| Logic Bugs | 1 | 1 | 1 | 0 | 3 |
| Error Handling | 0 | 1 | 2 | 0 | 3 |
| Linter / Static Analysis | 0 | 2 | 3 | 1 | 6 |
| Dead Code | 0 | 0 | 2 | 1 | 3 |
| Spelling / Naming | 0 | 2 | 4 | 2 | 8 |
| Test Infrastructure | 1 | 2 | 3 | 1 | 7 |
| TODO / Tech Debt | 0 | 1 | 1 | 0 | 2 |
| Deprecated API | 0 | 1 | 0 | 0 | 1 |
| **Total** | **2** | **10** | **16** | **5** | **33** |

---

## BUG-001 — `append` result discarded in tag_attachment Delete (Logic Bug)

- **Severity**: Critical
- **File**: `zstack/provider/resource_zstack_tag_attachment.go:205-210`
- **Category**: Logic Bug / SA4010
- **Detected by**: golangci-lint (staticcheck SA4010)

```go
// Line 205-210
resourceUuids := make([]string, 0, len(state.ResourceUuids.Elements()))
for _, v := range state.ResourceUuids.Elements() {
    resourceUuids = append(resourceUuids, v.(types.String).ValueString())
}

err := r.client.DetachTagFromResources(state.TagUuid.ValueString(), param.DeleteModePermissive)
```

**Problem**: `resourceUuids` is populated but **never passed** to `DetachTagFromResources`. The function call on line 210 only passes `TagUuid` and `DeleteModePermissive` — the specific resource UUIDs to detach from are ignored. This means **Delete detaches the tag from ALL resources** instead of just the ones in state, or the SDK ignores the missing parameter entirely.

**Fix**: Pass `resourceUuids` to the API call, or verify the SDK's `DetachTagFromResources` signature. The `resourceUuids` variable is computed but never used — this is almost certainly a bug where the intended behavior was to detach from specific resources.

---

## BUG-002 — Empty branch in instance resource (Logic Bug)

- **Severity**: Medium
- **File**: `zstack/provider/resource_zstack_instance.go:1143-1146`
- **Category**: Logic Bug / SA9003
- **Detected by**: golangci-lint (staticcheck SA9003)

```go
if dataDiskCephPoolName != "" {
    // Note: Pools field no longer available in PrimaryStorageInventoryView in SDK v2
    // Pool name validation is skipped; the API will validate on the server side
}
```

**Problem**: Empty `if` branch does nothing. The validation was removed during SDK v2 migration but the conditional structure was left behind.

**Fix**: Remove the empty `if` block entirely, keeping only the comment if documentation is needed.

---

## BUG-003 — `strings.Title` deprecated since Go 1.18 (Deprecated API)

- **Severity**: High
- **File**: `zstack/utils/filter.go:38`
- **Category**: Deprecated API / SA1019
- **Detected by**: golangci-lint (staticcheck SA1019)

```go
fieldName := strings.Title(apiFieldName)
```

**Problem**: `strings.Title` has been deprecated since Go 1.18 because it doesn't handle Unicode punctuation correctly. The Go team recommends `golang.org/x/text/cases` instead.

**Fix**: Replace with:
```go
import "golang.org/x/text/cases"
import "golang.org/x/text/language"

caser := cases.Title(language.Und)
fieldName := caser.String(apiFieldName)
```

---

## BUG-004 — Test helper panics on login failure (Test Infrastructure)

- **Severity**: Critical
- **File**: `zstack/provider/provider_test.go:104-113`
- **Category**: Test Infrastructure
- **Detected by**: Code review

```go
func testAccClientLoggedIn() *client.ZSClient {
    cli := testAccClient()
    if os.Getenv("ZSTACK_ACCESS_KEY_ID") == "" {
        if _, err := cli.Login(context.Background()); err != nil {
            panic(fmt.Sprintf("testAccClientLoggedIn: login failed: %v", err))
        }
    }
    return cli
}
```

**Problem**: `panic()` in test helpers crashes the **entire test binary** — not just the current test. This makes CI debugging extremely difficult: all subsequent tests are killed, stack traces are unhelpful, and the Go test harness cannot produce proper failure reports.

**Fix**: Change signature to accept `*testing.T` and use `t.Fatalf` or `t.Skip`:
```go
func testAccClientLoggedIn(t *testing.T) *client.ZSClient {
    t.Helper()
    cli := testAccClient()
    if os.Getenv("ZSTACK_ACCESS_KEY_ID") == "" {
        if _, err := cli.Login(context.Background()); err != nil {
            t.Skipf("testAccClientLoggedIn: login failed (skipping): %v", err)
        }
    }
    return cli
}
```

---

## BUG-005 — `TestQueryEnvironment` uses `t.Fatal` for missing env vars (Test Infrastructure)

- **Severity**: High
- **File**: `zstack/provider/query_env_test.go:28-30`
- **Category**: Test Infrastructure
- **Detected by**: `go test` (confirmed FAIL)

```go
if akID == "" || akSecret == "" {
    t.Fatal("ZSTACK_ACCESS_KEY_ID and ZSTACK_ACCESS_KEY_SECRET must be set")
}
```

**Problem**: `t.Fatal` causes CI to report FAIL when env vars are not set. This is the **only test that fails** in the current suite (`go test -short` reports 1 FAIL / 53 PASS). Acceptance tests should `t.Skip` when environment is unavailable.

**Fix**: Change `t.Fatal` to `t.Skip`:
```go
if akID == "" || akSecret == "" {
    t.Skip("ZSTACK_ACCESS_KEY_ID and ZSTACK_ACCESS_KEY_SECRET must be set")
}
```

---

## BUG-006 — Hardcoded default credentials in test helper (Test Infrastructure)

- **Severity**: Medium
- **File**: `zstack/provider/provider_test.go:98-101`
- **Category**: Test Infrastructure / Security

```go
return client.NewZSClient(client.NewZSConfig(host, port, "zstack").LoginAccount(
    getEnvOrDefault("ZSTACK_ACCOUNT_NAME", "admin"),
    getEnvOrDefault("ZSTACK_ACCOUNT_PASSWORD", "password"),
).ReadOnly(true).Debug(false))
```

**Problem**: Falls back to `admin/password` when env vars are missing, which may cause unintended authentication attempts against real environments. Tests should fail or skip explicitly when credentials are not provided.

**Fix**: Remove default credentials; require explicit env vars or skip.

---

## BUG-007 — Unchecked error returns (Error Handling)

- **Severity**: High
- **File**: Multiple files
- **Category**: Error Handling / errcheck
- **Detected by**: golangci-lint (errcheck)

| File | Line | Unchecked Call |
|------|------|----------------|
| `resource_antipattern_test.go` | 62 | `defer f.Close()` |
| `resource_antipattern_test.go` | 233 | `defer f.Close()` |
| `resource_zstack_virtual_router_instance_test.go` | 30 | `json.NewEncoder(w).Encode(res)` |
| `resource_zstack_virtual_router_instance_test.go` | 46 | `json.NewEncoder(w).Encode(res)` |
| `resource_zstack_virtual_router_instance_test.go` | 64 | `json.NewEncoder(w).Encode(res)` |
| `resource_zstack_virtual_router_instance_test.go` | 72 | `w.Write([]byte(\`{}\`))` |

**Problem**: Return values of `f.Close()`, `json.Encode()`, and `w.Write()` are silently discarded. While often benign in tests, `f.Close()` errors can indicate data loss and `Encode` errors can cause silent test failures.

**Fix**: For `f.Close()` in production-like code, check error. For test mock handlers, either ignore with `_ =` to make intent explicit, or log the error.

---

## BUG-008 — Redundant embedded field selector (Error Handling)

- **Severity**: Low
- **File**: `zstack/provider/resource_zstack_networking_secgroup.go:167`
- **Category**: Code Quality / QF1008
- **Detected by**: golangci-lint (staticcheck QF1008)

```go
p.BaseParam.SystemTags = []string{"SdnControllerUuid::" + u}
```

**Problem**: `BaseParam` is an embedded field — the selector `p.BaseParam.SystemTags` can be simplified to `p.SystemTags`.

**Fix**: `p.SystemTags = []string{"SdnControllerUuid::" + u}`

---

## BUG-009 — Could use tagged switch (Code Quality)

- **Severity**: Low
- **File**: `zstack/utils/filter.go:65`
- **Category**: Code Quality / QF1003
- **Detected by**: golangci-lint (staticcheck QF1003)

```go
if key == "memory_size" {
    fieldValue = fmt.Sprintf("%d", BytesToMB(field.Int()))
} else if key == "disk_size" {
    ...
```

**Problem**: Chain of `if key == "..."` should be a `switch key {` for readability.

**Fix**: Refactor to `switch key { case "memory_size": ... case "disk_size": ... }`.

---

## BUG-010 — Unused function `copyStringValues` (Dead Code)

- **Severity**: Medium
- **File**: `zstack/provider/state_helpers.go:52`
- **Category**: Dead Code
- **Detected by**: golangci-lint (unused)

```go
func copyStringValues(values []types.String) []types.String {
    if len(values) == 0 {
        return nil
    }
    copied := make([]types.String, len(values))
    copy(copied, values)
    return copied
}
```

**Problem**: Function is defined but never called anywhere in the codebase.

**Fix**: Delete the function or add a usage. If it was intended for a future feature, convert to a TODO with issue reference.

---

## BUG-011 — Unused function `envInt` (Dead Code)

- **Severity**: Medium
- **File**: `zstack/provider/test_env_loader_test.go:175`
- **Category**: Dead Code
- **Detected by**: golangci-lint (unused)

```go
func envInt(m map[string]interface{}, key string) int {
```

**Problem**: Function is defined but never called.

**Fix**: Delete or use.

---

## BUG-012 — Misspelled file name `virtural` → `virtual` (Spelling)

- **Severity**: High
- **File**: `zstack/provider/data_source_zstack_virtural_router_images.go`
- **Category**: Spelling / File Name

**Problem**: File name contains `virtural` (typo for `virtual`). This affects maintainability, grep-ability, and developer confidence.

**Fix**: Rename to `data_source_zstack_virtual_router_images.go`. Update any references in `provider.go` registration or imports. Also rename the corresponding test file if one exists.

---

## BUG-013 — Misspelled struct field `BackupStorges` → `BackupStorages` (Spelling)

- **Severity**: High
- **File**: `zstack/provider/data_source_zstack_backup_storages.go:40,144`
- **Category**: Spelling / Identifier

```go
BackupStorges []backupStorage `tfsdk:"backup_storages"`
// ...
state.BackupStorges = append(state.BackupStorges, backupStorageState)
```

**Problem**: Go field name `BackupStorges` is misspelled (missing "a"). The `tfsdk` tag is correct (`backup_storages`), so Terraform users won't see the typo, but it affects code readability and internal consistency.

**Fix**: Rename `BackupStorges` → `BackupStorages` across both occurrences.

---

## BUG-014 — Misspelled type name `gpuDeviceTyp` → `gpuDeviceType` (Spelling)

- **Severity**: Medium
- **File**: `zstack/provider/resource_zstack_instance.go:30,50-51`
- **Category**: Spelling / Identifier

```go
type gpuDeviceTyp string
// ...
mdevDevice gpuDeviceTyp = "mdevDevice"
pciDevice  gpuDeviceTyp = "pciDevice"
```

**Problem**: Type name `gpuDeviceTyp` is missing the trailing "e" — should be `gpuDeviceType`.

**Fix**: Rename to `gpuDeviceType` (3 occurrences in the same file).

---

## BUG-015 — Truncated schema description `"ype of the cluster"` (Spelling)

- **Severity**: Medium
- **File**: `zstack/provider/data_source_zstack_clusters.go:110`
- **Category**: Spelling / User-facing

```go
Description: "ype of the cluster",
```

**Problem**: Missing leading "T" — should be `"Type of the cluster"`. This is **user-facing** in Terraform registry documentation and CLI output.

**Fix**: Change to `"Type of the cluster"`.

---

## BUG-016 — Inconsistent factory/type naming in provider registration (Naming)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: Medium
- **File**: `zstack/provider/provider.go:221-265`
- **Category**: Naming Consistency

```go
// Examples of inconsistent capitalization:
ZStackvmsDataSource          // lowercase "vms"
ZStackl3NetworkDataSource    // lowercase "l3"
ZStackmnNodeDataSource       // lowercase "mn"
ZStackvrouterDataSource      // "vrouter" not "VirtualRouter"
```

Compare with properly capitalized names:
```go
ZStackVirtualRouterImageDataSource  // correct
ZStackAutoScalingGroupDataSource    // correct
```

**Problem**: Exported factory function names had inconsistent capitalization of acronyms and abbreviations.

**Current fix**: Standardized to consistent PascalCase in both declarations and provider registration: `ZStackVMsDataSource`, `ZStackL3NetworkDataSource`, `ZStackMNNodeDataSource`, `ZStackVRouterDataSource`.

---

## BUG-017 — Type/factory name mismatch: `vmResource` vs `InstanceResource()` (Naming)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: Low
- **File**: `zstack/provider/resource_zstack_instance.go:32,73`
- **Category**: Naming Consistency

```go
type vmResource struct { ... }           // internal type name
func InstanceResource() resource.Resource { return &vmResource{} }  // factory name
```

**Problem**: Factory returned `&vmResource{}` but was called `InstanceResource()`. Other resources are consistent (e.g., `imageResource` → `ImageResource()`). This made it harder to navigate the codebase.

**Current fix**: Renamed `vmResource` → `instanceResource` so the internal type now matches the exported factory naming pattern.

---

## BUG-018 — Mixed acronym casing: UUID vs Uuid, IP vs Ip (Naming)

- **Severity**: Low
- **File**: Multiple files across `zstack/provider/`
- **Category**: Naming Consistency

```go
// Examples found in various files:
L3NetworkUUID types.String  // "UUID" all-caps
PeerL3NetworkUuids types.String  // "Uuids" mixed case
Iprange []ipRangeModel  // "Ip" not "IP"
AttachedClusterUuids  // "Uuids" not "UUIDs"
```

**Problem**: Go convention recommends consistent acronym casing (all-caps `UUID`, `IP`, `VM`, or all-lower in unexported). The codebase mixes both styles.

**Fix**: Decide on convention (recommended: `UUID`, `IP`, `VM` for Go fields) and apply repo-wide. Note: this is a large refactor — lower priority than functional bugs.

---

## BUG-019 — TODO: Delete re-queries VM instance (Tech Debt)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: High
- **File**: `zstack/provider/resource_zstack_instance.go:1053`
- **Category**: TODO / Design Debt

```go
//TODO: query vm instance again in delete function is not smart.
// Update vm instance's data disk state in read function is a better way
```

**Problem**: The Delete function re-queried the VM to get data disk state. This was acknowledged as suboptimal — data disk state should be populated in Read() so Delete() can use the state directly.

**Current fix**: `data_disks` state now persists the data volume UUIDs during Create/Read, and Delete uses those persisted state values instead of issuing a fresh `GetVmInstance()` call.

---

## BUG-020 — TODO: modify mapping tools (Tech Debt)

- **Severity**: Medium
- **File**: `zstack/provider/data_source_zstack_hook_scripts.go:44`
- **Category**: TODO / Tech Debt

```go
// todo modify mapping tools
```

**Problem**: Vague TODO without owner, issue link, or action plan.

**Fix**: Either implement the mapping change, or convert to `// TODO(owner): description — see issue #N`.

---

## BUG-021 — Antipattern test uses relative glob (Test Infrastructure)

- **Severity**: Medium
- **File**: `zstack/provider/resource_antipattern_test.go:27-35`
- **Category**: Test Infrastructure

```go
files, err := filepath.Glob("resource_zstack_*.go")
if err != nil {
    t.Fatalf("failed to glob resource files: %v", err)
}
if len(files) == 0 {
    t.Fatal("no resource files found — test may be running from wrong directory")
}
```

**Problem**: Relative `filepath.Glob` depends on the working directory. Running `go test` from the module root vs package dir produces different results. The test `t.Fatal`s instead of skipping.

**Fix**: Use `runtime.Caller(0)` to compute absolute path (like `test_env_loader_test.go` does):
```go
_, filename, _, _ := runtime.Caller(0)
dir := filepath.Dir(filename)
files, err := filepath.Glob(filepath.Join(dir, "resource_zstack_*.go"))
```

---

## BUG-022 — Antipattern test uses naive string scanning (Test Infrastructure)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: Medium
- **File**: `zstack/provider/resource_antipattern_test.go`
- **Category**: Test Infrastructure

**Problem**: The antipattern scanner searched for code patterns using simple substring matching (`strings.Contains`) on raw source text. This approach:
- Produces **false positives** from comments, string literals, or multiline formatting
- Produces **false negatives** when code is formatted differently than expected

**Current fix**: Replaced string scanning with AST-based scanning using `go/ast` and `go/parser`, plus regression tests covering commented code vs real violations.

---

## BUG-023 — Hardcoded test resource names (Test Infrastructure)

- **Status**: ✅ Fixed (2026-04-21)

- **Severity**: Medium
- **File**: Multiple acceptance test files (e.g., `resource_zstack_volume_test.go:48`)
- **Category**: Test Infrastructure

```go
name = "acc-test-volume"
```

**Problem**: Fixed resource names can cause collisions when tests run concurrently or in shared environments.

**Current fix**: Added shared `testAccName()` helper in `provider_test.go` and switched the affected acceptance tests to randomized `acc-test-<name>-<suffix>` resource names.

Example helper pattern:
```go
name := fmt.Sprintf("acc-test-volume-%s", acctest.RandString(8))
```

---

## BUG-024 — Acceptance tests lack update steps (Test Infrastructure)

- **Status**: ✅ Fixed (2026-05-03)

- **Severity**: Low
- **File**: All acceptance tests
- **Category**: Test Coverage Gap

**Problem**: All existing acceptance tests only exercise Create → (optional Import) → Destroy. None test attribute updates. If Update logic has bugs, they go undetected.

**Resolution**: Added explicit create → update → import coverage to `TestAccZoneResource`, `TestAccClusterResource`, `TestAccCreateImageResource`, `TestAccVolumeResource`, `TestAccBackupStorageResource`, `TestAccLogServerResource`, `TestAccL3NetworkResource`, `TestAccL2VlanNetworkResource`, `TestAccLoadBalancerResource`, `TestAccInstanceResource`, `TestAccLoadBalancerListenerResource`, `TestAccPortForwardingRuleResource`, `TestAccUserResource`, `TestAccRoleResource`, `TestAccIAM2ProjectResource`, `TestAccIAM2VirtualIDResource`, `TestAccIAM2OrganizationResource`, `TestAccCertificateResource`, `TestAccSchedulerJobResource`, `TestAccSchedulerTriggerResource`, `TestAccWebhookResource`, `TestAccAffinityGroupResource`, `TestAccAutoScalingGroupResource`, `TestAccGlobalConfigResource`, `TestAccSNSTopicResource`, `TestAccSNSHttpEndpointResource`, `TestAccSshKeyPairResource`, `TestAccDiskOfferingResource`, `TestAccEIPResource`, and `TestAccVipResource`. Covered update paths include image/VM metadata and guest OS fields, network/load-balancer/listener/PF metadata, IAM2/user/role/certificate/scheduler/webhook metadata, auto scaling sizing/removal policy, global config value, SNS endpoint fields, SSH key pair metadata, disk offering/EIP/VIP metadata, cluster metadata, and zone name/description/state.

**Excluded from update-step coverage**: Resources whose provider/SDK semantics are explicitly no-update or no-op (`access_key`, `guest_tool_attachment`, `instance_offering`, `instance_scripts_execution`, `l2vxlan_network`, `policy`, `reserved_ip`, `sns_email_endpoint`, `subnet_ip_range`, `tag_attachment`, `virtual_router_image`, `virtual_router_offering`, `vm_nic`, `volume_backup`, plus similar create/delete-only resources), and `account` where current environment reports `UpdateAccount` as 404. These are not acceptance update coverage gaps.

---

## BUG-025 — Commented-out code blocks across many files (Dead Code)

- **Severity**: Low
- **File**: 19+ files in `zstack/provider/`
- **Category**: Dead Code

**Files with significant `/* ... */` commented blocks**:
- `data_source_zstack_instances.go`
- `data_source_zstack_tags.go`
- `data_source_zstack_backup_storages.go`
- `data_source_zstack_disk_offers.go`
- `data_source_zstack_zone.go`
- `data_source_zstack_virtural_router_images.go`
- `data_source_zstack_l3networks.go`
- `data_source_zstack_hosts.go`
- `data_source_zstack_images.go`
- `data_source_zstack_l2networks.go`
- `data_source_zstack_vips.go`
- `data_source_zstack_clusters.go`
- `resource_zstack_image.go` (commented `removeStringFromSlice` function)
- `resource_zstack_disk_offering.go`
- `provider.go`

**Problem**: Large commented code blocks add noise, confuse diffs, and suggest incomplete cleanups.

**Fix**: Delete stale blocks. If code is needed later, it's in git history. For planned-but-deferred work, replace with a TODO linking to an issue.

---

## Static Analysis Summary (golangci-lint)

Full output from `golangci-lint run ./...`:

| # | Rule | File | Line | Message |
|---|------|------|------|---------|
| 1 | errcheck | resource_antipattern_test.go | 62 | `f.Close()` unchecked |
| 2 | errcheck | resource_antipattern_test.go | 233 | `f.Close()` unchecked |
| 3 | errcheck | resource_zstack_virtual_router_instance_test.go | 30 | `Encode()` unchecked |
| 4 | errcheck | resource_zstack_virtual_router_instance_test.go | 46 | `Encode()` unchecked |
| 5 | errcheck | resource_zstack_virtual_router_instance_test.go | 64 | `Encode()` unchecked |
| 6 | errcheck | resource_zstack_virtual_router_instance_test.go | 72 | `w.Write()` unchecked |
| 7 | SA9003 | resource_zstack_instance.go | 1143 | Empty branch |
| 8 | QF1008 | resource_zstack_networking_secgroup.go | 167 | Redundant embedded field selector |
| 9 | **SA4010** | **resource_zstack_tag_attachment.go** | **207** | **Append result never used** |
| 10 | SA1019 | utils/filter.go | 38 | `strings.Title` deprecated |
| 11 | QF1003 | utils/filter.go | 65 | Could use tagged switch |
| 12 | unused | state_helpers.go | 52 | `copyStringValues` never used |
| 13 | unused | test_env_loader_test.go | 175 | `envInt` never used |

**`go vet ./...`**: Clean (0 issues)  
**`go build ./...`**: Clean (0 errors)  
**`go test -short`**: 53 PASS, 1 FAIL (`TestQueryEnvironment`), 76 SKIP (acceptance tests)

---

## Test Results Summary

| Status | Count | Notes |
|--------|-------|-------|
| PASS | 53 | All schema/metadata unit tests + 1 mock test + 3 antipattern tests |
| FAIL | 1 | `TestQueryEnvironment` — uses `t.Fatal` for missing env vars (BUG-005) |
| SKIP | 76 | Acceptance tests — skipped without `TF_ACC` or env data |
| **Total** | **130** | |

---

## Priority Fix Order

### Immediate (P0) — Fix before next release

| Bug | Description | Effort |
|-----|-------------|--------|
| BUG-001 | `append` result discarded in tag_attachment Delete — likely functional bug | 30 min |
| BUG-004 | Test helper `panic()` crashes test binary | 15 min |
| BUG-005 | `TestQueryEnvironment` fails CI with `t.Fatal` | 5 min |
| BUG-003 | `strings.Title` deprecated — may break on future Go versions | 30 min |

### Short-term (P1) — Fix within 1 sprint

| Bug | Description | Effort |
|-----|-------------|--------|
| BUG-012 | Rename `virtural` → `virtual` in filename | 15 min |
| BUG-013 | Fix `BackupStorges` → `BackupStorages` | 10 min |
| BUG-014 | Fix `gpuDeviceTyp` → `gpuDeviceType` | 10 min |
| BUG-015 | Fix truncated description `"ype of the cluster"` | 5 min |
| BUG-007 | Unchecked error returns (6 occurrences) | 30 min |
| BUG-010 | Remove unused `copyStringValues` | 5 min |
| BUG-011 | Remove unused `envInt` | 5 min |
| BUG-021 | Fix antipattern test relative glob | 15 min |

### Medium-term (P2) — Technical debt cleanup

| Bug | Description | Effort |
|-----|-------------|--------|
| BUG-002 | Remove empty `if` branch in instance resource | 5 min |
| BUG-006 | Remove hardcoded default credentials in tests | 15 min |
| BUG-008 | Simplify embedded field selector | 5 min |
| BUG-009 | Refactor if-chain to switch | 10 min |
| BUG-016 | Standardize factory naming in provider.go | 1 hr |
| BUG-019 | Implement TODO: move disk state to Read() | 2-3 hrs |
| BUG-020 | Resolve or track TODO in hook_scripts | 15 min |
| BUG-022 | Replace naive string scanning with AST | 2-3 hrs |
| BUG-023 | Randomize test resource names | 1 hr |
| BUG-025 | Clean up commented-out code blocks (19+ files) | 2-3 hrs |

### Low Priority (P3) — Style consistency

| Bug | Description | Effort |
|-----|-------------|--------|
| BUG-017 | Align vmResource/InstanceResource naming | 30 min |
| BUG-018 | Standardize acronym casing (UUID/Uuid/IP/Ip) | 3-4 hrs |
| BUG-024 | Add update steps to acceptance tests | Done |
| BUG-041–052 | File name ↔ TypeName drift (12 处批量) | 2-3 hrs |

---

## BUG-040 — virtual_routers data source 未注册到 provider

- **Severity**: P1
- **File**: `zstack/provider/data_source_zstack_virtual_routers.go:87` (TypeName 定义)
- **Category**: Provider Bug / 未注册
- **Detected by**: Acceptance test TestAccZStackVirtualRoutersDataSource (2026-04-24, 172.24.189.211)

**Problem**: 存在 `data_source_zstack_virtual_routers.go` 实现和 `data_source_zstack_virtual_routers_test.go` 测试文件，但 data source `zstack_virtual_routers` 未在 `provider.go` 的 `DataSources()` 方法中注册。Terraform 报 "The provider hashicorp/zstack does not support data source zstack_virtual_routers"。

**Root cause（批量扫描确认）**: `ZStackVRouterDataSource` factory 实际已在 `provider.go:232` 注册，但其 `Metadata()` 方法内的 TypeName 是 `zstack_virtual_router_instances`（见 `data_source_zstack_virtual_routers.go:87`），而测试/examples 文件命名暗示应为 `zstack_virtual_routers`。属于"文件名/TypeName 漂移"这一类的最严重案例（同类还有 12 处，见 BUG-041–052）。

**Evidence**:
- `go build ./...` ✅ 通过（编译无错误，纯运行时 TypeName mismatch）
- `provider.go:232` 注册了 `ZStackVRouterDataSource`
- 测试 HCL 里用的是 `data "zstack_virtual_routers" "test" {}`
- examples/docs 里用的是 `zstack_virtual_router_instances`（与实际 TypeName 一致）

**Fix**: 两种方案二选一：
1. **推荐**：改 `data_source_zstack_virtual_routers.go:87` 的 TypeName 为 `"_virtual_routers"`，并同步更新 examples/data-sources/virtual_router_instances/ → virtual_routers/ 和 docs。
2. 改测试 HCL 为 `data "zstack_virtual_router_instances"`，保留 examples/docs 不变（变更面小但与文件名不一致的味道保留）。

---

## BUG-041 — DataSource 文件名 ↔ TypeName 漂移：backup_storages

- **Severity**: P3
- **File**: `zstack/provider/data_source_zstack_backup_storages.go`
- **Category**: Naming Consistency
- **Detected by**: 批量扫描 (2026-04-24)

**Problem**: 文件名 `data_source_zstack_backup_storages.go` 暗示 TypeName 应为 `zstack_backup_storages`，但实际是 `zstack_backupstorages`（少了下划线）。examples/data-sources/backupstorages/data-source.tf 遵循实际 TypeName。

**Fix**: 统一改为 `zstack_backup_storages` + 同步 examples/docs（推荐），或接受漂移并在 docs 中说明。

---

## BUG-042 — DataSource 文件名 ↔ TypeName 漂移：instance_guest_tools

- **Severity**: P3
- **File**: `zstack/provider/data_source_zstack_instance_guest_tools.go`
- **Category**: Naming Consistency

**Problem**: 文件名暗示 TypeName 应为 `zstack_instance_guest_tools`，实际是 `zstack_guest_tools`。

**Fix**: 重命名文件为 `data_source_zstack_guest_tools.go`，或将 TypeName 改为 `zstack_instance_guest_tools`（会破坏 examples/docs）。推荐前者。

---

## BUG-043 — DataSource 文件名 ↔ TypeName 漂移：instance_scripts

- **Severity**: P3
- **File**: `zstack/provider/data_source_zstack_instance_scripts.go`
- **Category**: Naming Consistency

**Problem**: 文件名暗示 `zstack_instance_scripts`，实际 `zstack_scripts`。和 BUG-048（resource 同根问题）为同源，建议统一修复。

**Fix**: 讨论后确定 scripts/instance_scripts 哪个更准确（`instance_scripts` 更具描述性），然后同步 TypeName + filename + examples + docs。

---

## BUG-044 — DataSource 文件名 ↔ TypeName 漂移：mn_nodes

- **Severity**: P3
- **File**: `zstack/provider/data_source_zstack_mn_nodes.go`
- **Category**: Naming Consistency

**Problem**: 文件名 `mn_nodes`，TypeName `zstack_mnnodes`（无下划线）。examples/mnnodes/ 遵循实际 TypeName。

**Fix**: 推荐 TypeName 改为 `zstack_mn_nodes`（更规范），同步文件夹/docs。

---

## BUG-045 — DataSource 文件名 ↔ TypeName 漂移：sdn_controllers

- **Severity**: P3
- **File**: `zstack/provider/data_source_zstack_sdn_controllers.go`
- **Category**: Naming Consistency

**Problem**: 文件名暗示 `zstack_sdn_controllers`，实际 `zstack_networking_sdn_controllers`。Resource 同名 TypeName 是 `zstack_sdn_controller`（单数），说明 data source TypeName 莫名加了 `networking_` 前缀。

**Fix**: 删掉 `networking_` 前缀（与 resource 一致），同步 examples/docs。

---

## BUG-046 — DataSource 文件名 ↔ TypeName 漂移：zone

- **Severity**: P3
- **File**: `zstack/provider/data_source_zstack_zone.go`
- **Category**: Naming Consistency

**Problem**: 文件名单数 `zone`，TypeName 复数 `zstack_zones`。

**Fix**: 重命名文件为 `data_source_zstack_zones.go`（推荐，和其他复数 data source 一致）。

---

## BUG-047 — Resource 文件名 ↔ TypeName 漂移：disk_offering

- **Severity**: P3
- **File**: `zstack/provider/resource_zstack_disk_offering.go`
- **Category**: Naming Consistency

**Problem**: 文件名 `disk_offering`，TypeName `zstack_disk_offer`（截断）。examples/resources/disk_offer/ 遵循实际 TypeName。

**Fix**: 推荐 TypeName 改为 `zstack_disk_offering`（完整英文词），同步 examples/docs。对应 data source 已是 `zstack_disk_offers`，也应改为 `zstack_disk_offerings`。

---

## BUG-048 — Resource 文件名 ↔ TypeName 漂移：guest_tool_attachment

- **Severity**: P3
- **File**: `zstack/provider/resource_zstack_guest_tool_attachment.go`
- **Category**: Naming Consistency

**Problem**: 文件名单数 `guest_tool_attachment`，TypeName 复数 `zstack_guest_tools_attachment`。

**Fix**: 重命名文件为 `resource_zstack_guest_tools_attachment.go`（和 TypeName 一致）。

---

## BUG-049 — Resource 文件名 ↔ TypeName 漂移：instance_offering

- **Severity**: P3
- **File**: `zstack/provider/resource_zstack_instance_offering.go`
- **Category**: Naming Consistency

**Problem**: 文件名 `instance_offering`，TypeName `zstack_instance_offer`（截断）。与 BUG-047 同模式。

**Fix**: TypeName 改为 `zstack_instance_offering`，和 BUG-047 一起处理。

---

## BUG-050 — Resource 文件名 ↔ TypeName 漂移：instance_scripts_execution

- **Severity**: P3
- **File**: `zstack/provider/resource_zstack_instance_scripts_execution.go`
- **Category**: Naming Consistency

**Problem**: 文件名 `instance_scripts_execution`，TypeName `zstack_script_execution`（缺 `instance_` 前缀且 `scripts` 变单数）。

**Fix**: 与 BUG-043 / BUG-051 统一讨论 scripts 命名空间，推荐 TypeName 改为 `zstack_instance_scripts_execution`。

---

## BUG-051 — Resource 文件名 ↔ TypeName 漂移：instance_scripts

- **Severity**: P3
- **File**: `zstack/provider/resource_zstack_instance_scripts.go`
- **Category**: Naming Consistency

**Problem**: 文件名 `instance_scripts`（复数），TypeName `zstack_script`（无前缀且单数）。和 BUG-043 同源。

**Fix**: 与 BUG-043 / BUG-050 统一为 scripts 三件套：`zstack_instance_scripts`（data source 复数列表）、`zstack_instance_script`（resource 单数实体）、`zstack_instance_script_execution`（resource 动作）。

---

## BUG-052 — Resource 文件名 ↔ TypeName 漂移：virtual_router_offering

- **Severity**: P3
- **File**: `zstack/provider/resource_zstack_virtual_router_offering.go`
- **Category**: Naming Consistency

**Problem**: 文件名 `virtual_router_offering`，TypeName `zstack_virtual_router_offer`（截断）。与 BUG-047/049 同模式。

**Fix**: TypeName 改为 `zstack_virtual_router_offering`，和 BUG-047/049 一起处理。

---

## 批量修复建议（BUG-041–052）

这 12 个命名漂移最好按三类分组一次性修复，避免 N 次 PR 污染 changelog：

| Group | Scope | Bugs | 推荐一次性 PR |
|---|---|---|---|
| A. 截断词补全 (`_offer`→`_offering`) | resource + ds 双份 | BUG-047, BUG-049, BUG-052 + `zstack_disk_offers` / `zstack_instance_offers` / `zstack_virtual_router_offers` data source | 同 PR |
| B. 下划线/复数统一 | DataSource | BUG-041, BUG-044, BUG-046 | 同 PR |
| C. scripts 命名空间清理 | resource + ds | BUG-043, BUG-050, BUG-051 | 同 PR |
| D. 其他孤立 | — | BUG-042, BUG-045, BUG-048 | 可单独 |

每组修复 = 改 TypeName + 改文件名（可选） + 同步 examples 目录名 + 同步 docs 文件名 + 更新 CHANGELOG（破坏性变更）。

**注意**：这些都是 Terraform Registry 已发布的 TypeName，生产用户升级会 break，需要：
1. 提供 state migration 或 alias（Plugin Framework 不直接支持 alias，需要保留旧 TypeName 并标注 deprecated）；
2. 在 CHANGELOG 明确 breaking change + provider major version bump；
3. 或者反向决策——接受漂移，把文件名改为匹配 TypeName（零破坏但文件名更丑）。

建议先在 sprint 规划会上决策走向，再统一执行。

---

# 生产级功能 Bug (2026-04-24 审核 FIXING 项并归并)

> 以下条目源自 `.worktrees/qa-tracker-repo/_bmad-output/bug-tracker.md` 的 QA 发现，
> 于 2026-04-24 对代码进行实际扫描验证，确认仍未修复。重新编号为 BUG-053+ 以纳入主 tracker。

## BUG-053 — iam2_project Delete 缺 Expunge 调用

- **Severity**: P0 (CRITICAL)
- **Status**: ✅ Fixed (2026-04-24, commit `c45c082`)；后续随 SDK v0.0.5 `DeleteAndExpungeIAM2Project` 进一步收敛为单步调用
- **File**: `zstack/provider/resource_zstack_iam2_project.go:226-243`
- **Evidence (2026-04-24 扫描)**: Delete 方法只调 `DeleteIAM2Project(uuid, DeleteModePermissive)`，无 `ExpungeIAM2Project` 后续调用
- **Category**: Missing Expunge / Soft-delete residual
- **Related Story**: Story-04

**Problem**: ZStack IAM2 Project 在软删除后保留在回收站，导致同名 project 无法重新创建。

**Fix**:
```go
if err := r.client.DeleteIAM2Project(...); err != nil { ... }
if err := r.client.ExpungeIAM2Project(state.Uuid.ValueString()); err != nil {
    // log warning but don't fail — resource is already soft-deleted
}
```

SDK 方法已就位：`~/go/pkg/mod/github.com/zstackio/zstack-sdk-go-v2@v0.0.4/pkg/client/iam2project_actions.go:56`

---

## BUG-054 — policy Create 硬编码空 Statements 数组

- **Severity**: P0 (CRITICAL)
- **Status**: ✅ Fixed (2026-04-25, commit `6588547`)
- **File**: `zstack/provider/resource_zstack_policy.go:114-122`
- **Evidence**:
```go
// Note: Statements are required by the API but we're skipping them for now
// as per the requirement to manage basic CRUD only
createParam := param.CreatePolicyParam{
    Params: param.CreatePolicyParamDetail{
        ...
        Statements:  []param.PolicyStatementParam{},
    },
}
```
- **Category**: Missing Feature / Hardcoded Value
- **Related Story**: Story-05

**Problem**: 当前硬编码空数组意味着创建的 Policy 毫无效力（无 Allow/Deny 规则），整个 policy 资源**无实际用途**。

**Fix**: 在 Schema 增加 `statements` 块（`SchemaList` of nested block with `effect`/`actions`/`resources`），Create/Read/Update 都要映射。注意避免之前评论里"修复引入 god-mode"的陷阱——**不要**给默认 allow-all statements。

---

## BUG-055 — BUG-5 遗漏：3 处 Optional+Computed 仍缺 IsUnknown guard

- **Severity**: P1 (HIGH)
- **Status**: ✅ Fixed (2026-04-24, commit `c45c082`)
- **Evidence (antipattern scanner output 2026-04-24)**:
  - `resource_zstack_virtual_router_offering.go:100`  (`IsDefault.IsNull()` — 字段实际是 Optional 无 Computed，可能是误报)
  - `resource_zstack_virtual_router_offering.go:~97`  (`description` Optional+Computed+`UseStateForUnknown` → 真正 BUG-5 实例)
  - `resource_zstack_vpc_shared_qos.go:140`  (`description` Optional+Computed+`UseStateForUnknown`)
  - `resource_zstack_vpc_shared_qos.go:236`  (`bandwidth` Int64 Optional+Computed → **Int64 零值 = 0 可能导致 API 错误**)
- **Category**: 模式 2：Optional+Computed 缺 IsUnknown 检查（BUG-5 的收尾）

**Problem**: BUG-5 主体已修，但上述 3 处遗漏。Scanner 报告的 `provider.go:133` 是误报（provider 配置块，没有 UseStateForUnknown）。

**Fix** 每处:
```go
// Before:
if !plan.Bandwidth.IsNull() { ... }
// After:
if !plan.Bandwidth.IsNull() && !plan.Bandwidth.IsUnknown() { ... }
```

---

## BUG-056 — 系统性：SDK `PutWithRespKey` 返回空 struct

- **Severity**: P0 (CRITICAL)
- **Status**: ✅ Fixed (2026-04-25, commit `6588547`) — provider 侧走方案 B（Update 后 re-query）；SDK v0.0.5 已修底层（[`SDK-BUG-001`](https://github.com/zstackio/zstack-sdk-go-v2/blob/master/pkg/docs/SDK-BUG-001-PutWithRespKey-Empty-Envelope.md)），provider re-query 现降级为一致性兜底
- **Scope**: `instance` / `vm_cdrom` / `l3network` Update 方法；可能还有 `global_config`
- **Category**: 模式 1（跨资源系统性 SDK bug）

**Problem**: 这几个资源的 Update 方法调用 SDK `Put*` API，SDK 响应解析层缺少 `responseKey` 参数，导致 client 返回空 struct。Provider 把空 struct 写回 state → state 所有字段变成零值。

**Fix 方案**:
- **A (推荐)**: SDK 侧修响应解析 — 一次性覆盖所有受影响 API
- **B (provider 绕过)**: Update 方法在 Put 之后立刻 `Query*` 重新读取，像 Read 方法一样填充 state
- **C (最差但最快)**: 仅在空响应时 fallback 到 Query

跟 SDK 团队确认能否走 A；否则挨个资源走 B。

---

## BUG-057 — l2vxlan_network `vni` 缺 RequiresReplace

- **Severity**: P0 (CRITICAL — 静默数据丢失)
- **Status**: ✅ Fixed (2026-04-24, commit `c45c082`)
- **File**: `zstack/provider/resource_zstack_l2vxlan_network.go` (Schema 中 `vni` 字段)
- **Category**: Plan Modifier Missing

**Problem**: 用户改 vni 时 `terraform apply` 显示 Update，但 ZStack API 不支持改 vni。导致 Terraform state 更新但实际资源未变，**静默数据丢失**。

**Fix**: 给 `vni` 加 `int64planmodifier.RequiresReplace()`。

---

## BUG-058 — l3network Create: `ipVersion=0` 被 API 拒绝

- **Severity**: P0 (CRITICAL)
- **Status**: ✅ Fixed (确认代码中已有 IsUnknown guard, line 185-187, 2026-04-24)
- **File**: `zstack/provider/resource_zstack_l3network.go` Create
- **Category**: 模式 2 的具体实例（Int64 零值）

**Problem**: 用户省略 `ip_version` 时 plan 值为 Unknown，若无 IsUnknown guard 直接 `ValueInt64()` 取到零值 0，ZStack API 拒绝 ipVersion=0。

**Fix**: 同 BUG-055 模式——加 IsUnknown guard；或 Schema 设 Default(4)。

---

## BUG-059 — l3network Delete URL 缺 UUID

- **Severity**: P0 (CRITICAL)
- **Status**: ✅ Closed as Not Reproducible (2026-05-03)
- **File**: `zstack/provider/resource_zstack_l3network.go` Delete
- **Category**: Invalid / stale tracker entry

**Original problem**: `DeleteL3Network` 调用时 URL 模板拼接没把 UUID 放对位置，导致 DELETE 请求打到错误 endpoint。

**Disposition**: Re-checked against the current SDK code in use by this provider. `pkg/client/l3network_actions.go` calls `cli.Delete("v1/l3-networks", uuid, string(deleteMode))`; `pkg/client/http_client.go` routes that through `getDeleteURL()` → `getURL()`, and `getURL()` appends `/{resourceId}` when `resourceId` is non-empty. The resulting endpoint is `/v1/l3-networks/{uuid}?deleteMode=...`, so the tracker entry is stale and no provider or SDK fix is required for this claim.

---

## BUG-060 — instance Create: stringPtr 传空字符串

- **Severity**: P1 (HIGH)
- **Status**: ✅ Fixed (2026-04-24, commit `849b93f`)；后续 BUG-065 把 `stringValueOrNull()` sweep 范围扩到 24 处
- **File**: `zstack/provider/resource_zstack_instance.go` Create (多处)
- **Category**: Nil vs Empty String

**Problem**: 用户省略可选 string 字段时，代码做 `stringPtr(plan.Field.ValueString())` 得到 `*string` 指向空字符串。ZStack API 对某些字段不接受空字符串（需要 null/omit）。

**Fix**: 用 `stringPtrOrNil` 辅助（已存在于 `resource_zstack_policy.go:120`），空字符串返回 nil：
```go
func stringPtrOrNil(s string) *string {
    if s == "" { return nil }
    return &s
}
```

---

## BUG-061 — global_config Read 后 state 全空

- **Severity**: P1 (HIGH)
- **Status**: ✅ Fixed (2026-04-25, commit `6588547`) — provider 走方案 B（QueryGlobalConfig 重读）；底层随 SDK v0.0.5 SDK-BUG-001 一并修复
- **File**: `zstack/provider/resource_zstack_global_config.go`
- **Category**: 模式 1 的亲戚（SDK `PutWithSpec` responseKey 缺失）
- **Related**: BUG-056

**Problem**: 与 BUG-056 同根源，只是走 `PutWithSpec` 分支（不是 `PutWithRespKey`）。表现一致：响应解析得到空 struct。

**Fix**: 跟 BUG-056 一起修，或走 provider 侧 fallback 到 `QueryGlobalConfig`。

---

## BUG-062 — instance_scripts 三处缺陷

- **Severity**: P2 (MEDIUM)
- **Status**: ✅ Fixed (2026-04-24, commit `849b93f`)
- **File**: `zstack/provider/resource_zstack_instance_scripts.go`
- **Category**: 混合（BUG-5 + BUG-63 的 Read 字段遗漏）

**子缺陷**:
1. **a (Create/Update)**: `description` 字段处理不完整
2. **b (Create/Update)**: `script_timeout` (Int64 Optional+Computed) 缺 IsUnknown guard → BUG-5 家族
3. **c (Read)**: `encoding_type` 从 API 返回但未写回 state → 模式 3 实例

**Fix**: 按子缺陷分别处理——a/b 是 guard 问题；c 是 Read 方法字段映射补齐。

---

## BUG-063 — subnet_ip_range Read 后 `ip_range_type` 丢失

- **Severity**: P2 (MEDIUM)
- **Status**: ✅ Fixed (2026-04-24, commit `849b93f`)
- **File**: `zstack/provider/resource_zstack_subnet_ip_range.go` Read
- **Category**: 模式 3 — Read 方法字段遗漏

**Problem**: `ip_range_type` 在 Schema 和 Create 中都有，但 Read 方法没把 API 返回的值写回 state。第二次 `terraform plan` 会显示 drift。

**Fix**: 在 Read 方法里加 `state.IpRangeType = types.StringValue(inv.IpRangeType)` 之类的映射。

---

## BUG-064 — `reserved_ips` ZQL envelope 解码错位 + generator 缺 length guard

- **Severity**: P1 (HIGH)
- **Status**: ✅ Fixed (2026-04-26, commit `1275c54`)
- **File**:
  - `zstack/provider/data_source_zstack_reserved_ips.go` Read
  - `zstack/provider/resource_zstack_reserved_ip.go` Read
  - `zstack/provider/testdata/generate_tf.go`
- **Category**: SDK Usage Bug / Read Failure

**Problem**: `data.zstack_reserved_ips` 与 `resource zstack_reserved_ip` 的 Read 都调用 `Zql(ctx, q, &dst, "inventories")`。SDK `Zql` 把尾部 varargs 当作 unmarshal-key 钻进 JSON envelope，但 ZQL 真实响应是：

```json
{ "results": [ { "inventories": [ ... ] } ] }
```

顶层没有 `"inventories"`，导致每次 Read 报 `"Get: key not found"` (`jsonutils.ErrJsonDictKeyNotFound`)，资源完全无法刷新。

**Fix**:
1. Read 方法改为先反序列化整个 envelope 到一个反映响应结构的 struct，再 flatten `Results[*].Inventories`（与 ZStack 内部其它 ZQL 调用相同模式）。
2. `testdata/generate_tf.go` 给若干 `length()` 调用包 `try(length(...), 0)`，避免空集合导致 fixture 渲染失败。

---

## BUG-065 — Optional+Computed 字符串清洗 + instance Update 后 re-read（24 处 sweep）

- **Severity**: P1 (HIGH)
- **Status**: ✅ Fixed (2026-04-26, commit `d3aafc6`)
- **File**:
  - `zstack/provider/resource_zstack_instance.go`（含两个新辅助 `normalizeNetworkInterfacesFromVM` / `buildUpdatedStateFromVM` + `resource_zstack_instance_bug019_test.go`）
  - `zstack/provider/resource_zstack_l3network.go`
  - `zstack/provider/resource_zstack_sdn_controller.go`
  - `zstack/provider/resource_zstack_aliyun_proxy_vswitch.go`
  - 共 24 处 Optional+Computed 字符串字段（description / dns_domain / vendor_version / status …）
- **Category**: 模式 1 扩展（SDK-WA-001 显式落地） + Empty/Null Drift Cleanup
- **Related**: BUG-019 (Delete 已用 state 而非 re-query), BUG-060 (`stringPtrOrNil`), SDK-WA-001

**Problem**:
1. 大量 Optional+Computed 字符串字段在 ZStack API 返回空字符串时被写回 state 为 `""`，与 HCL 中省略字段的语义（null）不一致，每次 plan 都出现 `"" → null` 噪音 drift。
2. `instance` 资源 Update 路径直接信任 SDK 返回的 (空) struct，state 与 ZStack 实际持久化的内容偶发不符（与 BUG-056/SDK-WA-001 同根，但 SDK v0.0.5 修复前需要 provider 显式兜底）。

**Fix**:
1. 引入 `stringValueOrNull(s string) types.String` 辅助：空字符串 → `types.StringNull()`；24 处调用点切换。
2. `instance.Update` 在成功调用 `UpdateVmInstance` 之后再走 `findResourceByGet(GetVmInstance)` 重读，并通过两个新辅助从 fresh inventory 重建 state（包含 data-disk 同步，把 BUG-019 的 Delete 改造延伸到 Update 路径）。
3. 新增两个单元测试覆盖辅助函数。
4. `testdata/generate_tf.go -only=datasources` 现在为全部 42 个 list-shaped data source 生成 fixture（all / by_name / by_uuid / by_name_pattern / by_nonexistent_uuid 五变体 + by_name/by_uuid 一致性断言），并给每个 data source 按 schema 真实可过滤字段裁剪（如 `mn_nodes` 只支持 `all`、`networking_secgroup_rules` 只有 `priority`、`license_authorized_capacity` 是 singleton）。

**Verified**: `go build ./...` / `go vet ./...` / `go test ./zstack/provider/` 全通过；42/42 生成 fixture 通过 `terraform validate`（dev_overrides）。

---

## 系统性模式汇总（BUG-053–065）

| 模式 | 影响 Bug | 根因 |
|---|---|---|
| **模式 1**：SDK Put 响应解析缺 responseKey | BUG-056 / BUG-061 / BUG-065 | zstack-sdk-go-v2 响应层；上游已在 v0.0.5 修复（[`SDK-BUG-001`](https://github.com/zstackio/zstack-sdk-go-v2/blob/master/pkg/docs/SDK-BUG-001-PutWithRespKey-Empty-Envelope.md)） |
| **模式 2**：Optional+Computed 缺 IsUnknown guard | BUG-055 / BUG-058 / BUG-062b | provider 侧可单独扫修 |
| **模式 3**：Read 方法字段遗漏 | BUG-062c / BUG-063 | 需对 30+ 资源做 Schema↔Read 字段对齐检查 |
| **模式 4**：Empty string 空响应应映射为 null | BUG-060 / BUG-065 | `stringValueOrNull` / `stringPtrOrNil` sweep |
| **模式 5**：ZQL envelope 路径错位 | BUG-064 | SDK `Zql` 当 varargs 钻 envelope，但 ZQL 响应是嵌套 `results.inventories` |

### 优先级建议

1. **BUG-053 / BUG-054**：功能性 critical，单资源独立修，预计各 1-2 小时
2. **BUG-056 系列（056/061）**：等 SDK 修响应层是**最干净的路径**；若等不及走 provider re-read
3. **BUG-055 / BUG-058 / BUG-062b**：antipattern scanner 已标出位置，一个 PR 全修
4. **BUG-057 / BUG-059**：单点功能 bug，独立修
5. **BUG-060**：引入 `stringPtrOrNil` 后一次 sweep 全部 stringPtr 调用
6. **BUG-062c / BUG-063**：和 Read 字段对齐 sweep 一起做

---

# SDK Workaround Registry (provider 侧绕过的 SDK Bug)

> **目的**：所有 provider 代码里"绕过 SDK bug"的地方集中登记。每条都标识：(1) SDK 哪里坏了；(2) provider 怎么绕；(3) 上游修好后这里要回收的代码位置。
>
> **生成日期**: 2026-04-25（SDK v0.0.4 基线）  
> **更新**: 2026-05-06 — provider 升级到 SDK v0.0.8；BUG-087 / BUG-088 的 event response unwrap bug 已由 SDK 修复。
>
> **当前 provider SDK pin**: `github.com/zstackio/zstack-sdk-go-v2 v0.0.8`
> **当前登记 7 类 SDK bug/workaround，含已修复、已回收和仍 Open 项**

| 编号 | 上游 SDK 状态 | Provider 当前状态 |
|---|---|---|
| **SDK-WA-001** | ✅ Fixed (v0.0.5, [`SDK-BUG-001`](https://github.com/zstackio/zstack-sdk-go-v2/blob/master/pkg/docs/SDK-BUG-001-PutWithRespKey-Empty-Envelope.md)) | 23 资源 re-query 仍在；建议升 SDK 后保留 re-query 当一致性兜底，仅清 `// SDK bug:` 注释 |
| **SDK-WA-002** | 🔲 Open ([`SDK-BUG-002`](https://github.com/zstackio/zstack-sdk-go-v2/blob/master/pkg/docs/SDK-BUG-002-ZSClient-Post-URL-Template.md)) | 7 资源已绕过；待 SDK 修后回收 |
| **SDK-WA-003** | ✅ Fixed (v0.0.5, [`SDK-BUG-003`](https://github.com/zstackio/zstack-sdk-go-v2/blob/master/pkg/docs/SDK-BUG-003-IAM2Project-Soft-Delete.md)) — 新增 `DeleteAndExpungeIAM2Project` | provider 仍走两步；升 SDK 后改单步调用（PR-B 一并做） |
| **SDK-WA-004** | — (设计差异，非 bug) | 保留作 adapter 层 |
| **SDK-WA-005** | ✅ Fixed in SDK v0.0.6 (2026-04-27) — Provider workaround N/A (vip_qos was skipped, 可重新启用) | 13-misc-ops 跳过 vip_qos 的注释可去除，加回真实 Create+Destroy 用例 |
| **SDK-WA-006** | ✅ Fixed in SDK v0.0.6 (2026-04-27) — Provider workaround **REMOVED** | `isSDKTimeParseError` helper + `resource_zstack_instance.go` / `resource_zstack_image.go` Update 路径的 swallow-then-`findResourceByGet` 分支已清除；保留 `findResourceByGet(GetVmInstance)` / `findResourceByGet(GetImage)` 作 state 一致性兜底 |
| **SDK-WA-007** | ✅ Fixed in SDK v0.0.8 (2026-05-06) — `AddVmNicToSecurityGroup` / `AttachTagToResources` unwrap event responses without `inventory` | Provider upgraded to v0.0.8; no provider-side create workaround needed. `zstack_tag_attachment` Read also hardened to query user-tags by `tagPatternUuid` and match configured `resource_uuids` |

---

## SDK-WA-001 — `PutWithRespKey` / `PutWithSpec` 返回空 struct

### 上游状态
✅ **Fixed in SDK v0.0.5** — `PutWithAsync` / `PostWithAsync` 在 `responseKey == ""` 时已落地 `inventory` envelope 兜底。详见 [`SDK-BUG-001`](https://github.com/zstackio/zstack-sdk-go-v2/blob/master/pkg/docs/SDK-BUG-001-PutWithRespKey-Empty-Envelope.md)。

### SDK Bug
SDK 中所有 Update 方法（如 `UpdateVmInstance`、`UpdateL3Network`、`UpdateAlarm` 等）调用 `PutWithRespKey(path, uuid, "", params, &resp)`，第三个参数 `responseKey=""`。
HTTP 客户端逻辑（`pkg/client/http_client.go` 第 374 行）：
```go
if len(responseKey) == 0 {
    return location, resp.Unmarshal(retVal)  // ← 不指定 key，整个 resp 反序列化
}
return location, resp.Unmarshal(retVal, responseKey)  // ← 应该走这条
```
ZStack API 的 PUT 响应实际是 `{"inventory": {...}}`，缺少 `"inventory"` key 导致 `Unmarshal(retVal)` 没匹配字段，返回**全零值 struct**。

### 影响范围
**23 个 provider 资源**已经在 Update（部分含 Create）后做 re-query 绕过：

| Provider 文件 | SDK 方法 | 绕过模式 |
|---|---|---|
| `resource_zstack_account.go` | `UpdateAccount` | `_, err := Update(...)` → `GetAccount(uuid)` |
| `resource_zstack_affinity_group.go` | `UpdateAffinityGroup` | `_, err :=` → `GetAffinityGroup` |
| `resource_zstack_alarm.go` | `UpdateAlarm` | `_, err :=` → `findResourceByQuery(QueryAlarm)` |
| `resource_zstack_auto_scaling_group.go` | `UpdateAutoScalingGroup` | `_, err :=` → `GetAutoScalingGroup` |
| `resource_zstack_backup_storage.go` | `UpdateBackupStorage` | `_, err :=` |
| `resource_zstack_cluster.go` | `UpdateCluster` | `_, err :=` |
| `resource_zstack_global_config.go` | `UpdateGlobalConfig` (Create+Update) | `_, err :=` → `findResourceByQuery(QueryGlobalConfig)` ★ 本轮修 |
| `resource_zstack_host.go` | `UpdateHost` / `UpdateKVMHost` | `_, err :=` |
| `resource_zstack_iam2_project.go` | `UpdateIAM2Project` | `_, err :=` |
| `resource_zstack_image_store_backup_storage.go` | `UpdateImageStoreBackupStorage` | `_, err :=` |
| `resource_zstack_instance.go` | `UpdateVmInstance` | `_, err :=` → `findResourceByGet(GetVmInstance)` |
| `resource_zstack_instance_scripts.go` | `UpdateGuestVmScript` | `_, err :=` |
| `resource_zstack_l2vlan_network.go` | `UpdateL2Network` | `_, err :=` |
| `resource_zstack_l3network.go` | `UpdateL3Network` | `_, err :=` → `findResourceByQuery(QueryL3Network)` ★ 本轮修 |
| `resource_zstack_load_balancer.go` | `UpdateLoadBalancer` | `_, err :=` → `GetLoadBalancer` |
| `resource_zstack_load_balancer_listener.go` | `UpdateLoadBalancerListener` | `_, err :=` → `GetLoadBalancerListener` |
| `resource_zstack_port_forwarding_rule.go` | `UpdatePortForwardingRule` | `_, err :=` |
| `resource_zstack_primary_storage.go` | `UpdatePrimaryStorage` | `_, err :=` |
| `resource_zstack_ssh_key_pair.go` | `UpdateSshKeyPair` | `_, err :=` |
| `resource_zstack_vm_cdrom.go` | `UpdateVmCdRom` | `_, err :=` → `findResourceByQuery(QueryVmCdRom)` ★ 本轮修 |
| `resource_zstack_volume.go` | `UpdateVolume` | `_, err :=` |
| `resource_zstack_volume_snapshot.go` | `UpdateVolumeSnapshot` | `_, err :=` |
| `resource_zstack_zone.go` | `UpdateZone` | `_, err :=` |

### 已记录文档
- `troubleshooting/SDK-BUG-UpdateAlarm-Empty-Response.md` — 完整原始报告（alarm 案例）
- 上游：[`SDK-BUG-001`](https://github.com/zstackio/zstack-sdk-go-v2/blob/master/pkg/docs/SDK-BUG-001-PutWithRespKey-Empty-Envelope.md)（v0.0.5 已修）

### SDK v0.0.5 升级后的回收策略（**已选定**）
2026-04-26 决策：**保留 23 处 re-query 当一致性兜底**，仅做最小回收。原因：
- v0.0.5 已让 SDK 自动从 `inventory` envelope 解码，re-query 已不必为正确性所需；
- 但保留 re-query 在 server-side compute（`lastOpDate`、动态字段）下仍能拿到最新值，且仅多一次网络往返，成本可忽略；
- 大规模回收 23 处 re-query = 大 diff + 回归风险，性价比低。

后续 PR-B（SDK v0.0.5 升级）会做：
1. `go get github.com/zstackio/zstack-sdk-go-v2@v0.0.5`；
2. 全仓清除 `// SDK bug:` / `// SDK-WA-001:` 注释（保留 re-query 代码本身）；
3. 不改 23 个资源的逻辑。

---

## SDK-WA-002 — `ZSClient.Post()` 不解析 URL 模板占位符

### 上游状态
✅ **Resolved / stale as of SDK v0.0.7 (2026-05-03)** — 当前 provider 已升级到 `github.com/zstackio/zstack-sdk-go-v2 v0.0.7`。`ZSClient.Post()` 的活跃实现已委托给 `ZSHttpClient.Post(context.Background(), ...)`，并且当前 SDK action 文件中未再发现 `cli.Post("v1/...{placeholder}...")` 形式的未替换模板调用。

### Original SDK Bug
SDK 嵌入了两个 HTTP client：
- `ZSHttpClient.Post(resource, params, retVal)` — 通过 `getPostURL()` 正确解析占位符
- `ZSClient.Post(path, params, result)` — 直接 `fmt.Sprintf("%s/%s", baseURL, path)`，**不替换 `{xxx}` 占位符**

`other_actions.go` 中 **101 个 action 函数（分布于 32 个文件）**调用 `cli.Post(...)` 时走的是被遮蔽的 `ZSClient.Post()`，URL 永远带 literal `{xxx}` → API 返回 404。

### 影响范围（部分确认）
| 文件 | 受影响方法 | Provider 状态 |
|---|---|---|
| `pkg/client/other_actions.go` | 全 101 个 | stale historical note; current SDK v0.0.7 no longer has active `cli.Post("v1/...{placeholder}...")` matches |
| `resource_zstack_primary_storage.go:218,234` | `AddLocalPrimaryStorage` / `AddNfsPrimaryStorage` | provider still directly calls `r.client.Post("v1/primary-storage/...", ...)`; this is no longer tied to the URL-template bug and can be revisited only if SDK typed methods become preferable |

### SDK 修复后的回收清单
SDK 给 `ZSClient.Post()` 加占位符替换后，provider 侧可：
1. 把直接 `r.client.Post("v1/primary-storage/...", ...)` 替换为 `r.client.AddLocalPrimaryStorage(...)` / `AddNfsPrimaryStorage(...)`（前提是 SDK 函数也修了入参）
2. 删 `// SDK bug:` 注释（`resource_zstack_primary_storage.go:218,234`）

---

## SDK-WA-003 — `IAM2Project` 软删除不释放 name

### 上游状态
✅ **Fixed in SDK v0.0.5** — 新增 `DeleteAndExpungeIAM2Project(uuid, deleteMode)` 复合方法（无 break change）。详见 [`SDK-BUG-003`](https://github.com/zstackio/zstack-sdk-go-v2/blob/master/pkg/docs/SDK-BUG-003-IAM2Project-Soft-Delete.md)。

### SDK Bug
`DeleteIAM2Project(uuid, mode)` 仅做软删除，project 进回收站；同名 project 无法重新创建。SDK 提供了独立的 `ExpungeIAM2Project(uuid)` 但没有"删除并清空回收站"的复合方法。

### 影响范围
| 文件 | 修复 |
|---|---|
| `resource_zstack_iam2_project.go:226-249` | Delete 之后调用 `ExpungeIAM2Project`；expunge 失败时降级为 warning（资源已经软删除）★ 本轮修 |

### SDK v0.0.5 升级后的回收清单（PR-B 一并做）
1. `go get github.com/zstackio/zstack-sdk-go-v2@v0.0.5`；
2. `resource_zstack_iam2_project.go` Delete 路径：
   - 删除两步调用（`DeleteIAM2Project` + `ExpungeIAM2Project` + warning 降级）；
   - 替换为单步 `r.client.DeleteAndExpungeIAM2Project(uuid, param.DeleteModePermissive)`；
3. 删除 `// SDK-BUG-003:` 注释。

---

## SDK-WA-004 — `findResourceByQuery` / `findResourceByGet` 兼容层

### Background（不算严格 SDK bug，更像设计差异）
不同资源用不同 SDK 方法读单条记录：
- 一些用 `Get<Resource>(uuid)` → 直接返回 struct 或 `ErrNotFound`
- 一些用 `Query<Resource>(*QueryParam)` → 返回 list，需自己过滤 + 处理空 list

Provider 内部用 `findResourceByGet` / `findResourceByQuery` 统一为"找单条记录或返回 `ErrResourceNotFound`"语义。

### 影响范围
**广泛使用**（`zstack/provider/utils.go` 或 helpers），不是 workaround 而是 adapter 层。
SDK 修不修都建议保留。

---

## SDK-WA-005 — `CreateVipQos` 响应 `lastOpDate` 解码失败

### 上游状态
✅ **Released in SDK v0.0.6 (2026-04-27)** — `pkg/client/http_client.go` 的 `PostWithAsync` (`:312`) 与 `PutWithAsync` (`:380`) 不再用 `json.Unmarshal([]byte(inventory.String()), retVal)`（stdlib 不识别 `ZStackTimeFormat`），改成 `resp.Unmarshal(retVal, responseKeyInventory)` 走 jsonutils 解码层。Provider 已升至 `github.com/zstackio/zstack-sdk-go-v2 v0.0.6`，`replace` 指令已撤除。Provider-side workaround N/A（vip_qos 走的是跳过路径，没有具体 swallow 代码可清；只需在 13-misc-ops 重新启用真实用例）。

### SDK Bug
ZStack 管理面 `CreateVipQos` 响应里 `lastOpDate` 用 `Apr 26, 2026 11:05:50 PM` 这种 Java `SimpleDateFormat` 风格，但 SDK 解码层硬绑 RFC3339 `2006-01-02T15:04:05Z07:00` → 直接解码崩。Create call 实际成功，但 SDK 抛错让 provider 走错误路径。

### 影响范围
| 文件 | 影响 |
|---|---|
| `resource_zstack_vip_qos.go` Create | Create 拿到响应解不开 → 报错 → provider 没机会写 state；vip_qos 资源对此 cluster 完全不可用 |

### Provider 现状
13-misc-ops 这次跑直接跳过了 vip_qos（log 里写了 BUG-068）。

### SDK 修复后回收
SDK-FIX-006 落地后，去掉 13-misc-ops 里跳过 vip_qos 的注释，加回真实 Create+Destroy 用例。

---

## SDK-WA-006 — `UpdateVmInstance` 响应 `lastOpDate` 解码失败

### 上游状态
✅ **Released in SDK v0.0.6 (2026-04-27)** — 同 SDK-WA-005，根因 `pkg/client/http_client.go` 的 `PostWithAsync` / `PutWithAsync` 在空 responseKey 兜底分支里用 stdlib `json.Unmarshal` 解 `inventory` 子树，丢失 `ZStackTimeFormat` 解析能力。改用 `resp.Unmarshal(retVal, responseKeyInventory)` 后 `Apr 27, 2026 3:10:09 PM` 这种格式可正确入 `time.Time`。Provider 已升至 `github.com/zstackio/zstack-sdk-go-v2 v0.0.6`，`replace` 指令已撤除。

### Provider-side workaround：✅ **REMOVED (2026-04-27)**
- 删除 `zstack/provider/resource_zstack_instance.go` 内 `isSDKTimeParseError` helper；
- 删除 `resource_zstack_instance.go` Update 路径的 swallow 分支（`if isSDKTimeParseError(err) { tflog.Warn(...) ... } else { resp.Diagnostics.AddError(...) }` 直接退化为 `if err != nil { return error }`）；
- 删除 `resource_zstack_image.go` Update 路径的同款 swallow 分支；
- 删除 `resource_zstack_instance.go` 不再使用的 `"strings"` 导入；
- 保留 `findResourceByGet(GetVmInstance)` / `findResourceByGet(GetImage)` 这一段 refresh，作为 Update / Read 之间 state-construction 的一致性兜底（成本忽略）。

### 验证（2026-04-27 v0.0.6 + workaround removed real-env 复测）
1. provider `go.mod` 升 `zstack-sdk-go-v2` v0.0.5 → v0.0.6，`go.sum` 经 `go mod tidy` 重生（`v0.0.6 h1:ZbIcdmEc6HvxXH54AdU2dMysdAYKl5omzOw6x9Ekvo8=`），原本指向 `/Users/.../zstack-sdk-go-v2` 的 `replace` 已删；
2. `go build ./...` + `go vet ./...` 全 clean；
3. **10-image-storage**（RUN_ID `r1444xxxx-img-006`）：Phase 1 Create 2 added → Phase 2a Update v1→v2 1 changed → Phase 2b Update v2→v1 1 changed → Phase 3 Destroy 2 destroyed，全程无 decode error；
4. **05-compute-vm**（同 RUN_ID）：Phase 1 Create 4 added → Phase 2a Update v2→v1 2 changed → Phase 2b Update v1→v2 2 changed → Phase 3 Destroy 4 destroyed，全程无 decode error；
5. server cross-check：cluster 内无残留 vm/image。

### SDK Bug
`UpdateVmInstance` 实际在服务端已经成功（curl 验证 VM 名/描述/guestOsType 全已更新），但 SDK 解响应时 `lastOpDate` 也是 `Apr 27, 2026 3:10:09 PM` 格式 → `parsing time "Apr 27, 2026 3:10:09 PM" as "2006-01-02T15:04:05Z07:00": cannot parse "Apr 27, 2026" as "2006"`，让 SDK 把成功的调用当失败抛出。

### Provider 侧 workaround（历史记录，已于 2026-04-27 移除）
~~`zstack/provider/resource_zstack_instance.go`~~：
1. ~~在 Update 里调 `r.client.UpdateVmInstance(...)` 后用 `isSDKTimeParseError(err)` helper 识别这个特定 decode error；~~
2. ~~是的话吞掉并 `tflog.Warn`，再走 `findResourceByGet(r.client.GetVmInstance, uuid)` 拿真实 state；~~
3. ~~不是的话照常走错误路径。~~

~~`isSDKTimeParseError` 直接定义在同文件末尾，只匹配 `"parsing time"` + `"2006-01-02T15:04:05Z07:00"` 两段子串，不会误吞别的解码错。~~

→ 2026-04-27 SDK v0.0.6 release 后已全部清除（详见 "Provider-side workaround：✅ REMOVED" 段落）。

### 影响范围
| 文件 | 影响 |
|---|---|
| `resource_zstack_instance.go` Update | 任何 in-place 修改（重命名 / description / platform / guest_os_type / cpu_num / memory_size 等）都会先撞上 SDK decode error，走 workaround 兜底 |

### 验证
real-env RUN_ID `r144465c0a` Phase 2：v2→v1 + v1→v2 双向 in-place Update 各 2 changed，0 errors；Phase 3 destroy 4 destroyed；server cross-check 无残留。

### SDK 修复后回收（已完成 2026-04-27）
SDK-FIX-006 v0.0.6 release 后已落地：
1. ✅ 删掉 `isSDKTimeParseError` helper 与 Update 里的 swallow 分支；
2. ✅ 让 Update 直接信任 SDK 错误码；
3. ✅ `findResourceByGet(GetVmInstance)` / `findResourceByGet(GetImage)` 保留当一致性兜底。

---

## SDK 修复跟进列表（请提给 SDK 团队）

| 编号 | SDK 修复内容 | 上游状态 | 影响 provider 文件数 |
|---|---|---|---|
| **SDK-FIX-001** | `PutWithRespKey` 在底层默认走 inventory envelope 兜底 | ✅ Fixed in v0.0.5（同时修了 `PostWithAsync`） | 23 |
| **SDK-FIX-002** | `ZSClient.Post()` 接管 URL 模板替换 | ✅ Resolved/stale in v0.0.7；当前 SDK action 中未发现活跃 `cli.Post("v1/...{placeholder}...")` 调用 | 0 confirmed active provider blockers |
| **SDK-FIX-003** | `DeleteIAM2Project` 加 purge 参数或新增 `DeleteAndExpungeIAM2Project` 一站式方法 | ✅ Fixed in v0.0.5（采用方案 B：新增 `DeleteAndExpungeIAM2Project`） | 1 |
| **SDK-FIX-004** | l3network Delete URL 修复（原 BUG-059） | ✅ Closed as not reproducible；当前 SDK `DeleteL3Network` 会发 `/v1/l3-networks/{uuid}?deleteMode=...` | 0 |
| **SDK-FIX-005** | directory Delete body 修复（原 BUG-076 / BUG-NEW-110） | ✅ Fixed in SDK v0.0.7；`DeleteDirectory` 走 `DeleteWithBody("v1/delete/directory", DeleteDirectoryParam{...})` | 0 |
| **SDK-FIX-006** | `pkg/client/http_client.go` `PostWithAsync` (:312) / `PutWithAsync` (:380) 在 `len(responseKey) == 0` 且响应包含 `inventory` 子树时，原本用 `json.Unmarshal([]byte(inventory.String()), retVal)`（stdlib），不识别 ZStack 自家 `ZStackTimeFormat`（`Apr 27, 2026 3:10:09 PM`），导致 `CreateVipQos`（BUG-068）/ `UpdateVmInstance`（BUG-084）/ `UpdateImage`（BUG-085）等响应解码崩。修复：替换为 `resp.Unmarshal(retVal, responseKeyInventory)` 走 jsonutils。同时移除文件里不再需要的 `encoding/json` 导入。| ✅ Released in v0.0.6 (2026-04-27)，provider 端 workaround 已于 2026-04-27 全部移除 | 3+（CreateVipQos / UpdateVmInstance / UpdateImage / 任何空 responseKey 走 inventory 兜底的资源） |
| **SDK-FIX-007** | `pkg/client/directory_actions.go:53` `DeleteDirectory` URL 用 `v1/delete/directory`，与同文件其他方法（`v1/directories`）和 SDK 主流 REST 风格不符；同类 RPC 风格遗留路径还有：`ssoclient_actions.go:16` `Post v1/delete/sso/client`、`ssoclient_actions.go:24` `Get v1/get/sso/client`、`ssoredirect_template_actions.go:{16,24,32}` `v1/{create,update,delete}/sso/...`、`other_actions.go:{4647,5930}` `v1/{add,remove}/resources/directory`、`oauth2token_actions.go:16` `v1/get/oauth2/token`、`saml2client_actions.go:{16,24}` `v1/{update,create}/saml2/client`、`cas_client_actions.go:{16,24}` `v1/{create,update}/cas/client`、`other_actions.go:{1466,4899}` `v1/{create,update}/oauth2/client`、`other_actions.go:{1843,7058}` `v1/get/vmSchedulingRules/...`。需要 SDK 团队对照 ZStack 官方 REST 路径核对。原 BUG-076 / BUG-NEW-110 提到的 directory 是该家族首例。| 🔲 待 SDK 团队按服务端实际 URL 校验 | 至少 11 处 RPC 风格 URL（仅 directory 在 real-env 触发） |

每条 SDK 修复对应的 provider 回收 PR 都很简单（删几行 re-query），但要注意：
- 不要在 SDK 上线**前**回收，否则破坏当前 provider 的正确性。
- 推荐：SDK release 后跟一个 `chore: clean up SDK workarounds (closes SDK-WA-001..003)` 的 PR。

### v0.0.5 联动落地计划

- **PR-A（本 PR）**：tracker hygiene + 入账 BUG-064/065 + 关闭 BUG-059 链至 SDK-BUG-004 + 标记 SDK-WA-001/003 上游已修。**纯文档**，不动 go.mod 也不动 provider 代码。
- **PR-B（后续）**：升 `zstack-sdk-go-v2` 到 v0.0.5；SDK-WA-003 改单步调用；全仓清 `// SDK bug:` / `// SDK-WA-001:` / `// SDK-BUG-003:` 注释；保留 SDK-WA-001 的 23 处 re-query 当一致性兜底（不动逻辑）。
- **PR-C（2026-04-27 落地，本 session）**：BUG-069 / BUG-071 修复 + `zstack_instance.network_interfaces` Schema 重设计（default_l3 / static_ip 改 Optional+Computed，多 NIC 校验 + 自动选首 NIC）+ 顶层 `platform` / `guest_os_type` / `architecture` 增列（前两走 Update，后者 RequiresReplace）+ SDK-WA-006（UpdateVmInstance 时间戳解码兜底）+ BUG-085 修复（`zstack_image` Update 路径打开，lift RequiresReplace 4 字段 + 实装 Update，复用 SDK-WA-006 的 `isSDKTimeParseError` swallow workaround）。real-env RUN_ID `r144465c0a` Create + Update（双向）+ Destroy 全过（05-compute-vm 4+2+2+4，10-image-storage 2+1+1+2）。
