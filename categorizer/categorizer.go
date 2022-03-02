package categorizer

import (
	"strings"

	budget "github.com/jbleduigou/budgetcategorizer"
	"go.uber.org/zap"
)

// Categorizer provides and interface for assigning a category to a transaction
type Categorizer interface {
	Categorize(t budget.Transaction) budget.Transaction
}

// NewCategorizer will provide an instance of a Categorizer, implementation is not exposed
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
			zap.L().Info("Assigning category to transaction",
				zap.String("category", value),
				zap.String("transaction-description", t.Description))
			return output
		}
	}
	zap.L().Warn("No matching categories found for transaction",
		zap.String("transaction-description", t.Description))
	return output
}
