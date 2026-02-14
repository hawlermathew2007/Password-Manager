package data

import (
	"github.com/google/uuid"
)

type Storage struct {
	ID 				uuid.UUID
	Username	string
	Domain 		string
}

func NewStorage(username, domain string) *Storage {
    return &Storage{
        ID:       uuid.New(),
        Username: username,
        Domain:   domain,
    }
}

func (s *Storage) GetDomain() string {
    return s.Domain
}
