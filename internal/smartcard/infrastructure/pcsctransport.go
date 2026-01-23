package infrastructure

import (
	"strings"

	"github.com/andrei-dascalu/roeid-reader/internal/smartcard/domain"
	"github.com/ebfe/scard"
)

const LeaveCard scard.Disposition = scard.LeaveCard

// PCSCTransport implements smart card communication via PC/SC
type PCSCTransport struct {
	context        *scard.Context
	card           *scard.Card
	logger         *APDULogger
	selectedReader string
}

// NewPCSCTransport creates a new PC/SC transport
func NewPCSCTransport() *PCSCTransport {
	return &PCSCTransport{}
}

// SetLogger sets the logger for transport events
func (t *PCSCTransport) SetLogger(logger *APDULogger) {
	t.logger = logger
}

// Connect establishes a PC/SC context and connects to a card
func (t *PCSCTransport) Connect() error {
	ctx, err := scard.EstablishContext()
	if err != nil {
		return domain.NewTransportError(domain.ErrNoContext,
			"failed to establish PC/SC context", err)
	}
	t.context = ctx

	// List available readers
	readers, err := ctx.ListReaders()
	if err != nil {
		t.context.Release()
		t.context = nil
		return domain.NewTransportError(domain.ErrNoReaders,
			"failed to list readers", err)
	}

	if len(readers) == 0 {
		t.context.Release()
		t.context = nil
		return domain.NewTransportError(domain.ErrNoReaders,
			"no smart card readers found", nil)
	}

	// Select reader (prefer "Generic EMV" readers)
	t.selectedReader = t.selectReader(readers)
	if t.logger != nil {
		t.logger.LogReaderSelected(t.selectedReader, readers)
	}

	// Connect to card
	card, err := ctx.Connect(t.selectedReader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		t.context.Release()
		t.context = nil
		return domain.NewTransportError(domain.ErrNoCard,
			"failed to connect to card", err)
	}
	t.card = card

	// Log successful connection with protocol info
	if t.logger != nil {
		status, err := t.Status()
		if err == nil {
			t.logger.LogConnect(t.selectedReader, status.ActiveProtocol)
			t.logger.LogATR(status.ATR)
		}
	}

	return nil
}

// Transmit sends an APDU and receives a response
func (t *PCSCTransport) Transmit(apdu *domain.APDU) (*domain.Response, error) {
	if t.card == nil {
		return nil, domain.NewTransportError(domain.ErrNoCard,
			"not connected to card", nil)
	}

	// Log outgoing command
	apduBytes := apdu.Bytes()
	if t.logger != nil {
		t.logger.LogCommand(apduBytes)
	}

	responseData, err := t.card.Transmit(apduBytes)
	if err != nil {
		transportErr := domain.NewTransportError(domain.ErrTransmissionFailed,
			"APDU transmission failed", err)
		if t.logger != nil {
			t.logger.LogError(transportErr)
		}
		return nil, transportErr
	}

	// Log incoming response
	if t.logger != nil {
		t.logger.LogResponse(responseData)
	}

	return domain.NewResponse(responseData), nil
}

// Disconnect closes the card connection
func (t *PCSCTransport) Disconnect() error {
	if t.logger != nil {
		t.logger.LogDisconnect()
	}
	if t.card != nil {
		t.card.Disconnect(LeaveCard)
		t.card = nil
	}
	if t.context != nil {
		t.context.Release()
		t.context = nil
	}
	t.selectedReader = ""
	return nil
}

// Status returns current card status (ATR and active protocol)
func (t *PCSCTransport) Status() (*domain.CardStatus, error) {
	if t.card == nil {
		return nil, domain.NewTransportError(domain.ErrNoCard,
			"not connected to card", nil)
	}

	status, err := t.card.Status()
	if err != nil {
		return nil, domain.NewTransportError(domain.ErrConnectionLost,
			"failed to get card status", err)
	}

	protocol := "T=0"
	if status.ActiveProtocol == scard.ProtocolT1 {
		protocol = "T=1"
	}

	return &domain.CardStatus{
		ATR:            status.Atr,
		ActiveProtocol: protocol,
		Reader:         t.selectedReader,
	}, nil
}

// selectReader chooses a reader, preferring "Generic EMV" readers
func (t *PCSCTransport) selectReader(readers []string) string {
	for _, reader := range readers {
		if strings.Contains(reader, "Generic EMV") {
			return reader
		}
	}
	// Fallback to first reader
	return readers[0]
}

// ListReaders returns available smart card readers
func (t *PCSCTransport) ListReaders() ([]string, error) {
	if t.context == nil {
		return nil, domain.NewTransportError(domain.ErrNoContext,
			"not connected to PC/SC", nil)
	}
	return t.context.ListReaders()
}
