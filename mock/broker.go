package mock

import (
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/stretchr/testify/mock"
)

// NewBroker provides a mock instance of a Broker
func NewBroker() *Broker {
	return &Broker{}
}

// Broker is an implementation of the Broker interface with a mock, use for testing not for production
type Broker struct {
	mock.Mock
}

// Send is the method for sending a message to the broker
func (_m *Broker) Send(_a0 budget.Transaction) error {
	ret := _m.Called(_a0)
	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(error)
}
