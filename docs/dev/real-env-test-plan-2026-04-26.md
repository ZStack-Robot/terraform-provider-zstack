# Real-Env Resource Test Plan — terraform-provider-zstack

> **目的**：用真实 ZStack 环境（来自 `.env.test`，**不要把里面的 host/port/access key 写进任何文档或代码**）跑 `terraform apply` / `destroy`，覆盖 111 个 resource，验证 provider 的端到端正确性。
>
> **形态**：本计划 + 落在 `zstack/provider/testdata/terraform/real-env-test/<category>/` 下的可直接 `terraform init && terraform apply` 的骨架。**所有 .tf 骨架都在 gitignored 路径下**（`.gitignore:11`），不会被提交，仅供本地跑测。
>
> **隔离原则**（最重要）：所有 create 都用 `acc-test-<run-id>-<resource>` 命名前缀；**不修改、不删除任何已存在的资源**。`terraform destroy` 后用 grep 校验残留作兜底。
>
> **平面（Plane）维度**：本计划 v2 在分类（Category）之外新增「Plane」维度，区分 **User-plane**（用户面，如 vm/image/volume，普通租户日常操作）和 **Admin-plane**（管理面，如 host/cluster/l2/l3/zone/primary_storage 等基础设施，**一旦操作出错具有破坏性**）。Admin-plane 测试要求更严的前置审批和更强的回滚兜底。

---

## 0. 全局配置

### 0.1 环境变量（来自 `.env.test`）

> ⚠️ 实际 host / port / access key / secret **只放在本地 `.env.test`**（已在 `.gitignore`）。本计划不写出任何具体值。运行测试前 `source .env.test` 即可加载下列变量：
>
> - `ZSTACK_HOST`
> - `ZSTACK_PORT`
> - `ZSTACK_ACCESS_KEY_ID`
> - `ZSTACK_ACCESS_KEY_SECRET`

跑法（每类目录里都一样）：

```bash
source .env.test
cd zstack/provider/testdata/terraform/real-env-test/<category>
terraform init -upgrade
terraform plan -var "run_id=$(date +%s)$(openssl rand -hex 2)" -out=tfplan
terraform apply tfplan
# 业务验证 / 截图 / 写报告
terraform destroy -auto-approve -var "run_id=..."
# 兜底：检查 acc-test-<run-id> 前缀残留
bash ../_shared/postcheck.sh "$RUN_ID"
```

### 0.2 共享脚手架

`zstack/provider/testdata/terraform/real-env-test/_shared/`：

- `provider.tf` — provider 配置（读 env）
- `variables.tf` — `run_id` (string, default = timestamp+random)、`zone_uuid`、`cluster_uuid` 等只读引用
- `locals.tf` — `name_prefix = "acc-test-${var.run_id}"`
- `postcheck.sh` — destroy 后扫前缀残留的 stub（待补 SDK 调用）

每个 category 目录里 `provider.tf` / `variables.tf` / `locals.tf` 都是 symlink 到 `_shared/`。

### 0.3 命名约定

| 资源 | 命名 |
|---|---|
| 自创 | `${local.name_prefix}-<resource>-<seq>` 例：`acc-test-2604261830-instance-01` |
| 引用现有 | 通过 `var.<thing>_uuid` 显式注入，不在 .tf 内 hard-code UUID |
| Tag | `terraform_acc_test = "true"` + `run_id = "${var.run_id}"`（便于残留扫描） |

### 0.4 残留扫描（Post-destroy)

```bash
# 用 zstack-cli 或 SDK 扫所有 ${name_prefix} 前缀的残留
make test-real-env-postcheck RUN_ID=$RUN_ID
# 或手动：
for kind in vm-instance volume l3network secgroup ...; do
  zstack-cli QueryXxx name~="${RUN_ID}%"
done
```

如有残留 → **failed clean-up**，记录到本计划末尾的「残留事故簿」并人工清理。

---

## 1. Plane × Category 总览

