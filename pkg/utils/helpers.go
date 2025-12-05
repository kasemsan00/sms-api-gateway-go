package utils

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

const (
	// DefaultCharset is the default charset for random string generation
	DefaultCharset = "ABCDEFGHIJKLMOPQRSTUVWXYZabcdefghijklmopqrstuvwxyz"
	// AlphanumericCharset includes letters and numbers
	AlphanumericCharset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandomString generates a random string of specified length using the given charset
func GenerateRandomString(length int, charset string) string {
	if charset == "" {
		charset = DefaultCharset
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GenerateRoomName generates a random 6-character alphabetic room name
func GenerateRoomName() string {
	return GenerateRandomString(6, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

// GenerateLinkID generates a random 6-character link ID
func GenerateLinkID(charset string) string {
	if charset == "" {
		charset = DefaultCharset
	}
	return GenerateRandomString(6, charset)
}

// GenerateIdentity generates a 10-character identity string
func GenerateIdentity() string {
	return GenerateRandomString(10, AlphanumericCharset+"$")
}

// GenerateViewerIdentity generates a viewer identity
func GenerateViewerIdentity() string {
	return "viewer_" + GenerateIdentity()
}

// GenerateGuestName generates a guest username
func GenerateGuestName() string {
	return fmt.Sprintf("Guest-%d", rand.Intn(100))
}

// GenerateUserName generates a username
func GenerateUserName() string {
	return fmt.Sprintf("User-%d", rand.Intn(100))
}

// MD5Hash creates an MD5 hash of a string
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// FormatDateTime formats time to MySQL datetime format
func FormatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// FormatDateTimeNow returns current time in MySQL datetime format
func FormatDateTimeNow() string {
	return FormatDateTime(time.Now())
}

// AddDays adds days to a time and returns formatted string
func AddDays(t time.Time, days int) string {
	return FormatDateTime(t.AddDate(0, 0, days))
}

// ParseDateTime parses a MySQL datetime string
func ParseDateTime(s string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", s)
}

// NullStringValue returns the value of a sql.NullString or empty string
func NullStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

// NullIntValue returns the value of a sql.NullInt32 or 0
func NullIntValue(ni sql.NullInt32) int32 {
	if ni.Valid {
		return ni.Int32
	}
	return 0
}

// NullInt64Value returns the value of a sql.NullInt64 or 0
func NullInt64Value(ni sql.NullInt64) int64 {
	if ni.Valid {
		return ni.Int64
	}
	return 0
}

// NullFloat64Value returns the value of a sql.NullFloat64 or 0
func NullFloat64Value(nf sql.NullFloat64) float64 {
	if nf.Valid {
		return nf.Float64
	}
	return 0
}

// NullBoolValue returns the value of a sql.NullBool or false
func NullBoolValue(nb sql.NullBool) bool {
	if nb.Valid {
		return nb.Bool
	}
	return false
}

// NullTimeValue returns the value of a sql.NullTime or zero time
func NullTimeValue(nt sql.NullTime) time.Time {
	if nt.Valid {
		return nt.Time
	}
	return time.Time{}
}

// ToNullString converts a string to sql.NullString
func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

// ToNullInt32 converts an int to sql.NullInt32
func ToNullInt32(i int) sql.NullInt32 {
	return sql.NullInt32{
		Int32: int32(i),
		Valid: true,
	}
}

// ToNullFloat64 converts a float64 to sql.NullFloat64
func ToNullFloat64(f float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: f,
		Valid:   true,
	}
}

// ToNullBool converts a bool to sql.NullBool
func ToNullBool(b bool) sql.NullBool {
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}

// ToNullTime converts a time.Time to sql.NullTime
func ToNullTime(t time.Time) sql.NullTime {
	return sql.NullTime{
		Time:  t,
		Valid: !t.IsZero(),
	}
}

// Pointer helpers

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to an int
func IntPtr(i int) *int {
	return &i
}

// Int32Ptr returns a pointer to an int32
func Int32Ptr(i int32) *int32 {
	return &i
}

// Float64Ptr returns a pointer to a float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}
