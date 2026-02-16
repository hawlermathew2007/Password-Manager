package main

import (
	"fmt"
	// "os"
	//"log"
  // "encoding/json"
	"tools/account"
	"tools/tree"
	"tools/data"
	"tools/ui"
	"github.com/joho/godotenv"
	"github.com/google/uuid"
	"github.com/rivo/tview"
  "github.com/gdamore/tcell/v2"
)

// How to make data in Storage hard to dump the hash?
// Domain should be unique

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
	app 				= tview.NewApplication()
	logText  		= tview.NewTextView()
	// Some Global list should be here (Domains)
) 

// Global func
// func RemoveDuplicates[T comparable](array []T) []T {
// 	seen := make(map[T]bool)
// 	result := []T{}
// 	for _, item := range array {
// 		if !seen[item] {
// 			seen[item] = true
// 			result = append(result, item)
// 		}
// 	}
// 	return result
// }   

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

func main() {

	root := tview.NewTreeNode("Root")

	// Load ENV VAR 
	err := godotenv.Load()
	if err != nil {
		Log("Error loading .emv file.")
	}

	// Load Data here
	IDz := uuid.New()
	IDm := uuid.New()
	IDa := uuid.New()

	// There should be an addedList to track and help update storage here

	// creds := [...]data.Credential{
	// 	{
	// 		ID: IDz,
	// 		Password: os.Getenv("ZAX_PASS"),
	// 	},
	// 	{
	// 		ID: IDm,
	// 		Password: os.Getenv("MATHEW_PASS"),
	// 	},
	// 	{
	// 		ID: IDa,
	// 		Password: os.Getenv("ALEX_PASS"),
	// 	},
	// }
	// This should be passed to security modules

	dataLoaded := [...]data.Storage{
		{
			ID: IDz,
			Username: "Zax",
			Notes: "My account for Active Directory Practice",
			Domain: "adds.com",
		},
		{
			ID: IDm,
			Username: "Mathew",
			Notes: "My account for Active Directory Practice",
			Domain: "adds.com",
		},
		{
			ID: IDa,
			Username: "Alex",
			Notes: "My Social Media Account",
			Domain: "facebook.com",
		},
	}

	_data := dataLoaded[:]

	treeHelper := tree.Tree{
		Root: root,
		NodesList: make(map[string]*tview.TreeNode),
		ChildsList: make(map[string][]string),
		ChildNodeList: _data,
		LogFunc: Log,
	}

	accountHelper := account.Account{
		Tree: &treeHelper,
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

	// Boxes
	uiHolder := ui.Item{}
	treeHelper.LoadTree()
	tree := uiHolder.NewTree(&treeHelper)
	addAccForm := ui.NewAddAccForm(&accountHelper)
	accDetails := ui.NewAccDetails(&accountHelper) // Sample (will be removed once the Storage is refined)

	// Add Pages for Auth, TCROSS dashboard, (Leaked Password Table)
	
	// Grid Layout
	grid := tview.NewGrid().
		SetRows(4, 0, 5).
		SetColumns(35,0,0).
		AddItem(tree, 1, 0, 1, 1, 25, 20, false).
		AddItem(addAccForm, 1, 1, 1, 2, 25, 50, true).
		AddItem(tcross, 0, 0, 1, 3, 3, 0, false).
		AddItem(text, 2, 0, 1, 2, 3, 0, false).
		AddItem(log, 2, 2, 1, 1, 3, 0, false)

	uiHolder.Grid = grid
	uiHolder.UsedItem = addAccForm
	uiHolder.CurrentItem = addAccForm
	uiHolder.AddAccount = addAccForm
	uiHolder.AccountDetails = accDetails

	// Vim keys
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			// Start to Store data to file here
			// Problem: Cannot type q when entering text in field
			app.Stop()
			return nil
		case '0':
			app.SetFocus(tree)
			return nil
		case '1':
			app.SetFocus(uiHolder.CurrentItem)
			return nil
		}
		return event
	})
	
	if err := app.SetRoot(grid, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