### 1.1 Plane（平面）

| Plane | 含义 | 风险特征 | 测试纪律 |
|---|---|---|---|
| 👤 **User** | 终端租户日常 CRUD：vm、image、volume、snapshot、ssh_key、tag、alarm 订阅… | 影响范围限于自己创建的对象；做错可重做 | 标准隔离前缀 + destroy 兜底即可 |
| 🛡️ **Admin** | 基础设施侧：zone、cluster、host、l2、l3、primary_storage、backup_storage、global_config、license、role/policy、sdn_controller… | **一旦写错就是破坏性**：删错 cluster 把上面的 VM 全杀；改错 global_config 全局生效；host add/remove 触碰物理机 | 默认 **🚫 Skip**；只在用户独占测试管控节点时白名单放开个别只读/创建独立子树的项 |
| 🔁 **Mixed** | 用户面/管理面在同一资源混合：vpc 既是租户操作也是基础设施；secgroup 创建是用户但绑到现网 VM 是管理 | 需要逐子场景拆分 | apply 前白名单审批，attachment 类只允许挂自创资源 |

> **关键判据**：操作失败时是否会让"现网正在跑业务的 VM/网络/存储"出问题？是 → Admin-plane。

### 1.2 Plane × Category 矩阵

| # | Category | 主 Plane | Resource 数 | 可测 | Skip | 优先级 |
|---|---|---|---|---|---|---|
| 00 | iam | 🛡️ Admin (account/role/policy) + 👤 User (ssh_key_pair, iam2_virtual_id) | 8 | 8 | 0 | P0 |
| 01 | infra-readonly (zone) | 🛡️ Admin | 1 | 0 (data only) | 1 | P0 |
| 02 | compute-base | 🛡️ Admin | 5 | 4 | 1 (host) | P0 |
| 03 | network-base | 🛡️ Admin | 7 | 7 | 0 | P0 |
| 04 | storage-base | 🛡️ Admin | 4 | 1 | 3 (ceph/nvme/iscsi) | P0 |
| 05 | compute-vm | 👤 User | 6 | 6 | 0 | P0 |
| 06 | network-advanced | 🔁 Mixed (eip/pf=User; vrouter/policy_route=Admin) | 12 | 9 | 3 (multicast/flow*) | P1 |
| 07 | network-vpc | 🔁 Mixed | 7 | 4 | 3 (aliyun_*) | P1 |
| 08 | secgroup | 👤 User | 4 | 4 | 0 | P1 |
| 09 | loadbalancer | 👤 User | 4 | 4 | 0 | P1 |
| 10 | image-storage | 👤 User (image/backup/cdp/dataset) | 5 | 5 | 0 | P1 |
| 11 | monitoring | 👤 User (alarm/sns/scheduler) + 🛡️ Admin (log_server/snmp_agent/email_media) | 13 | 13 | 0 | P2 |
| 12 | iam2-extras | 🛡️ Admin (iam2_organization/directory) + 👤 User (tag/tag_attachment) | 4 | 4 | 0 | P2 |
| 13 | misc-ops | 🛡️ Admin (global_config, vcenter, ipsec) + 👤 User (stack_template/dataset/scripts/price_table) | 12 | 10 | 2 (vcenter/ipsec) | P2 |
| 99 | skip-list | — | 18 | 0 | 18 | 🚫 |

**覆盖小计**：可测 = 111 − 18 = **93**。其中

- 👤 User-plane: ~50 resources — **跑通门槛低**，直接做 P0 主线
- 🛡️ Admin-plane: ~30 resources — **风险高**，需要分两阶段：先 only-create-fresh（不动现网）；再 with-cleanup（先确认事先创建的子树还能完整 destroy）
- 🔁 Mixed: ~13 resources — 拆子场景

### 1.3 Admin-plane 额外纪律

