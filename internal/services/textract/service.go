package textract

import (
	"github.com/google/uuid"
	"github.com/hectorvent/cloudstack/internal/services/textract/model"
)

type TextractService struct{}

func NewTextractService() *TextractService {
	return &TextractService{}
}

func (s *TextractService) DetectDocumentText() map[string]any {
	wordID := uuid.New().String()
	return map[string]any{
		"DocumentMetadata": map[string]any{"Pages": 1},
		"Blocks": []model.Block{
			{
				BlockType:  "WORD",
				Confidence: 99.9,
				Text:       "CloudStack",
				ID:         wordID,
				Geometry: &model.Geometry{
					BoundingBox: &model.BoundingBox{Width: 0.1, Height: 0.05, Left: 0.1, Top: 0.1},
					Polygon: []*model.Point{
						{X: 0.1, Y: 0.1}, {X: 0.2, Y: 0.1}, {X: 0.2, Y: 0.15}, {X: 0.1, Y: 0.15},
					},
				},
				Page: 1,
			},
		},
	}
}
