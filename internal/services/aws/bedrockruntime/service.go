package bedrockruntime

import (
	"fmt"
	"strings"
)

type BedrockRuntimeService struct{}

func NewBedrockRuntimeService() *BedrockRuntimeService {
	return &BedrockRuntimeService{}
}

func (s *BedrockRuntimeService) BuildConverseResponse(modelID string) map[string]any {
	return map[string]any{
		"output": map[string]any{
			"message": map[string]any{
				"role": "assistant",
				"content": []map[string]any{
					{"text": "CloudStack stub response for model=" + modelID},
				},
			},
		},
		"stopReason": "end_turn",
		"usage": map[string]any{
			"inputTokens":  10,
			"outputTokens": 12,
			"totalTokens":  22,
		},
		"metrics": map[string]any{
			"latencyMs": 1,
		},
	}
}

func (s *BedrockRuntimeService) BuildInvokeModelResponse(modelID string) []byte {
	lower := strings.ToLower(modelID)
	if strings.Contains(lower, "anthropic") {
		return []byte(fmt.Sprintf(`{"id":"msg_stub","type":"message","role":"assistant","content":[{"type":"text","text":"CloudStack stub response"}],"model":"%s","stop_reason":"end_turn","usage":{"input_tokens":10,"output_tokens":12}}`, modelID))
	}
	return []byte(`{"outputs":[{"text":"CloudStack stub response"}]}`)
}
