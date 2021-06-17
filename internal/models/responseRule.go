package models

import (
	commonDTO "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
)

type RuleResponse struct {
	commonDTO.BaseResponse `json:",inline"`
	Rule                   Rule `json:"rule"`
}

func NewRuleResponse(requestId string, message string, statusCode int, rule Rule) RuleResponse {
	return RuleResponse{
		BaseResponse: commonDTO.NewBaseResponse(requestId, message, statusCode),
		Rule:         rule,
	}
}

type MultiRulesResponse struct {
	commonDTO.BaseResponse `json:",inline"`
	Rules                  []Rule `json:"rules"`
}

func NewMultiRulesResponse(requestId string, message string, statusCode int, rules []Rule) MultiRulesResponse {
	return MultiRulesResponse{
		BaseResponse: commonDTO.NewBaseResponse(requestId, message, statusCode),
		Rules:        rules,
	}
}
