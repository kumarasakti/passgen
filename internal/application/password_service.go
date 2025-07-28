package application

import (
	"github.com/kumarasakti/passgen/internal/domain/entities"
	"github.com/kumarasakti/passgen/internal/domain/services"
)

// GeneratePasswordRequest represents a request to generate passwords
type GeneratePasswordRequest struct {
	Config entities.PasswordConfig
}

// GeneratePasswordResponse represents the response from password generation
type GeneratePasswordResponse struct {
	Passwords []entities.Password
	Analyses  []services.PasswordAnalysis
}

// GenerateWordPasswordRequest represents a request to generate word-based passwords
type GenerateWordPasswordRequest struct {
	Word       string
	Strategy   entities.TransformationStrategy
	Complexity entities.ComplexityLevel
	Count      int
}

// GenerateWordPasswordResponse represents the response from word-based password generation
type GenerateWordPasswordResponse struct {
	Passwords []string
	Analyses  []services.PasswordAnalysis
	Pattern   entities.WordPattern
}

// CheckPasswordRequest represents a request to check password strength
type CheckPasswordRequest struct {
	Password string
}

// CheckPasswordResponse represents the response from password strength checking
type CheckPasswordResponse struct {
	Result services.StrengthCheckResult
}

// PasswordService orchestrates password-related operations
type PasswordService struct {
	generator             *services.PasswordGenerator
	analyzer              *services.PasswordAnalyzer
	strengthChecker       *services.PasswordStrengthChecker
	wordPasswordGenerator *services.WordPasswordGenerator
}

// NewPasswordService creates a new PasswordService instance
func NewPasswordService() *PasswordService {
	analyzer := services.NewPasswordAnalyzer()
	return &PasswordService{
		generator:             services.NewPasswordGenerator(),
		analyzer:              analyzer,
		strengthChecker:       services.NewPasswordStrengthChecker(),
		wordPasswordGenerator: services.NewWordPasswordGenerator(analyzer),
	}
}

// GeneratePasswords generates passwords and provides analysis
func (ps *PasswordService) GeneratePasswords(req GeneratePasswordRequest) (GeneratePasswordResponse, error) {
	if err := req.Config.Validate(); err != nil {
		return GeneratePasswordResponse{}, err
	}

	passwords, err := ps.generator.GenerateMultiplePasswords(req.Config)
	if err != nil {
		return GeneratePasswordResponse{}, err
	}

	analyses := make([]services.PasswordAnalysis, len(passwords))
	for i, password := range passwords {
		analyses[i] = ps.analyzer.AnalyzePassword(password, req.Config)
	}

	return GeneratePasswordResponse{
		Passwords: passwords,
		Analyses:  analyses,
	}, nil
}

// CheckPasswordStrength checks the strength of a given password
func (ps *PasswordService) CheckPasswordStrength(req CheckPasswordRequest) CheckPasswordResponse {
	password := entities.NewPassword(req.Password)
	result := ps.strengthChecker.CheckPasswordStrength(password)

	return CheckPasswordResponse{
		Result: result,
	}
}

// GeneratePresetPassword generates a password using predefined presets
func (ps *PasswordService) GeneratePresetPassword(presetType string) (GeneratePasswordResponse, error) {
	config, err := ps.getPresetConfig(presetType)
	if err != nil {
		return GeneratePasswordResponse{}, err
	}

	return ps.GeneratePasswords(GeneratePasswordRequest{Config: config})
}

// getPresetConfig returns configuration for predefined presets
func (ps *PasswordService) getPresetConfig(presetType string) (entities.PasswordConfig, error) {
	switch presetType {
	case "secure":
		return entities.PasswordConfig{
			Length: 16, IncludeLower: true, IncludeUpper: true,
			IncludeNumbers: true, IncludeSymbols: true, Count: 1,
		}, nil
	case "simple":
		return entities.PasswordConfig{
			Length: 12, IncludeLower: true, IncludeUpper: true,
			IncludeNumbers: true, IncludeSymbols: false, Count: 1,
		}, nil
	case "pin":
		return entities.PasswordConfig{
			Length: 6, IncludeLower: false, IncludeUpper: false,
			IncludeNumbers: true, IncludeSymbols: false, Count: 1,
		}, nil
	case "alphanumeric":
		return entities.PasswordConfig{
			Length: 12, IncludeLower: true, IncludeUpper: true,
			IncludeNumbers: true, IncludeSymbols: false, Count: 1,
		}, nil
	default:
		return entities.PasswordConfig{}, entities.NewPasswordError("unknown preset: " + presetType)
	}
}

// GenerateWordPasswords generates word-based passwords
func (ps *PasswordService) GenerateWordPasswords(req GenerateWordPasswordRequest) (GenerateWordPasswordResponse, error) {
	// Create word pattern
	pattern := entities.NewWordPattern(req.Word)

	// Set strategy if provided
	if req.Strategy != "" {
		pattern.SetStrategy(req.Strategy)
	}

	// Set complexity if provided
	if req.Complexity != "" {
		pattern.SetComplexity(req.Complexity)
	}

	// Default count to 1 if not specified
	count := req.Count
	if count <= 0 {
		count = 1
	}

	// Generate passwords
	var passwords []string
	var err error

	if count == 1 {
		password, genErr := ps.wordPasswordGenerator.GenerateWordPassword(pattern)
		if genErr != nil {
			return GenerateWordPasswordResponse{}, genErr
		}
		passwords = []string{password}
	} else {
		passwords, err = ps.wordPasswordGenerator.GenerateMultipleWordPasswords(pattern, count)
		if err != nil {
			return GenerateWordPasswordResponse{}, err
		}
	}

	// Analyze each password
	analyses := make([]services.PasswordAnalysis, len(passwords))
	for i, password := range passwords {
		analysis, analyzeErr := ps.wordPasswordGenerator.AnalyzeWordPassword(password, req.Word)
		if analyzeErr != nil {
			return GenerateWordPasswordResponse{}, analyzeErr
		}
		analyses[i] = *analysis
	}

	return GenerateWordPasswordResponse{
		Passwords: passwords,
		Analyses:  analyses,
		Pattern:   *pattern,
	}, nil
}
