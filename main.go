package main

import (
	"fmt"
	// "os"
	"log"
  // "encoding/json"
	"tools/account"
	"tools/tree"
	"tools/data"
	"github.com/joho/godotenv"
	"github.com/google/uuid"
	"github.com/rivo/tview"
  "github.com/gdamore/tcell/v2"
)

// How to make data in Storage hard to dump the hash?

// Func orders:
// - RecognizeDomain
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
// account 
// = Manage accounts
// tree
// = Manage tree UI
// security - crypto, scan, auth 
// = Help with Cryptography, Scan Weak/Leaked Password, Authentication
// store - store, local 
// = store help with deal with storing data to file & local deal with local memory

// Global Func: RemoveDuplicates, RecognizeDomain, Log
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

// Global Var
var (
	app 			= tview.NewApplication()
	logText  	= tview.NewTextView()
	form 			= tview.NewForm()
	root 			= tview.NewTreeNode("Root")
	// Some Global list should be here (Domains)
) 

// Global func
func RemoveDuplicates[T comparable](array []T) []T {
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

func Log(log string) {
	go func() {
		app.QueueUpdateDraw(func() {
			fmt.Fprintln(logText, log)
		})
	}()
	// Remember to store the log to .log
}

func UpdateAll(){
	// Update the Value of Current seen Element if Info in Storage change
	// Use the Added/Deleted List to guide
}

func ChangeBox() {
	// Change to Account Details in DomainName/CategoryName
	// Change to Add Account
	// Change to List Domain in Root (1st)
	// Change to List Account in Domain 
	// Change to Table of Leaked and Weak Password
}


func main() {
	// Load ENV VAR 
	err := godotenv.Load()
	if err != nil {
			log.Fatal("Error loading .env file")
	}

	// Load Data here
	// password: os.Getenv("ZAX_PASS"),
	data := [...]data.Storage{
		{
			ID: uuid.New(),
			Username: "Zax",
			Domain: "adds.com",
		},
		{
			ID: uuid.New(),
			Username: "Mathew",
			Domain: "adds.com",
		},
		{
			ID: uuid.New(),
			Username: "Alex",
			Domain: "facebook.com",
		},
	}

	slice := data[:]

	// There should be an addedList to track and help update storage here

	domains := []string{}
	for _, c := range slice {
		domains = append(domains, c.Domain)
	}
	
	domains = RemoveDuplicates(domains) 
	
	accountHelper := account.Account{LogFunc: Log}
	treeHelper := tree.Tree{
		Root: root,
		NodesList: make(map[string]*tview.TreeNode),
		ChildNodeList: slice,
		LogFunc: Log,
	}
	// dataHelper

	// Banner
	tcross :=  tview.NewTextView().
	   SetText("TCROSS\n\n[aa] Add Account 	[ad] Add Domain 	[ac] Add Category		[sp] Scan A Weak Password	 [spa] Scan Weak All Passwords	 [l] View logs") // Should be change to Flex

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
	treeHelper.LoadTree() // Convert fixed [3]Storage to *[]Storage   
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
			accountHelper.AddAccount(&treeHelper, domain, account)
		})

	form.
		SetBorder(true).
		SetTitle(" Add Account [1] ").
		SetTitleAlign(tview.AlignLeft)
	
	// Account Details Box (Box 2)
	
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
			// Start to Store data here
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