1. **仅在自创子树内操作**：例如 cluster 测试就 create 一个全新 cluster，不挂 host、不绑 ps；测完直接 destroy。**绝不在已存在的 cluster 上 update/attach 任何东西**。
2. **Update / Delete 二阶段批准**：每个 admin-plane category 跑前必须本地 `terraform plan` 出来人工读一遍，确认 plan 里没有 `~` (update) / `-` (delete) 现网资源。看到任何一个就停。
3. **存档 plan**：`terraform show -json tfplan > plan-${RUN_ID}.json` 留档，便于事后审计。
4. **destroy 失败必须人工兜底**：admin-plane 残留比 user-plane 严重得多——一个孤儿 cluster 比一个孤儿 VM 影响大得多。
5. **数据库类全部默认 Skip**：`zbox_backup` / `database_backup` / `cdp_task` 在 admin-plane 模式下默认不跑，要跑必须先确认目标无业务数据。

---

## 2. 每类详细计划

### 2.1 Category 00 — IAM（独立运行，无外部依赖）

**Plane**：🛡️ Admin（account/role/policy 影响认证授权链路）+ 👤 User（ssh_key_pair / iam2_virtual_id）
**目录**：`zstack/provider/testdata/terraform/real-env-test/00-iam/`
**资源**：`account`, `user`, `role`, `policy`, `access_key`, `ssh_key_pair`, `iam2_project`, `iam2_virtual_id`

> **Admin 子项纪律**：account / role / policy 创建必须命名前缀严格匹配；测试完后立即 destroy。**不要修改 admin 账号** (`admin` 本人) 或现有任何 system role / 内置 policy。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `account` | — | create + import + destroy；`stateName` 字段非空 |
| basic | `user` | account | create+attach role；rotate password |
| basic | `role` | — | create + attach policy；list permissions |
| basic | `policy` | — | create + statements 不为空 |
| basic | `access_key` | account/user | create → secret 在 state 中只出现一次 |
| basic | `ssh_key_pair` | — | create → public_key 反射回 state |
| basic | `iam2_project` | — | create + Delete + Expunge → name 释放（验证 SDK-WA-003）|
| basic | `iam2_virtual_id` | iam2_project | attach to project + detach |

**风险**：

- iam2_project name 是 ZStack 全局唯一；用 `name_prefix` 即可避免撞名
- access_key 的 `secret` 仅 create 时返回；Read 之后 secret 字段会被 redact

**Skip 子场景**：

- `account` 的 password rotation（影响后续登录链路）

---

### 2.2 Category 01 — Infra Read-only

**Plane**：🛡️ Admin（read-only）
**目录**：`zstack/provider/testdata/terraform/real-env-test/01-infra-readonly/`
**资源**：`zone`

> **注意**：本环境只有一个 zone，且不是测试 zone。**禁止** create/delete。
> 这一类只跑 `data "zstack_zones"` 验证读路径；resource 形式直接放在 Skip List。

---

### 2.3 Category 02 — Compute Base

**Plane**：🛡️ Admin（cluster / host），👤 User-ish（instance_offering / disk_offering / affinity_group 是租户能见的"规格目录"）
**目录**：`02-compute-base/`
**资源**：`cluster`, `host`, `instance_offering`, `disk_offering`, `affinity_group`

> **Admin 警戒**：cluster create 出来后**不要**做 `attach host` 或 `attach primary_storage`；只跑空 cluster 闭环。host **完全不测**（见下表）。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `instance_offering` | — | create CPU=1 / Memory=512MB（不挂任何 VM） |
| basic | `disk_offering` | — | create 1GB |
| basic | `affinity_group` | — | create 空组 |
| basic | `cluster` | zone | create empty cluster；不加 host |
| **🚫** | `host` | cluster | **Skip**：真物理 host add/remove 不可逆 |

**注意**：`host` 一旦 add 进集群且 connect 失败，会留下"挂起 host" → **完全不测**。

---

### 2.4 Category 03 — Network Base

