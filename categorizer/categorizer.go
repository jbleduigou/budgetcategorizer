package categorizer

import (
	"fmt"
	"strings"

	budget "github.com/jbleduigou/budgetcategorizer"
)

type Categorizer interface {
	Categorize(t budget.Transaction) budget.Transaction
}

func NewCategorizer(libelles map[string]string) Categorizer {
	return &categorizerImpl{libelles: libelles}
}

type categorizerImpl struct {
	libelles map[string]string
}

func (c *categorizerImpl) Categorize(t budget.Transaction) budget.Transaction {
	output := budget.NewTransaction(t.Date, t.Description, t.Comment, "???", t.Value)
	for key, value := range c.libelles {
		if strings.Contains(t.Description, key) {
			output.Category = value
			fmt.Printf("Assigning category '%s' to transaction with description '%s'\n", value, t.Description)
			return *output
		}
	}
	fmt.Printf("No matching categories found for transaction with description '%s'\n", t.Description)
	return *output
}
