package client

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
