// Copyright (c) HashiCorp, Inc.

package view

import "time"

type SessionView struct {
	UUID        string    `json:"uuid"`        //资源的UUID，唯一标示该资源
	AccountUuid string    `json:"accountUuid"` //账户UUID
	UserUuid    string    `json:"userUuid"`    //用户UUID
	ExpiredDate time.Time `json:"expiredDate"` //会话过期日期
	CreateDate  time.Time `json:"createDate"`  //创建时间
}

type WebUISessionView struct {
	SessionId       string `json:"sessionId"`   //资源的UUID
	AccountUuid     string `json:"accountUuid"` //账户UUID
	UserUuid        string `json:"userUuid"`    //用户UUID
	UserName        string `json:"username"`    //用户名
	LoginType       string `json:"loginType"`
	CurrentIdentity string `json:"currentIdentity"`
	ZSVersion       string `json:"zsVersion"` //ZStack Cloud详细版本
}
