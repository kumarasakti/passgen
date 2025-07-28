package entities

import (
	"testing"
)

func TestTransformationStrategy_String(t *testing.T) {
	tests := []struct {
		strategy TransformationStrategy
		expected string
	}{
		{StrategyLeetspeak, "leetspeak"},
		{StrategyMixedCase, "mixed-case"},
		{StrategySuffix, "suffix"},
		{StrategyPrefix, "prefix"},
		{StrategyInsert, "insert"},
		{StrategyHybrid, "hybrid"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := string(tt.strategy); got != tt.expected {
				t.Errorf("TransformationStrategy = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestComplexityLevel_String(t *testing.T) {
	tests := []struct {
		level    ComplexityLevel
		expected string
	}{
		{ComplexityLow, "low"},
		{ComplexityMedium, "medium"},
		{ComplexityHigh, "high"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := string(tt.level); got != tt.expected {
				t.Errorf("ComplexityLevel = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestWordPattern_Creation(t *testing.T) {
	tests := []struct {
		name     string
		word     string
		expected string
	}{
		{
			name:     "simple word",
			word:     "password",
			expected: "password",
		},
		{
			name:     "word with spaces",
			word:     "  test  ",
			expected: "test",
		},
		{
			name:     "mixed case word",
			word:     "TeSt",
			expected: "test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp := NewWordPattern(tt.word)
			if wp.Word != tt.expected {
				t.Errorf("NewWordPattern().Word = %v, want %v", wp.Word, tt.expected)
			}
			
			// Check default values
			if wp.Strategy != StrategyHybrid {
				t.Errorf("Expected default strategy to be hybrid")
			}
			if wp.Complexity != ComplexityMedium {
				t.Errorf("Expected default complexity to be medium")
			}
		})
	}
}
