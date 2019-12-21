package mock

import (
	"io"

	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/stretchr/testify/mock"
)

func NewParser() *MockParser {
	return &MockParser{}
}

type MockParser struct {
	mock.Mock
}

func (_m *MockParser) ParseTransactions(_a0 io.Reader) ([]*budget.Transaction, error) {
	ret := _m.Called(_a0)
	if ret.Get(1) != nil {
		return ret.Get(0).([]*budget.Transaction), ret.Get(1).(error)
	}
	return ret.Get(0).([]*budget.Transaction), nil
}
