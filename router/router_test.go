package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExactSameMatch(t *testing.T) {
	assert := assert.New(t)

	output, _ := Match("/", "/")
	assert.Equal(true, output, "match / /")
}

func TestSingleParamMatch(t *testing.T) {
	assert := assert.New(t)

	output, _ := Match("/:id", "/")
	assert.Equal(true, output, "match /:id /")
}

func TestMultipleParamMatch(t *testing.T) {
	assert := assert.New(t)

	output, _ := Match("/:id/:name", "/")
	assert.Equal(false, output, "match /:id/:name /")
}
