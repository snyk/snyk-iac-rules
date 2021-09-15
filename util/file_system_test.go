package util

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateDirectory(t *testing.T) {
	oldCreateDirectoryIfNotExists := createDirectoryIfNotExists
	defer func() {
		createDirectoryIfNotExists = oldCreateDirectoryIfNotExists
	}()

	oldComputePath := computePath
	defer func() {
		computePath = oldComputePath
	}()
	computePath = func(string, string) string {
		return "test"
	}

	t.Run("Returns the directory path when it succeeds", func(t *testing.T) {
		createDirectoryIfNotExists = func(string, bool) error {
			return nil
		}

		dirPath, err := CreateDirectory("../nested/workingDirectory", "name", true)
		assert.Nil(t, err)
		assert.Equal(t, "test", dirPath)
	})

	t.Run("Returns the directory path when it fails", func(t *testing.T) {
		createDirectoryIfNotExists = func(string, bool) error {
			return errors.New("Test")
		}

		dirPath, err := CreateDirectory("../nested/workingDirectory", "name", true)
		assert.NotNil(t, err)
		assert.Equal(t, "Test", err.Error())
		assert.Equal(t, "test", dirPath)
	})
}

func TestCreateFileAndClose(t *testing.T) {
	oldCreateFileAndCloseAndClose := createFileAndClose
	defer func() {
		createFileAndClose = oldCreateFileAndCloseAndClose
	}()
	createFileAndClose = func(string) error {
		return nil
	}

	oldComputePath := computePath
	defer func() {
		computePath = oldComputePath
	}()
	computePath = func(string, string) string {
		return "test"
	}

	t.Run("Returns the file path when it succeeds", func(t *testing.T) {
		createFileAndClose = func(string) error {
			return nil
		}

		dirName, err := CreateFile("../nested/workingDirectory", "name")
		assert.Nil(t, err)
		assert.Equal(t, "test", dirName)
	})

	t.Run("Returns the file path when it fails", func(t *testing.T) {
		createFileAndClose = func(string) error {
			return errors.New("Test")
		}

		dirName, err := CreateFile("../nested/workingDirectory", "name")
		assert.NotNil(t, err)
		assert.Equal(t, "Test", err.Error())
		assert.Equal(t, "test", dirName)
	})
}

func TestCreateDirectoryIfNotExists(t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err)

	t.Run("Creates a folder if it doesn't already exist", func(t *testing.T) {
		err := createDirectoryIfNotExists(path.Join(dir, "folder"), true)
		assert.Nil(t, err)
		_, err = os.Stat("./folder")
		assert.Nil(t, err)
		err = os.Remove("./folder")
		assert.Nil(t, err)
	})

	t.Run("Returns an error if folder already exists and it's a strict call", func(t *testing.T) {
		err := createDirectoryIfNotExists(dir, true)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "Directory already exists at path")
	})

	t.Run("Returns nothing error if folder already exists and it's not a strict call", func(t *testing.T) {
		err := createDirectoryIfNotExists(dir, false)
		assert.Nil(t, err)
	})
}

func TestComputePath(t *testing.T) {
	assert.Equal(t, "nested/workingDirectory/name", computePath("./nested/workingDirectory", "name"))
	assert.Equal(t, "workingDirectory/name", computePath("./workingDirectory", "name"))
	assert.Equal(t, "../workingDirectory/name", computePath("../workingDirectory", "name"))
}
