package utils

import "testing"

func TestPasswordHash(t *testing.T) {
	hash, err := HashPassword("correct-password")
	if err != nil {
		t.Fatal(err)
	}
	if hash == "correct-password" {
		t.Fatal("password was stored as plain text")
	}
	if !CheckPasswordHash("correct-password", hash) {
		t.Fatal("correct password rejected")
	}
	if CheckPasswordHash("wrong-password", hash) {
		t.Fatal("wrong password accepted")
	}
}
