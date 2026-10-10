package service

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWTSecretFromEnv(t *testing.T) {
	const secret = "smoke-jwt-secret"
	t.Setenv("JWT_SECRET", secret)
	if string(currentJWTSecret()) != secret {
		t.Fatal("JWT_SECRET was not applied")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": "smoke-user",
		"group":  "user",
	})
	signed, err := token.SignedString(currentJWTSecret())
	if err != nil {
		t.Fatal(err)
	}
	uid, group, err := ParseJWT(signed)
	if err != nil {
		t.Fatal(err)
	}
	if uid != "smoke-user" || group != "user" {
		t.Fatalf("claims uid=%s group=%s", uid, group)
	}

	other := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"userId": "x", "group": "user"})
	bad, err := other.SignedString([]byte("different-secret"))
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := ParseJWT(bad); err == nil {
		t.Fatal("token signed with another secret was accepted")
	}
}
