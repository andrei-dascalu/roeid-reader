package domain

import "math/big"

// Point represents a point on an elliptic curve
type Point struct {
	X *big.Int
	Y *big.Int
}

// NewPoint creates a new point
func NewPoint(x, y *big.Int) *Point {
	return &Point{X: x, Y: y}
}

// IsPointAtInfinity checks if this is the point at infinity
func (p *Point) IsPointAtInfinity() bool {
	return p == nil || (p.X == nil && p.Y == nil)
}

// EllipticCurve defines the interface for elliptic curve operations
type EllipticCurve interface {
	Name() string
	P() *big.Int
	A() *big.Int
	B() *big.Int
	G() *Point
	Order() *big.Int
	ScalarMult(k *big.Int, p *Point) *Point
	Add(p1, p2 *Point) *Point
}

// AESKey represents a symmetric AES key with secure memory handling
type AESKey struct {
	key []byte
}

// NewAESKey creates an AES key from bytes
func NewAESKey(keyBytes []byte) *AESKey {
	key := make([]byte, len(keyBytes))
	copy(key, keyBytes)
	return &AESKey{key: key}
}

// Bytes returns the key material (creates a copy for safety)
func (k *AESKey) Bytes() []byte {
	result := make([]byte, len(k.key))
	copy(result, k.key)
	return result
}

// Clear securely erases the key material
func (k *AESKey) Clear() {
	for i := range k.key {
		k.key[i] = 0
	}
}

// KDF defines the interface for key derivation
type KDF interface {
	Derive(context []byte, counter uint32, length int) ([]byte, error)
}
