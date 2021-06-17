// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

// This package provides a implementation of a ProtocolDriver interface.
//
package driver

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	dsModels "github.com/edgexfoundry/device-sdk-go/v2/pkg/models"
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/http"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/interfaces"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/models"

	"github.com/rddigital/device-scenario/internal/config"
)

var once sync.Once
var driver *ScenarioDriver

type ScenarioDriver struct {
	lc            logger.LoggingClient
	serviceConfig *config.ServiceConfig
	commandClient interfaces.CommandClient
}

type ScenarioContent struct {
	DeviceName  string
	CommandName string
	BodyMap     map[string]string
}

func NewProtocolDriver() dsModels.ProtocolDriver {
	once.Do(func() {
		driver = new(ScenarioDriver)
	})
	return driver
}

func (d *ScenarioDriver) DisconnectDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.lc.Infof("ScenarioDriver.DisconnectDevice: device-scenario driver is disconnecting to %s", deviceName)
	return nil
}

func (d *ScenarioDriver) Initialize(lc logger.LoggingClient, asyncCh chan<- *dsModels.AsyncValues, deviceCh chan<- []dsModels.DiscoveredDevice) error {
	d.lc = lc
	d.serviceConfig = &config.ServiceConfig{}

	ds := service.RunningService()

	if err := ds.LoadCustomConfig(d.serviceConfig, ServiceCustomConfigName); err != nil {
		return fmt.Errorf("unable to load '%s' custom configuration: %s", ServiceCustomConfigName, err.Error())
	}

	if err := d.serviceConfig.ServiceCustomConfig.Validate(); err != nil {
		return fmt.Errorf("'%s' custom configuration validation failed: %s", ServiceCustomConfigName, err.Error())
	}

	urlCoreCommand := d.serviceConfig.ServiceCustomConfig.CommandClientInfo.Url()
	d.commandClient = http.NewCommandClient(urlCoreCommand)

	return nil
}

func (d *ScenarioDriver) HandleReadCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest) (res []*dsModels.CommandValue, err error) {
	return nil, fmt.Errorf("ScenarioDriver.HandleReadCommands; read commands not supported")
}

func (d *ScenarioDriver) HandleWriteCommands(deviceName string, protocols map[string]models.ProtocolProperties, reqs []dsModels.CommandRequest,
	params []*dsModels.CommandValue) error {

	for _, param := range params {
		switch param.DeviceResourceName {
		case "TriggerScenario":
			content, ok := protocols[ContentPropertyName]
			if !ok {
				d.lc.Debugf("No content in Scenario: %s", deviceName)
				return nil
			}
			arrAction := parseContent(content)
			if arrAction == nil {
				d.lc.Debugf("No content in Scenario: %s", deviceName)
				return nil
			}

			ctx := context.Background()
			arrError := make([]string, 0)
			for _, action := range arrAction {
				_, errAction := d.commandClient.IssueSetCommandByName(ctx, action.DeviceName, action.CommandName, action.BodyMap)
				if errAction != nil {
					arrError = append(arrError, errAction.Message())
					d.lc.Debugf("Send command '%s' to device '%s' failed", action.CommandName, action.DeviceName)
				} else {
					d.lc.Debugf("Send command '%s' to device '%s' successed", action.CommandName, action.DeviceName)
				}
			}

			if len(arrError) > 0 {
				errStr := strings.Join(arrError, ";")
				return fmt.Errorf("ScenarioDriver.HandleWriteCommands: Some actions errored: %s", errStr)
			}
		default:
			return fmt.Errorf("ScenarioDriver.HandleWriteCommands: there is no matched device resource for %s", param.String())
		}
	}

	return nil
}

func (d *ScenarioDriver) Stop(force bool) error {
	d.lc.Info("ScenarioDriver.Stop: device-scenario driver is stopping...")
	return nil
}

func (d *ScenarioDriver) AddDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.lc.Debugf("a new Device is added: %s", deviceName)
	return nil
}

func (d *ScenarioDriver) UpdateDevice(deviceName string, protocols map[string]models.ProtocolProperties, adminState models.AdminState) error {
	d.lc.Debugf("Device %s is updated", deviceName)
	return nil
}

func (d *ScenarioDriver) RemoveDevice(deviceName string, protocols map[string]models.ProtocolProperties) error {
	d.lc.Debugf("Device %s is removed", deviceName)
	return nil
}

func parseContent(content map[string]string) []ScenarioContent {
	if content == nil {
		return nil
	}
	if len(content) == 0 {
		return nil
	}

	arrContent := make([]ScenarioContent, 0, len(content))
	for nameParam, bodyParam := range content {
		deviceName, commandName, ok := parseName(nameParam)
		if !ok {
			continue
		}
		bodyMap, err := parseBody(bodyParam)
		if err != nil {
			continue
		}

		contentElement := ScenarioContent{
			DeviceName:  deviceName,
			CommandName: commandName,
			BodyMap:     bodyMap,
		}
		arrContent = append(arrContent, contentElement)
	}

	return arrContent
}

func parseName(params string) (deviceName string, commandName string, ok bool) {
	ok = false
	arrStr := strings.Split(params, "/")
	if len(arrStr) < 2 {
		return
	}
	deviceName = arrStr[0]
	commandName = arrStr[1]
	ok = true
	return
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
