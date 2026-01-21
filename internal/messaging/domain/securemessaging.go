package domain

// SendSequenceCounter tracks the message sequence number for secure messaging
type SendSequenceCounter struct {
	value uint64
}

// NewSendSequenceCounter initializes an SSC at 0
func NewSendSequenceCounter() *SendSequenceCounter {
	return &SendSequenceCounter{value: 0}
}

// Increment advances the counter (call before encryption)
func (s *SendSequenceCounter) Increment() {
	s.value++
}

// Bytes returns the SSC as 8-byte big-endian (for CMAC computation)
func (s *SendSequenceCounter) Bytes() [8]byte {
	var result [8]byte
	for i := 0; i < 8; i++ {
		result[i] = byte((s.value >> (56 - (i * 8))) & 0xFF)
	}
	return result
}

// Value returns the current counter value
func (s *SendSequenceCounter) Value() uint64 {
	return s.value
}

// SecureMessage represents an encrypted APDU with CMAC authentication
type SecureMessage struct {
	EncryptedData []byte // Encrypted APDU body
	CMAC          []byte // Authentication tag (8 bytes for AES-CMAC)
}

// NewSecureMessage creates a secure message
func NewSecureMessage(encrypted []byte, cmac []byte) *SecureMessage {
	return &SecureMessage{
		EncryptedData: encrypted,
		CMAC:          cmac,
	}
}
