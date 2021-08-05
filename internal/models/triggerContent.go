package models

type ContentTrigger struct {
	TriggerIndex *int  `json:"triggerIndex" validate:"required"`
	TriggerState *bool `json:"triggerState" validate:"required"`
}
