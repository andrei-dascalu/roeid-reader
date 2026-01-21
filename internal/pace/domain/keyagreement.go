package domain

import "math/big"

// KeyAgreement represents ECDH key agreement state
type KeyAgreement struct {
	EphemeralPrivateKey *big.Int // Ks
	EphemeralPublicKey  *Point   // Qs
	CardPublicKey       *Point   // Qc
	SharedSecret        *big.Int // Z
}

// AuthenticationTag represents proof of session key possession
type AuthenticationTag struct {
	Tag []byte // CMAC or hash
}

// AuthenticationContext holds both sides' authentication tags
type AuthenticationContext struct {
	TerminalTag *AuthenticationTag
	CardTag     *AuthenticationTag
}

// IsAuthenticated checks if both tags match (mutual authentication)
func (ac *AuthenticationContext) IsAuthenticated() bool {
	if ac.TerminalTag == nil || ac.CardTag == nil {
		return false
	}
	if len(ac.TerminalTag.Tag) != len(ac.CardTag.Tag) {
		return false
	}
	// Constant-time comparison would be better
	for i := range ac.TerminalTag.Tag {
		if ac.TerminalTag.Tag[i] != ac.CardTag.Tag[i] {
			return false
		}
	}
	return true
}