**Plane**：🛡️ Admin（l2/l3/ip_range/reserved_ip）+ 👤 User（vip / certificate）
**目录**：`03-network-base/`
**资源**：`l2vlan_network`, `l2vxlan_network`, `l3network`, `subnet_ip_range`, `reserved_ip`, `vip`, `certificate`

> **Admin 警戒**：l2 vlan id 必须从 `var.test_vlan_id` 中取一个**确认空闲**的；用错号会和现网撞到现网 vlan 把租户网络切断。**绝不允许** attach 现有 cluster 之外的任何已用 cluster；本计划默认只 attach `var.cluster_uuid`。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `l2vlan_network` | zone, cluster, vlan_id (var) | create + attach cluster |
| basic | `l2vxlan_network` | zone, vxlan_pool | create |
| basic | `l3network` | l2vlan_network | create + 立即 destroy（验证 SDK-BUG-004 已修：URL 携带 UUID） |
| basic | `subnet_ip_range` | l3network | create startIp/endIp |
| basic | `reserved_ip` | l3network, ip_range | create + verify ZQL envelope decode（验证 BUG-064 已修） |
| basic | `vip` | l3network | acquire vip → release |
| basic | `certificate` | — | create with PEM cert string |

**风险**：

- VLAN ID 必须在 `var.vlan_pool` 内（默认 4000-4090，需用户确认无冲突）
- `reserved_ip` 的 IP 必须在 ip_range 范围内、且未被分配

---

### 2.5 Category 04 — Storage Base

**Plane**：🛡️ Admin（全部）
**目录**：`04-storage-base/`
**资源**：`primary_storage` (LocalStorage 模式), `image_store_backup_storage`, `ceph_pool`, `iscsi_server`

> **Admin 警戒**：primary_storage / backup_storage 一旦 attach 错 cluster 或 url 错都会触发现网存储重连/重挂；本类只允许 create + 立即 destroy 闭环，**不允许** attach 现有 cluster。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `primary_storage` | zone, cluster | create LocalStorage（不挂载真实盘）|
| basic | `image_store_backup_storage` | — | create with mock URL |
| basic | `ceph_pool` | ceph_primary_storage (existing) | 仅 Query 现有 pool；不创建 |
| **🚫** | `ceph_primary_storage` / `ceph_backup_storage` | — | **Skip**：需要真实 mon URL，配错会破坏现网 |
| **🚫** | `nvme_server` / `iscsi_server` | — | **Skip**：需要真实硬件端点 |

---

### 2.6 Category 05 — Compute VM（耗时最长）

**Plane**：👤 User（这是最典型的用户面 — 创建虚机/盘/快照都是租户日常）
**目录**：`05-compute-vm/`
**资源**：`instance`, `vm_nic`, `vm_cdrom`, `volume`, `volume_snapshot`, `volume_backup`

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `instance` | image, l3, instance_offering | create + Update name (验证 BUG-065 stringValueOrNull 修复)+ destroy |
| compose | `vm_nic` | instance, l3 | attach extra NIC + detach |
| compose | `vm_cdrom` | instance | attach iso + detach |
| basic | `volume` | disk_offering, primary_storage | create → attach to test instance → detach |
| basic | `volume_snapshot` | volume | create snapshot + revert |
| basic | `volume_backup` | volume, backup_storage | create backup |

**风险**：

- VM 启动需要 image 已经 Ready；用 `data "zstack_images"` 选一个 mini cloud-image
- 测试镜像 UUID 通过 `var.test_image_uuid` 注入，避免动到生产镜像

---

### 2.7 Category 06 — Network Advanced

**Plane**：🔁 Mixed
- 👤 User: `eip`, `port_forwarding_rule`
- 🛡️ Admin: `virtual_router_offering`, `virtual_router_image`, `virtual_router_instance`, `vrouter_route_table`, `vrouter_route_entry`, `policy_route_rule_set`, `policy_route_rule`
- 🚫 Skip: `multicast_router`, `flow_collector`, `flow_meter`

