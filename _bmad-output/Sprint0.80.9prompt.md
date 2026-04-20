Story 08: 补全不完整测试维度

  执行 Story 08: 补全 secgroup_attachment Import 和 virtual_router_instance CheckDestroy。

  ## 背景

  项目有 2 个"不完整级"资源测试需要补全:
  1. secgroup_attachment — 缺 Import Step（且资源代码未实现 ImportState）
  2. virtual_router_instance — 缺 CheckDestroy

  reserved_ip 已确认完整（有 Import Step + CheckDestroy），无需处理。

  ## Task 1: secgroup_attachment ImportState 实现 + Import Step

  资源文件: zstack/provider/resource_zstack_networking_secgroup_attachment.go
  测试文件: zstack/provider/resource_zstack_networking_secgroup_attachment_test.go

  现状:
  - 资源 **未实现 ResourceWithImportState 接口**
  - 资源无 UUID 字段，使用 composite key: secgroup_uuid + nic_uuid
  - 测试只有 Create Step，无 Import Step

  步骤:
  1. 在资源代码中实现 ImportState 方法（composite ID 格式建议: secgroup_uuid:nic_uuid，用 ":" 分隔）
  2. 在 var 块中添加 `_ resource.ResourceWithImportState = &securityGroupAttachmentResource{}`
  3. 在测试中添加 Import Step + importStateIdFunc
  4. ImportStateVerifyIdentifierAttribute 需要设为 "id"（因为该资源无 uuid 属性，主键是 id）

  参考:
  - 已有 composite ID import 实现: resource_zstack_scheduler_job.go（用 ":"）、resource_zstack_global_config.go（用 "/"）
  - 已有测试模式: resource_zstack_scheduler_job_test.go 的 importStateIdSchedulerJob

  ## Task 2: virtual_router_instance CheckDestroy

  资源文件: zstack/provider/resource_zstack_virtual_router_instance.go
  测试文件: zstack/provider/resource_zstack_virtual_router_instance_test.go

  现状:
  - 测试使用 mock server + resource.UnitTest（非真实 API）
  - 无 CheckDestroy 函数
  - Delete mock 只返回 200 + {}

  步骤:
  1. 在 check_destroy_test.go 新增 testAccCheckVirtualRouterInstanceDestroy
  2. 在测试中引用（注意: UnitTest 模式下 CheckDestroy 仍会执行，mock DELETE 端点需正确响应 404 on re-read）
  3. 需要在 mock 中添加状态跟踪: DELETE 后 GET 应返回空/404，否则 CheckDestroy 会误判资源仍存在

  ## 验证

  - go build ./... 编译通过
  - go test ./zstack/provider/ -run TestAccZStackSecurityGroupAttachment -v（需环境）
  - go test ./zstack/provider/ -run TestAccZStackVirtualRouterInstance -v（mock，可本地跑）

  完成后标注每个 Task 的实际结果。

  ---
  Story 09: 批量 Disappears 测试

  执行 Story 09: 为现有资源批量补全 Disappears 测试。

  ## 背景

  当前 Disappears 测试覆盖: 6 个（zone, account, ssh_key_pair, affinity_group, cluster, iam2_project）
  目标: 新增约 26 个，覆盖所有已有验收测试的资源。

  ## 框架

  通用基础设施已就绪:
  - check_disappears_test.go 已有 stateCheckDisappears 通用工厂（第59行）
  - 模式: stateCheck<X>Disappears 工厂函数 → 在 _test.go 中写 TestAcc<X>Resource_disappears 测试

  ## 每个 disappears 测试的固定模式

  1. 在 check_disappears_test.go 添加工厂函数:
  ```go
  func stateCheck<X>Disappears(resourceAddress string) statecheck.StateCheck {
      return stateCheckDisappears(resourceAddress, func(cli *client.ZSClient, id string) error {
          return cli.Delete<X>(id, param.DeleteModePermissive)
      })
  }

  2. 在对应 _test.go 中添加测试:
  func TestAcc<X>Resource_disappears(t *testing.T) {
      // 复用已有 Create config
      resource.ParallelTest(t, resource.TestCase{
          ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
          CheckDestroy:             testAccCheck<X>Destroy,
          Steps: []resource.TestStep{
              {
                  Config: <same create config>,
                  ConfigStateChecks: []statecheck.StateCheck{
                      stateCheck<X>Disappears("zstack_<x>.test"),
                  },
                  ExpectNonEmptyPlan: true,
              },
          },
      })
  }

  资源清单（按批次）

  Batch 1: 有 Get 方法的资源（简单）

  image, volume, l2vlan_network, load_balancer, load_balancer_listener,
  virtual_router_image, virtual_router_offering, auto_scaling_group,
  port_forwarding_rule, instance_offering, disk_offering, tag,
  security_group, security_group_rule

  Batch 2: 用 Query 方法的资源（Delete 函数需查 SDK）

  certificate, webhook, vip, reserved_ip, user, role, policy,
  iam2_virtual_id, iam2_organization, sns_topic, scheduler_trigger,
  scheduler_job, alarm

  跳过（不适用）

  - sns_email_endpoint — API 不支持删除
  - sns_http_endpoint — API 不支持删除
  - global_config — 语义为"重置"非"删除"
  - secgroup_attachment — 无独立删除语义
  - tag_attachment — 关联删除

  注意事项

  1. Delete 方法名不一定是 Delete，需查各资源的 Delete 方法确认 SDK 调用
  2. 部分资源的 Delete 使用 Query-style 而非 Get-style（参考 check_destroy_test.go 现有模式）
  3. reserved_ip 的删除方法特殊（可能是 releaseReservedIp，不是 delete），需查资源代码
  4. 有回收站的资源（volume）用 DeleteModePermissive 确保直接删除

  验证

  - go build ./... 编译通过
  - 抽样运行 5 个: zone, account, image, volume, policy
