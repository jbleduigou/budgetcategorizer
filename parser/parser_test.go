package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTransactions(t *testing.T) {
	p := NewParser("1a5931ca-dd5d-11ea-90cb-3822e2348205")
	transactions, _ := p.ParseTransactions(bytes.NewBufferString(content))
	assert.Equal(t, 5, len(transactions))
	debit := transactions[0]
	assert.Equal(t, debit.Date, "18/12/2019")
	assert.Equal(t, debit.Description, "Paiement Par Carte Brulerie Des Capuci Brest 15/02 ")
	assert.Equal(t, debit.Comment, "")
	assert.Equal(t, debit.Category, "Courses Alimentation")
	assert.Equal(t, debit.Value, 3.18)
	credit := transactions[1]
	assert.Equal(t, credit.Date, "18/12/2019")
	assert.Equal(t, credit.Description, "3446335 Remise De Cheque Ref: 3446335 ")
	assert.Equal(t, credit.Comment, "")
	assert.Equal(t, credit.Category, "")
	assert.Equal(t, credit.Value, -30.13)
	cheque1 := transactions[2]
	assert.Equal(t, cheque1.Date, "28/02/2020")
	assert.Equal(t, cheque1.Description, "Cheque Emis 8936392")
	assert.Equal(t, cheque1.Comment, "")
	assert.Equal(t, cheque1.Category, "Courses Alimentation")
	assert.Equal(t, cheque1.Value, 118.8)
	cheque2 := transactions[3]
	assert.Equal(t, cheque2.Date, "29/02/2020")
	assert.Equal(t, cheque2.Description, "Cheque Emis 5423696")
	assert.Equal(t, cheque2.Comment, "")
	assert.Equal(t, cheque2.Category, "Courses Alimentation")
	assert.Equal(t, cheque2.Value, 39.0)
	grand := transactions[4]
	assert.Equal(t, grand.Date, "01/01/2020")
	assert.Equal(t, grand.Description, "Virement En Votre Faveur Youpi Youpi ")
	assert.Equal(t, grand.Comment, "")
	assert.Equal(t, grand.Category, "")
	assert.Equal(t, grand.Value, -1620.0)
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
18/12/2019;"PAIEMENT PAR CARTE      
BRULERIE DES CAPUCI BREST  15/02  



";3,18;;
18/12/2019;"PEOPLE ACTING TOGETHER AS A GROUP 
CAN ACCOMPLISH THINGS WHICH NO INDIVIDUAL 
ACTING ALONE COULD EVER 
HOPE TO BRING ABOUT.
";not-a-number;;
18/12/2019;"3446335 
REMISE DE CHEQUE REF: �3446335
¤£¥€¡¶€";;30,13;
18/12/2019;"3446535 
REMISE DE CHEQUE 
";;not-a-number;
28/02/2020;"8936392 
CHEQUE EMIS 
";118,80;;
29/02/2020;"5423696/0000000/000000000 
CHEQUE EMIS 
";39,00;;
01/01/2020;"VIREMENT EN VOTRE FAVEUR
Youpi  


Youpi             ";;1�620,00;
  
Mr Jb Le Duigou     
VISA DUAL BZH DI no  4533 07xx xxxx xx60;

;Pas d'opération pour cette carte entre le 18/12/2019 et 18/12/2019;
`
