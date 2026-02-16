package data

import (
	"github.com/google/uuid"
)

// Deal with Local variable

type Credential struct {
	ID 				uuid.UUID
	Password 	string
}

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
type Storage struct {
	ID 				uuid.UUID
	Username	string
	Notes			string
	Domain 		string
}

func NewStorage(username, domain string) *Storage {
	return &Storage{
		ID:       uuid.New(),
		Username: username,
		Notes:		"Test Note bla bla",
		Domain:   domain,
	}
}

