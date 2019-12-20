package exporter

import (
	"fmt"
	"testing"

	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	transactions := []*budget.Transaction{}
	transactions = append(transactions, budget.NewTransaction("18/12/2019", "<libelle>", "<commentaire>", "<category>", 13.37))
	transactions = append(transactions, budget.NewTransaction("18/12/2019", "<libelle>", "<commentaire>", "<category>", -19.84))

	e := NewExporter()

	result, _ := e.Export(transactions)
	expected :=
		`18/12/2019,<libelle>,<commentaire>,<category>,13.37
18/12/2019,<libelle>,<commentaire>,<category>,-19.84
`
	assert.Equal(t, expected, string(result))
}

type mockWriterWithError struct {
}

func (m *mockWriterWithError) Write(record []string) error {
	return fmt.Errorf("Error for unit tests")
}

func TestExportWithError(t *testing.T) {
	transactions := []*budget.Transaction{}
	transactions = append(transactions, budget.NewTransaction("18/12/2019", "<libelle>", "<commentaire>", "<category>", 13.37))
	transactions = append(transactions, budget.NewTransaction("18/12/2019", "<libelle>", "<commentaire>", "<category>", -19.84))

	e := &csvExporter{}

	err := e.exportToCSV(transactions, &mockWriterWithError{})

	assert.Equal(t, err.Error(), "Error for unit tests")
}
