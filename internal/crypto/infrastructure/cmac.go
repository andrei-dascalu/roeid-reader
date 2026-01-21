package infrastructure

// CMACProvider wraps AES-CMAC operations
// Uses golang.org/x/crypto/cmac
type CMACProvider struct {
	// Key will be stored here during session
}

// NewCMACProvider creates a new CMAC provider
func NewCMACProvider() *CMACProvider {
	return &CMACProvider{}
}

// TODO: Implement Compute(key, data) -> []byte
