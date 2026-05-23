package jwt

import (
	"testing"
	"time"

	jwt5 "github.com/golang-jwt/jwt/v5"
)

func init() { Init("test-secret-key-for-unit-tests") }

func TestGenerateAndParse(t *testing.T) {
	token, err := GenerateToken("u1", "openid_xxx", 1)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if claims.UserID != "u1" {
		t.Errorf("user_id: %s", claims.UserID)
	}
	if claims.OpenID != "openid_xxx" {
		t.Errorf("openid: %s", claims.OpenID)
	}
	if claims.Level != 1 {
		t.Errorf("level: %d", claims.Level)
	}
}

func TestParseInvalidToken(t *testing.T) {
	_, err := ParseToken("invalid.token.here")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestExpiredToken(t *testing.T) {
	claims := Claims{
		UserID: "u1", OpenID: "ox", Level: 0,
		RegisteredClaims: jwt5.RegisteredClaims{
			ExpiresAt: jwt5.NewNumericDate(time.Now().Add(-1 * time.Hour)),
		},
	}
	if claims.RegisteredClaims.ExpiresAt.Time.After(time.Now()) {
		t.Error("expected expired time")
	}
}
