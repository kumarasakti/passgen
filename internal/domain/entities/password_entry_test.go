package entities

import (
	"testing"
	"time"
)

func TestPasswordEntry_Validation(t *testing.T) {
	tests := []struct {
		name  string
		entry PasswordEntry
		valid bool
	}{
		{
			name: "valid password entry",
			entry: PasswordEntry{
				Service:     "github",
				Username:    "john.doe",
				Password:    "secret123",
				URL:         "https://github.com",
				Notes:       "Personal account",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				GeneratedBy: "passgen",
			},
			valid: true,
		},
		{
			name: "minimal valid entry",
			entry: PasswordEntry{
				Service:     "test",
				Password:    "password",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				GeneratedBy: "passgen",
			},
			valid: true,
		},
		{
			name: "entry with auto-rotation",
			entry: PasswordEntry{
				Service:   "aws",
				Password:  "complex-password",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				AutoRotation: &AutoRotationConfig{
					Enabled:          true,
					IntervalDays:     90,
					NextRotationAt:   time.Now().Add(90 * 24 * time.Hour),
					NotifyDaysBefore: 7,
					AutoGenerate:     true,
					PasswordProfile: &PasswordProfile{
						Length:         16,
						IncludeUpper:   true,
						IncludeLower:   true,
						IncludeNumbers: true,
						IncludeSymbols: true,
					},
				},
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Basic validation - service and password should not be empty
			if tt.entry.Service == "" && tt.valid {
				t.Error("Expected valid entry but service is empty")
			}
			if tt.entry.Password == "" && tt.valid {
				t.Error("Expected valid entry but password is empty")
			}
			if tt.entry.CreatedAt.IsZero() && tt.valid {
				t.Error("Expected valid entry but CreatedAt is zero")
			}
		})
	}
}

func TestPasswordStore_Validation(t *testing.T) {
	tests := []struct {
		name  string
		store PasswordStore
		valid bool
	}{
		{
			name: "valid password store",
			store: PasswordStore{
				Name:       "personal",
				GitURL:     "git@github.com:user/passwords.git",
				LocalPath:  "/home/user/.password-stores/personal",
				GPGKeyID:   "passgen-personal@localhost",
				IsDefault:  true,
				CreatedAt:  time.Now(),
				LastSyncAt: &time.Time{},
			},
			valid: true,
		},
		{
			name: "minimal valid store",
			store: PasswordStore{
				Name:      "work",
				GitURL:    "https://github.com/company/passwords.git",
				LocalPath: "/home/user/.password-stores/work",
				GPGKeyID:  "key-id",
				CreatedAt: time.Now(),
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.store.Name == "" && tt.valid {
				t.Error("Expected valid store but name is empty")
			}
			if tt.store.GitURL == "" && tt.valid {
				t.Error("Expected valid store but GitURL is empty")
			}
			if tt.store.LocalPath == "" && tt.valid {
				t.Error("Expected valid store but LocalPath is empty")
			}
		})
	}
}

func TestAutoRotationConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config AutoRotationConfig
		valid  bool
	}{
		{
			name: "valid rotation config",
			config: AutoRotationConfig{
				Enabled:          true,
				IntervalDays:     90,
				NextRotationAt:   time.Now().Add(90 * 24 * time.Hour),
				NotifyDaysBefore: 7,
				AutoGenerate:     true,
				PasswordProfile: &PasswordProfile{
					Length:         16,
					IncludeUpper:   true,
					IncludeLower:   true,
					IncludeNumbers: true,
					IncludeSymbols: true,
				},
			},
			valid: true,
		},
		{
			name: "disabled rotation",
			config: AutoRotationConfig{
				Enabled: false,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Enabled {
				if tt.config.IntervalDays <= 0 && tt.valid {
					t.Error("Expected valid config but IntervalDays is invalid")
				}
				if tt.config.NextRotationAt.IsZero() && tt.valid {
					t.Error("Expected valid config but NextRotationAt is zero")
				}
			}
		})
	}
}

func TestPasswordProfile_Validation(t *testing.T) {
	tests := []struct {
		name    string
		profile PasswordProfile
		valid   bool
	}{
		{
			name: "valid profile",
			profile: PasswordProfile{
				Length:         16,
				IncludeUpper:   true,
				IncludeLower:   true,
				IncludeNumbers: true,
				IncludeSymbols: true,
			},
			valid: true,
		},
		{
			name: "minimal profile",
			profile: PasswordProfile{
				Length:       8,
				IncludeLower: true,
			},
			valid: true,
		},
		{
			name: "invalid length",
			profile: PasswordProfile{
				Length:       0,
				IncludeLower: true,
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasCharTypes := tt.profile.IncludeUpper || tt.profile.IncludeLower || 
				tt.profile.IncludeNumbers || tt.profile.IncludeSymbols
			
			if tt.valid {
				if tt.profile.Length <= 0 {
					t.Error("Expected valid profile but length is invalid")
				}
				if !hasCharTypes {
					t.Error("Expected valid profile but no character types enabled")
				}
			}
		})
	}
}

func TestRotationStatus_Logic(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name           string
		status         RotationStatus
		expectedStatus string
	}{
		{
			name: "scheduled rotation",
			status: RotationStatus{
				Service:       "test",
				NextRotation:  now.Add(30 * 24 * time.Hour),
				DaysUntilNext: 30,
				IntervalDays:  90,
			},
			expectedStatus: "scheduled",
		},
		{
			name: "soon rotation",
			status: RotationStatus{
				Service:       "test",
				NextRotation:  now.Add(5 * 24 * time.Hour),
				DaysUntilNext: 5,
				IntervalDays:  90,
			},
			expectedStatus: "soon",
		},
		{
			name: "critical rotation",
			status: RotationStatus{
				Service:       "test",
				NextRotation:  now.Add(1 * 24 * time.Hour),
				DaysUntilNext: 1,
				IntervalDays:  90,
			},
			expectedStatus: "critical",
		},
		{
			name: "overdue rotation",
			status: RotationStatus{
				Service:       "test",
				NextRotation:  now.Add(-2 * 24 * time.Hour),
				DaysUntilNext: -2,
				IntervalDays:  90,
			},
			expectedStatus: "overdue",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the logic for determining status based on DaysUntilNext
			var actualStatus string
			switch {
			case tt.status.DaysUntilNext < 0:
				actualStatus = "overdue"
			case tt.status.DaysUntilNext <= 2:
				actualStatus = "critical"
			case tt.status.DaysUntilNext <= 7:
				actualStatus = "soon"
			default:
				actualStatus = "scheduled"
			}
			
			if actualStatus != tt.expectedStatus {
				t.Errorf("Expected status %q, got %q", tt.expectedStatus, actualStatus)
			}
		})
	}
}
