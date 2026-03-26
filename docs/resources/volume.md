# zstack_volume 资源文档

## 资源说明
zstack_volume 用于管理 ZStack 云平台中的数据盘资源，支持创建、挂载到虚拟机、卸载、扩容、删除等操作。

## 参数
- name (必填)：卷名称
- description (可选)：描述
- disk_offering_uuid (可选)：磁盘规格UUID
- disk_size (可选)：磁盘大小（字节）
- primary_storage_uuid (可选)：主存储UUID
- resource_uuid (可选)：资源UUID
- tag_uuids (可选)：标签UUID列表
- vm_instance_uuid (可选)：要挂载的虚拟机UUID

## 属性
- uuid：卷UUID
- 其余同参数

## 示例
```hcl
resource "zstack_volume" "example" {
  name                 = "data-volume-1"
  disk_offering_uuid   = "uuid-xxxx"
  disk_size            = 107374182400 # 100G
  primary_storage_uuid = "uuid-ps"
  vm_instance_uuid     = "uuid-vm"
}
```

---
# zstack_volume_snapshot 资源文档

## 资源说明
zstack_volume_snapshot 用于管理数据盘快照，支持创建、删除、导入等操作。

## 参数
- name (必填)：快照名称
- description (可选)：描述
- volume_uuid (必填)：所属数据盘UUID

## 属性
- uuid：快照UUID
- 其余同参数

## 示例
```hcl
resource "zstack_volume_snapshot" "snap1" {
  name        = "snap-1"
  volume_uuid = zstack_volume.example.uuid
}
```

---
# zstack_volumes 数据源

## 用法
```hcl
data "zstack_volumes" "all" {}
```

## 属性
- volumes：卷列表（包含 uuid、name、description、disk_offering_uuid、disk_size、primary_storage_uuid、vm_instance_uuid 等）

---
# zstack_volume_snapshots 数据源

## 用法
```hcl
data "zstack_volume_snapshots" "all" {}
```

## 属性
- snapshots：快照列表（包含 uuid、name、description、volume_uuid 等）
