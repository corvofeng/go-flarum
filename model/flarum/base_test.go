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
	testEmail := "hello@flarum"
	data := (curUser.Attributes).(*CurrentUser)
	data.Email = testEmail

	if curUser.Attributes.GetType() != "users" {
		t.Errorf("Get wrong type CurrentUser")
	}
	if curUser.GetType() != "users" {
		t.Errorf("Get wrong type CurrentUser")
	}
	// attr, err := curUser.Attributes.GetAttributes()
	attr, err := curUser.GetAttributes()
	v := attr["attributes"].(map[string]interface{})
	if err != nil {
		t.Error(err)
	}
	if v["email"] != testEmail {
		t.Errorf("Get wrong attributes")
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
