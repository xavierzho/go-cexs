package constants

// AuthType represents the authentication level
type AuthType int

const (
	None AuthType = iota
	Keyed
	Signed
)
