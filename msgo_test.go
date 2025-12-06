package msgo_test

import (
	"strings"
	"testing"
	"time"

	"github.com/0x4A756E65/msgo"
)

func TestParse(t *testing.T) {
	year := 24 * time.Hour * 36525 / 100
	month := year / 12
	ms := func(v float64) time.Duration {
		return time.Duration(v * float64(time.Millisecond))
	}

	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{"hours short", "2h", 2 * time.Hour, false},
		{"hours long", "2 hours", 2 * time.Hour, false},
		{"year", "1y", year, false},
		{"week", "1w", 7 * 24 * time.Hour, false},
		{"month", "2mo", 2 * month, false},
		{"days long", "2 days", 2 * 24 * time.Hour, false},
		{"months long", "1 month", month, false},
		{"decimal hours", "1.5h", 90 * time.Minute, false},
		{"leading decimal", ".5h", 30 * time.Minute, false},
		{"decimal ms", ".5ms", ms(0.5), false},
		{"bare number", "500", 500 * time.Millisecond, false},
		{"upper minutes", "3MINS", 3 * time.Minute, false},
		{"negative", "-1.5h", -90 * time.Minute, false},
		{"negative decimal", "-.5h", -30 * time.Minute, false},
		{"empty", "", 0, true},
		{"too long", strings.Repeat("▲", 101), 0, true},
		{"nonsense", "abc", 0, true},
		{"nonsense symbol", "☃", 0, true},
		{"bad number", "10-.5", 0, true},
		{"unknown unit literal", "ms", 0, true},
		{"unknown unit", "10 lightyears", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := msgo.Parse(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("Parse(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Fatalf("Parse(%q) unexpected error: %v", tt.input, err)
			}

			if got != tt.want {
				t.Fatalf("Parse(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormatShort(t *testing.T) {
	year := 24 * time.Hour * 36525 / 100
	month := year / 12
	ms := func(v float64) time.Duration {
		return time.Duration(v * float64(time.Millisecond))
	}

	tests := []struct {
		input time.Duration
		want  string
	}{
		{2 * time.Hour, "2h"},
		{90 * time.Second, "2m"},
		{15 * time.Minute, "15m"},
		{36 * time.Hour, "2d"},
		{2 * month, "2mo"},
		{ms(500), "500ms"},
		{-ms(500), "-500ms"},
		{ms(1000), "1s"},
		{ms(10000), "10s"},
		{ms(60 * 1000), "1m"},
		{ms(60 * 10000), "10m"},
		{ms(60 * 60 * 1000), "1h"},
		{ms(60 * 60 * 10000), "10h"},
		{ms(24 * 60 * 60 * 1000), "1d"},
		{ms(24 * 60 * 60 * 6000), "6d"},
		{ms(1 * 7 * 24 * 60 * 60 * 1000), "1w"},
		{ms(2 * 7 * 24 * 60 * 60 * 1000), "2w"},
		{ms(30.4375 * 24 * 60 * 60 * 1000), "1mo"},
		{ms(30.4375 * 24 * 60 * 60 * 1200), "1mo"},
		{ms(30.4375 * 24 * 60 * 60 * 10000), "10mo"},
		{ms(365.25*24*60*60*1000 + 1), "1y"},
		{ms(365.25*24*60*60*1200 + 1), "1y"},
		{ms(365.25*24*60*60*10000 + 1), "10y"},
		{-2 * time.Hour, "-2h"},
		{-ms(1000), "-1s"},
		{-ms(10000), "-10s"},
		{-ms(60 * 1000), "-1m"},
		{-ms(60 * 10000), "-10m"},
		{-ms(60 * 60 * 1000), "-1h"},
		{-ms(60 * 60 * 10000), "-10h"},
		{-ms(24 * 60 * 60 * 1000), "-1d"},
		{-ms(24 * 60 * 60 * 6000), "-6d"},
		{-ms(1 * 7 * 24 * 60 * 60 * 1000), "-1w"},
		{-ms(2 * 7 * 24 * 60 * 60 * 1000), "-2w"},
		{-ms(30.4375 * 24 * 60 * 60 * 1000), "-1mo"},
		{-ms(30.4375 * 24 * 60 * 60 * 1200), "-1mo"},
		{-ms(30.4375 * 24 * 60 * 60 * 10000), "-10mo"},
		{-ms(365.25*24*60*60*1000 + 1), "-1y"},
		{-ms(365.25*24*60*60*1200 + 1), "-1y"},
		{-ms(365.25*24*60*60*10000 + 1), "-10y"},
		{ms(234234234), "3d"},
		{-ms(234234234), "-3d"},
	}

	for _, tt := range tests {
		if got := msgo.FormatShort(tt.input); got != tt.want {
			t.Fatalf("FormatShort(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatLong(t *testing.T) {
	year := 24 * time.Hour * 36525 / 100
	month := year / 12
	ms := func(v float64) time.Duration {
		return time.Duration(v * float64(time.Millisecond))
	}

	tests := []struct {
		input time.Duration
		want  string
	}{
		{2 * time.Hour, "2 hours"},
		{90 * time.Second, "2 minutes"},
		{time.Hour, "1 hour"},
		{36 * time.Hour, "2 days"},
		{1500 * time.Millisecond, "2 seconds"},
		{500 * time.Millisecond, "500 ms"},
		{-2 * time.Hour, "-2 hours"},
		{2 * month, "2 months"},
		{ms(1200), "1 second"},
		{ms(10000), "10 seconds"},
		{ms(60 * 1000), "1 minute"},
		{ms(60 * 1200), "1 minute"},
		{ms(60 * 10000), "10 minutes"},
		{ms(60 * 60 * 1000), "1 hour"},
		{ms(60 * 60 * 1200), "1 hour"},
		{ms(60 * 60 * 10000), "10 hours"},
		{ms(24 * 60 * 60 * 1000), "1 day"},
		{ms(1 * 24 * 60 * 60 * 1200), "1 day"},
		{ms(6 * 24 * 60 * 60 * 1000), "6 days"},
		{ms(1 * 7 * 24 * 60 * 60 * 1000), "1 week"},
		{ms(2 * 7 * 24 * 60 * 60 * 1000), "2 weeks"},
		{ms(30.4375 * 24 * 60 * 60 * 1000), "1 month"},
		{ms(30.4375 * 24 * 60 * 60 * 1200), "1 month"},
		{ms(30.4375 * 24 * 60 * 60 * 10000), "10 months"},
		{ms(365.25*24*60*60*1000 + 1), "1 year"},
		{ms(365.25*24*60*60*1200 + 1), "1 year"},
		{ms(365.25*24*60*60*10000 + 1), "10 years"},
		{ms(234234234), "3 days"},
		{-ms(1200), "-1 second"},
		{-ms(10000), "-10 seconds"},
		{-ms(60 * 1000), "-1 minute"},
		{-ms(60 * 1200), "-1 minute"},
		{-ms(60 * 10000), "-10 minutes"},
		{-ms(60 * 60 * 1000), "-1 hour"},
		{-ms(60 * 60 * 1200), "-1 hour"},
		{-ms(60 * 60 * 10000), "-10 hours"},
		{-ms(1 * 24 * 60 * 60 * 1000), "-1 day"},
		{-ms(1 * 24 * 60 * 60 * 1200), "-1 day"},
		{-ms(6 * 24 * 60 * 60 * 1000), "-6 days"},
		{-ms(1 * 7 * 24 * 60 * 60 * 1000), "-1 week"},
		{-ms(2 * 7 * 24 * 60 * 60 * 1000), "-2 weeks"},
		{-ms(30.4375 * 24 * 60 * 60 * 1000), "-1 month"},
		{-ms(30.4375 * 24 * 60 * 60 * 1200), "-1 month"},
		{-ms(30.4375 * 24 * 60 * 60 * 10000), "-10 months"},
		{-ms(365.25*24*60*60*1000 + 1), "-1 year"},
		{-ms(365.25*24*60*60*1200 + 1), "-1 year"},
		{-ms(365.25*24*60*60*10000 + 1), "-10 years"},
		{-ms(234234234), "-3 days"},
	}

	for _, tt := range tests {
		if got := msgo.FormatLong(tt.input); got != tt.want {
			t.Fatalf("FormatLong(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	d := 90 * time.Second
	if got := msgo.Format(d, false); got != "2m" {
		t.Fatalf("Format short = %q, want %q", got, "2m")
	}
	if got := msgo.Format(d, true); got != "2 minutes" {
		t.Fatalf("Format long = %q, want %q", got, "2 minutes")
	}
}

func TestFormatInvalidLength(t *testing.T) {
	if _, err := msgo.Parse(""); err == nil {
		t.Fatalf("Parse empty should error")
	}
	if _, err := msgo.Parse(strings.Repeat("x", 101)); err == nil {
		t.Fatalf("Parse overly long string should error")
	}
}
