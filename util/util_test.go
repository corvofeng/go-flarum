package util

import (
	"testing"
)

func TestToken(t *testing.T) {
	t1 := GetNewToken()
	t2 := GetNewToken()

	if VerifyToken(t1, t2) {
		t.Error("Get wrong csrf")
	}

	if !VerifyToken(t1, t1) {
		t.Error("Get wrong csrf")
	}
}
