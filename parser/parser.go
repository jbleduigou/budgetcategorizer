package parser

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	budget "github.com/jbleduigou/budgetcategorizer"
)

type Parser interface {
	ParseTransactions(r io.Reader) ([]*budget.Transaction, error)
}

type csvParser struct {
}

func NewParser() Parser {
	return &csvParser{}
}

func (c *csvParser) ParseTransactions(r io.Reader) ([]*budget.Transaction, error) {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.Comma = ';'
	reader.FieldsPerRecord = 4
	reader.FieldsPerRecord = -1
	return c.parse(reader)
}

func (c *csvParser) parse(reader *csv.Reader) (transactions []*budget.Transaction, err error) {
	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, each := range rawCSVdata {
		if len(each) == 4 {
			date := each[0]
			libelle := c.sanitizeDescription(each[1])
			debit, err := c.parseAmount(each[2])
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			t := budget.NewTransaction(date, libelle, "", "Courses Alimentation", debit)
			transactions = append(transactions, t)
		}
		if len(each) == 5 {
			date := each[0]
			if "Date" != date {
				libelle := c.sanitizeDescription(each[1])
				credit, err := c.parseAmount(each[3])
				if err != nil {
					fmt.Printf("%v\n", err)
				}
				t := budget.NewTransaction(date, libelle, "", "", -credit)
				transactions = append(transactions, t)
			}
		}
	}
	fmt.Printf("Found %v transactions\n", len(transactions))
	return transactions, nil
}

func (c *csvParser) sanitizeDescription(d string) string {
	libelle := []byte(d)
	{
		re := regexp.MustCompile(`\n`)
		libelle = re.ReplaceAll(libelle, []byte(" "))
	}
	{
		re := regexp.MustCompile(`[\s]+`)
		libelle = re.ReplaceAll(libelle, []byte(" "))
	}
	return string(libelle)
}

func (c *csvParser) parseAmount(a string) (float64, error) {
	creditStr := []byte(a)
	{
		re := regexp.MustCompile(`,`)
		creditStr = re.ReplaceAll(creditStr, []byte("."))
	}
	credit, err := strconv.ParseFloat(string(creditStr), 64)
	return credit, err
}
