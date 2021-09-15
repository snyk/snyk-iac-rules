package util

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartProgress(t *testing.T) {
	rescueStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	err := StartProgress("Start progress", func() error {
		return nil
	})
	assert.Nil(t, err)

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		_, err := io.Copy(&buf, r)
		assert.Nil(t, err)
		outC <- buf.String()
	}()

	w.Close()
	os.Stderr = rescueStderr
	out := <-outC

	assert.Equal(t, "[ ] Start progress[/] Start progress[x] Start progress[ ] Start progress[/] Start progress[x] Start progress", out)
}
