package infrastructure

import (
	"fmt"
	"log"
	"strings"

	"github.com/andrei-dascalu/roeid-reader/internal/smartcard/domain"
	"github.com/ebfe/scard"
)

const LeaveCard scard.Disposition = scard.LeaveCard

// PCSCTransport implements smart card communication via PC/SC
type PCSCTransport struct {
	context *scard.Context
	card    *scard.Card
}

// NewPCSCTransport creates a new PC/SC transport
func NewPCSCTransport() *PCSCTransport {
	return &PCSCTransport{}
}

// Connect establishes a PC/SC context and connects to a card
func (t *PCSCTransport) Connect() error {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return fmt.Errorf("failed to establish PC/SC context: %w", err)
	}
	t.context = ctx

	// List available readers
	readers, err := ctx.ListReaders()
	if err != nil {
		t.context.Release()
		return fmt.Errorf("failed to list readers: %w", err)
	}

	if len(readers) == 0 {
		t.context.Release()
		return fmt.Errorf("no smart card readers found")
	}

	// Select reader (prefer "Generic EMV" readers)
	selectedReader := t.selectReader(readers)

	// Connect to card
	card, err := ctx.Connect(selectedReader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		t.context.Release()
		return fmt.Errorf("failed to connect to card: %w", err)
	}

	t.card = card
	return nil
}

// Transmit sends an APDU and receives a response
func (t *PCSCTransport) Transmit(apdu *domain.APDU) (*domain.Response, error) {
	if t.card == nil {
		return nil, fmt.Errorf("not connected to card")
	}

	responseData, err := t.card.Transmit(apdu.Bytes())
	if err != nil {
		return nil, fmt.Errorf("transmission failed: %w", err)
	}

	return domain.NewResponse(responseData), nil
}

// Disconnect closes the card connection
func (t *PCSCTransport) Disconnect() error {
	if t.card != nil {
		t.card.Disconnect(LeaveCard)
		t.card = nil
	}
	if t.context != nil {
		t.context.Release()
		t.context = nil
	}
	return nil
}

// Status returns current card status (ATR and active protocol)
func (t *PCSCTransport) Status() (*domain.CardStatus, error) {
	if t.card == nil {
		return nil, fmt.Errorf("not connected to card")
	}

	status, err := t.card.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get card status: %w", err)
	}

	return &domain.CardStatus{
		ATR:            status.Atr,
		ActiveProtocol: fmt.Sprintf("T=%d", status.ActiveProtocol),
	}, nil
}

// selectReader chooses a reader, preferring "Generic EMV" readers
func (t *PCSCTransport) selectReader(readers []string) string {
	for _, reader := range readers {
		if strings.Contains(reader, "Generic EMV") {
			log.Printf("Selected reader: %s", reader)
			return reader
		}
	}
	// Fallback to first reader
	log.Printf("No Generic EMV reader found. Using: %s", readers[0])
	return readers[0]
}

// ListReaders returns available smart card readers
func (t *PCSCTransport) ListReaders() ([]string, error) {
	if t.context == nil {
		return nil, fmt.Errorf("not connected to PC/SC")
	}
	return t.context.ListReaders()
}
