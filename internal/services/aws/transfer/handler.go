package transfer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type TransferJsonHandler struct {
	service *TransferService
}

func NewTransferJsonHandler(service *TransferService) *TransferJsonHandler {
	return &TransferJsonHandler{
		service: service,
	}
}

func (h *TransferJsonHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	ctx := context.Background()

	switch action {
	case "CreateServer":
		var req struct {
			Protocols []string `json:"Protocols"`
		}
		json.Unmarshal(request, &req)
		server, err := h.service.CreateServer(ctx, req.Protocols)
		if err != nil {
			return nil, err
		}
		return map[string]string{"ServerId": server.ServerId}, nil

	case "ListServers":
		servers, _ := h.service.ListServers(ctx)
		return map[string]any{"Servers": servers}, nil

	default:
		return nil, fmt.Errorf("UnknownOperationException: Operation %s is not supported", action)
	}
}

func (h *TransferJsonHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("Transfer does not support Query protocol")
}
