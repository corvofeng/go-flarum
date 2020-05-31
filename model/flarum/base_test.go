package flarum

import (
	"testing"
)

func TestCreateResources(t *testing.T) {

	curUser := NewResource(ECurrentUser)
	if curUser.(*CurrentUser).Type != "user" {
		t.Errorf("Get wrong type CurrentUser")
	}
}
