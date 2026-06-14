package model

type Block struct {
	BlockType     string          `json:"BlockType"`
	Confidence    float64         `json:"Confidence"`
	Text          string          `json:"Text,omitempty"`
	ID            string          `json:"Id"`
	Relationships []*Relationship `json:"Relationships,omitempty"`
	Geometry      *Geometry       `json:"Geometry"`
	Page          int             `json:"Page,omitempty"`
}

type Relationship struct {
	Type string   `json:"Type"`
	Ids  []string `json:"Ids"`
}

type Geometry struct {
	BoundingBox *BoundingBox `json:"BoundingBox"`
	Polygon     []*Point     `json:"Polygon"`
}

type BoundingBox struct {
	Width  float64 `json:"Width"`
	Height float64 `json:"Height"`
	Left   float64 `json:"Left"`
	Top    float64 `json:"Top"`
}

type Point struct {
	X float64 `json:"X"`
	Y float64 `json:"Y"`
}
