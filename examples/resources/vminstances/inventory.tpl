[masters]
%{ for ip in master ~}
${ip} internal_ip=${ip}  zone=zone1 
%{ endfor ~} 

[masters-shadow]
%{ for ip in mastershadow ~}
${ip} internal_ip=${ip}  zone=zone1 
%{ endfor ~}


[slaves]


[masters:vars]
control_plane=k8smaster.qfusion.irds
role=master

[masters-shadow:vars]
role=shadow

[slaves:vars]
role=slave

[all:vars]
ansible_ssh_user=root
ansible_ssh_port=22
ansible_ssh_extra_args='-o StrictHostKeyChecking=no'
plat_info=irds
# registry存储目录,默认/gaea
registry_data_dir=/gaea
# 本地存储目录
data_dir=/opt/qfusion
ntp_server=PROFILE_NTP
# CPU超分比
cpu_overrate=2
# 内存超分比
memory_overrate=1
# 是否开启强制反亲和
cluster_forced_anti_affinity=true
# 开启本地存储
enable_local_storage=true
# 本地存储超分比
local_storage_overrate=2
# 本地存储的IOPS
localpv_iops=1000000
# 卷大小，单位GB
localpv_storage_size=100000
# 本地存储默认介质，仅支持SSD,NVMe,HDD
localpv_storage_medium=SSD
# 是否启用分布式存储
enable_linstor_storage=false
# 必须与 linstor-installer 所配置的值一致
pool_name=data_pool
# 为 elasticsearch 设置单独的数据目录，不填时使用默认目录 `{data_dir}/elasticsearch`
es_data_path=
# 为监控组件设置单独的数据目录，不填时使用默认目录 `{data_dir}/monitor`
monitor_data_path=
# 预置的数据库镜像类型，不填时仅包含基础镜像
database_modules=all
only_run_master=true
# 是否开启用户备份接口权限
enable_user_backup_storage_interface=false
# 是否禁用kube-proxy，内核版本必须大于等于4.19.57才能禁用
disable_kube_proxy=true
# 是否使用crd，用于固定ip需求
ipam_use_crd=false
# 产品部署模式，qfusion: qfusion模式; qdataOnly: 仅安装qdata部分; qdataFull: 全量安装
deploy_mode=qfusion
# 安装时所有节点的cpu架构 可选 [amd64, arm64, multiple] multiple混合部署
cpuPlatform=amd64

