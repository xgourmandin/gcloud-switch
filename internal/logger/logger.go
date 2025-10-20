package logger

import (
	"fmt"
	"os"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// Info logs an informational message in cyan
func Info(msg string, args ...any) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, "%s%s%s\n", colorCyan, msg, colorReset) //nolint:errcheck
	} else {
		// Format key-value pairs
		_, _ = fmt.Fprintf(os.Stdout, "%s%s%s", colorCyan, msg, colorReset) //nolint:errcheck
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				_, _ = fmt.Fprintf(os.Stdout, " %s%v%s=%s%v%s", colorGray, args[i], colorReset, colorBold, args[i+1], colorReset) //nolint:errcheck
			}
		}
		_, _ = fmt.Fprintln(os.Stdout) //nolint:errcheck
	}
}

// Error logs an error message in red
func Error(msg string, args ...any) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "%s✗ %s%s\n", colorRed, msg, colorReset) //nolint:errcheck
	} else {
		_, _ = fmt.Fprintf(os.Stderr, "%s✗ %s%s", colorRed, msg, colorReset) //nolint:errcheck
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				_, _ = fmt.Fprintf(os.Stderr, " %s%v%s=%s%v%s", colorGray, args[i], colorReset, colorBold, args[i+1], colorReset) //nolint:errcheck
			}
		}
		_, _ = fmt.Fprintln(os.Stderr) //nolint:errcheck
	}
}

// Debug logs a debug message in gray
func Debug(msg string, args ...any) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, "%s[DEBUG] %s%s\n", colorGray, msg, colorReset) //nolint:errcheck
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "%s[DEBUG] %s%s", colorGray, msg, colorReset) //nolint:errcheck
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				_, _ = fmt.Fprintf(os.Stdout, " %v=%v", args[i], args[i+1]) //nolint:errcheck
			}
		}
		_, _ = fmt.Fprintln(os.Stdout) //nolint:errcheck
	}
}

// Success logs a success message with a checkmark in green
func Success(msg string, args ...any) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, "%s✓ %s%s\n", colorGreen, msg, colorReset) //nolint:errcheck
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "%s✓ %s%s", colorGreen, msg, colorReset) //nolint:errcheck
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				_, _ = fmt.Fprintf(os.Stdout, " %s%v%s=%s%v%s", colorGray, args[i], colorReset, colorGreen, args[i+1], colorReset) //nolint:errcheck
			}
		}
		_, _ = fmt.Fprintln(os.Stdout) //nolint:errcheck
	}
}

// Warning logs a warning message in yellow
func Warning(msg string, args ...any) {
	if len(args) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, "%s⚠ %s%s\n", colorYellow, msg, colorReset) //nolint:errcheck
	} else {
		_, _ = fmt.Fprintf(os.Stdout, "%s⚠ %s%s", colorYellow, msg, colorReset) //nolint:errcheck
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				_, _ = fmt.Fprintf(os.Stdout, " %s%v%s=%s%v%s", colorGray, args[i], colorReset, colorYellow, args[i+1], colorReset) //nolint:errcheck
			}
		}
		_, _ = fmt.Fprintln(os.Stdout) //nolint:errcheck
	}
}

// Plain logs a plain message without any color or prefix
func Plain(msg string, args ...any) {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(os.Stdout, msg) //nolint:errcheck
	} else {
		_, _ = fmt.Fprint(os.Stdout, msg) //nolint:errcheck
		for i := 0; i < len(args); i += 2 {
			if i+1 < len(args) {
				_, _ = fmt.Fprintf(os.Stdout, " %v: %v", args[i], args[i+1]) //nolint:errcheck
			}
		}
		_, _ = fmt.Fprintln(os.Stdout) //nolint:errcheck
	}
}
