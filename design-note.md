```
[Device.Id]
[Device.Name]
[Device.Descriptor]
[Device.AdminState]
[Device.Protocols]
    [actions]
        "{DeviceName}/{CommandName}" : "{JSON Body}"
    [notify]
        "notify": "{true/false}"
    [conditions]
        "{index}": "{
            'Logic':'{and/or}',
            'Type':'{schedule/threshold}',
            'StartTime': '',
            'EndTime': '',
            'DeviceThreshold': '',
            'OperatorThreshold': '',
            'ResourceThreshold': '',
            'ValueThreshold': ''
        }"
```

> In Schedule, Kuiper service: `Rule.Id = Interval.Name = IntervalAction.Name = "_" + Rule.Id + "_" + "{index}"`

> The body of Kuiper action and IntervalAction: `{"triggerIndex":"{index}", "TriggerState":"true/false"}`