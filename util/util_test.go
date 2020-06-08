package util

import (
	"testing"
)

func TestToken(t *testing.T) {
	t1 := GetNewToken()

	if !VerifyToken(t1, t1) {
		t.Error("Get wrong csrf")
	}
}
