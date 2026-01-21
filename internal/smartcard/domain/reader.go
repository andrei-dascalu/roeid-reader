package domain

// Reader represents a smart card reader
type Reader interface {
	// Name returns the reader name
	Name() string

	// IsConnected checks if a card is present
	IsConnected() bool
}

// ReaderList holds multiple readers
type ReaderList []Reader
