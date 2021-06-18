package application

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/http"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	ctCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/requests"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	ctModels "github.com/edgexfoundry/go-mod-core-contracts/v2/models"

	"github.com/google/uuid"
	"github.com/rddigital/device-scenario/internal/cache"
	"github.com/rddigital/device-scenario/internal/client"
	cm "github.com/rddigital/device-scenario/internal/common"
	"github.com/rddigital/device-scenario/internal/models"
)

var (
	host                 string
	port                 int
	lc                   logger.LoggingClient
	commandClient        interfaces.CommandClient
	intervalClient       interfaces.IntervalClient
	intervalActionClient interfaces.IntervalActionClient
	notificationClient   interfaces.NotificationClient
	ruleEngineClient     client.RuleEngineClient
)

const (
	StreamName         = "events"
	StreamSQLTemplate  = string(`{"sql":"create stream %s () WITH (FORMAT=\"JSON\", TYPE=\"edgex\")"}`)
	AddRuleSQLTemplate = string(`{"id":"%s","sql":"SELECT (collect(%s)[1] %s %s) as v FROM %s GROUP BY PADCOUNTWINDOW(2,1) FILTER(WHERE meta(deviceName) = \"%s\") HAVING collect(%s)[0]  %s %s OR collect(%s)[1] %s %s",
	"actions": [{
		"rest": {
			"url": "http://%s:%d/api/v2/rule/id/%s",
			"method": "post",
			"dataTemplate": "{\"triggerState\":{{.v}},\"triggerIndex\":%d}",
			"sendSingle": true
		  }
		}
	]}`)
	UpdateRuleSQLTemplate = string(`{"sql":"SELECT (collect(%s)[1] %s %s) as v FROM %s GROUP BY PADCOUNTWINDOW(2,1) FILTER(WHERE meta(deviceName) = \"%s\") HAVING collect(%s)[0]  %s %s OR collect(%s)[1] %s %s",
	"actions": [{
		"rest": {
			"url": "http://%s:%d/api/v2/rule/id/%s",
			"method": "post",
			"dataTemplate": "{\"triggerState\":{{.v}},\"triggerIndex\":%d}",
			"sendSingle": true
		  }
		}
	]}`)
)

func InitRuleApplication(l logger.LoggingClient, portService int, hostService, urlCoreCommand, urlNotification, urlSchduler, urlRuleEngine string) error {
	lc = l
	host = hostService
	port = portService
	commandClient = http.NewCommandClient(urlCoreCommand)
	notificationClient = http.NewNotificationClient(urlNotification)
	intervalClient = http.NewIntervalClient(urlSchduler)
	intervalActionClient = http.NewIntervalActionClient(urlNotification)
	ruleEngineClient = client.NewKuiperRuleClient(urlRuleEngine)

	_, err := ruleEngineClient.DescribeStream(StreamName)
	if err != nil {
		lc.Errorf(err.Error())
		streamStr := fmt.Sprintf(StreamSQLTemplate, StreamName)
		_, err = ruleEngineClient.CreateStream(streamStr)
		if err != nil {
			return err
		}
	}
	cache.InitCache()
	sysnRule()

	return nil
}

// remove all intervals, interval actions, rule engines that do not belong any rules
func sysnRule() {
	ctx := context.Background()
	// sysn intervals
	intervals, err := intervalClient.AllIntervals(ctx, 0, cm.DefaultLimit)
	if err != nil {
		for _, i := range intervals.Intervals {
			if id := parseName(i.Name); id != "" {
				if !cache.Rules().CheckExistsById(id) {
					intervalClient.DeleteIntervalByName(ctx, i.Name)
				}
			}
		}
	}

	// sysn interval actions
	intervalActions, err := intervalActionClient.AllIntervalActions(ctx, 0, cm.DefaultLimit)
	if err != nil {
		for _, i := range intervalActions.Actions {
			if id := parseName(i.Name); id != "" {
				if !cache.Rules().CheckExistsById(id) {
					intervalActionClient.DeleteIntervalActionByName(ctx, i.Name)
				}
			}
		}
	}

	// sysn rule engines
	ruleEngines, e := ruleEngineClient.ShowRules()
	if e != nil {
		for _, r := range ruleEngines {
			if id := parseName(r["id"].(string)); id != "" {
				if !cache.Rules().CheckExistsById(id) {
					ruleEngineClient.DropRule(r["id"].(string))
				}
			}
		}
	}
}

func parseName(name string) string {
	arrStr := strings.SplitAfter(name, cm.CharacterGenName)
	if len(arrStr) < 2 {
		return ""
	}
	return arrStr[0]
}

