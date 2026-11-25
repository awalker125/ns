package main

import (
	"os"
	"testing"
)

func TestFormatMessage(t *testing.T) {
	previousValue, wasSet := os.LookupEnv("NO_COLOR")
	t.Cleanup(func() {
		if wasSet {
			_ = os.Setenv("NO_COLOR", previousValue)
		} else {
			_ = os.Unsetenv("NO_COLOR")
		}
	})

	cases := []struct {
		name    string
		noColor *string
		level   string
		color   string
		want    string
	}{
		{
			name:  "colors only the info prefix when NO_COLOR is unset",
			level: "info",
			color: colorGreen,
			want:  colorGreen + "info" + colorReset + ": message",
		},
		{
			name:  "colors only the error prefix when NO_COLOR is unset",
			level: "error",
			color: colorRed,
			want:  colorRed + "error" + colorReset + ": message",
		},
		{
			name:    "disables color when NO_COLOR is empty",
			noColor: stringPointer(""),
			level:   "info",
			color:   colorGreen,
			want:    "info: message",
		},
		{
			name:    "disables color when NO_COLOR has a value",
			noColor: stringPointer("1"),
			level:   "error",
			color:   colorRed,
			want:    "error: message",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.noColor == nil {
				if err := os.Unsetenv("NO_COLOR"); err != nil {
					t.Fatalf("Unsetenv(NO_COLOR): %v", err)
				}
			} else if err := os.Setenv("NO_COLOR", *tc.noColor); err != nil {
				t.Fatalf("Setenv(NO_COLOR): %v", err)
			}

			if got := formatMessage(tc.level, tc.color, "message"); got != tc.want {
				t.Errorf("formatMessage() = %q, want %q", got, tc.want)
			}
		})
	}
}

func stringPointer(value string) *string {
	return &value
}
