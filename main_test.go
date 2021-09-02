package main

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/snyk/snyk-iac-custom-rules/cmd"
)

func Test(t *testing.T) {
	cmd := cmd.RootCommand
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	err := cmd.Execute()
	assert.Nil(t, err)

	out, err := ioutil.ReadAll(b)
	assert.Nil(t, err)

	assert.Contains(t, string(out), "An SDK to write, debug, test, and bundle custom rules for Snyk IaC.")
}
