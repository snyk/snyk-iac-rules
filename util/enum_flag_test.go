package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnumFlag(t *testing.T) {
	enumFlag := NewEnumFlag("default", []string{"option1", "option2"})

	assert.Equal(t, "{option1,option2}", enumFlag.Type())
	assert.Equal(t, false, enumFlag.IsSet())
	assert.Equal(t, "default", enumFlag.String())

	err := enumFlag.Set("option1")
	assert.Nil(t, err)
	assert.Equal(t, true, enumFlag.IsSet())
	assert.Equal(t, "option1", enumFlag.String())

	err = enumFlag.Set("wrong")
	assert.NotNil(t, err)
	assert.Equal(t, true, enumFlag.IsSet())
	assert.Equal(t, "option1", enumFlag.String())
}
