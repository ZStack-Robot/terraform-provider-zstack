package view

import "time"

type ManagementNodeInventoryView struct {
	UUID      string    `json:"uuid"`      //资源的UUID，唯一标示该资源
	HostName  string    `json:"hostName"`  //宿主机名称
	JoinDate  time.Time `json:"joinDate"`  //加入时间
	HeartBeat time.Time `json:"heartBeat"` //心跳时间
}
