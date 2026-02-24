package data

import (
	"github.com/google/uuid"
)

// Deal with Local variable

type Credential struct {
	ID 				uuid.UUID
	Password 	string
} // Used for Retreiving encrypted password using DecryptPass (temp)

type AccountDetails struct {
	ID 				uuid.UUID // User ID (temp)
	Username	string
	Domain 		string
	Notes			string
}

type PartialDomainDetail struct {
	ID 				uuid.UUID // Also User ID
	Username	string
}

type DomainDetails struct {
	ID 					uuid.UUID // Domain ID (may not be needed as domain will be unique)
	DomainName 	string // Validation needed (check if it actually a domain)
	Usernames	 	[]PartialDomainDetail
} // Needed for Tree

// Should be refined
// {
// 	IDs: DJHFEDOMAINSHH,
// 	DomainName: "adds.com",
// 	Accounts: [
// 		{
// 			ID:       DJHFEUSERSSHH,
// 			Username: Mathew,
// 			Notes:		"Test Note bla bla",
// 		}
// 	]
// }

