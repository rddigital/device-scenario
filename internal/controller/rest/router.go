package rest

import (
	"encoding/json"
	"net/http"

	"github.com/edgexfoundry/device-sdk-go/v2/pkg/service"
	contractsCommon "github.com/edgexfoundry/go-mod-core-contracts/v2/common"
	commonDTO "github.com/edgexfoundry/go-mod-core-contracts/v2/dtos/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/errors"

	"github.com/rddigital/device-scenario/internal/common"
)

func InitRuleServer() {
	ds := service.RunningService()
	ds.AddRoute(common.ApiRuleRoute, AddRuleHander, http.MethodPost)
	ds.AddRoute(common.ApiAllRuleRoute, GetAllRuleHander, http.MethodGet)
	ds.AddRoute(common.ApiRuleByNameRoute, GetRuleByNameHander, http.MethodGet)
	ds.AddRoute(common.ApiRuleByNameRoute, UpdateRuleByNameHander, http.MethodPut)
	ds.AddRoute(common.ApiRuleByNameRoute, DeleteRuleByNameHander, http.MethodDelete)
	ds.AddRoute(common.ApiRuleTriggerByIdRoute, TriggerRuleByIdHander, http.MethodPost)
}

// SendResponse puts together the response packet for the V2 API
func SendResponse(
	writer http.ResponseWriter,
	request *http.Request,
	response interface{},
	statusCode int) {

	ds := service.RunningService()
	lc := ds.LoggingClient

	correlationID := request.Header.Get(contractsCommon.CorrelationHeader)

	writer.Header().Set(contractsCommon.CorrelationHeader, correlationID)
	writer.Header().Set(contractsCommon.ContentType, contractsCommon.ContentTypeJSON)
	writer.WriteHeader(statusCode)

	if response != nil {
		data, err := json.Marshal(response)
		if err != nil {
			lc.Error("Unable to marshal response", "error", err.Error(), contractsCommon.CorrelationHeader, correlationID)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = writer.Write(data)
		if err != nil {
			lc.Error("Unable to write %s response", "error", err.Error(), contractsCommon.CorrelationHeader, correlationID)
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func SendEdgexError(
	writer http.ResponseWriter,
	request *http.Request,
	err errors.EdgeX) {

	ds := service.RunningService()
	lc := ds.LoggingClient
	correlationID := request.Header.Get(contractsCommon.CorrelationHeader)

	lc.Error(err.Error(), contractsCommon.CorrelationHeader, correlationID)
	lc.Debug(err.DebugMessages(), contractsCommon.CorrelationHeader, correlationID)
	response := commonDTO.NewBaseResponse("", err.Error(), err.Code())
	SendResponse(writer, request, response, err.Code())
}
