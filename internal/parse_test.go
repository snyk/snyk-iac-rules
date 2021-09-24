package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/snyk/snyk-iac-rules/util"
	"github.com/stretchr/testify/assert"
)

func mockParseParams() *ParseCommandParams {
	return &ParseCommandParams{
		Format: util.NewEnumFlag(HCL2, []string{HCL2, YAML}),
	}
}

func TestParse(t *testing.T) {
	input := []byte("{\"foo\": \"bar\"}")
	expected := `{
	"foo": "bar"
}
`

	oldReadFile := readFile
	defer func() {
		readFile = oldReadFile
	}()
	readFile = func(string) ([]byte, error) {
		return []byte("Test"), nil
	}

	t.Run("Prints out the JSON format of YAML", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, err := os.Pipe()
		assert.Nil(t, err)
		os.Stdout = w

		oldParseYAML := parseYAML
		defer func() {
			parseYAML = oldParseYAML
		}()
		parseYAML = func(content []byte, parsedInput interface{}) error {
			assert.Equal(t, "Test", string(content))
			return json.Unmarshal(input, parsedInput)
		}

		parseParams := mockParseParams()
		err = parseParams.Format.Set(YAML)
		assert.Nil(t, err)

		err = RunParse([]string{"test"}, parseParams)
		assert.Nil(t, err)

		outC := make(chan string)
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func() {
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			assert.Nil(t, err)
			outC <- buf.String()
		}()

		w.Close()
		os.Stdout = rescueStdout
		out := <-outC

		assert.Equal(t, expected, out)
	})

	t.Run("Prints out the JSON format of HCL2", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, err := os.Pipe()
		assert.Nil(t, err)
		os.Stdout = w

		oldParseHCL2 := parseHCL2
		defer func() {
			parseHCL2 = oldParseHCL2
		}()
		parseHCL2 = func(content []byte, parsedInput interface{}) error {
			assert.Equal(t, "Test", string(content))
			return json.Unmarshal(input, parsedInput)
		}

		parseParams := mockParseParams()
		err = parseParams.Format.Set(HCL2)
		assert.Nil(t, err)

		err = RunParse([]string{"test"}, parseParams)
		assert.Nil(t, err)

		outC := make(chan string)
		// copy the output in a separate goroutine so printing can't block indefinitely
		go func() {
			var buf bytes.Buffer
			_, err = io.Copy(&buf, r)
			assert.Nil(t, err)
			outC <- buf.String()
		}()

		w.Close()
		os.Stdout = rescueStdout
		out := <-outC

		assert.Equal(t, expected, out)
	})
}
