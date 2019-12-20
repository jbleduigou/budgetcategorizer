package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetResultFileName(t *testing.T) {
	c := &command{}

	output := c.getResultFileName("CA20191220_1142.CSV")
	assert.Equal(t, "CA20191220_1142-result.txt", output)
}
