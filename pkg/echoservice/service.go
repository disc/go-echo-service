package echoservice

import "errors"

type EchoService interface {
	Echo(string) (string, error)
}

type echoService struct{}

func NewEchoService() EchoService {
	return echoService{}
}

func (echoService) Echo(s string) (string, error) {
	if s == "" {
		return "", ErrEmpty
	}
	return s, nil
}

var (
	// ErrEmpty is returned when provided string si empty
	ErrEmpty = errors.New("Empty string")
)

// ServiceMiddleware is a chain-able behavior modifier for EchoService.
type ServiceMiddleware func(EchoService) EchoService
