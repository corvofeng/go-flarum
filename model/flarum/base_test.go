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
	curUser := NewResource(ECurrentUser)
	if curUser.Attributes.GetType() != "users" {
		t.Errorf("Get wrong type CurrentUser")
	}
	if curUser.GetType() != "users" {
		t.Errorf("Get wrong type CurrentUser")
	}
}

func TestBindRelation(t *testing.T) {
	diss := NewResource(EDiscussion)
	testType := "simple type"
	testID := uint64(993)

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
