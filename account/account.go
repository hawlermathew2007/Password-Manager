package account

import (
	"fmt"
	"tools/tree"
)

type Account struct {
	LogFunc 				func(string)
}

func (acc *Account) ProvideAccsInDomain() {

}

func (acc *Account) ProvideDetails() {

}

func (acc *Account) AddAccount(tree *tree.Tree, domainName string, accountName string) {
	// Store Account
	// Recognize the Domain to auto categorize
	// Add account to Domain
	if !tree.HasDomain(domainName) {
		tree.AddNodeNChild(domainName, accountName)
	} else{
		tree.AddChild(domainName, accountName)
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
