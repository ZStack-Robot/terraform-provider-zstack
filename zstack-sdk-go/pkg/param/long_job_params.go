// Copyright (c) ZStack.io, Inc.

package param

type LongJobState string

const (
	Waiting   LongJobState = "Waiting"
	Suspended LongJobState = "Suspended"
	Running   LongJobState = "Running"
	Succeeded LongJobState = "Succeeded"
	Canceling LongJobState = "Canceling"
	Canceled  LongJobState = "Canceled"
	Failed    LongJobState = "Failed"
)

type SubmitLongJobParam struct {
	BaseParam
	Params SubmitLongJobDetailParam `json:"params"`
}
type SubmitLongJobDetailParam struct {
	Name               *string  `json:"name"`
	Description        *string  `json:"description"`
	JobName            string   `json:"jobName" validate:"required"`
	JobData            string   `json:"jobData" validate:"required"`
	ResourceUuid       *string  `json:"resourceUuid" `
	TargetResourceUuid *string  `json:"targetResourceUuid" `
	TagUuids           []string `json:"tagUuids" `
}

type UpdateLongJobParam struct {
	BaseParam
	UpdateLongJob UpdateLongJobDetailParam `json:"updateLongJob"`
}
type UpdateLongJobDetailParam struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