**目录**：`06-network-advanced/`
**资源**：`eip`, `port_forwarding_rule`, `virtual_router_offering`, `virtual_router_image`, `virtual_router_instance`, `vrouter_route_table`, `vrouter_route_entry`, `policy_route_rule_set`, `policy_route_rule`, `multicast_router`, `flow_collector`, `flow_meter`

> **Admin 警戒**：vrouter 实例上线后会自动接管它所在 l3 的 DHCP/DNS/SNAT；本类的 vrouter 必须挂在**自创的 l3** (`zstack_l3network.test`) 上，**不允许**绑定到现有 l3。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `eip` | l3network (public), vip | acquire + bind + unbind |
| basic | `port_forwarding_rule` | vip, instance | create rule + ping through |
| basic | `virtual_router_offering` | l3network | create offering |
| basic | `virtual_router_image` | image (existing) | register image as router image |
| basic | `virtual_router_instance` | offering, image, l3 | spin up vrouter |
| compose | `vrouter_route_table` + `vrouter_route_entry` | vrouter_instance | create table + add entry |
| basic | `policy_route_rule_set` + `policy_route_rule` | vrouter | basic policy chain |
| **🚫** | `multicast_router` / `flow_collector` / `flow_meter` | — | **Skip**：依赖外部多播 / NetFlow 收集器 |

---

### 2.8 Category 07 — Network VPC

**Plane**：🔁 Mixed（vpc 主体是租户层，但 vpc_firewall / vpc_shared_qos 是管理面策略）
**目录**：`07-network-vpc/`
**资源**：`vpc`, `vpc_firewall`, `vpc_ha_group`, `vpc_shared_qos`, `aliyun_proxy_vpc`, `aliyun_proxy_vswitch`, `aliyun_nas_access_group`

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `vpc` | zone | create empty vpc |
| basic | `vpc_firewall` | vpc | create rule set |
| basic | `vpc_ha_group` | vpc, two l3 | create HA group |
| basic | `vpc_shared_qos` | vpc | create policy |
| **🚫** | `aliyun_proxy_*` / `aliyun_nas_*` | — | **Skip**：依赖阿里云账号 + RAM |

---

### 2.9 Category 08 — Security Group

**Plane**：👤 User（standard tenant security policy）
**目录**：`08-secgroup/`
**资源**：`networking_secgroup`, `networking_secgroup_rule`, `networking_secgroup_attachment`, `access_control_list`

> **Mixed 子项纪律**：`networking_secgroup_attachment` 必须挂到**自创的 acc-test-${RUN_ID}-instance**，**绝对不允许**挂到现网 332 个 VM 中的任何一个。一旦发现 plan 里出现现网 nic_uuid，立即 abort。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `networking_secgroup` | — | create empty group |
| basic | `networking_secgroup_rule` | secgroup | add ingress rule（22/tcp） |
| compose | `networking_secgroup_attachment` | secgroup, **自创 instance** | attach + detach（用 05 类创建的 acc-test instance，不动现网 VM） |
| basic | `access_control_list` | — | create ACL |

---

### 2.10 Category 09 — Load Balancer

**Plane**：👤 User
**目录**：`09-loadbalancer/`
**资源**：`load_balancer`, `load_balancer_listener`, `lb_server_group`, `vip_qos`

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `load_balancer` | vip | create LB |
| basic | `load_balancer_listener` | load_balancer | create TCP/80 listener (验证 SDK-BUG-002 URL template 已修) |
| basic | `lb_server_group` | load_balancer | create + add backend |
| basic | `vip_qos` | vip | bandwidth limit |

---

### 2.11 Category 10 — Image / Storage（带数据）

**Plane**：👤 User（image / dataset 是租户上传/管理）；🛡️ Admin 警戒（database_backup / cdp_task 涉及 mn 数据）
**目录**：`10-image-storage/`
**资源**：`image`, `database_backup`, `cdp_policy`, `cdp_task`, `dataset`

