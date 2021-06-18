package cache

import (
	"sync"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	ctModels "github.com/edgexfoundry/go-mod-core-contracts/v2/models"
	"github.com/rddigital/device-scenario/internal/common"
	"github.com/rddigital/device-scenario/internal/models"
)

type RuleCache interface {
	CheckExistsById(id string) bool
	ForId(id string) (models.Rule, bool)
	ForName(name string) (models.Rule, bool)
	All() []models.Rule
	Add(rule models.Rule)
	Update(rule models.Rule)
	RemoveByName(name string)
	UpdateStateRule(id string, index int, state bool)
	GetStateRule(id string, index int) bool
}

type ruleCache struct {
	ruleMap   map[string]models.Rule // key is rule id
	nameIdMap map[string]string
	stateMap  map[string][]bool
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
		if d.ProfileName == common.AutoScenarioProfile {
			scenarios = append(scenarios, d)
		}
	}

	sizeMap := len(scenarios) + 1 // minimum = 1
	ruleMap := make(map[string]models.Rule, sizeMap)
	nameIdMap := make(map[string]string, sizeMap)
	stateMap := make(map[string][]bool, sizeMap)

	for _, s := range scenarios {
		if rule, ok := models.RuleFromDevice(s); ok {
			ruleMap[rule.Id] = rule
			nameIdMap[rule.Name] = rule.Id
			stateMap[rule.Id] = make([]bool, len(rule.Conditions))
		} else {
			// Remove scenarios are not valid
			ds.RemoveDeviceByName(s.Name)
		}
	}

	rc = &ruleCache{
		ruleMap:   ruleMap,
		nameIdMap: nameIdMap,
		stateMap:  stateMap,
	}
}

func (rc *ruleCache) CheckExistsById(id string) bool {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	_, ok := rc.ruleMap[id]
	return ok
}

func (rc *ruleCache) ForId(id string) (models.Rule, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	rule, ok := rc.ruleMap[id]
	return rule, ok
}

func (rc *ruleCache) ForName(name string) (models.Rule, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	id, ok := rc.nameIdMap[name]
	if !ok {
		return models.Rule{}, false
	}
	rule, ok := rc.ruleMap[id]
	return rule, ok
}

func (rc *ruleCache) All() []models.Rule {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	rules := make([]models.Rule, 0, len(rc.ruleMap))
	for _, r := range rc.ruleMap {
		rules = append(rules, r)
	}

	return rules
}

func (rc *ruleCache) Add(rule models.Rule) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	rc.update(rule)
}

func (rc *ruleCache) Update(rule models.Rule) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	rc.update(rule)
}

func (rc *ruleCache) update(rule models.Rule) {
	rc.delete(rule.Name)

	rc.nameIdMap[rule.Name] = rule.Id
	rc.ruleMap[rule.Id] = rule
	rc.stateMap[rule.Id] = make([]bool, len(rule.Conditions))
}

func (rc *ruleCache) delete(name string) {
	id, ok := rc.nameIdMap[name]
	if !ok {
		return
	}

	delete(rc.nameIdMap, name)
	delete(rc.ruleMap, id)
	delete(rc.stateMap, id)
}

func (rc *ruleCache) RemoveByName(name string) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	rc.delete(name)
}

func (rc *ruleCache) UpdateStateRule(id string, index int, state bool) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	if len(rc.stateMap[id]) <= index {
		return
	}
	rc.stateMap[id][index] = state
}

func (rc *ruleCache) GetStateRule(id string, index int) bool {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	if len(rc.stateMap[id]) <= index {
		return false
	}
	return rc.stateMap[id][index]
}
