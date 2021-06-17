package cache

import (
	"sync"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	ctModels "github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/rddigital/device-scenario/internal/models"
)

type RuleCache interface {
	ForId(id string) (models.Rule, bool)
	ForName(name string) (models.Rule, bool)
	All() []models.Rule
	Add(rule models.Rule) error
	Update(rule models.Rule) error
	RemoveByName(name string) error
}

type ruleCache struct {
	ruleMap   map[string]models.Rule // key is rule id
	nameIdMap map[string]string
	mutex     sync.RWMutex
}

var (
	rc *ruleCache
)

func Rules() RuleCache {
	return rc
}

// InitCache Init basic state for cache
func InitCache() {
	ds := service.RunningService()
	devices := ds.Devices()
	scenarios := make([]ctModels.Device, 0, len(devices))

	for _, d := range devices {
		if d.ProfileName == "AutoScenario" {
			scenarios = append(scenarios, d)
		}
	}

	sizeMap := len(scenarios) + 1 // minimum = 1
	ruleMap := make(map[string]models.Rule, sizeMap)
	nameIdMap := make(map[string]string, sizeMap)
	for _, s := range scenarios {
		if rule, ok := models.RuleFromDevice(s); ok {
			ruleMap[rule.Id] = rule
			nameIdMap[rule.Name] = rule.Id
		}
	}

	rc = &ruleCache{
		ruleMap:   ruleMap,
		nameIdMap: nameIdMap,
	}
}

func (rc *ruleCache) ForId(id string) (models.Rule, bool) {

	return models.Rule{}, false
}

func (rc *ruleCache) ForName(name string) (models.Rule, bool) {

	return models.Rule{}, false
}

func (rc *ruleCache) All() []models.Rule {

	return nil
}

func (rc *ruleCache) Add(rule models.Rule) error {

	return nil
}

func (rc *ruleCache) Update(rule models.Rule) error {

	return nil
}

func (rc *ruleCache) RemoveByName(name string) error {

	return nil
}