func AddRule(rule models.Rule) errors.EdgeX {
	if rule.Id == "" {
		id, _ := uuid.NewUUID()
		rule.Id = id.String()
	}

	if _, ok := cache.Rules().ForName(rule.Name); ok {
		err := fmt.Errorf("rule '%s' already exists", rule.Name)
		lc.Error(err.Error())
		return errors.NewCommonEdgeXWrapper(err)
	}

	// Alway unlock rule when change conditions
	rule.AdminState = ctModels.Unlocked

	lc.Debugf("adding rule: %s", rule.Name)

	err := addRuleConditions(rule)
	if err != nil {
		err = fmt.Errorf("add rule condition '%s' error: %s", rule.Name, err.Error())
		lc.Error(err.Error())
		return errors.NewCommonEdgeXWrapper(err)
	}
	lc.Debugf("add rule conditions '%s' success", rule.Name)

	properties := models.RuleToProperties(rule)
	var device = ctModels.Device{
		Id:             rule.Id,
		Name:           rule.Name,
		Description:    rule.Description,
		AdminState:     rule.AdminState,
		OperatingState: ctModels.Up,
		Protocols:      properties,
		ProfileName:    cm.AutoScenarioProfile,
		ServiceName:    cm.DeviceServiceName,
	}

	ds := service.RunningService()
	_, err = ds.AddDevice(device)
	if err != nil {
		err = fmt.Errorf("add rule '%s' to database error: %s", rule.Name, err.Error())
		lc.Error(err.Error())
		return errors.NewCommonEdgeXWrapper(err)
	}

	cache.Rules().Add(rule)
	lc.Debugf("add rule '%s' success", rule.Name)

	return nil
}

func GetAllRule() []models.Rule {
	return cache.Rules().All()
}

func GetRuleByName(name string) (models.Rule, errors.EdgeX) {
	rule, ok := cache.Rules().ForName(name)
	if !ok {
		return models.Rule{}, errors.NewCommonEdgeX(errors.KindEntityDoesNotExist, fmt.Sprintf("rule '%s' does not exists", name), nil)
	}
	return rule, nil
}

func UpdateRuleByName(name string, rule models.Rule) errors.EdgeX {
	oldRule, ok := cache.Rules().ForName(name)
	if !ok {
		err := fmt.Errorf("rule with id '%s' does not exists", oldRule.Id)
		lc.Error(err.Error())
		return errors.NewCommonEdgeXWrapper(err)
	}

	rule.Id = oldRule.Id
	if rule.Name == "" {
		rule.Name = oldRule.Name
	}
	if rule.Description == "" {
		rule.Description = oldRule.Description
	}
	if rule.AdminState == "" {
		rule.AdminState = oldRule.AdminState
	}
	if rule.NotifyEnable == "" {
		rule.NotifyEnable = oldRule.NotifyEnable
	}
	if len(rule.Actions) == 0 {
		rule.Actions = oldRule.Actions
	}
	if len(rule.Conditions) == 0 {
		rule.Conditions = oldRule.Conditions
	}

	lc.Debugf("updating rule with id '%s'", rule.Id)

	var forceUpdate bool = false
	if !reflect.DeepEqual(rule.Conditions, oldRule.Conditions) {
		forceUpdate = true
		// Alway unlock rule when change conditions
		rule.AdminState = ctModels.Unlocked
	}
	if rule.AdminState != oldRule.AdminState {
		forceUpdate = true
	}

	if forceUpdate {
		err := updateRuleConditions(rule, oldRule)
		if err != nil {
			err = fmt.Errorf("update rule condition with id '%s' error: %s", rule.Id, err.Error())
			lc.Error(err.Error())
			return errors.NewCommonEdgeXWrapper(err)
		}
		lc.Debugf("update rule conditions with id '%s' success", rule.Id)
	}

	properties := models.RuleToProperties(rule)
	var device = ctModels.Device{
		Id:             rule.Id,
		Name:           rule.Name,
		Description:    rule.Description,
		AdminState:     rule.AdminState,
		OperatingState: ctModels.Up,
		Protocols:      properties,
		ProfileName:    cm.AutoScenarioProfile,
		ServiceName:    cm.DeviceServiceName,
	}

	ds := service.RunningService()
	err := ds.UpdateDevice(device)
	if err != nil {
		err = fmt.Errorf("update rule with id '%s' in database error: %s", rule.Id, err.Error())
		lc.Error(err.Error())
		return errors.NewCommonEdgeXWrapper(err)
	}

	cache.Rules().Update(rule) // update rule and reset states
	lc.Debugf("update rule with id '%s' success", rule.Id)

	return nil
}

