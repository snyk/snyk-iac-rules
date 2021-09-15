package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepeatedStringFlag(t *testing.T) {
	repeatedStringFlag := NewRepeatedStringFlag("default")

	assert.Equal(t, "string", repeatedStringFlag.Type())
	assert.Equal(t, []string{"default"}, repeatedStringFlag.Strings())
	assert.Equal(t, true, repeatedStringFlag.IsSet())
	assert.Equal(t, "default", repeatedStringFlag.String())

	err := repeatedStringFlag.Set("option1")
	assert.Nil(t, err)
	assert.Equal(t, []string{"option1"}, repeatedStringFlag.Strings())
	assert.Equal(t, true, repeatedStringFlag.IsSet())
	assert.Equal(t, "option1", repeatedStringFlag.String())

	err = repeatedStringFlag.Set("option2")
	assert.Nil(t, err)
	assert.Equal(t, []string{"option1", "option2"}, repeatedStringFlag.Strings())
	assert.Equal(t, true, repeatedStringFlag.IsSet())
	assert.Equal(t, "option1,option2", repeatedStringFlag.String())
}
