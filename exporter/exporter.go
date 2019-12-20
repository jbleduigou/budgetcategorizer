package exporter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	budget "github.com/jbleduigou/budgetcategorizer"
)

type Exporter interface {
	Export(transactions []*budget.Transaction) ([]byte, error)
}

type csvwriter interface {
	Write(record []string) error
}

type csvExporter struct {
}

func NewExporter() Exporter {
	return &csvExporter{}
}

func (e *csvExporter) Export(transactions []*budget.Transaction) ([]byte, error) {
	var b bytes.Buffer
	writer := csv.NewWriter(&b)

	err := e.exportToCSV(transactions, writer)

	writer.Flush()
	return b.Bytes(), err
}

func (e *csvExporter) exportToCSV(transactions []*budget.Transaction, writer csvwriter) error {
	for _, value := range transactions {
		tunasse := strconv.FormatFloat(value.Value, 'f', -1, 64)

		line := []string{value.Date, value.Description, value.Comment, value.Category, string(tunasse)}
		err := writer.Write(line)
		if err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
	}
	return nil
}
