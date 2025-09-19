package helpers

import (
	"os"
	"testing"
)

// TestResolveConfigPath_Simple проверяет простую функцию ResolveConfigPath
func TestResolveConfigPath_Simple(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		defaultPath string
		expected    string
	}{
		{
			name:        "default path when no env",
			defaultPath: "config/default.yaml",
			expected:    "config/default.yaml",
		},
		{
			name:        "env variable takes precedence over default",
			envValue:    "config/from-env.yaml",
			defaultPath: "config/default.yaml",
			expected:    "config/from-env.yaml",
		},
		{
			name:        "empty env uses default",
			envValue:    "",
			defaultPath: "config/default.yaml",
			expected:    "config/default.yaml",
		},
		{
			name:        "env with whitespace is trimmed",
			envValue:    "  config/from-env.yaml  ",
			defaultPath: "config/default.yaml",
			expected:    "config/from-env.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем оригинальное значение
			originalEnv := os.Getenv("CONFIG_PATH")
			defer func() {
				if originalEnv != "" {
					os.Setenv("CONFIG_PATH", originalEnv)
				} else {
					os.Unsetenv("CONFIG_PATH")
				}
			}()

			// Устанавливаем тестовое значение
			if tt.envValue != "" {
				os.Setenv("CONFIG_PATH", tt.envValue)
			} else {
				os.Unsetenv("CONFIG_PATH")
			}

			result := ResolveConfigPath(tt.defaultPath)
			if result != tt.expected {
				t.Errorf("ResolveConfigPath() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// BenchmarkResolveConfigPath_Simple бенчмарк простой функции
func BenchmarkResolveConfigPath_Simple(b *testing.B) {
	os.Unsetenv("CONFIG_PATH")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ResolveConfigPath("config/default.yaml")
	}
}
