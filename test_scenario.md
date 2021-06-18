## Run edgex in a terminal
1. Download
> `cd $GOPATH/src/github.com/edgexfoundry/ && git clone https://github.com/edgexfoundry/edgex-go.git`
2. Build
> `echo 2.0.0 > $GOPATH/src/github.com/edgexfoundry/edgex-go/VERSION`
> 
> `cd $GOPATH/src/github.com/edgexfoundry/edgex-go && make build && cd`
1. Run
> `sudo service redis-server start`

> `cd $GOPATH/src/github.com/edgexfoundry/edgex-go && make run`

## Run random service in a terminal
1. Download
> `cd $GOPATH/src/github.com/edgexfoundry/ && git clone https://github.com/edgexfoundry/device-random.git`
2. Build & Run
> `echo 2.0.0 > $GOPATH/src/github.com/edgexfoundry/device-random/VERSION`
> 
> `cd $GOPATH/src/github.com/edgexfoundry/device-random && make && cd`

> `export EDGEX_SECURITY_SECRET_STORE=false && cd $GOPATH/src/github.com/edgexfoundry/device-random/cmd/ && ./device-random`

1. Test
> `curl -X PUT -d '{"Min_Int8": "0"}' localhost:59882/api/v2/device/name/Random-Integer-Generator01/Min_Int8`
> 
> `curl -X PUT -d '{"Max_Int8": "1"}' localhost:59882/api/v2/device/name/Random-Integer-Generator01/Max_Int8`

> `curl -X PUT -d '{"Min_Int8": "-100"}' localhost:59882/api/v2/device/name/Random-Integer-Generator01/Min_Int8`

> `curl -X GET localhost:59882/api/v2/device/name/Random-Integer-Generator01/RandomValue_Int8 | json_pp`

## Run kuiper in a terminal
1. Download
> `cd $GOPATH/src/github.com/rddigital/ && git clone https://github.com/rddigital/ekuiper.git`
2. Build & Run
> `cd $GOPATH/src/github.com/rddigital/ekuiper/ && make && cd`
> 
> `cd $GOPATH/src/github.com/rddigital/ekuiper/_build/*/bin/ && ./kuiperd`

## Run scenario service in a terminal
1. Build & Run
> `cd $GOPATH/src/github.com/rddigital/device-scenario && make && cd`

> `export EDGEX_SECURITY_SECRET_STORE=false && cd $GOPATH/src/github.com/rddigital/device-scenario/cmd/ && ./device-scenario`

## Setup config in a terminal
1. Get all rule
> `curl --location --request GET 'http://localhost:59990/api/v2/rule/all' | json_pp`
2. Add new rule and start
> `curl --location --request POST 'http://localhost:59990/api/v2/rule'
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
}' | json_pp`
3. Get rule
> `curl --location --request GET 'http://localhost:59990/api/v2/rule/name/auto1' | json_pp`
4. Start rule
> `curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
--header 'Content-Type: application/json'
--data-raw '{
    "adminState": "UNLOCKED"
}' | json_pp`
5. Stop rule
> `curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
--header 'Content-Type: application/json'
--data-raw '{
    "adminState": "LOCKED"
}' | json_pp`
6. Update actions rule
> `curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
--header 'Content-Type: application/json'
--data-raw '{
    "apiVersion":"v2",
    "rule: {
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
}' | json_pp`
7. Update "EnableNotify"
> `curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
--header 'Content-Type: application/json'
--data-raw '{
    "apiVersion": "v2",
    "rule":{
        "notifyEnable": "false"
    }
}' | json_pp`
8. Update conditions
> `curl --location --request PUT 'http://localhost:59990/api/v2/rule/name/auto1'
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
}' | json_pp`
9. Delete rule
> `curl --location --request DELETE 'http://localhost:59990/api/v2/rule/name/auto1'  | json_pp`

## Test rule

### Before update conditions: `0 <= value < 10`
Expected result: no trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 20, "name" : "dev1"}' -t reading`

Expected result: no trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 40, "name" : "dev2"}' -t reading`

Expected result: trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 40, "name" : "dev1"}' -t reading`

Expected result: no trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 20, "name" : "dev1"}' -t reading`

Expected result: trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 20, "name" : "dev2"}' -t reading`

Expected result: no trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 40, "name" : "dev2"}' -t reading`

### After update conditions: `10 <= value < 20`
Expected result: no trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 20, "name" : "dev1"}' -t reading`

Expected result: no trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 40, "name" : "dev2"}' -t reading`

Expected result: no trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 40, "name" : "dev1"}' -t reading`

Expected result: trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 20, "name" : "dev2"}' -t reading`

Expected result: no trigger

`mosquitto_pub -h broker.emqx.io -m '{"temperature": 20, "name" : "dev1"}' -t reading`