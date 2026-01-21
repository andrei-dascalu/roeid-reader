package domain

// Identity represents a personal identity record from CEI
type Identity struct {
	Surname        string
	GivenNames     string
	CNP            string // Cod Numeric Personal (national ID)
	DateOfBirth    string
	PlaceOfBirth   string
	Citizenship    string
	Gender         string
	DocumentNumber string
	Series         string
	IssueDate      string
	IssuePlace     string
	ExpiryDate     string
}

// IdentityRepository defines the interface for identity data access
type IdentityRepository interface {
	// Read fetches identity data from card via secure messaging
	Read() (*Identity, error)

	// Save persists identity data locally (optional)
	Save(identity *Identity) error
}
