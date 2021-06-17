package models

type ContentTrigger struct {
	TriggerIndex string `json:"triggerIndex" validate:"required"`
	TriggerState string `json:"triggerState" validate:"required"`
}
