package flarum

import (
	"reflect"
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

func TestSetData(t *testing.T) {
	testID := uint64(993)
	diss := NewResource(EDiscussion, testID)
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "empty Resource slice",
			input:    []Resource{},
			expected: []Resource{},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "string input",
			input:    "test string",
			expected: "test string",
		},
		{
			name:     "integer input",
			input:    123,
			expected: 123,
		},
		{
			name:     "non-empty Resource slice",
			input:    []Resource{diss},
			expected: []Resource{diss},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiDoc := &APIDoc{}
			apiDoc.SetData(tt.input)
			if !reflect.DeepEqual(apiDoc.Data, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, apiDoc.Data)
			}
		})
	}
}
