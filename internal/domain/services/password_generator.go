package services

import (
	"crypto/rand"
	"math/big"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

// Provides cryptographically secure password generation with configurable character sets
type PasswordGenerator struct {
	charsetManager *entities.CharacterSet
}

// Initializes secure password generation with character set management
func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{
		charsetManager: entities.NewCharacterSet(),
	}
}

// Creates cryptographically random password using specified configuration parameters
func (pg *PasswordGenerator) GeneratePassword(config entities.PasswordConfig) (entities.Password, error) {
	if err := config.Validate(); err != nil {
		return entities.Password{}, err
	}

	charset, err := pg.charsetManager.BuildCharset(config)
	if err != nil {
		return entities.Password{}, err
	}

	passwordBytes := make([]byte, config.Length)
	for i := range passwordBytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return entities.Password{}, entities.NewPasswordError("failed to generate random number: " + err.Error())
		}
		passwordBytes[i] = charset[num.Int64()]
	}

	return entities.NewPassword(string(passwordBytes)), nil
}

// Produces collection of unique secure passwords for batch generation needs
func (pg *PasswordGenerator) GenerateMultiplePasswords(config entities.PasswordConfig) ([]entities.Password, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	passwords := make([]entities.Password, config.Count)
	for i := 0; i < config.Count; i++ {
		password, err := pg.GeneratePassword(config)
		if err != nil {
			return nil, err
		}
		passwords[i] = password
	}

	return passwords, nil
}
