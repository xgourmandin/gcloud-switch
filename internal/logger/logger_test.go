package logger

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close() //nolint:errcheck
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r) //nolint:errcheck
	return buf.String()
}

func captureErrorOutput(f func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	_ = w.Close() //nolint:errcheck
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r) //nolint:errcheck
	return buf.String()
}

func TestInfo(t *testing.T) {
	output := captureOutput(func() {
		Info("Test message")
	})

	if !strings.Contains(output, "Test message") {
		t.Errorf("Expected output to contain 'Test message', got: %s", output)
	}
}

func TestInfoWithArgs(t *testing.T) {
	output := captureOutput(func() {
		Info("Config added", "name", "dev-config", "project_id", "dev-123")
	})

	if !strings.Contains(output, "Config added") {
		t.Errorf("Expected output to contain 'Config added', got: %s", output)
	}
	if !strings.Contains(output, "dev-config") {
		t.Errorf("Expected output to contain 'dev-config', got: %s", output)
	}
	if !strings.Contains(output, "dev-123") {
		t.Errorf("Expected output to contain 'dev-123', got: %s", output)
	}
}

func TestSuccess(t *testing.T) {
	output := captureOutput(func() {
		Success("Operation successful")
	})

	if !strings.Contains(output, "✓") {
		t.Errorf("Expected output to contain '✓', got: %s", output)
	}
	if !strings.Contains(output, "Operation successful") {
		t.Errorf("Expected output to contain 'Operation successful', got: %s", output)
	}
}

func TestSuccessWithArgs(t *testing.T) {
	output := captureOutput(func() {
		Success("Configuration switched", "name", "prod")
	})

	if !strings.Contains(output, "✓") {
		t.Errorf("Expected output to contain '✓', got: %s", output)
	}
	if !strings.Contains(output, "Configuration switched") {
		t.Errorf("Expected output to contain 'Configuration switched', got: %s", output)
	}
	if !strings.Contains(output, "prod") {
		t.Errorf("Expected output to contain 'prod', got: %s", output)
	}
}

func TestError(t *testing.T) {
	output := captureErrorOutput(func() {
		Error("Error occurred")
	})

	if !strings.Contains(output, "✗") {
		t.Errorf("Expected output to contain '✗', got: %s", output)
	}
	if !strings.Contains(output, "Error occurred") {
		t.Errorf("Expected output to contain 'Error occurred', got: %s", output)
	}
}

func TestWarning(t *testing.T) {
	output := captureOutput(func() {
		Warning("Warning message")
	})

	if !strings.Contains(output, "⚠") {
		t.Errorf("Expected output to contain '⚠', got: %s", output)
	}
	if !strings.Contains(output, "Warning message") {
		t.Errorf("Expected output to contain 'Warning message', got: %s", output)
	}
}

func TestDebug(t *testing.T) {
	output := captureOutput(func() {
		Debug("Debug info")
	})

	if !strings.Contains(output, "[DEBUG]") {
		t.Errorf("Expected output to contain '[DEBUG]', got: %s", output)
	}
	if !strings.Contains(output, "Debug info") {
		t.Errorf("Expected output to contain 'Debug info', got: %s", output)
	}
}

func TestPlain(t *testing.T) {
	output := captureOutput(func() {
		Plain("Plain text message")
	})

	if !strings.Contains(output, "Plain text message") {
		t.Errorf("Expected output to contain 'Plain text message', got: %s", output)
	}
	// Plain should not have any special prefix
	if strings.Contains(output, "✓") || strings.Contains(output, "✗") || strings.Contains(output, "⚠") {
		t.Errorf("Expected output to not contain special symbols, got: %s", output)
	}
}

func TestPlainWithArgs(t *testing.T) {
	output := captureOutput(func() {
		Plain("Status", "active", "true")
	})

	if !strings.Contains(output, "Status") {
		t.Errorf("Expected output to contain 'Status', got: %s", output)
	}
	if !strings.Contains(output, "active") {
		t.Errorf("Expected output to contain 'active', got: %s", output)
	}
	if !strings.Contains(output, "true") {
		t.Errorf("Expected output to contain 'true', got: %s", output)
	}
}

func TestColorCodes(t *testing.T) {
	// Test that color codes are defined
	if colorReset == "" {
		t.Error("colorReset should not be empty")
	}
	if colorGreen == "" {
		t.Error("colorGreen should not be empty")
	}
	if colorRed == "" {
		t.Error("colorRed should not be empty")
	}
	if colorCyan == "" {
		t.Error("colorCyan should not be empty")
	}
	if colorYellow == "" {
		t.Error("colorYellow should not be empty")
	}
}
