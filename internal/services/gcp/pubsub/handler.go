package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hectorvent/cloudstack/internal/services/gcp/pubsub/model"
)

type PubSubHandler struct {
	service *PubSubService
}

func NewPubSubHandler(service *PubSubService) *PubSubHandler {
	return &PubSubHandler{
		service: service,
	}
}

func (h *PubSubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	ctx := context.Background()

	// Pattern: /v1/projects/{project}/topics
	if strings.Contains(path, "/topics") {
		if strings.HasSuffix(path, "/topics") {
			if r.Method == "GET" {
				h.handleListTopics(w, r, ctx)
				return
			}
		}
		
		if strings.Contains(path, ":publish") {
			topicName := strings.TrimSuffix(strings.TrimPrefix(path, "/v1/"), ":publish")
			h.handlePublish(w, r, ctx, topicName)
			return
		}

		// Create/Get/Delete Topic
		topicName := strings.TrimPrefix(path, "/v1/")
		if r.Method == "PUT" {
			h.handleCreateTopic(w, r, ctx, topicName)
			return
		}
		if r.Method == "GET" {
			h.handleGetTopic(w, r, ctx, topicName)
			return
		}
		if r.Method == "DELETE" {
			h.handleDeleteTopic(w, r, ctx, topicName)
			return
		}
	}

	// Pattern: /v1/projects/{project}/subscriptions
	if strings.Contains(path, "/subscriptions") {
		if strings.HasSuffix(path, "/subscriptions") {
			if r.Method == "GET" {
				h.handleListSubscriptions(w, r, ctx)
				return
			}
			if r.Method == "POST" {
				// GCP sometimes uses POST to root for create if name is in body
				h.handleCreateSubscription(w, r, ctx, "")
				return
			}
		}
		
		if strings.Contains(path, ":pull") {
			subName := strings.TrimSuffix(strings.TrimPrefix(path, "/v1/"), ":pull")
			h.handlePull(w, r, ctx, subName)
			return
		}
		if strings.Contains(path, ":acknowledge") {
			subName := strings.TrimSuffix(strings.TrimPrefix(path, "/v1/"), ":acknowledge")
			h.handleAcknowledge(w, r, ctx, subName)
			return
		}

		// Create/Get/Delete Subscription
		subName := strings.TrimPrefix(path, "/v1/")
		if r.Method == "PUT" {
			h.handleCreateSubscription(w, r, ctx, subName)
			return
		}
		if r.Method == "GET" {
			h.handleGetSubscription(w, r, ctx, subName)
			return
		}
		if r.Method == "DELETE" {
			h.handleDeleteSubscription(w, r, ctx, subName)
			return
		}
	}

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, `{"error": {"code": 501, "message": "PubSub path not implemented: %s %s"}}`, r.Method, path)
}

func (h *PubSubHandler) handleListTopics(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	topics, _ := h.service.ListTopics(ctx)
	res := model.TopicsList{
		Topics: topics,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *PubSubHandler) handleCreateTopic(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	topic, err := h.service.CreateTopic(ctx, name)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(topic)
}

func (h *PubSubHandler) handleGetTopic(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	// For now just return if it exists
	topics, _ := h.service.ListTopics(ctx)
	for _, t := range topics {
		if t.Name == name {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(t)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *PubSubHandler) handleDeleteTopic(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteTopic(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}

func (h *PubSubHandler) handleListSubscriptions(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	subs, _ := h.service.ListSubscriptions(ctx)
	res := model.SubscriptionsList{
		Subscriptions: subs,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *PubSubHandler) handleCreateSubscription(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	var req struct {
		Name  string `json:"name"`
		Topic string `json:"topic"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if name == "" { name = req.Name }
	
	sub, err := h.service.CreateSubscription(ctx, name, req.Topic)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *PubSubHandler) handleGetSubscription(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	subs, _ := h.service.ListSubscriptions(ctx)
	for _, s := range subs {
		if s.Name == name {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(s)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (h *PubSubHandler) handleDeleteSubscription(w http.ResponseWriter, r *http.Request, ctx context.Context, name string) {
	h.service.DeleteSubscription(ctx, name)
	w.WriteHeader(http.StatusNoContent)
}

func (h *PubSubHandler) handlePublish(w http.ResponseWriter, r *http.Request, ctx context.Context, topicName string) {
	var req struct {
		Messages []*model.Message `json:"messages"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	ids, err := h.service.Publish(ctx, topicName, req.Messages)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"messageIds": ids})
}

func (h *PubSubHandler) handlePull(w http.ResponseWriter, r *http.Request, ctx context.Context, subName string) {
	var req struct {
		MaxMessages int `json:"maxMessages"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.MaxMessages == 0 { req.MaxMessages = 1 }
	
	received, err := h.service.Pull(ctx, subName, req.MaxMessages)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"receivedMessages": received})
}

func (h *PubSubHandler) handleAcknowledge(w http.ResponseWriter, r *http.Request, ctx context.Context, subName string) {
	var req struct {
		AckIds []string `json:"ackIds"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	h.service.Acknowledge(ctx, subName, req.AckIds)
	w.WriteHeader(http.StatusNoContent)
}
