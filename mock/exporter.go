package mock

import (
	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/stretchr/testify/mock"
)

func NewExporter() *MockExporter {
	return &MockExporter{}
}

type MockExporter struct {
	mock.Mock
}

func (_m *MockExporter) Export(_a0 []*budget.Transaction) ([]byte, error) {
	ret := _m.Called(_a0)
	if ret.Get(1) != nil {
		return ret.Get(0).([]byte), ret.Get(1).(error)
	}
	return ret.Get(0).([]byte), nil
}
