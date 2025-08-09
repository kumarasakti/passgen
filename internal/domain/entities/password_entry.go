package entities

import (
	"time"
)

// PasswordEntry represents a stored password with metadata
type PasswordEntry struct {
	Service     string            `json:"service"`
	Username    string            `json:"username,omitempty"`
	Password    string            `json:"password"`
	URL         string            `json:"url,omitempty"`
	Notes       string            `json:"notes,omitempty"`
	Metadata    map[string]string `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	GeneratedBy string            `json:"generated_by"`

	// Auto-rotation features (optional)
	AutoRotation    *AutoRotationConfig `json:"auto_rotation,omitempty"`
	RotationHistory []RotationRecord    `json:"rotation_history,omitempty"`
}

// AutoRotationConfig defines automatic password rotation settings
type AutoRotationConfig struct {
	Enabled          bool             `json:"enabled"`
	IntervalDays     int              `json:"interval_days"`     // e.g., 30, 60, 90
	NextRotationAt   time.Time        `json:"next_rotation_at"`
	NotifyDaysBefore int              `json:"notify_days_before"` // e.g., 7 days warning
	AutoGenerate     bool             `json:"auto_generate"`      // true = auto-generate new password
	PasswordProfile  *PasswordProfile `json:"password_profile,omitempty"`
}

// PasswordProfile defines custom password generation rules for auto-rotation
type PasswordProfile struct {
	Length         int    `json:"length"`
	IncludeUpper   bool   `json:"include_upper"`
	IncludeLower   bool   `json:"include_lower"`
	IncludeNumbers bool   `json:"include_numbers"`
	IncludeSymbols bool   `json:"include_symbols"`
	CustomRules    string `json:"custom_rules,omitempty"` // e.g., "no-ambiguous"
}

// RotationRecord tracks password rotation history
type RotationRecord struct {
	RotatedAt    time.Time `json:"rotated_at"`
	PreviousHash string    `json:"previous_hash"` // SHA256 of old password (for audit)
	Reason       string    `json:"reason"`        // "auto-rotation", "manual", "breach"
	GeneratedBy  string    `json:"generated_by"`
}

// PasswordMetadata represents password information for display (no actual password)
type PasswordMetadata struct {
	Service      string            `json:"service"`
	Username     string            `json:"username,omitempty"`
	URL          string            `json:"url,omitempty"`
	Notes        string            `json:"notes,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	AutoRotation *AutoRotationInfo `json:"auto_rotation,omitempty"` // Only if enabled
	StrengthInfo string            `json:"strength_info"`
}

// AutoRotationInfo represents rotation information for display
type AutoRotationInfo struct {
	Enabled       bool      `json:"enabled"`
	IntervalDays  int       `json:"interval_days"`
	NextRotation  time.Time `json:"next_rotation"`
	DaysUntilNext int       `json:"days_until_next"`
}
