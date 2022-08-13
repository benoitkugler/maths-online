package client

import "encoding/json"

const (
	Invalid Binary = iota // Invalide
	And                   // Et
	Or                    // Ou
)

type SignSymbol uint8

const (
	Nothing        SignSymbol = iota //
	Zero                             // 0
	ForbiddenValue                   // ||
)

// By default a slice of SignSymbol is marshalled as string
// by Go, which is not recognized by Dart
func (s SignSymbol) MarshalJSON() ([]byte, error) { return json.Marshal(uint8(s)) }
