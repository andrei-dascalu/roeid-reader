package domain

// Card represents a connected smart card instance
type Card interface {
	// Transmit sends an APDU command and receives response
	Transmit(apdu *APDU) (*Response, error)

	// Disconnect closes the card connection
	Disconnect() error

	// Status returns current card status
	Status() (*CardStatus, error)
}

// CardStatus holds information about a connected card
type CardStatus struct {
	ATR            []byte // Answer To Reset
	ActiveProtocol string // T=0, T=1, etc.
	Reader         string // Reader name
}
