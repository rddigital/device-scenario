package models

import (
	"strconv"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/rddigital/device-scenario/internal/common"
)

type Rule struct {
	Id           string            `json:"id,omitempty" validate:"omitempty,uuid"`
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	AdminState   models.AdminState `json:"adminState,omitempty" validate:"omitempty,oneof='UNLOCKED' 'LOCKED'"`
	Actions      []Action          `json:"actions,omitempty"`
	NotifyEnable string            `json:"notifyEnable,omitempty" validate:"omitempty,oneof='true' 'false'"`
	Conditions   []Condition       `json:"conditions,omitempty"`
}

func RuleToProperties(rule Rule) map[string]models.ProtocolProperties {
	var protocol = make(map[string]models.ProtocolProperties)

	actionsProperty := ActionsToProperties(rule.Actions)
	if len(actionsProperty) > 0 {
		protocol[common.ActionsProperty] = actionsProperty
	}

	notifyEnableProperty := make(map[string]string)
	notifyEnableProperty[common.NotifyEnableProperty] = rule.NotifyEnable
	protocol[common.NotifyEnableProperty] = notifyEnableProperty

	conditionsProperty := ConditionsToProperties(rule.Conditions)
	if len(conditionsProperty) > 0 {
		protocol[common.ConditionsProperty] = conditionsProperty
	}

	return protocol
}

func RuleFromDevice(d models.Device) (rule Rule, ok bool) {
	rule.Id = d.Id
	rule.Name = d.Name
	rule.Description = d.Description
	rule.AdminState = d.AdminState

	if pp, ok := d.Protocols[common.ConditionsProperty]; ok {
		rule.Conditions = ConditionsFromProperties(pp)
	} else {
		return Rule{}, false
	}

	if pp, ok := d.Protocols[common.NotifyEnableProperty]; ok {
		rule.NotifyEnable = pp[common.NotifyEnableProperty]
	}

	if pp, ok := d.Protocols[common.ActionsProperty]; ok {
		rule.Actions = ActionsFromProperties(pp)
	} else {
		// Require (Actions != nil) or (Actions = nil and NotifyEnable = true)
		if enable, _ := strconv.ParseBool(rule.NotifyEnable); !enable {
			return Rule{}, false
		}
	}

	return rule, true
}
