package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
)

// The CLI ignores the first input from os.Args so we must mock it
var mock = "test"

func Test(t *testing.T) {
	result, stdout, _ := os.Pipe()

	app := cli.NewApp()
	app.Writer = stdout

	err := app.Run([]string{mock})
	assert.Nil(t, err)
	stdout.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, result)
	assert.Nil(t, err)

	assert.Contains(t, buf.String(), "NAME")
	assert.Contains(t, buf.String(), "USAGE")
	assert.Contains(t, buf.String(), "COMMANDS")
	assert.Contains(t, buf.String(), "GLOBAL OPTIONS")
	assert.Contains(t, buf.String(), "--help, -h")
}
