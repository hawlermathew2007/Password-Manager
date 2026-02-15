package account

import (
	"fmt"
	"tools/tree"
)

type Account struct {
	Tree 			*tree.Tree	
	LogFunc 	func(string)
}

func (acc *Account) ProvideAccsInDomain() {

}

func (acc *Account) ProvideDetails() {

}

func (acc *Account) AddAccount(domainName string, accountName string) {
	// Store Account
	// Recognize the Domain to auto categorize
	// Add account to Domain
	if !acc.Tree.HasDomain(domainName) {
		acc.Tree.AddNodeNChild(domainName, accountName)
	} else{
		acc.Tree.AddChild(domainName, accountName)
	}

	// Store Password it is encrypted ofc in passwordList along with ID

	acc.LogFunc(fmt.Sprintf("Successfully added \"%s\" to %s domain", accountName, domainName))
}

func (acc *Account) DeleteAccount()  {
	
}

func (acc *Account) UpdateAccount() {

}

func (acc *Account) MoveAccount() {
	// Accept x & p vim keys
	// Buffer to store file - Potential for Pwn?
	// Delete & Add Account func in here
}
