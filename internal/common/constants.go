package common

import contractsCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/common"

const (
	ActionsProperty      = "actions"
	NotifyEnableProperty = "notify"
	ConditionsProperty   = "conditions"

	ScheduleRuleType  = "schedule"
	ThresholdRuleType = "threshold"
)

// Constants related to defined routes in the v2 service APIs
const (
	ApiRuleRoute            = contractsCommon.ApiBase + "/" + Rule                                          // POST
	ApiAllRuleRoute         = ApiRuleRoute + "/" + contractsCommon.All                                      // GET
	ApiRuleByNameRoute      = ApiRuleRoute + "/" + contractsCommon.Name + "/{" + contractsCommon.Name + "}" // GET, PUT
	ApiRuleTriggerByIdRoute = ApiRuleRoute + "/" + contractsCommon.Id + "/{" + contractsCommon.Id + "}"     // POST
)

// Constants related to defined url path names and parameters in the v2 service APIs
const (
	Rule = "rule"
)

// Constants related to defined profiles and device service
const (
	ManualScenarioProfile = "ManualScenario"
	AutoScenarioProfile   = "AutoScenario"

	DeviceServiceName = "scenario"
)

// Constants related to defined logic type
const (
	AndLogic = "and"
	OrLogic  = "or"
)

const (
	CharacterGenName = "_"
	DefaultLimit     = 1000
)
