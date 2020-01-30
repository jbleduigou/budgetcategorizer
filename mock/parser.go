package mock

import (
	"io"

	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/stretchr/testify/mock"
)

// NewParser provides a mock instance of a Parser
func NewParser() *Parser {
	return &Parser{}
}

// Parser is an implementation of the Parser interface with a mock, use for testing not for production
type Parser struct {
	mock.Mock
}

// ParseTransactions parses an array of Transaction from a reader
func (_m *Parser) ParseTransactions(_a0 io.Reader) ([]budget.Transaction, error) {
	ret := _m.Called(_a0)
	if ret.Get(1) != nil {
		return ret.Get(0).([]budget.Transaction), ret.Get(1).(error)
	}
	return ret.Get(0).([]budget.Transaction), nil
}