func DeleteRuleByName(name string) errors.EdgeX {
	rule, ok := cache.Rules().ForName(name)
	if !ok {
		err := fmt.Errorf("rule '%s' does not exists", name)
		lc.Error(err.Error())
		// return errors.NewCommonEdgeXWrapper(err)
	}

	lc.Debugf("deleting rule: %s", name)

	err := deleteRuleConditions(rule)
	if err != nil {
		err = fmt.Errorf("delete rule condition '%s' error: %s", rule.Name, err.Error())
		lc.Error(err.Error())
		// return errors.NewCommonEdgeXWrapper(err)
	}
	lc.Debugf("delete rule conditions '%s' success", rule.Name)

	ds := service.RunningService()
	err = ds.RemoveDeviceByName(name)
	if err != nil {
		err = fmt.Errorf("delete rule '%s' in database error: %s", rule.Name, err.Error())
		lc.Error(err.Error())
		return errors.NewCommonEdgeXWrapper(err)
	}

	cache.Rules().RemoveByName(name)
	lc.Debugf("delete rule '%s' success", rule.Name)
	return nil
}

func TriggerRuleById(id string, contentTrigger models.ContentTrigger) {
	rule, ok := cache.Rules().ForId(id)
	if !ok {
		lc.Errorf("rule with id '%s' does not exists", id)
		return
	}

	if rule.AdminState != ctModels.Unlocked {
		lc.Debugf("rule '%s' locked -> no excute actions", rule.Name)
		return
	}

	index := *contentTrigger.TriggerIndex
	if index >= len(rule.Conditions) {
		lc.Errorf("rule '%s' do not have the %d-rd condition", rule.Name, index)
		return
	}

	newState := *contentTrigger.TriggerState
	cache.Rules().UpdateStateRule(id, index, newState)
	defer func() {
		if rule.Conditions[index].Type == cm.ScheduleRuleType {
			cache.Rules().UpdateStateRule(id, index, false)
		}
	}()

	if checkRuleConditions(id) {
		lc.Infof("rule '%s' triggered", rule.Name)
		triggerRule(rule.Name)
	}
}

func generateName(baseName string, index int) string {
	return fmt.Sprintf("%s%s%d", baseName, cm.CharacterGenName, index)
}

func _addIntervalScheduleAdd(rule models.Rule, index int) error {
	ctx := context.Background()
	name := generateName(rule.Id, index)

	interval := dtos.Interval{
		Name:     name,
		Start:    rule.Conditions[index].StartTime,
		End:      rule.Conditions[index].EndTime,
		Interval: rule.Conditions[index].IntervalTime,
	}
	reqs := make([]requests.AddIntervalRequest, 1)
	reqs[0] = requests.AddIntervalRequest{
		BaseRequest: common.BaseRequest{
			Versionable: common.NewVersionable(),
		},
		Interval: interval,
	}
	_, err := intervalClient.Add(ctx, reqs)

	return err
}

func _updateIntervalSchedule(rule models.Rule, index int) error {
	ctx := context.Background()
	name := generateName(rule.Id, index)

	interval := dtos.UpdateInterval{
		Name:     &name,
		Start:    &rule.Conditions[index].StartTime,
		End:      &rule.Conditions[index].EndTime,
		Interval: &rule.Conditions[index].IntervalTime,
	}
	reqs := make([]requests.UpdateIntervalRequest, 1)
	reqs[0] = requests.UpdateIntervalRequest{
		BaseRequest: common.BaseRequest{
			Versionable: common.NewVersionable(),
		},
		Interval: interval,
	}
	_, err := intervalClient.Update(ctx, reqs)

	return err
}

func _addIntervalActionSchedule(rule models.Rule, index int) error {
	ctx := context.Background()

	name := generateName(rule.Id, index)
	action := dtos.IntervalAction{
		Name:         name,
		IntervalName: name,
		AdminState:   string(rule.AdminState),
		Address: dtos.Address{
			Type: "REST",
			Host: host,
			Port: port,
			RESTAddress: dtos.RESTAddress{
				Path:       "/rule/id/" + name,
				HTTPMethod: "POST",
			},
		},
		ContentType: ctCommon.ContentTypeJSON,
		Content:     fmt.Sprintf("{\"triggerState\":true, \"triggerIndex\":%d}", index),
	}
	reqs := make([]requests.AddIntervalActionRequest, 1)
	reqs[0] = requests.AddIntervalActionRequest{
		BaseRequest: common.BaseRequest{
			Versionable: common.NewVersionable(),
		},
		Action: action,
	}
	intervalActionClient.Add(ctx, reqs)

	return nil
}

