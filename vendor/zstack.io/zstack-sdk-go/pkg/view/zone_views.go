package view

type ZoneView struct {
	BaseInfoView
	BaseTimeView

	State string `json:"state"` //状态
	Type  string `json:"type"`  //类型
}
