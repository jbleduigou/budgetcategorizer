package mock

import (
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/stretchr/testify/mock"
)

func NewBroker() *MockBroker {
	return &MockBroker{}
}

type MockBroker struct {
	mock.Mock
}

func (_m *MockBroker) Send(_a0 budget.Transaction) error {
	ret := _m.Called(_a0)
	if ret.Get(0) == nil {
		return nil
	}
	return ret.Get(0).(error)
}
