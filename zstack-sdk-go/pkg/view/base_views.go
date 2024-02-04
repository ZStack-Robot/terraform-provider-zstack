package view

import "time"

type BaseInfoView struct {
	UUID        string `json:"uuid"`        //资源的UUID，唯一标示该资源
	Name        string `json:"name"`        //资源名称
	Description string `json:"description"` //资源的详细描述
}

type BaseTimeView struct {
	CreateDate time.Time `json:"createDate"` //创建时间
	LastOpDate time.Time `json:"lastOpDate"` //最后一次修改时间
}

type ErrorCodeView struct {
	Code        string                 `json:"code"`        //错误码号，错误的全局唯一标识，例如SYS.1000, HOST.1001
	Description string                 `json:"description"` // 	错误的概要描述
	Details     string                 `json:"details"`     //错误的详细信息
	Elaboration string                 `json:"elaboration"` //保留字段，默认为null
	Location    string                 `json:"location"`
	Cost        string                 `json:"cost"`
	Opaque      map[string]interface{} `json:"opaque"` //保留字段，默认为null
	Cause       *ErrorCodeView         `json:"cause"`  //根错误，引发当前错误的源错误，若无原错误，该字段为null
}
