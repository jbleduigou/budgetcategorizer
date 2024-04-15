package categorizer

import (
	"log/slog"
	"strings"

	budget "github.com/jbleduigou/budgetcategorizer"
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
		if strings.Contains(strings.ToLower(t.Description), strings.ToLower(key)) {
			output.Category = value
			slog.Info("Assigning category to transaction",
				slog.String("category", value),
				slog.String("transaction-description", t.Description))
			return output
		}
	}
	slog.Warn("No matching categories found for transaction",
		slog.String("transaction-description", t.Description))
	return output
}
