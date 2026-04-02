package utils

import "database/sql"

// NullString returns an interface that is nil if the string is empty,
// which gets correctly handled by pgx as a NULL value in the database.
func NullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// ToNullString converts a string to sql.NullString
func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// StringValue returns the value of a string pointer, or an empty string if nil.
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
