package output

import (
	"testing"
)

func TestSetGetFormat(t *testing.T) {
	SetFormat(Table)
	if GetFormat() != Table {
		t.Errorf("GetFormat() = %d, want Table(%d)", GetFormat(), Table)
	}

	SetFormat(JSON)
	if GetFormat() != JSON {
		t.Errorf("GetFormat() = %d, want JSON(%d)", GetFormat(), JSON)
	}

	SetFormat(Table)
}

func TestFormatPnL(t *testing.T) {
	tests := []struct {
		name string
		val  float64
		want string
	}{
		{"positive", 123.45, "+123.45"},
		{"negative", -50.00, "-50.00"},
		{"zero", 0, "0.00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatPnL(tt.val)
			if len(got) == 0 {
				t.Error("FormatPnL returned empty string")
			}
			// ANSI color codes make exact match hard; check value is present
			if tt.val > 0 && !containsSubstring(got, "123.45") {
				t.Errorf("FormatPnL(%f) should contain '123.45', got %q", tt.val, got)
			}
		})
	}
}

func TestFormatPercent(t *testing.T) {
	tests := []struct {
		val      float64
		contains string
	}{
		{5.25, "5.25%"},
		{-3.10, "3.10%"},
		{0, "0.00%"},
	}

	for _, tt := range tests {
		got := FormatPercent(tt.val)
		if !containsSubstring(got, tt.contains) {
			t.Errorf("FormatPercent(%f) should contain %q, got %q", tt.val, tt.contains, got)
		}
	}
}

func TestFormatMoney(t *testing.T) {
	tests := []struct {
		val  float64
		want string
	}{
		{1500000, "$1.5M"},
		{50000, "$50.0K"},
		{999, "$999.00"},
		{0, "$0.00"},
	}

	for _, tt := range tests {
		got := FormatMoney(tt.val)
		if got != tt.want {
			t.Errorf("FormatMoney(%f) = %q, want %q", tt.val, got, tt.want)
		}
	}
}

func TestFormatBool(t *testing.T) {
	trueResult := FormatBool(true, "yes", "no")
	if !containsSubstring(trueResult, "yes") {
		t.Errorf("FormatBool(true) should contain 'yes', got %q", trueResult)
	}

	falseResult := FormatBool(false, "yes", "no")
	if !containsSubstring(falseResult, "no") {
		t.Errorf("FormatBool(false) should contain 'no', got %q", falseResult)
	}
}

func TestBoldGreenRedYellowCyan(t *testing.T) {
	if Bold("test") == "" {
		t.Error("Bold returned empty")
	}
	if Green("test") == "" {
		t.Error("Green returned empty")
	}
	if Red("test") == "" {
		t.Error("Red returned empty")
	}
	if Yellow("test") == "" {
		t.Error("Yellow returned empty")
	}
	if Cyan("test") == "" {
		t.Error("Cyan returned empty")
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && findSubstring(s, sub))
}

func findSubstring(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
