package dispatch

import "context"

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Dispatch(ctx context.Context) error {
	return nil
}
