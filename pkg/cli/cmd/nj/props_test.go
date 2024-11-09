// File: cmd/properties_test.go
package cmd

import (
	"bytes"
	"errors"
	"testing"

	"radius/clierrors"
)

func TestRunProperties(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		mockProps     map[string]string
		mockError     error
		expectedOut   string
		expectedError string
	}{
		{
			name: "Successful retrieval",
			args: []string{"resource1"},
			mockProps: map[string]string{
				"Name":   "Example Resource",
				"Type":   "Service",
				"Status": "Active",
			},
			mockError:   nil,
			expectedOut: "Name: Example Resource\nType: Service\nStatus: Active\n",
		},
		{
			name:          "Resource not found",
			args:          []string{"invalid-resource"},
			mockProps:     nil,
			mockError:     clierrors.NotFoundError,
			expectedError: `The resource "invalid-resource" could not be found.`,
		},
		{
			name:          "Unexpected error",
			args:          []string{"resource2"},
			mockProps:     nil,
			mockError:     errors.New("database connection failed"),
			expectedError: `error retrieving properties for resource "resource2": database connection failed`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock getResourceProperties
			originalFunc := getResourceProperties
			getResourceProperties = func(id string) (map[string]string, error) {
				return tt.mockProps, tt.mockError
			}
			defer func() { getResourceProperties = originalFunc }()

			// Capture output
			var out bytes.Buffer
			rootCmd.SetOut(&out)
			rootCmd.SetErr(&out)
			rootCmd.SetArgs(append([]string{"properties"}, tt.args...))

			// Execute command
			err := rootCmd.Execute()

			if tt.expectedError != "" {
				if err == nil || err.Error() != tt.expectedError {
					t.Errorf("expected error '%s', got '%v'", tt.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if out.String() != tt.expectedOut {
					t.Errorf("expected output '%s', got '%s'", tt.expectedOut, out.String())
				}
			}
		})
	}
}
