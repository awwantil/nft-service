package helpers

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StrToBool force convert string to boolean
func StrToBool(s string) bool {
	if s == "да" {
		return true
	}

	if v, err := strconv.ParseBool(s); err == nil {
		return v
	}

	return false
}

// StrToInt force convert string to int
func StrToInt(s string) int {
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return 0
}

// StrToInt64 force convert string to int64
func StrToInt64(s string) int64 {
	return int64(StrToInt(s))
}

// StrToUInt force convert string to uint
func StrToUInt64(s string) uint64 {
	return uint64(StrToInt(s))
}

// StrToFloat force convert string to float64
func StrToFloat(s string) float64 {
	if n, err := strconv.ParseFloat(s, 64); err == nil {
		return n
	}
	return float64(0.0)
}

// Abs return abs int value
func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

// Find returns the smallest index i at which x == a[i],
// or len(a) if there is no such index.
func Find(a []string, x string) int {
	for i, n := range a {
		if x == n {
			return i
		}
	}
	return -1
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

// Md5Hash calculate md5 hash of string
func Md5Hash(s string) string {
	data := []byte(s)
	return fmt.Sprintf("%x", md5.Sum(data))
}

// ParseTime parse string date with format
func ParseTime(value, format string) time.Time {
	res, _ := time.Parse(format, value)
	return res
}

// ParseNullTime parse string date with format and return sql.NullTime value
func ParseNullTime(value, format string) sql.NullTime {
	res, err := time.Parse(format, value)
	if err != nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Valid: true, Time: res}
}

// CreateNullTime parse string date with format and return sql.NullTime value
func CreateNullTime(value time.Time) sql.NullTime {
	return sql.NullTime{Time: value, Valid: true}
}

// CreateNullInt64 parse string date with format and return sql.NullTime value
func CreateNullInt64(value int) sql.NullInt64 {
	return sql.NullInt64{Int64: int64(value), Valid: true}
}

// CompareTimes compare two times string and return -1 0 1:
// time1 < time2 - -1
// time1 == time2 - 0
// time1 > time2 - 1
func CompareTimes(time1, time2 string) int {
	time1 = strings.Split(time1, ".")[0]
	time1Parts, err := StringArrayToInt(strings.Split(time1, ":"))
	if err != nil {
		return -2
	}
	time2 = strings.Split(time2, ".")[0]
	time2Parts, err := StringArrayToInt(strings.Split(time2, ":"))
	if err != nil {
		return -2
	}

	return CompareIntArrays(time1Parts, time2Parts)
}

// CompareIntArrays comparing array of integers
func CompareIntArrays(value1, value2 []int) int {
	compare := func(v1, v2 int) int {
		if v1 < v2 {
			return -1
		} else if v1 > v2 {
			return 1
		}
		return 0
	}

	for i := 0; i < len(value1); i++ {
		if v := compare(value1[i], value2[i]); v != 0 {
			return v
		}
	}
	return 0
}

// StringArrayToInt convert array of strings to array of ints
func StringArrayToInt(values []string) ([]int, error) {
	result := make([]int, 0)
	for _, i := range values {
		j, err := strconv.Atoi(i)
		if err != nil {
			return nil, err
		}
		result = append(result, j)
	}
	return result, nil
}

// MakeSearchKey remove from the string non-letter symbole
func MakeSearchKey(value string) string {
	re, err := regexp.Compile(`[^\p{Cyrillic}\W]`)
	if err != nil {
		return ""
	}
	result := re.ReplaceAllString(strings.ToLower(value), "")

	re, err = regexp.Compile(`[\W]`)
	if err != nil {
		return result
	}
	return re.ReplaceAllString(strings.ToLower(result), "")
}

// UniqueIds returns a new slice containing only unique elements from the input slice.
func UniqueIds(ids []int64) []int64 {
	seen := make(map[int64]struct{})
	result := make([]int64, 0, len(ids))
	for _, v := range ids {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
