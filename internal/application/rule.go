package application

import (
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	"github.com/rddigital/device-scenario/internal/models"
)

func AddRule(rule models.Rule) errors.EdgeX {
	return nil
}

func GetAllRule() ([]models.Rule, errors.EdgeX) {
	return []models.Rule{}, nil
}

func GetRuleByName(name string) (models.Rule, errors.EdgeX) {
	return models.Rule{}, nil
}

func UpdateRuleByName(name string, rule models.Rule) errors.EdgeX {
	return nil
}

func DeleteRuleByName(name string) errors.EdgeX {
	return nil
}

func TriggerRuleById(id string, contentTrigger models.ContentTrigger) {
}
