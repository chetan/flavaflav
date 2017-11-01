package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandURL(t *testing.T) {
	u := "https://t.co/5C1hXyTnnR"
	e, err := ExpandURL(u)
	assert.NoError(t, err)
	println(e)
	assert.Equal(t, "https://twitter.com/KellyannePolls/status/917865369964023808", e)
}
