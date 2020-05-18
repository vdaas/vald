package usecase

import "strings"

type Operation uint8

const (
	Unknown Operation = iota
	Insert
	Search
)

func Atoo(s string) Operation {
	switch strings.ToLower(s) {
	case "insert": return Insert
	case "search": return Search
	default: return Unknown
	}
}

func (o Operation) String() string {
	switch o {
	case Insert: return "insert"
	case Search: return "search"
	default: return "unknown operation"
	}
}
