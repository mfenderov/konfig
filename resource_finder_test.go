package konfig

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindRootPath(t *testing.T) {
	// Save the original rootInstance to restore it after the test
	originalRootInstance := rootInstance
	defer func() {
		rootInstance = originalRootInstance
	}()

	// Create a new rootInstance for this test
	rootInstance = &root{}

	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "konfig-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a resources directory in the temp directory
	resourcesDir := filepath.Join(tempDir, "resources")
	err = os.Mkdir(resourcesDir, 0755)
	assert.NoError(t, err)

	// Change to a subdirectory of the temp directory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	assert.NoError(t, err)

	// Save current directory to restore it later
	currentDir, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(currentDir)

	// Change to the subdirectory
	err = os.Chdir(subDir)
	assert.NoError(t, err)

	// Test finding the root path
	path, err := findRootPath()
	assert.NoError(t, err)

	// Evaluate symlinks to get the real path
	realTempDir, err := filepath.EvalSymlinks(tempDir)
	assert.NoError(t, err)
	assert.Equal(t, realTempDir, path)
}

func TestFindRootPath_NoResourcesDir(t *testing.T) {
	// Save the original rootInstance to restore it after the test
	originalRootInstance := rootInstance
	defer func() {
		rootInstance = originalRootInstance
	}()

	// Create a new rootInstance for this test
	rootInstance = &root{}

	// Create a temporary directory without a resources subdirectory
	tempDir, err := os.MkdirTemp("", "konfig-test-no-resources")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Save current directory to restore it later
	currentDir, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(currentDir)

	// Change to the temp directory
	err = os.Chdir(tempDir)
	assert.NoError(t, err)

	// Test finding the root path - should fail
	_, err = findRootPath()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to find root path with 'resources' folder")
}

func TestFindRootPath_Caching(t *testing.T) {
	// Save the original rootInstance to restore it after the test
	originalRootInstance := rootInstance
	defer func() {
		rootInstance = originalRootInstance
	}()

	// Create a new rootInstance for this test
	rootInstance = &root{}

	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "konfig-test-caching")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a resources directory in the temp directory
	resourcesDir := filepath.Join(tempDir, "resources")
	err = os.Mkdir(resourcesDir, 0755)
	assert.NoError(t, err)

	// Save current directory to restore it later
	currentDir, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(currentDir)

	// Change to the temp directory
	err = os.Chdir(tempDir)
	assert.NoError(t, err)

	// Evaluate symlinks to get the real path
	realTempDir, err := filepath.EvalSymlinks(tempDir)
	assert.NoError(t, err)

	// First call to findRootPath
	path1, err1 := findRootPath()
	assert.NoError(t, err1)
	assert.Equal(t, realTempDir, path1)

	// Remove the resources directory to verify caching
	err = os.RemoveAll(resourcesDir)
	assert.NoError(t, err)

	// Second call to findRootPath should return the cached result
	path2, err2 := findRootPath()
	assert.NoError(t, err2)
	assert.Equal(t, realTempDir, path2)
}
