package interfaces

// StorageBackend defines the interface for different storage backends
type StorageBackend interface {
	// File operations
	SaveFile(key string, data []byte) error
	LoadFile(key string) ([]byte, error)
	ListFiles(prefix string) ([]string, error)
	DeleteFile(key string) error
	FileExists(key string) (bool, error)

	// Storage operations
	Initialize(storeName string) error
	IsInitialized(storeName string) (bool, error)

	// Sync operations (for backends that support it)
	Sync() error

	// Backend information
	GetBackendType() string
	GetConnectionInfo() map[string]string
}

// StorageConfig represents the configuration for a storage backend
type StorageConfig struct {
	Type       string            `yaml:"type" json:"type"`
	Settings   map[string]string `yaml:"settings" json:"settings"`
	Encryption string            `yaml:"encryption" json:"encryption"`
}