> **Admin 警戒**：`database_backup` 备份的是 ZStack 管理节点自己的 DB；不要在生产 mn 上跑 restore；本类只测 backup 路径，restore 完全 Skip。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `image` | backup_storage | create from mini-image URL（公开 cirros） |
| basic | `database_backup` | — | create immediate backup |
| basic | `cdp_policy` | — | create policy |
| basic | `cdp_task` | cdp_policy, volume | create task |
| basic | `dataset` | — | create empty dataset |

**风险**：

- `image` 下载耗时（5-15 min）；用 cirros (12MB) 即可
- `cdp_task` 会触发真实 IO；只跑 dry-run 模式

---

### 2.12 Category 11 — Monitoring / Alarm / SNS

**Plane**：🔁 Mixed
- 👤 User: `alarm`, `sns_topic`, `sns_email_endpoint`, `sns_http_endpoint`, `webhook`, `monitor_group`, `monitor_template`, `scheduler_job`, `scheduler_trigger`, `instance_scripts`
- 🛡️ Admin: `email_media` (改 SMTP 全局配置), `log_server` (rsyslog/fluentd 全局), `snmp_agent` (SNMP 全局)

**目录**：`11-monitoring/`
**资源**：`alarm`, `monitor_group`, `monitor_template`, `sns_topic`, `sns_email_endpoint`, `sns_http_endpoint`, `email_media`, `webhook`, `log_server`, `snmp_agent`, `scheduler_job`, `scheduler_trigger`, `instance_scripts`

> **Admin 警戒**：email_media 配错会让现网告警邮件投不出去；log_server 配错会切断日志收集；snmp_agent 配错暴露管理端口。本类的 admin 子项默认在 `terraform plan` 阶段先 dry-run，确认 read 路径正常后再 apply create；apply 后立即 destroy。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `sns_topic` | — | create |
| basic | `sns_email_endpoint` | sns_topic | subscribe email |
| basic | `sns_http_endpoint` | sns_topic | subscribe https://example.com |
| basic | `email_media` | — | create SMTP config (mock SMTP 不会真发) |
| basic | `alarm` | sns_topic | create + Update name (验证 SDK-BUG-001 已修) |
| basic | `monitor_template` / `monitor_group` | — | create empty |
| basic | `webhook` | sns_topic | create |
| basic | `log_server` | — | create syslog target |
| basic | `snmp_agent` | — | create v2c agent |
| basic | `scheduler_job` + `scheduler_trigger` | — | create cron job |
| basic | `instance_scripts` | — | create script content |

---

### 2.13 Category 12 — IAM2 Extras

**Plane**：🛡️ Admin（iam2_organization / directory）+ 👤 User（tag / tag_attachment）
**目录**：`12-iam2-extras/`
**资源**：`iam2_organization`, `directory`, `tag`, `tag_attachment`

> **Mixed 纪律**：`tag_attachment` 同 secgroup_attachment 一样，**只能挂到自创 acc-test instance**。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `iam2_organization` | iam2_project | create |
| basic | `directory` | iam2_organization | create dept |
| basic | `tag` | — | create user tag |
| compose | `tag_attachment` | tag, **自创 instance** | attach + detach |

---

### 2.14 Category 13 — Misc Ops

**Plane**：🛡️ Admin（global_config / vcenter / ipsec / port_mirror / pci_device_offering）+ 👤 User（stack_template / resource_stack / dataset / scripts / price_table / preconfiguration_template）
**目录**：`13-misc-ops/`
**资源**：`global_config`, `stack_template`, `resource_stack`, `preconfiguration_template`, `pci_device_offering`, `instance_scripts_execution`, `guest_tool_attachment`, `port_mirror`, `port_mirror_session`, `price_table`, `ipsec_connection`, `vcenter`

