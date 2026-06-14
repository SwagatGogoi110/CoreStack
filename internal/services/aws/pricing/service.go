package pricing

type PricingService struct{}

func NewPricingService() *PricingService {
	return &PricingService{}
}

func (s *PricingService) DescribeServices() map[string]any {
	return map[string]any{
		"FormatVersion": "aws_v1",
		"Services":      []any{},
	}
}

func (s *PricingService) GetProducts() map[string]any {
	return map[string]any{
		"FormatVersion": "aws_v1",
		"PriceList":     []string{},
	}
}
