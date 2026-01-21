package application

import (
	"fmt"

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

// DiagnoseCard reads and displays card information
func (s *SmartCardService) DiagnoseCard() error {
	status, err := s.transport.Status()
	if err != nil {
		return fmt.Errorf("failed to get card status: %w", err)
	}

	fmt.Printf("Card ATR: %02X\n", status.ATR)
	fmt.Printf("Active Protocol: %s\n", status.ActiveProtocol)
	return nil
}

// SelectApplication sends SELECT APDU to activate an application
func (s *SmartCardService) SelectApplication(aid []byte) error {
	apdu := &domain.APDU{
		CLA:  0x00,
		INS:  0xA4,
		P1:   0x04,
		P2:   0x00,
		Data: aid,
		Le:   0x00,
	}

	resp, err := s.transport.Transmit(apdu)
	if err != nil {
		s.logger.LogError(err)
		return err
	}

	if !resp.IsSuccess() {
		err := domain.NewStatusError(resp)
		s.logger.LogError(err)
		return err
	}

	fmt.Printf("SELECT Application Response: %02X\n", resp.Data)
	return nil
}

// VerifyPIN sends VERIFY APDU to authenticate with PIN
func (s *SmartCardService) VerifyPIN(pin []byte) error {
	if len(pin) > 8 {
		pin = pin[:8]
	}

	apdu := &domain.APDU{
		CLA:  0x00,
		INS:  0x20,
		P1:   0x00,
		P2:   0x01,
		Data: pin,
	}

	resp, err := s.transport.Transmit(apdu)
	if err != nil {
		s.logger.LogError(err)
		return err
	}

	if !resp.IsSuccess() {
		err := domain.NewStatusError(resp)
		s.logger.LogError(err)
		return err
	}

	fmt.Printf("PIN Verification Response: %02X\n", resp.Data)
	return nil
}

// Transmit sends a raw APDU
func (s *SmartCardService) Transmit(apdu *domain.APDU) (*domain.Response, error) {
	s.logger.LogCommand(apdu.Bytes())
	resp, err := s.transport.Transmit(apdu)
	if err != nil {
		s.logger.LogError(err)
		return nil, err
	}
	fullResp := append(resp.Data, resp.SW1, resp.SW2)
	s.logger.LogResponse(fullResp)
	return resp, nil
}
