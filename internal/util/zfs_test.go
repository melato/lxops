package util

import "testing"

func TestFormatZfsSize(t *testing.T) {
	tests := []struct {
		name    string
		gb      float32
		want    string
		wantErr bool
	}{
		{"whole", 1, "1GB", false},
		{"whole with trailing zeros", 10, "10GB", false},
		{"hundreds", 100, "100GB", false},
		{"half", 1.5, "1536MB", false},
		{"quarter", 0.25, "256MB", false},
		{"fraction above 1", 2.125, "2176MB", false},
		{"one MB", 1.0 / 1024, "1MB", false},
		{"rounds to nearest MB", 0.001, "1MB", false},
		{"less than 1GB", 0.5, "512MB", false},
		{"1023 MB", 1023.0 / 1024, "1023MB", false},
		{"negative", -1, "", true},
		{"small negative", -0.001, "", true},
		{"zero", 0, "", true},
		{"below 1 MB", 0.0004, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FormatZfsSize(tt.gb)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FormatZfsSize(%v) error = %v, wantErr %v", tt.gb, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("FormatZfsSize(%v) = %q, want %q", tt.gb, got, tt.want)
			}
		})
	}
}
