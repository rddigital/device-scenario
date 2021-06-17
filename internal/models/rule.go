package models

import (
	"strconv"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/rddigital/device-scenario/internal/common"
)

type Rule struct {
	Id           string            `json:"id,omitempty" validate:"omitempty,uuid"`
	Name         string            `json:"name" validate:"required"`
	Description  string            `json:"description,omitempty"`
	AdminState   models.AdminState `json:"adminState,omitempty" validate:"omitempty,oneof='ENABLE' 'DISABLE'"`
	Actions      []Action          `json:"actions,omitempty"`
	NotifyEnable bool              `json:"notifyEnable,omitempty"`
	Conditions   []Condition       `json:"conditions" validate:"required,gt=0"`
}

func RuleToProperties(rule Rule) map[string]models.ProtocolProperties {
	var protocol = make(map[string]models.ProtocolProperties)

	actionsProperty := ActionsToProperties(rule.Actions)
	if len(actionsProperty) > 0 {
		protocol[common.ActionsProperty] = actionsProperty
	}

	notifyEnableProperty := make(map[string]string)
	notifyEnableProperty[common.NotifyEnableProperty] = strconv.FormatBool(rule.NotifyEnable)
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
		strBool := pp[common.ConditionsProperty]
		rule.NotifyEnable, _ = strconv.ParseBool(strBool)
	}

	if pp, ok := d.Protocols[common.ActionsProperty]; ok {
		rule.Actions = ActionsFromProperties(pp)
	} else {
		// Require (Actions != nil) or (Actions = nil and NotifyEnable = true)
		if !rule.NotifyEnable {
			return Rule{}, false
		}
	}

	return rule, false
}
