package models

import "strings"

type Action struct {
	DeviceName  string `json:"deviceName" validate:"required"`
	CommandName string `json:"commandName" validate:"required"`
	Body        string `json:"body" validate:"required"`
}

func ActionsToProperties(actions []Action) map[string]string {
	properties := make(map[string]string, len(actions))
	for _, action := range actions {
		key := action.DeviceName + "/" + action.CommandName
		properties[key] = action.Body
	}
	return properties
}

func ActionsFromProperties(properties map[string]string) []Action {
	actions := make([]Action, 0, len(properties))
	for key, value := range properties {
		arrStr := strings.Split(key, "/")
		if len(arrStr) < 2 {
			continue
		}
		var action = Action{
			DeviceName:  arrStr[0],
			CommandName: arrStr[1],
			Body:        value,
		}
		actions = append(actions, action)
	}

	return actions
}