> **Admin 警戒（global_config 单独说一下）**：这是**最危险**的 admin-plane 资源，因为它没有 create/delete 概念，只有 update。改错一个 key 全集群生效。本类只允许把 `value` 改成"原值"作 idempotent 测试，**禁止**改成新值。挑一个无副作用的 key（建议 `agent.deployTimeout` 之类）做"原值→原值" Update 验证 SDK-BUG-001 修复，跑完即停。

| 子场景 | 资源 | 前置 | 验证点 |
|---|---|---|---|
| basic | `global_config` | — | Update only（不能 create / delete）；测 BUG-061 修复 |
| basic | `stack_template` | — | create CFN-style template |
| basic | `resource_stack` | stack_template | apply stack |
| basic | `preconfiguration_template` | — | create cloud-init template |
| basic | `pci_device_offering` | — | create offering（不挂载真实 PCI） |
| basic | `instance_scripts_execution` | instance_scripts, instance | execute on test VM |
| basic | `guest_tool_attachment` | instance | attach guest tool |
| basic | `port_mirror` + `port_mirror_session` | l2 | create mirror |
| basic | `price_table` | — | create table |
| **🚫** | `ipsec_connection` | — | **Skip**：需对端实体设备 |
| **🚫** | `vcenter` | — | **Skip**：需真 vCenter 凭据 + 库存影响巨大 |

---

## 3. 🚫 Skip List（18 个，明确不在真实环境跑）

> 触发的副作用 ≥ 收益的资源；保留单元测试 + schema 测试即可。

| 资源 | 跳过原因 |
|---|---|
| `host` | 加错或删错会让现网 host 失联 |
| `ceph_primary_storage` | 需要真实 mon URLs；配错破坏现网存储 |
| `ceph_backup_storage` | 同上 |
| `nvme_server` | 需要真实 NVMe-oF target |
| `iscsi_server` | 需要真实 iSCSI portal |
| `vcenter` | 需 vCenter 凭据；触发库存同步 |
| `baremetal_chassis` | 需真实 IPMI；写错触发 BMC 异常 |
| `baremetal_instance` | 同上 |
| `baremetal_pxe_server` | 需 DHCP/TFTP 控制权 |
| `aliyun_proxy_vpc` | 需阿里云 RAM |
| `aliyun_proxy_vswitch` | 同上 |
| `aliyun_nas_access_group` | 同上 |
| `v2v_conversion_host` | 需 V2V 服务后端 |
| `sdn_controller` | 影响整个 SDN 控制面 |
| `multicast_router` | 需多播订阅 |
| `flow_collector` | 需 NetFlow collector 端 |
| `flow_meter` | 同上 |
| `ipsec_connection` | 需对端 IPsec 设备 |
| `license` / `*_security_machine` 系 | 一次性激活 / 安全机制不可逆 |
| `zbox_backup` | 需 ZBox 真实存储 |
| `container_management_endpoint` | 需 K8s 后端 |
| `info_sec_security_machine` 等 5 个安全机 | 真实安全设备 |

完整 list 以本节为准；旧的 `_bmad-output/test-status-overview.md` 已作为过期过程文档移除。

---

## 4. 执行流程

### 4.1 单类跑通

```bash
source .env.test
cd examples/real-env-test/00-iam
RUN_ID=$(date +%s)$(openssl rand -hex 2)
terraform init -upgrade
terraform plan -var run_id=$RUN_ID -out=tfplan
terraform apply tfplan
# 业务验证（手动 / 截图）
terraform destroy -auto-approve -var run_id=$RUN_ID
# 残留扫描
bash ../_shared/postcheck.sh $RUN_ID
```

### 4.2 全量跑（执行顺序按 Category # 由小到大）

依赖链：`00 → 02 → 03 → 04 → 05 → 06,07,08,09 → 10 → 11 → 12 → 13`

`01-infra-readonly` 在任意时机都可单独跑（只读）。

### 4.3 失败定位

每个 category 目录里跑出来的：

