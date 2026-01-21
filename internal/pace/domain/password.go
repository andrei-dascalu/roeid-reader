package domain

// Password represents a user-provided PIN or CAN
type Password struct {
	value []byte
	typ   PasswordType
}

// PasswordType enum
type PasswordType int

const (
	PasswordTypePIN PasswordType = iota
	PasswordTypeCSAN
)

// NewPassword creates a Password
func NewPassword(value []byte, typ PasswordType) *Password {
	return &Password{
		value: make([]byte, len(value)),
		typ:   typ,
	}
}

// Bytes returns the raw password bytes
func (p *Password) Bytes() []byte {
	return p.value
}

// Clear securely erases the password from memory
func (p *Password) Clear() {
	for i := range p.value {
		p.value[i] = 0
	}
}

// Nonce represents a card-provided random value for PACE mapping
type Nonce struct {
	data []byte
}

// NewNonce creates a Nonce
func NewNonce(data []byte) *Nonce {
	return &Nonce{
		data: make([]byte, len(data)),
	}
}

// Bytes returns nonce raw data
func (n *Nonce) Bytes() []byte {
	return n.data
}
