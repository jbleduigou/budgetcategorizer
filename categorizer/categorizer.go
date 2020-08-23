package categorizer

import (
	"fmt"
	"strings"

	budget "github.com/jbleduigou/budgetcategorizer"
)

// Categorizer provides and interface for assigning a category to a transaction
type Categorizer interface {
	Categorize(t budget.Transaction) budget.Transaction
}

// NewCategorizer will provide an instance of a Categorizer, implementation is not exposed
func NewCategorizer(libelles map[string]string, awsRequestID string) Categorizer {
	return &categorizerImpl{libelles: libelles, requestID: awsRequestID}
}

type categorizerImpl struct {
	libelles  map[string]string
	requestID string
}

func (c *categorizerImpl) Categorize(t budget.Transaction) budget.Transaction {
	output := budget.NewTransaction(t.Date, t.Description, t.Comment, "???", t.Value)
	for key, value := range c.libelles {
		if strings.Contains(t.Description, key) {
			output.Category = value
			fmt.Printf("[INFO] %v Assigning category '%s' to transaction with description '%s'\n", c.requestID, value, t.Description)
			return output
		}
	}
	fmt.Printf("[WARN] %v No matching categories found for transaction with description '%s'\n", c.requestID, t.Description)
	return output
}
