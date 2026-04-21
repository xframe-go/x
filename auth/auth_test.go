package auth

import (
	"testing"
)

func TestManager_GenerateToken(t *testing.T) {
	config := Config{
		Secret:     "test-secret-key",
		Expiration: 24,
	}
	manager := NewManager(config)

	token, err := manager.GenerateToken("user-code-123")
	if err != nil {
		t.Errorf("GenerateToken() error = %v", err)
	}
	if token == "" {
		t.Error("GenerateToken() returned empty token")
	}

	claims, err := manager.ParseToken(token)
	if err != nil {
		t.Errorf("ParseToken() error = %v", err)
	}
	if claims.ID != "user-code-123" {
		t.Errorf("ParseToken() UserID = %v, want %v", claims.ID, "user-code-123")
	}
}

func TestManager_ParseToken_Invalid(t *testing.T) {
	config := Config{
		Secret:     "test-secret-key",
		Expiration: 24,
	}
	manager := NewManager(config)

	_, err := manager.ParseToken("invalid-token")
	if err == nil {
		t.Error("ParseToken() should return error for invalid token")
	}
}
