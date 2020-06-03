package flarum

import (
	"testing"
)

func assertEqual(t *testing.T, a, b interface{}) {
	t.Helper()
	if a != b {
		t.Errorf("Not Equal. %d %d", a, b)
	}
}

func TestCreateResources(t *testing.T) {
	testID := uint64(993)
	curUser := NewResource(ECurrentUser, testID)
	if curUser.Attributes.GetType() != "users" {
		t.Errorf("Get wrong type CurrentUser")
	}
	if curUser.GetType() != "users" {
		t.Errorf("Get wrong type CurrentUser")
	}
}

func TestBindRelation(t *testing.T) {
	testID := uint64(993)
	diss := NewResource(EDiscussion, testID)
	testType := "simple type"

	posts := RelationArray{
		Data: []BaseRelation{
			InitBaseResources(testID, testType),
		},
	}
	diss.BindRelations("Posts", posts)
	newPosts := diss.Relationships.(*DiscussionRelations).Posts

	assertEqual(t, newPosts.Data[0].GetID(), testID)
	assertEqual(t, newPosts.Data[0].Type, testType)
}
