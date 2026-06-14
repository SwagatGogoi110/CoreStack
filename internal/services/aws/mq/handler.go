package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hectorvent/cloudstack/internal/core/common"
)

type MqHandler struct {
	service *MqService
}

func NewMqHandler(service *MqService) *MqHandler {
	return &MqHandler{
		service: service,
	}
}

func (h *MqHandler) HandleJSON(action string, request json.RawMessage, rc *common.RequestContext) (any, error) {
	return nil, fmt.Errorf("MQ does not support standard JSON dispatcher")
}

func (h *MqHandler) HandleQuery(action string, params url.Values, rc *common.RequestContext) (string, error) {
	return "", fmt.Errorf("MQ does not support Query protocol")
}

func (h *MqHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v1")
	ctx := context.Background()

	if path == "/brokers" && r.Method == "POST" {
		var req struct {
			BrokerName     string `json:"brokerName"`
			EngineType     string `json:"engineType"`
			DeploymentMode string `json:"deploymentMode"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		broker, _ := h.service.CreateBroker(ctx, req.BrokerName, req.EngineType, req.DeploymentMode)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"brokerId": broker.BrokerId, "brokerArn": broker.BrokerArn})
		return
	}

	if path == "/brokers" && r.Method == "GET" {
		brokers, _ := h.service.ListBrokers(ctx)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"brokerSummaries": brokers})
		return
	}

	if strings.HasPrefix(path, "/brokers/") {
		brokerId := strings.TrimPrefix(path, "/brokers/")
		if r.Method == "GET" {
			broker, err := h.service.DescribeBroker(ctx, brokerId)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(broker)
			return
		}
		if r.Method == "DELETE" {
			h.service.DeleteBroker(ctx, brokerId)
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]any{})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
