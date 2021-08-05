// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018-2021 IOTech Ltd
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/edgexfoundry/device-sdk-go/v2/pkg/startup"

	device_scenario "github.com/rddigital/device-scenario"
	"github.com/rddigital/device-scenario/driver"
)

const (
	serviceName string = "scenario"
)

func main() {
	d := driver.NewProtocolDriver()
	startup.Bootstrap(serviceName, device_scenario.Version, d)
}
