package models

import (
	commonDTO "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
)

type RuleRequest struct {
	commonDTO.BaseRequest `json:",inline"`
	Rule                  Rule `json:"rule"`
}
