## Models

1. Rule

```
type Rule struct {
	Id           string            `json:"id,omitempty" validate:"omitempty,uuid"`
	Name         string            `json:"name,omitempty"`
	Description  string            `json:"description,omitempty"`
	AdminState   models.AdminState `json:"adminState,omitempty" validate:"omitempty,oneof='UNLOCKED' 'LOCKED'"`
	Actions      []Action          `json:"actions,omitempty"`
	NotifyEnable string            `json:"notifyEnable,omitempty" validate:"omitempty,oneof='true' 'false'"`
	Conditions   []Condition       `json:"conditions,omitempty"`
}
```
2. Action

```
type Action struct {
	DeviceName  string `json:"deviceName" validate:"required"`
	CommandName string `json:"commandName" validate:"required"`
	Body        string `json:"body" validate:"required"`
}
```

3. Condition

```
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
```

> In Schedule, Kuiper service: `Rule.Id = Interval.Name = IntervalAction.Name = "_" + Rule.Id + "_" + "{index}"`

> The body of Kuiper action and IntervalAction: `{"triggerIndex":"{index}", "TriggerState":"true/false"}`

## API

1. `GET` `api/v2/rule/name/{rule-name}`

    - Get rule by name
2. `GET` `api/v2/rule/all`

    - Get all rule
3. `POST` `api/v2/rule`

   - Add new rule
    
Example:

```
curl --location --request POST 'http://localhost:59990/api/v2/rule'
--header 'Content-Type: application/json'
--data-raw '{
    "apiVersion":"v2",
    "rule":{
        "name": "auto2",
        "description": "this is auto trigger scenario",
        "adminState": "UNLOCKED",
        "actions": [
            {
                "deviceName": "Random-Integer-Generator01",
                "commandName": "Min_Int8",
                "body": "{\"Min_Int8\": \"0\"}"
            },
            {
                "deviceName": "Random-Integer-Generator01",
                "commandName": "Max_Int8",
                "body": "{\"Max_Int8\": \"10\"}"
            }
        ],
        "notifyEnable": "false",
        "conditions": [
            {
                "logic": "and",
                "type": "threshold",
                "deviceThreshold": "dev1",
                "resourceThreshold": "temperature",
                "operatorThreshold": ">=",
                "valueThreshold": "30"
            },
            {
                "logic": "or",
                "type": "threshold",
                "deviceThreshold": "dev2",
                "resourceThreshold": "temperature",
                "operatorThreshold": "<=",
                "valueThreshold": "30"
            }
        ]
    }
}'
```
4. `PUT` `api/v2/rule/name/{rule-name}`

    - Start/stop rule

Example: start

```
curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
--header 'Content-Type: application/json'
--data-raw '{
    "adminState": "UNLOCKED"
}'
```

Example: stop

```
curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
--header 'Content-Type: application/json'
--data-raw '{
    "adminState": "LOCKED"
}'
```

5. `PUT` `api/v2/rule/name/{rule-name}`

    - Update rule
    - If update conditions, rule alway restart

Example: update multi-info

```
curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
--header 'Content-Type: application/json'
--data-raw '{
    "apiVersion":"v2",
    "rule": {
        "name": "auto1",
        "description": "this is auto scenario",
        "adminState": "UNLOCKED",
        "actions": [
            {
                "deviceName": "Random-Integer-Generator01",
                "commandName": "Min_Int8",
                "body": "{\"Min_Int8\": \"10\"}"
            },
            {
                "deviceName": "Random-Integer-Generator01",
                "commandName": "Max_Int8",
                "body": "{\"Max_Int8\": \"20\"}"
            }
        ],
        "notifyEnable": "false"
    }    
}'
```

Example: update only EnableNotify
    
```
curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
--header 'Content-Type: application/json'
--data-raw '{
    "apiVersion": "v2",
    "rule":{
        "notifyEnable": "false"
    }
}
```

Example: update conditions
    
```
curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
--header 'Content-Type: application/json'
--data-raw '{
    "apiVersion": "v2",
    "rule":{
        "conditions": [
        {
            "logic": "and",
            "type": "threshold",
            "deviceThreshold": "dev1",
            "resourceThreshold": "temperature",
            "operatorThreshold": ">=",
            "valueThreshold": "30"
        },
        {
            "logic": "and",
            "type": "threshold",
            "deviceThreshold": "dev2",
            "resourceThreshold": "temperature",
            "operatorThreshold": "<=",
            "valueThreshold": "30"
        }
        ]
    }
}
```

6. `DELETE` `api/v2/rule/name/{rule-name}`

    - Delete rule

Example:

```
curl --location --request DELETE 'http://localhost:59990/api/v2/rule/name/auto1
```