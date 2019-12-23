package budget

// Transaction is a transaction in the budget (purchase, refund).
type Transaction struct {
	Date        string
	Description string
	Comment     string
	Category    string
	Value       float64
}

// NewTransaction creates a transaction.
func NewTransaction(date string, description string, comment string, category string, value float64) Transaction {
	return Transaction{Date: date, Description: description, Comment: comment, Category: category, Value: value}
}
