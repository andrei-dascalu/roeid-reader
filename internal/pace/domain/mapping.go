package domain

import "math/big"

// MappedDomain represents elliptic curve domain parameters after nonce mapping
type MappedDomain struct {
	P *big.Int      // Field prime
	A *big.Int      // Curve coefficient A
	B *big.Int      // Curve coefficient B
	G *Point        // Generator point
	N *big.Int      // Order of generator
	H *big.Int      // Cofactor
}

// Point represents an elliptic curve point
type Point struct {
	X *big.Int
	Y *big.Int
}

// IsPointAtInfinity checks if point is the identity element
func (p *Point) IsPointAtInfinity() bool {
	return p.X == nil && p.Y == nil
}

// MappingType identifies the nonce mapping algorithm
type MappingType int

const (
MappingTypeGM MappingType = iota // Generic Mapping
MappingTypeIM                      // Integrated Mapping
)

// Mapping represents the result of nonce-to-domain-parameters mapping
type Mapping struct {
	Type   MappingType
	Domain *MappedDomain
}