func _updateIntervalActionSchedule(rule models.Rule, index int) error {
	ctx := context.Background()

	name := generateName(rule.Id, index)
	content := fmt.Sprintf("{\"triggerState\":true, \"triggerIndex\":\"%d\"}", index)
	action := dtos.UpdateIntervalAction{
		Name:         &name,
		IntervalName: &name,
		AdminState:   (*string)(&rule.AdminState),
		Address: &dtos.Address{
			Type: "REST",
			Host: host,
			Port: port,
			RESTAddress: dtos.RESTAddress{
				Path:       "/rule/id/" + name,
				HTTPMethod: "POST",
			},
		},
		Content: &content,
	}
	reqs := make([]requests.UpdateIntervalActionRequest, 1)
	reqs[0] = requests.UpdateIntervalActionRequest{
		BaseRequest: common.BaseRequest{
			Versionable: common.NewVersionable(),
		},
		Action: action,
	}
	intervalActionClient.Update(ctx, reqs)

	return nil
}

func updateScheduleState(rule models.Rule, index int) error {
	return _updateIntervalActionSchedule(rule, index)
}

func _addRuleEngine(rule models.Rule, index int) error {
	name := generateName(rule.Id, index)

	c := rule.Conditions[index]
	ruleStr := fmt.Sprintf(AddRuleSQLTemplate, name, c.ResourceThreshold, c.OperatorThreshold, c.ValueThreshold, StreamName, c.DeviceThreshold,
		c.ResourceThreshold, c.OperatorThreshold, c.ValueThreshold,
		c.ResourceThreshold, c.OperatorThreshold, c.ValueThreshold,
		host, port, rule.Id, index)
	_, err := ruleEngineClient.CreateRule(ruleStr)
	if err != nil {
		return err
	}

	if rule.AdminState == ctModels.Locked {
		_, err = ruleEngineClient.StopRule(name)
	}
	return err
}

func _updateRuleEngine(rule models.Rule, index int) error {
	name := generateName(rule.Id, index)

	c := rule.Conditions[index]
	ruleStr := fmt.Sprintf(UpdateRuleSQLTemplate, c.ResourceThreshold, c.OperatorThreshold, c.ValueThreshold, StreamName, c.DeviceThreshold,
		c.ResourceThreshold, c.OperatorThreshold, c.ValueThreshold,
		c.ResourceThreshold, c.OperatorThreshold, c.ValueThreshold,
		host, port, rule.Id, index)

	_, err := ruleEngineClient.UpdateRule(name, ruleStr)
	return err
}

func updateRuleEngineState(rule models.Rule, index int) error {
	name := generateName(rule.Id, index)

	var err error
	if rule.AdminState == ctModels.Unlocked {
		_, err = ruleEngineClient.RestartRule(name)
	} else {
		_, err = ruleEngineClient.StopRule(name)
	}
	return err
}

func _addRuleElement(rule models.Rule, index int) error {
	if rule.Conditions[index].Type == cm.ThresholdRuleType {
		return _addRuleEngine(rule, index)
	}

	if err := _addIntervalScheduleAdd(rule, index); err != nil {
		return err
	}
	if err := _addIntervalActionSchedule(rule, index); err != nil {
		return err
	}
	return nil
}

func _updateRuleElement(newRule models.Rule, oldRule models.Rule, index int) error {
	var updateContent, updateState bool
	var err error
	// TODO: check if condition exists in database

	// if update Kuiper
	if newRule.Conditions[index].Type == cm.ThresholdRuleType {
		if !reflect.DeepEqual(newRule.Conditions[index], oldRule.Conditions[index]) {
			updateContent = true
		}
		if newRule.AdminState != oldRule.AdminState {
			updateState = true
			if newRule.AdminState == ctModels.Unlocked && updateContent {
				// Kuiper alway restart rule (state = Unlocked) after update
				updateState = false
			}
		}

		if updateContent {
			if err = _updateRuleEngine(newRule, index); err != nil {
				return err
			}
		}
		if updateState {
			if err = updateRuleEngineState(newRule, index); err != nil {
				return err
			}
		}

		return nil
	}

	// if update Schedule
	if !reflect.DeepEqual(newRule.Conditions[index], oldRule.Conditions[index]) {
		updateContent = true
	}
	if newRule.AdminState != oldRule.AdminState {
		updateState = true
	}

	if updateContent {
		if err = _updateIntervalSchedule(newRule, index); err != nil {
			return err
		}
	}
	if updateState {
		if err := updateScheduleState(newRule, index); err != nil {
			return err
		}
	}
	return nil
}

