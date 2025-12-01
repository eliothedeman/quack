package quack

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExistingFilePathValidate(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	require.Nil(t, err)
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath)
	tmpFile.Close()

	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testdir_*")
	require.Nil(t, err)
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name    string
		path    ExistingFilePath
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid_existing_file",
			path:    ExistingFilePath(tmpFilePath),
			wantErr: false,
		},
		{
			name:    "nonexistent_file",
			path:    ExistingFilePath("/nonexistent/path/to/file.txt"),
			wantErr: true,
			errMsg:  "not able to find file",
		},
		{
			name:    "directory_instead_of_file",
			path:    ExistingFilePath(tmpDir),
			wantErr: true,
			errMsg:  "is a directory. expected to be file",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.path.Validate()
			if test.wantErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), test.errMsg)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestExistingFilePathOpen(t *testing.T) {
	// Create a temporary file with some content
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	require.Nil(t, err)
	defer os.Remove(tmpFile.Name())

	content := "test content"
	_, err = tmpFile.WriteString(content)
	require.Nil(t, err)
	tmpFile.Close()

	t.Run("open_existing_file", func(t *testing.T) {
		efp := ExistingFilePath(tmpFile.Name())
		file := efp.Open()
		require.NotNil(t, file)
		defer file.Close()

		// Verify we can read from the file
		buf := make([]byte, len(content))
		n, err := file.Read(buf)
		assert.Nil(t, err)
		assert.Equal(t, len(content), n)
		assert.Equal(t, content, string(buf))
	})
}

func TestExistingFilePathOpenWith(t *testing.T) {
	// Create a temporary file
	tmpFile, err := os.CreateTemp("", "test_*.txt")
	require.Nil(t, err)
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	t.Run("open_file_read_only", func(t *testing.T) {
		efp := ExistingFilePath(tmpFile.Name())
		file := efp.OpenWith(os.O_RDONLY, 0)
		require.NotNil(t, file)
		defer file.Close()

		// File should be open and readable
		assert.NotNil(t, file)
	})

	t.Run("open_file_read_write", func(t *testing.T) {
		efp := ExistingFilePath(tmpFile.Name())
		file := efp.OpenWith(os.O_RDWR, 0)
		require.NotNil(t, file)
		defer file.Close()

		// Write to the file to verify it's writable
		content := "new content"
		n, err := file.WriteString(content)
		assert.Nil(t, err)
		assert.Equal(t, len(content), n)
	})

	t.Run("open_file_append", func(t *testing.T) {
		tmpFile2, err := os.CreateTemp("", "test2_*.txt")
		require.Nil(t, err)
		defer os.Remove(tmpFile2.Name())
		tmpFile2.Close()

		efp := ExistingFilePath(tmpFile2.Name())
		file := efp.OpenWith(os.O_APPEND|os.O_WRONLY, 0)
		require.NotNil(t, file)
		defer file.Close()

		// Write to the file in append mode
		content := "appended content"
		n, err := file.WriteString(content)
		assert.Nil(t, err)
		assert.Equal(t, len(content), n)
	})
}

func TestExistingFilePathString(t *testing.T) {
	path := ExistingFilePath("/path/to/file.txt")
	assert.Equal(t, "/path/to/file.txt", string(path))
}

func TestExistingFilePathRoundTrip(t *testing.T) {
	// Create a temporary file with some content
	tmpFile, err := os.CreateTemp("", "integration_test_*.txt")
	require.Nil(t, err)
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath)

	content := "integration test content"
	_, err = tmpFile.WriteString(content)
	require.Nil(t, err)
	tmpFile.Close()

	t.Run("round_trip_open_and_read", func(t *testing.T) {
		// Use absolute path for Open() to work properly
		efp := ExistingFilePath(tmpFilePath)

		// Open and read the file
		file := efp.Open()
		require.NotNil(t, file)
		defer file.Close()

		// Verify we can read the content
		buf := make([]byte, len(content))
		n, err := file.Read(buf)
		assert.Nil(t, err)
		assert.Equal(t, len(content), n)
		assert.Equal(t, content, string(buf))
	})
}
