package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMain handles setup and teardown for all tests
func TestMain(m *testing.M) {
	// Setup before all tests
	setup()

	// Run all tests
	exitCode := m.Run()

	// Cleanup after all tests
	teardown()

	os.Exit(exitCode)
}

// Test that the application can be built and run
func TestAppBuild(t *testing.T) {
	// This is a basic smoke test to ensure the app can be built
	// In a real CI pipeline, you would use build tags to exclude this from regular test runs
	cmd := exec.Command("go", "build", "-o", "termpilot.test")
	output, err := cmd.CombinedOutput()

	assert.NoError(t, err, "Build failed with error: %s", string(output))
	assert.FileExists(t, "termpilot.test")
}

// Setup test environment
func setup() {
	// Create any test files needed
	// Set environment variables for testing
	os.Setenv("TEST_MODE", "true")
}

// Cleanup test environment
func teardown() {
	// Remove test files
	os.Remove("termpilot.test")
	os.Remove("test.db")

	// Restore environment
	os.Unsetenv("TEST_MODE")
}