func _deleteRuleElement(rule models.Rule, index int) error {
	name := generateName(rule.Id, index)

	if rule.Conditions[index].Type == cm.ThresholdRuleType {
		_, err := ruleEngineClient.DropRule(name)
		return err
	}

	ctx := context.Background()

	if _, err := intervalActionClient.DeleteIntervalActionByName(ctx, name); err != nil {
		return err
	}
	if _, err := intervalClient.DeleteIntervalByName(ctx, name); err != nil {
		return err
	}
	return nil
}

func addRuleConditions(rule models.Rule) error {
	for index := range rule.Conditions {
		if err := _addRuleElement(rule, index); err != nil {
			return fmt.Errorf(err.Error())
		}
	}

	return nil
}

func updateRuleConditions(newRule models.Rule, oldRule models.Rule) error {
	if len(newRule.Conditions) < len(oldRule.Conditions) {
		for index := len(newRule.Conditions); index < len(oldRule.Conditions); index++ {
			if err := _deleteRuleElement(oldRule, index); err != nil {
				return err
			}
		}
	}

	for index := range newRule.Conditions {
		if index < len(oldRule.Conditions) {
			if newRule.Conditions[index].Type != oldRule.Conditions[index].Type {
				if err := _deleteRuleElement(oldRule, index); err != nil {
					return err
				}
				if err := _addRuleElement(newRule, index); err != nil {
					return err
				}
			} else {
				if err := _updateRuleElement(newRule, oldRule, index); err != nil {
					return err
				}
			}
		} else {
			if err := _addRuleElement(newRule, index); err != nil {
				return err
			}
		}
	}

	return nil
}

func deleteRuleConditions(rule models.Rule) error {
	for index := range rule.Conditions {
		if err := _deleteRuleElement(rule, index); err != nil {
			return fmt.Errorf(err.Error())
		}
	}

	return nil
}

func checkRuleConditions(id string) bool {
	rule, ok := cache.Rules().ForId(id)
	if !ok {
		return false
	}
	if len(rule.Conditions) <= 0 {
		return false
	}

	result := cache.Rules().GetStateRule(id, 0)
	for index := 1; index < len(rule.Conditions); index++ {
		state := cache.Rules().GetStateRule(id, index)
		if rule.Conditions[index].Logic == cm.AndLogic {
			result = result && state
		} else {
			result = result || state
		}
	}

	return result
}

func triggerRule(name string) {
	rule, ok := cache.Rules().ForName(name)
	if !ok {
		return
	}

	ctx := context.Background()
	if rule.Actions != nil {
		for index, action := range rule.Actions {
			bodyParam, err := parseBody(action.Body)
			if err != nil {
				lc.Errorf("Trigger rule '%s' error: parse content of action[%d] error: %s -> Abort action[%d]", name, index, err.Error(), index)
				continue
			}
			_, err = commandClient.IssueSetCommandByName(ctx, action.DeviceName, action.CommandName, bodyParam)
			if err != nil {
				lc.Errorf("Trigger rule '%s' error: execute action[%d] error:%s", name, index, err.Error())
			} else {
				lc.Debugf("Trigger rule '%s' execute action[%d] success", name, index)
			}
		}
	}

	if enable, _ := strconv.ParseBool(rule.NotifyEnable); enable {
		contentRequest := requests.NewAddNotificationRequest(
			dtos.Notification{
				Category: "trigger-event",
				Content:  fmt.Sprintf("auto-scenario '%s' triggered", name),
				Sender:   "scenario-service",
				Severity: ctModels.Normal,
				Status:   ctModels.New,
			},
		)
		_, err := notificationClient.SendNotification(ctx, []requests.AddNotificationRequest{contentRequest})
		if err != nil {
			lc.Errorf("Trigger rule '%s' error: send notification error:%s", name, err.Error())
		} else {
			lc.Debugf("Trigger rule '%s' send successful notification", name)
		}
	}
}

func parseBody(params string) (paramMap map[string]string, err error) {
	err = json.Unmarshal([]byte(params), &paramMap)
	if err != nil {
		return
	}

	if len(paramMap) == 0 {
		err = fmt.Errorf("no parameters specified")
		return
	}

	return
}
