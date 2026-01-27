package ws

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/q1bksuu/onebot-go-sdk/v11/dispatcher"
	"github.com/q1bksuu/onebot-go-sdk/v11/entity"
)

// HandleActionMessage parses an action request and returns the standardized response envelope.
func HandleActionMessage(
	ctx context.Context,
	data []byte,
	handler dispatcher.ActionRequestHandler,
	badRequestErr error,
) *entity.ActionResponseEnvelope {
	var reqEnv entity.ActionRequestEnvelope

	err := json.Unmarshal(data, &reqEnv)
	if err != nil {
		return invalidJSONResponse()
	}

	req := &entity.ActionRequest{Action: reqEnv.Action, Params: reqEnv.Params}

	resp, err := handler.HandleActionRequest(ctx, req)
	if err != nil {
		mapped := mapHandlerError(err, badRequestErr)

		return &entity.ActionResponseEnvelope{
			ActionRawResponse: *mapped,
			Echo:              reqEnv.Echo,
		}
	}

	if resp == nil {
		resp = &entity.ActionRawResponse{
			Status:  entity.StatusFailed,
			Retcode: -1,
			Message: "empty response",
		}
	}

	return &entity.ActionResponseEnvelope{
		ActionRawResponse: *resp,
		Echo:              reqEnv.Echo,
	}
}

func invalidJSONResponse() *entity.ActionResponseEnvelope {
	return &entity.ActionResponseEnvelope{
		ActionRawResponse: entity.ActionRawResponse{
			Status:  entity.StatusFailed,
			Retcode: entity.ActionResponseRetcode(1400),
			Message: "invalid json",
		},
	}
}

func mapHandlerError(err error, badRequestErr error) *entity.ActionRawResponse {
	switch {
	case errors.Is(err, dispatcher.ErrActionNotFound):
		return &entity.ActionRawResponse{
			Status:  entity.StatusFailed,
			Retcode: 1404,
			Message: err.Error(),
		}
	case badRequestErr != nil && errors.Is(err, badRequestErr):
		return &entity.ActionRawResponse{
			Status:  entity.StatusFailed,
			Retcode: 1400,
			Message: err.Error(),
		}
	default:
		return &entity.ActionRawResponse{
			Status:  entity.StatusFailed,
			Retcode: 1500,
			Message: err.Error(),
		}
	}
}
