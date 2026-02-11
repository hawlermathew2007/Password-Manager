package main

import (
	"fmt"
	"os"
	"log"
  // "encoding/json"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
  "github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// How to make Credential Storage hard to dump?

// Func orders:
// - recognizeDomain
// - addAccount
// - deleteAccount
// - moveAccount
// - storeInfo
// - updateInfo
// - deleteInfo
// - loadInfo
// - loadTree
// - expandTree & unexpandTree
// - changeBox
// - authentication
// - encryptInfo & decryptInfo

// Modules:
// authentication
// encryption
// account 
// root 
// scan weak & leaked password

// Global Func: removeDuplicates, recognizeDomain, addtoLog
// Class
// Tree 
// - func: Provide Current Nodes & Childs (with domain index - map), add - delete, move, update Nodes & Childs, loadTree
// - attr: Domains list, Account list (with domain index - map)
// Acc 
// - func: addAccount, deleteAccount, updateAccount (update Account Details),  provideDetails, decryptPass, scanWeakPass
// - attr: Details (username, password, domain)
// Info 
// - func: Add - Remove - Update Domain/AccountName/Pass (Should update the current list too), loadInfo
// - attr: Current Data Dumps

type Credential struct {
	id 				uuid.UUID
	username	string
	password 	string
	domain 		string
}

// Global Var
var (
	app 			= tview.NewApplication()
	logText  	= tview.NewTextView()
	form 			= tview.NewForm()
	root 			= tview.NewTreeNode("Root")
	// Some Global list should be here (Domains)
) 

// Global func
func removeDuplicates[T comparable](array []T) []T {
	seen := make(map[T]bool)
	result := []T{}
	for _, item := range array {
			if !seen[item] {
					seen[item] = true
					result = append(result, item)
			}
	}
	return result
}   

func addtoLog(log string) {
	go func() {
		app.QueueUpdateDraw(func() {
			fmt.Fprintln(logText, log)
		})
	}()
}

func recognizeDomain(nodes []string, domain string) bool  {
	// Check if Domain exist
	for _, n := range nodes {
		if n == domain {
			return true
		}
	}
	return false
}

func changeBox() {
	// Change to Account Details in DomainName/CategoryName
	// Change to Add Account
	// Change to List Domain in Root (1st)
	// Change to List Account in Domain 
	// Change to Table of Leaked and Weak Password
}

func authentication() {
	
}

type Account struct {

}

func provideDetails() {

}

func decryptPass() {

}

func scanWeakPass() {

}

func addAccount(domainName string, accountName string, curr_nodes []string) {
	// Store Account
	// Recognize the Domain to auto categorize
	// Add account to Domain
	if !recognizeDomain(curr_nodes, domainName) {
		domain := tview.NewTreeNode(domainName)
		domain.AddChild(tview.NewTreeNode(accountName))
		root.AddChild(domain)
	} else{
		var domain *tview.TreeNode
		root.Walk(func(node, parent *tview.TreeNode) bool {
			if node.GetText() == domainName {
					domain = node
					return false // Stop traversal
			}
			return true
		})
		domain.AddChild(tview.NewTreeNode(accountName))
	}

	// Store Password

	addtoLog(fmt.Sprintf("Successfully added \"%s\" to %s domain", accountName, domainName))
}

func deleteAccount()  {
	
}

func updateAccount() {

}

func moveAccount() {
	// Accept x & p vim keys
	// Buffer to store file - Potential for Pwn?
	// Delete & Add Account func in here
}

type Tree struct {
	ChildNodeList []string
	NodeList 			[]string
	ChildList 			[]string
}

func loadTree(data *[]Credential) {
	// Getting Domain, Account and Password through JSON file

	// Create array for storing Domain in memory

	// Create Tree branch
	for _, cred := range *data {
		domain := tview.NewTreeNode(cred.domain) // Add recognizeDomain here
		domain.AddChild(tview.NewTreeNode(cred.username))
		root.AddChild(domain)
	}
}

func expandTree() {
	
}

func unexpandTree() {

}

func deleteNode() {

}

type Data struct {

}

func loadInfo() {
	// Acname Domain, 
}

func storeInfo() {

}

func decryptInfo() {

}

func encryptInfo() {

}

func main() {
	// Load ENV VAR 
	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading .env file")
	}

	// Load Data here
	data := [...]Credential{
		{
			id: uuid.New(),
			username: "Zax",
			password: os.Getenv("ZAX_PASS"),
			domain: "adds.com",
		},
		{
			id: uuid.New(),
			username: "Mathew",
			password: os.Getenv("MATHEW_PASS"),
			domain: "adds.com",
		},
		{
			id: uuid.New(),
			username: "Alex",
			password: os.Getenv("ALEX_PASS"),
			domain: "facebook.com",
		},
	}

	nodes := []string{}
	for _, c := range data {
		nodes = append(nodes, c.domain)
	}
	
	nodes = removeDuplicates(nodes) 
	// Banner
	tcross :=  tview.NewTextView().
	   SetText("TCROSS\n\n[aa] Add Account 	[ac] Add Category		[sp] Scan A Weak Password	 [spa] Scan Weak All Passwords	 [l] View logs") // Should be change to Flex

	// Instruction
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetText("Navigate with h, j, k, l\nSelect an account or a Category and press d to delete.")

	// Log Fields
	logText.
		SetDynamicColors(true).
		SetScrollable(true).
		SetText("Your logs and errors will appear here.\n").
		SetChangedFunc(func() {
			app.Draw()
		})

	log := tview.NewGrid().
		SetBorders(true).
		SetRows(0).
		SetColumns(0).
		AddItem(logText, 0, 0, 1, 1, 0, 0, false)

	// Creating the Tree Root
	slice := data[:]
	loadTree(&slice) // Convert fixed []Credential to *[]Credential   
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	
	tree.
		SetBorder(true).
		SetTitle(" Accounts Overview [0]").
		SetTitleAlign(tview.AlignLeft)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {

	})
	
	// Add Account Section (Box 1)
	form.
		AddInputField("Domain", "", 30, nil, nil).
		AddInputField("Account", "", 30, nil, nil).
		AddPasswordField("Password", "", 30, '*', nil).
		AddButton("Add", func() {
			domain := form.GetFormItemByLabel("Domain").(*tview.InputField).GetText()
			account := form.GetFormItemByLabel("Account").(*tview.InputField).GetText()
			// password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			addAccount(domain, account, nodes)
		})

	form.
		SetBorder(true).
		SetTitle(" Add Account [1] ").
		SetTitleAlign(tview.AlignLeft)
	
	// Account Details Box
	
	// Grid Layout
	grid := tview.NewGrid().
		SetRows(4, 0, 5).
		SetColumns(35,0,0).
		AddItem(tree, 1, 0, 1, 1, 25, 20, false).
		AddItem(form, 1, 1, 1, 2, 25, 50, true).
		AddItem(tcross, 0, 0, 1, 3, 3, 0, false).
		AddItem(text, 2, 0, 1, 2, 3, 0, false).
		AddItem(log, 2, 2, 1, 1, 3, 0, false)

	// Vim keys
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			app.Stop()
			return nil
		case '0':
			app.SetFocus(tree)
			return nil
		case '1':
			app.SetFocus(form)
			return nil
		}
		return event
	})
	
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
