package categorizer

import (
	"fmt"
	"testing"

	budget "github.com/jbleduigou/budgetcategorizer"
	"github.com/stretchr/testify/assert"
)

func TestCategorizeFound(t *testing.T) {
	tr := budget.NewTransaction("18/12/2019", "Paiement Par Carte Express Proxi Saint Thonan 18/12", "", "", 13.37)
	l := make(map[string]string)
	l["Express Proxi Saint Thonan"] = "Courses Alimentation"
	c := NewCategorizer(l)

	expected := "Courses Alimentation"
	output := c.Categorize(tr)

	fmt.Println(output)
	assert.Equal(t, expected, string(output.Category))
}

func TestCategorizeNotFound(t *testing.T) {
	tr := budget.NewTransaction("18/12/2019", "Lorem ipsum dolor sit amet", "", "", 13.37)
	l := make(map[string]string)
	c := NewCategorizer(l)

	expected := "???"
	output := c.Categorize(tr)

	fmt.Println(output)
	assert.Equal(t, expected, string(output.Category))
}
