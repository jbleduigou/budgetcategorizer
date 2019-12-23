package categorizer

import (
	"strings"

	budget "github.com/jbleduigou/budgetcategorizer"
)

type Categorizer interface {
	Categorize(t budget.Transaction) (budget.Transaction, error)
}

func NewCategorizer() Categorizer {
	l := make(map[string]string)
	// TODO make it parametable
	l["Express Proxi Saint Thonan"] = "Courses Alimentation"
	return &categorizerImpl{libelles: l}
}

type categorizerImpl struct {
	libelles map[string]string
}

func (c *categorizerImpl) Categorize(t budget.Transaction) (budget.Transaction, error) {
	output := budget.NewTransaction(t.Date, t.Description, t.Comment, "???", t.Value)
	for key, value := range c.libelles {
		if strings.Contains(t.Description, key) {
			output.Category = value
			return *output, nil
		}
	}
	return *output, nil
}
