package services

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

// PasswordGenerator handles secure password generation
type PasswordGenerator struct {
	charsetManager *entities.CharacterSet
}

// NewPasswordGenerator creates a new PasswordGenerator instance
func NewPasswordGenerator() *PasswordGenerator {
	return &PasswordGenerator{
		charsetManager: entities.NewCharacterSet(),
	}
}

// GeneratePassword creates a cryptographically secure password.
//
// Behavior depends on config.NoRepeat:
//   - false (default): samples each position independently with replacement from the
//     full charset. This maximizes entropy (length * log2(charsetSize)) and is the
//     recommended mode for general use.
//   - true: guarantees every enabled character type appears at least once (when length
//     permits), samples without replacement (no duplicate characters), and applies a
//     cryptographically secure Fisher-Yates shuffle. This trades ~2 bits of entropy for
//     resistance to pattern-based cracking and stricter policy compliance.
func (pg *PasswordGenerator) GeneratePassword(config entities.PasswordConfig) (entities.Password, error) {
	if err := config.Validate(); err != nil {
		return entities.Password{}, err
	}

	charset, err := pg.charsetManager.BuildCharset(config)
	if err != nil {
		return entities.Password{}, err
	}

	if config.NoRepeat {
		return pg.generateNoRepeat(config, charset)
	}
	return pg.generateStandard(config, charset)
}

// generateStandard samples each position independently with replacement from charset.
// This is the default path and maximizes raw entropy.
func (pg *PasswordGenerator) generateStandard(config entities.PasswordConfig, charset string) (entities.Password, error) {
	passwordBytes := make([]byte, config.Length)
	charsetMax := big.NewInt(int64(len(charset)))

	for i := range passwordBytes {
		num, err := rand.Int(rand.Reader, charsetMax)
		if err != nil {
			return entities.Password{}, entities.NewPasswordError("failed to generate random number: " + err.Error())
		}
		passwordBytes[i] = charset[num.Int64()]
	}

	return entities.NewPassword(string(passwordBytes)), nil
}

// generateNoRepeat produces a password with guaranteed character-type coverage and no
// duplicate characters, then securely shuffles the result.
func (pg *PasswordGenerator) generateNoRepeat(config entities.PasswordConfig, charset string) (entities.Password, error) {
	categories, err := pg.charsetManager.BuildCategories(config)
	if err != nil {
		return entities.Password{}, err
	}

	// No-duplicate constraint: length must not exceed available unique characters
	if config.Length > len(charset) {
		return entities.Password{}, entities.NewPasswordError(fmt.Sprintf(
			"password length %d exceeds available unique characters %d (reduce length, enable more character types, or disable --no-repeat)",
			config.Length, len(charset)))
	}

	result := make([]byte, 0, config.Length)
	used := make(map[byte]bool, config.Length)

	// 1. Guarantee: pick one character from each enabled category (if length permits)
	if config.Length >= len(categories) {
		for _, category := range categories {
			char, err := pickUniqueChar(category, used)
			if err != nil {
				return entities.Password{}, err
			}
			result = append(result, char)
			used[char] = true
		}
	}

	// 2. Fill remaining positions from the full charset without replacement
	remaining := config.Length - len(result)
	for i := 0; i < remaining; i++ {
		char, err := pickUniqueChar(charset, used)
		if err != nil {
			return entities.Password{}, err
		}
		result = append(result, char)
		used[char] = true
	}

	// 3. Cryptographically secure Fisher-Yates shuffle.
	// Without this, guaranteed-category characters would always appear at the start,
	// making the password structure predictable.
	if err := secureShuffle(result); err != nil {
		return entities.Password{}, entities.NewPasswordError("failed to shuffle password: " + err.Error())
	}

	return entities.NewPassword(string(result)), nil
}

// GenerateMultiplePasswords generates multiple unique passwords based on the configuration
func (pg *PasswordGenerator) GenerateMultiplePasswords(config entities.PasswordConfig) ([]entities.Password, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	passwords := make([]entities.Password, 0, config.Count)
	seen := make(map[string]bool, config.Count)

	// Allow multiple attempts to handle rare duplicate collisions,
	// but cap to prevent infinite loops for tiny password spaces.
	maxAttempts := config.Count * 10
	attempts := 0

	for len(passwords) < config.Count && attempts < maxAttempts {
		attempts++
		password, err := pg.GeneratePassword(config)
		if err != nil {
			return nil, err
		}
		if !seen[password.Value] {
			seen[password.Value] = true
			passwords = append(passwords, password)
		}
	}

	if len(passwords) < config.Count {
		return nil, entities.NewPasswordError("failed to generate unique passwords (password space too small for requested count)")
	}

	return passwords, nil
}

// pickUniqueChar selects a random character from charset that has not been used yet.
// Uses crypto/rand for cryptographic security.
func pickUniqueChar(charset string, used map[byte]bool) (byte, error) {
	// Collect available (unused) characters
	available := make([]byte, 0, len(charset))
	for i := 0; i < len(charset); i++ {
		c := charset[i]
		if !used[c] {
			available = append(available, c)
		}
	}

	if len(available) == 0 {
		return 0, entities.NewPasswordError("no available unique characters remaining in character set")
	}

	idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(available))))
	if err != nil {
		return 0, entities.NewPasswordError("failed to generate random number: " + err.Error())
	}

	return available[idx.Int64()], nil
}

// secureShuffle performs a Fisher-Yates shuffle using crypto/rand.
// This ensures that guaranteed-category characters are randomly distributed
// throughout the password rather than clustered at the beginning.
func secureShuffle(arr []byte) error {
	for i := len(arr) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		arr[i], arr[j.Int64()] = arr[j.Int64()], arr[i]
	}
	return nil
}
