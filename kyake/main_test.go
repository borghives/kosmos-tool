package main

import (
	"testing"
)

func TestProcessTemplate(t *testing.T) {
	envMap := map[string]string{
		"SERVICE_NAME": "my-service",
		"PORT":         "8080",
		"EMPTY_VAR":    "",
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "Basic substitution",
			template: "Service: {{ .SERVICE_NAME }}",
			expected: "Service: my-service",
		},
		{
			name:     "Missing variable",
			template: "Host: {{ .MISSING_HOST }}",
			expected: "Host: ",
		},
		{
			name:     "Conditionals",
			template: "{{ if .PORT }}Port: {{ .PORT }}{{ end }}",
			expected: "Port: 8080",
		},
		{
			name:     "Conditionals missing",
			template: "{{ if .MISSING_HOST }}Host: {{ .MISSING_HOST }}{{ end }}",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ProcessTemplate(tt.template, envMap)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if actual != tt.expected {
				t.Errorf("ProcessTemplate() = %q, want %q", actual, tt.expected)
			}
		})
	}
}
