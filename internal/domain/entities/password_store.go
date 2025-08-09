package entities

import (
	"time"
)

// PasswordStore represents a password store configuration
type PasswordStore struct {
	Name       string     `yaml:"name"`
	GitURL     string     `yaml:"git_url"`
	LocalPath  string     `yaml:"local_path"`
	GPGKeyID   string     `yaml:"gpg_key_id"`
	IsDefault  bool       `yaml:"is_default"`
	CreatedAt  time.Time  `yaml:"created_at"`
	LastSyncAt *time.Time `yaml:"last_sync_at,omitempty"`
}

// StoreConfig represents the global store configuration
type StoreConfig struct {
	DefaultStore     string                   `yaml:"default_store"`
	Stores           map[string]PasswordStore `yaml:"stores"`
	ConfigPath       string                   `yaml:"-"`
	DefaultRotation  *DefaultRotationConfig   `yaml:"default_rotation,omitempty"`
	Notifications    *NotificationConfig      `yaml:"notifications,omitempty"`
}

// DefaultRotationConfig defines default rotation settings for new passwords
type DefaultRotationConfig struct {
	IntervalDays     int              `yaml:"interval_days"`
	NotifyDaysBefore int              `yaml:"notify_days_before"`
	AutoGenerate     bool             `yaml:"auto_generate"`
	PasswordProfile  *PasswordProfile `yaml:"password_profile,omitempty"`
}

// NotificationConfig defines notification settings
type NotificationConfig struct {
	Enabled bool   `yaml:"enabled"`
	Email   string `yaml:"email,omitempty"`
	Webhook string `yaml:"webhook,omitempty"`
}

// RotationStatus represents the status of password rotation for display
type RotationStatus struct {
	Service       string    `json:"service"`
	NextRotation  time.Time `json:"next_rotation"`
	DaysUntilNext int       `json:"days_until_next"`
	Status        string    `json:"status"` // "scheduled", "soon", "critical", "overdue"
	IntervalDays  int       `json:"interval_days"`
}
