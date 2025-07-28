package entities

import (
	"testing"
)

func TestPasswordConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  PasswordConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: PasswordConfig{
				Length:         12,
				IncludeLower:   true,
				IncludeUpper:   true,
				IncludeNumbers: true,
				IncludeSymbols: false,
				Count:          1,
			},
			wantErr: false,
		},
		{
			name: "minimum length",
			config: PasswordConfig{
				Length:         4,
				IncludeLower:   true,
				IncludeUpper:   false,
				IncludeNumbers: false,
				IncludeSymbols: false,
				Count:          1,
			},
			wantErr: false,
		},
		{
			name: "zero length",
			config: PasswordConfig{
				Length:         0,
				IncludeLower:   true,
				IncludeUpper:   false,
				IncludeNumbers: false,
				IncludeSymbols: false,
				Count:          1,
			},
			wantErr: true,
		},
		{
			name: "no character sets",
			config: PasswordConfig{
				Length:         12,
				IncludeLower:   false,
				IncludeUpper:   false,
				IncludeNumbers: false,
				IncludeSymbols: false,
				Count:          1,
			},
			wantErr: true,
		},
		{
			name: "zero count",
			config: PasswordConfig{
				Length:         12,
				IncludeLower:   true,
				IncludeUpper:   false,
				IncludeNumbers: false,
				IncludeSymbols: false,
				Count:          0,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("PasswordConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
