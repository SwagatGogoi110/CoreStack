package ce

type CostExplorerService struct{}

func NewCostExplorerService() *CostExplorerService {
	return &CostExplorerService{}
}

func (s *CostExplorerService) GetCostAndUsage() map[string]any {
	return map[string]any{
		"ResultsByTime": []any{},
		"GroupDefinitions": []any{},
		"DimensionValueAttributes": []any{},
	}
}

func (s *CostExplorerService) GetDimensionValues() map[string]any {
	return map[string]any{
		"DimensionValues": []any{},
		"ReturnSize": 0,
		"TotalSize": 0,
	}
}
