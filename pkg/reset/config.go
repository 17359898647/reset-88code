package reset

import (
	"fmt"
	"os"
)

// LoadConfig loads and validates configuration from environment variables
func LoadConfig() (*Config, error) {
	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, fmt.Errorf("未配置TOKEN环境变量，跳过重置操作")
	}

	return &Config{
		Token: token,
	}, nil
}

// MaskToken masks the token for logging (show first 8 chars only)
func MaskToken(token string) string {
	if len(token) <= 8 {
		return "****"
	}
	return token[:8] + "****"
}
