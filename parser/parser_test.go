package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTransactions(t *testing.T) {
	p := NewParser()
	transactions, _ := p.ParseTransactions(bytes.NewBufferString(content))
	assert.Equal(t, 2, len(transactions))
	debit := transactions[0]
	assert.Equal(t, debit.Date, "18/12/2019")
	assert.Equal(t, debit.Description, "Dans la vie on ne fait pas ce que l'on veut mais on est responsable de ce que l'on est. ")
	assert.Equal(t, debit.Comment, "")
	assert.Equal(t, debit.Category, "Courses Alimentation")
	assert.Equal(t, debit.Value, 3.18)
	credit := transactions[1]
	assert.Equal(t, credit.Date, "18/12/2019")
	assert.Equal(t, credit.Description, "3446335 Remise De Cheque Ref: 3446335 ")
	assert.Equal(t, credit.Comment, "")
	assert.Equal(t, credit.Category, "")
	assert.Equal(t, credit.Value, -30.13)
}

type mockReaderWithError struct {
}

func (m *mockReaderWithError) ReadAll() (records [][]string, err error) {
	return nil, fmt.Errorf("Error for unit tests")
}

func TestParseTransactionsWithError(t *testing.T) {
	p := csvParser{}
	transactions, err := p.parse(&mockReaderWithError{})
	assert.Nil(t, transactions)
	assert.Equal(t, err.Error(), "Error for unit tests")
}

var content = `
Téléchargement du  19/12/2019;

  
M. LE DUIGOU JEAN BAPTISTE    
CCHQ       no 62734091867;
Solde au 18/12/19 : 13,37 EUR

Liste des opérations du compte entre le 18/12/2019 et le 18/12/2019;

Date;Libellé;Débit Euros;Crédit Euros;
18/12/2019;"Dans la vie on ne fait pas 
ce que l'on veut 
mais on est responsable 
de ce que l'on est.
";3,18;
18/12/2019;"People acting together as a group 
can accomplish things which no individual 
acting alone could ever 
hope to bring about.
";not-a-number;
18/12/2019;"3446335 
Remise De Cheque Ref: �3446335
¤£¥€¡¶€";;30,13;
18/12/2019;"3446535 
Remise De Cheque 
";;not-a-number;

  
Mr Jb Le Duigou     
VISA DUAL BZH DI no  4533 07xx xxxx xx60;

;Pas d'opération pour cette carte entre le 18/12/2019 et 18/12/2019;
`
