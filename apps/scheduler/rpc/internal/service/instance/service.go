package instance

import "context"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ReportResult(ctx context.Context) error {
	return nil
}
