package account

import (
	// "fmt"
	"tools/tree"
	"tools/data"
	"tools/debug"
	// "github.com/google/uuid"
)

// Announce Manager and Contribute to the changes in UI

type LogFuncAcc func(logType debug.LogType, context debug.LogContext)

type Account struct {
	Tree 				*tree.Tree	
	DataManager *data.Manager
	LogFunc 		LogFuncAcc
}

func (acc *Account) ProvideAccsInDomain() {

}

func (acc *Account) ProvideDetails() {

}

func (acc *Account) AddAccount(accountName string, password string, domainName string, notes string) {
	// Store Account
	// Do something with data here (Update AccountLoaded,  AddedList)
	// Need Validation for Domain Name
	accountID := acc.DataManager.HandleNewAccDet(accountName, password, domainName, notes)
	if !acc.Tree.HasDomain(domainName) {
		acc.Tree.AddNodeNChild(domainName, accountName, accountID) 
	} else{
		acc.Tree.AddChild(domainName, accountName, accountID)
	}
	// Store Password it is encrypted ofc in passwordList along with ID

	acc.LogFunc(debug.MasterCreateAcc, debug.LogContext{
		Username: accountName,
		Domain: domainName,
	})
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
