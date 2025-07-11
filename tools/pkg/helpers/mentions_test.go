package helpers

import (
	"reflect"
	"testing"
)

func TestExtractNicknames(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Hello @user1 and @user2!", []string{"user1", "user2"}},
		{"@user1 @user1 @user2", []string{"user1", "user2"}},
		{"No nicknames here", []string{}},
		{"@123numeric", []string{"123numeric"}},
		{"Test enters\n@yser\n@alex", []string{"yser", "alex"}},
	}

	for _, test := range tests {
		result := ExtractNicknames(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ExtractNicknames(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestReplaceNicknamesWithUserIDs(t *testing.T) {
	nicknameToID := map[string]int64{"user1": 101, "user2": 202}
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello @user1 and @user2!", "Hello @{101} and @{202}!"},
		{"@user1 @unknown", "@{101} @unknown"},
		{"No replacements here", "No replacements here"},
	}

	for _, test := range tests {
		result := ReplaceNicknamesWithUserIDs(test.input, nicknameToID)
		if result != test.expected {
			t.Errorf("ReplaceNicknamesWithUserIDs(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestExtractUserIDs(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{"Hello @{101} and @{202}!", []int64{101, 202}},
		{"@{303} @{101} @{303}", []int64{303, 101}},
		{"No user IDs here", []int64{}},
	}

	for _, test := range tests {
		result := ExtractUserIDs(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("ExtractUserIDs(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestReplaceUserIDsWithNicknames(t *testing.T) {
	userIDToNickname := map[int64]string{101: "user1", 202: "user2"}
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello @{101} and @{202}!", "Hello [[@user1]] and [[@user2]]!"},
		{"@{101} @{999} @{202}", "[[@user1]] @{999} [[@user2]]"},
		{"No replacements here", "No replacements here"},
		{"Test enters\n@{101}\n@alex", "Test enters\n[[@user1]]\n@alex"},
	}

	for _, test := range tests {
		result := ReplaceUserIDsWithNicknames(test.input, userIDToNickname)
		if result != test.expected {
			t.Errorf("ReplaceUserIDsWithNicknames(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}
