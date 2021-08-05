package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	contractsCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	commonDTO "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"
	"github.com/gorilla/mux"

	"github.com/rddigital/device-scenario/internal/application"
	"github.com/rddigital/device-scenario/internal/common"
	"github.com/rddigital/device-scenario/internal/models"
)

func AddRuleHander(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer func() { _ = r.Body.Close() }()
	}

	var addRuleRequest models.RuleRequest
	err := json.NewDecoder(r.Body).Decode(&addRuleRequest)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to decode JSON", err)
		SendEdgexError(w, r, edgexErr)
		return
	}

	err = common.Validate(addRuleRequest)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to validation", err)
		SendEdgexError(w, r, edgexErr)
		return
	}

	if addRuleRequest.Rule.Name == "" {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "rule name empty", err)
		SendEdgexError(w, r, edgexErr)
		return
	}

	if enable, _ := strconv.ParseBool(addRuleRequest.Rule.NotifyEnable); (len(addRuleRequest.Rule.Actions) == 0) && (!enable) {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "no action rule", err)
		SendEdgexError(w, r, edgexErr)
		return
	}

	if len(addRuleRequest.Rule.Conditions) == 0 {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "rule condions empty", err)
		SendEdgexError(w, r, edgexErr)
		return
	}

	edgexErr := application.AddRule(addRuleRequest.Rule)
	if edgexErr == nil {
		correlationID := r.Header.Get(contractsCommon.CorrelationHeader)
		response := commonDTO.NewBaseResponse(correlationID, "", http.StatusOK)
		SendResponse(w, r, response, http.StatusOK)
	} else {
		SendEdgexError(w, r, edgexErr)
	}
}

func GetAllRuleHander(w http.ResponseWriter, r *http.Request) {
	rulesResponse := application.GetAllRule()
	correlationID := r.Header.Get(contractsCommon.CorrelationHeader)
	response := models.NewMultiRulesResponse(correlationID, "", http.StatusOK, rulesResponse)
	SendResponse(w, r, response, http.StatusOK)
}

func GetRuleByNameHander(w http.ResponseWriter, r *http.Request) {
	// URL parameters
	vars := mux.Vars(r)
	name := vars[contractsCommon.Name]

	ruleResponse, edgexErr := application.GetRuleByName(name)
	if edgexErr == nil {
		correlationID := r.Header.Get(contractsCommon.CorrelationHeader)
		response := models.NewRuleResponse(correlationID, "", http.StatusOK, ruleResponse)
		SendResponse(w, r, response, http.StatusOK)
	} else {
		SendEdgexError(w, r, edgexErr)
	}
}

func UpdateRuleByNameHander(w http.ResponseWriter, r *http.Request) {
	// URL parameters
	vars := mux.Vars(r)
	name := vars[contractsCommon.Name]

	if r.Body != nil {
		defer func() { _ = r.Body.Close() }()
	}

	var updateRuleRequest models.RuleRequest
	err := json.NewDecoder(r.Body).Decode(&updateRuleRequest)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to decode JSON", err)
		SendEdgexError(w, r, edgexErr)
		return
	}

	err = common.Validate(updateRuleRequest)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to validation", err)
		SendEdgexError(w, r, edgexErr)
		return
	}

	edgexErr := application.UpdateRuleByName(name, updateRuleRequest.Rule)
	if edgexErr == nil {
		correlationID := r.Header.Get(contractsCommon.CorrelationHeader)
		response := commonDTO.NewBaseResponse(correlationID, "", http.StatusOK)
		SendResponse(w, r, response, http.StatusOK)
	} else {
		SendEdgexError(w, r, edgexErr)
	}
}

func DeleteRuleByNameHander(w http.ResponseWriter, r *http.Request) {
	// URL parameters
	vars := mux.Vars(r)
	name := vars[contractsCommon.Name]

	edgexErr := application.DeleteRuleByName(name)
	if edgexErr == nil {
		correlationID := r.Header.Get(contractsCommon.CorrelationHeader)
		response := commonDTO.NewBaseResponse(correlationID, "", http.StatusOK)
		SendResponse(w, r, response, http.StatusOK)
	} else {
		SendEdgexError(w, r, edgexErr)
	}
}

func TriggerRuleByIdHander(w http.ResponseWriter, r *http.Request) {
	// URL parameters
	vars := mux.Vars(r)
	id := vars[contractsCommon.Id]

	if r.Body != nil {
		defer func() { _ = r.Body.Close() }()
	}

	var triggerContent models.ContentTrigger
	err := json.NewDecoder(r.Body).Decode(&triggerContent)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to decode JSON", err)
		SendEdgexError(w, r, edgexErr)
		return
	}

	err = common.Validate(triggerContent)
	if err != nil {
		edgexErr := errors.NewCommonEdgeX(errors.KindServerError, "failed to validation", err)
		SendEdgexError(w, r, edgexErr)
		return
	}

	go application.TriggerRuleById(id, triggerContent)

	correlationID := r.Header.Get(contractsCommon.CorrelationHeader)
	response := commonDTO.NewBaseResponse(correlationID, "", http.StatusAccepted)
	SendResponse(w, r, response, http.StatusAccepted)
}
