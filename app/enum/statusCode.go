package enum

type StatusCode int

const (
	// User not logged in
	UserNotLoggedIn  StatusCode = 10001 
	// Voter not logged in
	VoterNotLoggedIn StatusCode = 10002
)
