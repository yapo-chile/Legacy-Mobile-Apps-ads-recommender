package infrastructure

import (
	"time"

	"github.com/sony/gobreaker"
	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/interfaces/loggers"
	"gitlab.com/yapo_team/legacy/mobile-apps/ads-recommender/pkg/interfaces/repository"
)

var (
	// ErrTooManyRequests is returned when the CB state is half open and the requests count is over the cb maxRequests
	ErrTooManyRequests = gobreaker.ErrTooManyRequests

	// ErrOpenState is returned when the CB state is open
	ErrOpenState = gobreaker.ErrOpenState
)

// HTTPCircuitBreakerHandler struct to implements http repository operations with circuit breaker
type HTTPCircuitBreakerHandler struct {
	circuitBreaker CircuitBreaker
	logger         loggers.Logger
	httpHandler    repository.HTTPHandler
}

// NewHTTPCircuitBreakerHandler will create a new instance of a custom http request handler
func NewHTTPCircuitBreakerHandler(
	circuitBreaker CircuitBreaker,
	logger loggers.Logger,
	h repository.HTTPHandler,
) repository.HTTPHandler {
	return &HTTPCircuitBreakerHandler{
		circuitBreaker: circuitBreaker,
		logger:         logger,
		httpHandler:    h,
	}
}

// Send will execute the sending of a http request
// but in this case it will wait until it obtains a successful response
// in order to continue it's execution
func (h *HTTPCircuitBreakerHandler) Send(req repository.HTTPRequest) (interface{}, error) {
	h.logger.Debug("HTTP - %s - Sending HTTP with circuit breaker request to: %+v", req.GetMethod(), req.GetPath())

	var response interface{}
	var err error
	// do-while: try once or retry until circuit breaker closes
	for ok := true; ok; ok = (err == ErrOpenState || err == ErrTooManyRequests) {
		response, err = h.circuitBreaker.Execute(func() (interface{}, error) {
			return h.httpHandler.Send(req)
		})
	}

	return response, err
}

// NewRequest returns an initialized struct that can be used to make a http request
func (h *HTTPCircuitBreakerHandler) NewRequest() repository.HTTPRequest {
	return h.httpHandler.NewRequest()
}

// NewCircuitBreaker initializes circuit breaker wrapper
// name is the circuit breaker
// consecutiveFailures is the maximum of consecutive errors allowed before open state
// failureRatioTolerance is the maximum error ratio (errors vs requests qty) allowed before open state
// Interval is the cyclic period of the closed state for the CircuitBreaker to clear the internal Counts.
// If Interval is 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
// Timeout is the period of the open state, after which the state of the CircuitBreaker becomes half-open.
func NewCircuitBreaker(
	name string,
	consecutiveFailures uint32,
	failureRatioTolerance float64,
	timeout,
	interval int,
	logger loggers.Logger,
) CircuitBreaker {
	settings := gobreaker.Settings{
		Name:     name,
		Timeout:  time.Duration(timeout) * (time.Second),
		Interval: time.Duration(interval) * (time.Second),

		// If ReadyToTrip returns true, the CircuitBreaker will be placed into the open state
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			errorRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return errorRatio >= failureRatioTolerance || counts.ConsecutiveFailures > consecutiveFailures
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			logger.Error("CircuitBreaker: Changing status %+v to %+v", from.String(), to.String())
			if from == gobreaker.StateOpen { // represents Circuit breaker opened state
				logger.Error("CircuitBreaker: Waiting for closed state...")
			}
		},
	}

	return gobreaker.NewCircuitBreaker(settings)
}

// CircuitBreaker allows circuit breaker operations
type CircuitBreaker interface {
	// Execute wrapps a function. If the function returns too many errors, circuit breaker
	// will return "circuit breaker open" error
	Execute(req func() (interface{}, error)) (interface{}, error)
	// Name returns circuit breaker name
	Name() string
}