- `terraform.tfstate` (run-local，不入库)
- `apply.log` / `destroy.log`（重定向保存）
- `postcheck-residue.json`（残留扫描结果）

失败时 → 把这三个文件打包贴到对应 category 的 issue。

---

## 5. 进度跟踪表

| Category | 资源数 (含 skip) | 已写 .tf | 已 apply 通过 | 已 destroy 通过 | postcheck 干净 |
|---|---|---|---|---|---|
| 00-iam | 8 | 🔲 | 🔲 | 🔲 | 🔲 |
| 01-infra-readonly | 1 | 🔲 | 🔲 | 🔲 | 🔲 |
| 02-compute-base | 5 (1 skip) | 🔲 | 🔲 | 🔲 | 🔲 |
| 03-network-base | 7 | 🔲 | 🔲 | 🔲 | 🔲 |
| 04-storage-base | 4 (3 skip) | 🔲 | 🔲 | 🔲 | 🔲 |
| 05-compute-vm | 6 | 🔲 | 🔲 | 🔲 | 🔲 |
| 06-network-advanced | 12 (3 skip) | 🔲 | 🔲 | 🔲 | 🔲 |
| 07-network-vpc | 7 (3 skip) | 🔲 | 🔲 | 🔲 | 🔲 |
| 08-secgroup | 4 | 🔲 | 🔲 | 🔲 | 🔲 |
| 09-loadbalancer | 4 | 🔲 | 🔲 | 🔲 | 🔲 |
| 10-image-storage | 5 | 🔲 | 🔲 | 🔲 | 🔲 |
| 11-monitoring | 13 | 🔲 | 🔲 | 🔲 | 🔲 |
| 12-iam2-extras | 4 | 🔲 | 🔲 | 🔲 | 🔲 |
| 13-misc-ops | 12 (2 skip) | 🔲 | 🔲 | 🔲 | 🔲 |
| **合计可测** | **93** | — | — | — | — |
| Skip 列入 | 18 | — | — | — | — |

---

## 6. 残留事故簿（占位，跑出来才会有内容）

| 日期 | RUN_ID | Category | 残留资源 | 原因 | 处理 |
|---|---|---|---|---|---|
| _ | _ | _ | _ | _ | _ |

---

## 7. 与 SDK / 已知 BUG 的对照

跑通时也是这些 BUG 修复的最终验证：

| BUG | 验证点 | Category |
|---|---|---|
| SDK-BUG-001 (PutWithRespKey envelope) | Update 后 state 不归零 | 11-monitoring (alarm) / 03 (l3network) / 13 (global_config) / 05 (vm_cdrom) |
| SDK-BUG-002 (URL template) | listener Create 路径正确 | 09-loadbalancer |
| SDK-BUG-003 (DeleteIAM2Project soft) | Delete 后能立即重建同名 | 00-iam |
| SDK-BUG-004 (DeleteL3Network URL) | Delete 不再 404 | 03-network-base |
| BUG-064 (reserved_ips ZQL envelope) | List/Read 不报 "key not found" | 03-network-base |
| BUG-065 (stringValueOrNull sweep) | Update 后 description 字段 null/"" 一致 | 05-compute-vm |

---

## 8. 待用户决策的开放点

1. **VLAN ID 池**：用 `var.vlan_pool` 默认 `[4000, 4090]`，确认这些 VLAN 在物理网络上没人占？
2. **测试 image**：建议用 cirros-0.5.2 (12MB)；如果环境里已有合适小 image，给出 UUID 直接 reuse。
3. **acc-test 命名前缀长度**：当前 `acc-test-${run_id}-<resource>` 总长度可达 35+ 字符；ZStack VM name 上限 255 应该够用，但有些资源（如 SecGroup name）实际限制 64，需要逐一检查。
4. **是否允许并行 RUN_ID**：默认顺序跑；如果想 CI 多机并行，需要把 VLAN/IP 池切片成多份。

---

**作者**：claude (with sign-off) · **日期**：2026-04-26
