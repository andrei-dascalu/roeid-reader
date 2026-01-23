package application

import (
	"github.com/andrei-dascalu/roeid-reader/internal/smartcard/domain"
	"github.com/andrei-dascalu/roeid-reader/internal/smartcard/infrastructure"
)

// SmartCardService orchestrates smart card operations
type SmartCardService struct {
	transport *infrastructure.PCSCTransport
	logger    *infrastructure.APDULogger
}

// NewSmartCardService creates a new smart card service
func NewSmartCardService(
	transport *infrastructure.PCSCTransport,
	logger *infrastructure.APDULogger,
) *SmartCardService {
	// Wire up logger to transport for automatic APDU logging
	transport.SetLogger(logger)
	return &SmartCardService{
		transport: transport,
		logger:    logger,
	}
}

// Connect establishes connection to a smart card
func (s *SmartCardService) Connect() error {
	return s.transport.Connect()
}

// Disconnect closes the smart card connection
func (s *SmartCardService) Disconnect() error {
	return s.transport.Disconnect()
}

// Status returns current card status (ATR, protocol, reader)
func (s *SmartCardService) Status() (*domain.CardStatus, error) {
	return s.transport.Status()
}

// SelectApplication sends SELECT APDU to activate an application (ISO/IEC 7816-4)
func (s *SmartCardService) SelectApplication(aid []byte) (*domain.Response, error) {
	apdu := &domain.APDU{
		CLA:  0x00, // ISO/IEC 7816-4: Inter-industry command
		INS:  0xA4, // SELECT
		P1:   0x04, // Select by DF name (AID)
		P2:   0x00, // First or only occurrence
		Data: aid,
		Le:   0x00, // Accept any response length
	}

	resp, err := s.transport.Transmit(apdu)
	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return resp, domain.NewStatusError(resp)
	}

	return resp, nil
}

// VerifyPIN sends VERIFY APDU to authenticate with PIN (ISO/IEC 7816-4)
func (s *SmartCardService) VerifyPIN(pin []byte, pinRef byte) error {
	// PIN is typically padded/truncated to 8 bytes for CEI cards
	if len(pin) > 8 {
		pin = pin[:8]
	}

	apdu := &domain.APDU{
		CLA:  0x00,   // ISO/IEC 7816-4: Inter-industry command
		INS:  0x20,   // VERIFY
		P1:   0x00,   // No information given
		P2:   pinRef, // PIN reference (e.g., 0x01 for PIN1)
		Data: pin,
	}

	resp, err := s.transport.Transmit(apdu)
	if err != nil {
		return err
	}

	if !resp.IsSuccess() {
		return domain.NewStatusError(resp)
	}

	return nil
}

// Transmit sends a raw APDU command (logging handled by transport)
func (s *SmartCardService) Transmit(apdu *domain.APDU) (*domain.Response, error) {
	return s.transport.Transmit(apdu)
}

// TransmitBytes sends raw APDU bytes and returns the full response
func (s *SmartCardService) TransmitBytes(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return nil, domain.NewTransportError(domain.ErrTransmissionFailed,
			"APDU too short: minimum 4 bytes required", nil)
	}

	apdu := &domain.APDU{
		CLA: data[0],
		INS: data[1],
		P1:  data[2],
		P2:  data[3],
	}

	// Parse optional Lc/Data/Le fields
	if len(data) > 4 {
		lc := int(data[4])
		if len(data) > 5+lc {
			apdu.Data = data[5 : 5+lc]
			if len(data) > 5+lc {
				apdu.Le = data[5+lc]
			}
		} else if len(data) == 5 {
			// Just Le, no data
			apdu.Le = data[4]
		} else {
			apdu.Data = data[5:]
		}
	}

	resp, err := s.transport.Transmit(apdu)
	if err != nil {
		return nil, err
	}

	// Return full response including status words
	return append(resp.Data, resp.SW1, resp.SW2), nil
}
