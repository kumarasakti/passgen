package services

import (
	"testing"

	"github.com/kumarasakti/passgen/internal/domain/entities"
)

func TestWordPasswordGenerator_GenerateWordPassword(t *testing.T) {
	analyzer := NewPasswordAnalyzer()
	generator := NewWordPasswordGenerator(analyzer)

	tests := []struct {
		name       string
		word       string
		strategy   entities.TransformationStrategy
		complexity entities.ComplexityLevel
		wantErr    bool
	}{
		{
			name:       "valid word with leetspeak",
			word:       "password",
			strategy:   entities.StrategyLeetspeak,
			complexity: entities.ComplexityMedium,
			wantErr:    false,
		},
		{
			name:       "valid word with hybrid",
			word:       "security",
			strategy:   entities.StrategyHybrid,
			complexity: entities.ComplexityHigh,
			wantErr:    false,
		},
		{
			name:       "empty word",
			word:       "",
			strategy:   entities.StrategyMixedCase,
			complexity: entities.ComplexityLow,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern := entities.NewWordPattern(tt.word)
			pattern.Strategy = tt.strategy
			pattern.Complexity = tt.complexity

			password, err := generator.GenerateWordPassword(pattern)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateWordPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if password == "" {
					t.Error("Generated password should not be empty")
				}
				if password == tt.word {
					t.Error("Generated password should be different from original word")
				}
			}
		})
	}
}

func TestWordPasswordGenerator_GenerateMultipleWordPasswords(t *testing.T) {
	analyzer := NewPasswordAnalyzer()
	generator := NewWordPasswordGenerator(analyzer)

	word := "test"
	pattern := entities.NewWordPattern(word)
	pattern.Strategy = entities.StrategyHybrid
	pattern.Complexity = entities.ComplexityMedium
	count := 3

	passwords, err := generator.GenerateMultipleWordPasswords(pattern, count)
	if err != nil {
		t.Errorf("GenerateMultipleWordPasswords() error = %v", err)
		return
	}

	if len(passwords) != count {
		t.Errorf("Expected %d passwords, got %d", count, len(passwords))
	}

	// Check basic properties
	for i, password := range passwords {
		if password == "" {
			t.Errorf("Password %d should not be empty", i)
		}
		if password == word {
			t.Errorf("Password %d should be different from original word", i)
		}
	}
}

func TestWordPasswordGenerator_AnalyzeWordPassword(t *testing.T) {
	analyzer := NewPasswordAnalyzer()
	generator := NewWordPasswordGenerator(analyzer)

	word := "security"
	pattern := entities.NewWordPattern(word)
	pattern.Strategy = entities.StrategyHybrid
	pattern.Complexity = entities.ComplexityHigh

	password, err := generator.GenerateWordPassword(pattern)
	if err != nil {
		t.Errorf("GenerateWordPassword() error = %v", err)
		return
	}

	analysis, err := generator.AnalyzeWordPassword(password, word)
	if err != nil {
		t.Errorf("AnalyzeWordPassword() error = %v", err)
		return
	}

	if !analysis.WordBased {
		t.Error("Analysis should indicate word-based password")
	}

	if analysis.OriginalWord != word {
		t.Errorf("Original word = %v, want %v", analysis.OriginalWord, word)
	}

	if analysis.Password.Value != password {
		t.Errorf("Analysis password = %v, want %v", analysis.Password.Value, password)
	}

	if analysis.TransformationQuality == "" {
		t.Error("Transformation quality should not be empty")
	}
}
