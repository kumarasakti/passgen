package application

import (
	"testing"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

func TestPasswordService_CreatePasswordService(t *testing.T) {
	service := NewPasswordService()
	if service == nil {
		t.Error("PasswordService should not be nil")
	}
}

func TestGeneratePasswordRequest_ValidConfig(t *testing.T) {
	config := entities.PasswordConfig{
		Length:         12,
		IncludeLower:   true,
		IncludeUpper:   true,
		IncludeNumbers: true,
		IncludeSymbols: false,
		Count:          1,
	}

	request := GeneratePasswordRequest{
		Config: config,
	}

	if request.Config.Length != 12 {
		t.Errorf("Expected length 12, got %d", request.Config.Length)
	}

	if !request.Config.IncludeLower {
		t.Error("Expected IncludeLower to be true")
	}

	err := request.Config.Validate()
	if err != nil {
		t.Errorf("Valid config should not return error: %v", err)
	}
}

func TestGenerateWordPasswordRequest_ValidRequest(t *testing.T) {
	request := GenerateWordPasswordRequest{
		Word:       "security",
		Strategy:   entities.StrategyLeetspeak,
		Complexity: entities.ComplexityMedium,
		Count:      2,
	}

	if request.Word != "security" {
		t.Errorf("Expected word 'security', got %s", request.Word)
	}

	if request.Strategy != entities.StrategyLeetspeak {
		t.Errorf("Expected leetspeak strategy, got %s", request.Strategy)
	}

	if request.Complexity != entities.ComplexityMedium {
		t.Errorf("Expected medium complexity, got %s", request.Complexity)
	}

	if request.Count != 2 {
		t.Errorf("Expected count 2, got %d", request.Count)
	}
}

func TestGeneratePasswordResponse_Structure(t *testing.T) {
	password := entities.Password{Value: "test123"}
	passwords := []entities.Password{password}

	response := GeneratePasswordResponse{
		Passwords: passwords,
	}

	if len(response.Passwords) != 1 {
		t.Errorf("Expected 1 password, got %d", len(response.Passwords))
	}

	if response.Passwords[0].Value != "test123" {
		t.Errorf("Expected password 'test123', got %s", response.Passwords[0].Value)
	}
}

func TestGenerateWordPasswordResponse_Structure(t *testing.T) {
	passwords := []string{"s3cur1ty!", "p@ssw0rd123"}

	response := GenerateWordPasswordResponse{
		Passwords: passwords,
	}

	if len(response.Passwords) != 2 {
		t.Errorf("Expected 2 passwords, got %d", len(response.Passwords))
	}

	if response.Passwords[0] != "s3cur1ty!" {
		t.Errorf("Expected first password 's3cur1ty!', got %s", response.Passwords[0])
	}

	if response.Passwords[1] != "p@ssw0rd123" {
		t.Errorf("Expected second password 'p@ssw0rd123', got %s", response.Passwords[1])
	}
}
