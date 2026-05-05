package data

import (
	"github.com/google/uuid"
)

// Deal with Changes data locally

type Manager struct {
	ModifiedList 				[]string // Added, Deleted, Updated
	LoadedAccountsList	map[uuid.UUID]AccountDetails
	LoadedCredsList	 		map[uuid.UUID]string
}

func TrackLoadedList(accounts []AccountDetails, creds []Credential) *Manager {
	manager := &Manager{
		ModifiedList:       []string{}, // This is act as instruction for Store to act
		LoadedAccountsList: make(map[uuid.UUID]AccountDetails),
		LoadedCredsList:    make(map[uuid.UUID]string),
	}

	for _, cred := range creds {
		manager.LoadedCredsList[cred.ID] = cred.Password
	}

	for _, account := range accounts {
		manager.LoadedAccountsList[account.ID] = account
	}

	return manager
}

func (manager *Manager) HandleNewAccDet(username string, pwd string, domain string, notes string) uuid.UUID {
	accountID := uuid.New()
	manager.LoadedAccountsList[accountID] = AccountDetails{
		ID: accountID,
		Username: username,
		Domain: domain,
		Notes: notes,
	}
	manager.LoadedCredsList[accountID] = pwd
	return accountID
}
// Function that can identify data changes and update the ModifiedList and LoadedList
