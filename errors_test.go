package nordeasiirto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArgError(t *testing.T) {
	expected := "foo is invalid because bar"
	err := NewArgError("foo", "bar")

	assert.Equal(t, expected, err.Error())
}
