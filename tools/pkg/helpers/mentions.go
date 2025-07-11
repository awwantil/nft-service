package helpers

import (
	"fmt"
	"regexp"
	"strconv"
)

// ExtractNicknames extracts all nicknames from the text.
func ExtractNicknames(text string) []string {
	re := regexp.MustCompile(`@([a-zA-Z0-9_]+)`)
	matches := re.FindAllStringSubmatch(text, -1)

	nicknameSet := make(map[string]struct{})
	for _, match := range matches {
		nicknameSet[match[1]] = struct{}{}
	}

	nicknames := make([]string, 0, len(nicknameSet))
	for nickname := range nicknameSet {
		nicknames = append(nicknames, nickname)
	}

	return nicknames
}

// ReplaceNicknamesWithUserIDs replaces @nickname with @{userID}.
func ReplaceNicknamesWithUserIDs(text string, nicknameToID map[string]int64) string {
	re := regexp.MustCompile(`@([a-zA-Z0-9_]+)`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		nickname := match[1:] // Remove @
		if userID, ok := nicknameToID[nickname]; ok {
			return fmt.Sprintf("@{%d}", userID)
		}
		return match
	})
}

// ExtractUserIDs extracts all user IDs from the text.
func ExtractUserIDs(text string) []int64 {
	re := regexp.MustCompile(`@\{(\d+)\}`)
	matches := re.FindAllStringSubmatch(text, -1)

	userIDSet := make(map[int64]struct{})
	for _, match := range matches {
		num, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			continue
		}
		userIDSet[num] = struct{}{}
	}

	userIDs := make([]int64, 0, len(userIDSet))
	for userID := range userIDSet {
		userIDs = append(userIDs, userID)
	}

	return userIDs
}

// ReplaceUserIDsWithNicknames replaces @{userID} with @nickname.
func ReplaceUserIDsWithNicknames(text string, userIDToNickname map[int64]string) string {
	re := regexp.MustCompile(`@\{(\d+)\}`)
	return re.ReplaceAllStringFunc(text, func(match string) string {
		userID := match[2 : len(match)-1] // Extract user ID
		num, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			return match
		}

		if nickname, ok := userIDToNickname[num]; ok {
			return fmt.Sprintf("[[@%s]]", nickname)
		}
		return match
	})
}
