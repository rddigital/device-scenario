package models

import (
	"encoding/json"
	"strconv"

	"github.com/rddigital/device-scenario/internal/common"
)

type Condition struct {
	Logic string `json:"logic" validate:"required,oneof='and' 'or'"`
	Type  string `json:"type" validate:"required"`

	// Time condition
	StartTime    string `json:"startTime,omitempty"`
	EndTime      string `json:"endTime,omitempty"`
	IntervalTime string `json:"intervalTime,omitempty"`

	// Threshold condition
	DeviceThreshold   string `json:"deviceThreshold,omitempty"`
	OperatorThreshold string `json:"operatorThreshold,omitempty" validate:"omitempty,oneof='>' '<' '=' '>=' '<='"`
	ResourceThreshold string `json:"resourceThreshold,omitempty"`
	ValueThreshold    string `json:"valueThreshold,omitempty"`
}

func ConditionsToProperties(conditions []Condition) map[string]string {
	properties := make(map[string]string, len(conditions))
	for index, condition := range conditions {
		value, err := json.Marshal(condition)
		if err != nil {
			continue
		}
		key := strconv.Itoa(index)
		properties[key] = string(value)
	}
	return properties
}

func ConditionsFromProperties(properties map[string]string) []Condition {
	conditions := make([]Condition, len(properties))
	for key, value := range properties {
		var condition Condition
		err := json.Unmarshal([]byte(value), &condition)
		if err != nil {
			continue
		}
		if common.Validate(condition) != nil {
			continue
		}

		index, _ := strconv.Atoi(key)
		conditions[index] = condition
	}

	return conditions
}
