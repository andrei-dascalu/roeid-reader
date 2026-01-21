package application

import (
	"fmt"

	cryptoDomain "github.com/andrei-dascalu/roeid-reader/internal/crypto/domain"
	domainMsg "github.com/andrei-dascalu/roeid-reader/internal/messaging/domain"
	smartcardDomain "github.com/andrei-dascalu/roeid-reader/internal/smartcard/domain"
)

// SecureMessagingService handles encryption and decryption of APDUs
type SecureMessagingService struct {
	kEnc *cryptoDomain.AESKey
	kMac *cryptoDomain.AESKey
	ssc  *domainMsg.SendSequenceCounter
}

// NewSecureMessagingService creates a new secure messaging service
// Requires K_enc and K_mac from PACE key derivation
func NewSecureMessagingService(kEnc, kMac *cryptoDomain.AESKey) *SecureMessagingService {
	return &SecureMessagingService{
		kEnc: kEnc,
		kMac: kMac,
		ssc:  domainMsg.NewSendSequenceCounter(),
	}
}

// Encrypt encrypts an APDU and computes its CMAC
func (s *SecureMessagingService) Encrypt(apdu *smartcardDomain.APDU) (*domainMsg.SecureMessage, error) {
	// Increment SSC before encryption
	s.ssc.Increment()

	// TODO: Implement AES-CBC encryption with K_enc
	// TODO: Compute CMAC(SSC || Encrypted) with K_mac
	// TODO: Return SecureMessage with encrypted data and CMAC

	return nil, fmt.Errorf("encryption not yet implemented")
}

// Decrypt verifies CMAC and decrypts a secure message
func (s *SecureMessagingService) Decrypt(secure *domainMsg.SecureMessage) (*smartcardDomain.APDU, error) {
	// TODO: Verify CMAC(SSC || EncryptedData) with K_mac
	// TODO: Decrypt EncryptedData with K_enc
	// TODO: Increment SSC
	// TODO: Return plaintext APDU

	return nil, fmt.Errorf("decryption not yet implemented")
}
