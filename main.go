package main

import (
	"fmt"
	"os"
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
// Should add Error handling
// Domain ID is kinda redundant
// Ya need to add errors handler for data that might not always be available

// Objectives: Be able to add account and access the account accurately
// - Fix the tree
// - Change the NewAccDetails in UI Module
// - Be able to change the ModifiedList
// - decryptPass (can be simulated not actually decrypt lol)

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
	// .log has access.log for what the user done, errors.log for the program errors
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

	IDadds := uuid.New()
	IDface := uuid.New()

	// Temporary in Load all Mode (Default should be Lazy Load)
	// There should be an addedList to track and help update storage here
	// data := EncryptData(encryptionKey)

	// Data will manage this
	// This should not be exist here but leave as temp
	creds := [...]data.Credential{ 
		{
			ID: IDz,
			Password: os.Getenv("ZAX_PASS"),
		},
		{
			ID: IDm,
			Password: os.Getenv("MATHEW_PASS"),
		},
		{
			ID: IDa,
			Password: os.Getenv("ALEX_PASS"),
		},
	}
	// This should be passed to security modules

	// AccountDetails no need to have ID, the map in manage module should have
	// Will be refined
	accounts := [...]data.AccountDetails{
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

	// Tree will manage this
	domainLoaded:= [...]data.DomainDetails{
		{
			ID: IDadds,
			DomainName: "adds.com",
			Usernames: []data.PartialDomainDetail{
				{
					ID: IDz,
					Username: "Zax",
				},
				{
					ID: IDm,
					Username: "Mathew",
				},
			},
		},
		{
			ID: IDface,
			DomainName: "facebook.com",
			Usernames: []data.PartialDomainDetail{
				{
					ID: IDa,
					Username: "Alex",
				},
			},
		},
	}

	domain := domainLoaded[:]

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
	treeHelper := tree.Tree{
		Root: root,
		NodesList: make(map[string]*tview.TreeNode),
		ChildsList: make(map[*tview.TreeNode]uuid.UUID),
		ChildNodeList: domain,
		LogFunc: Log,
	}
	// dataHelper
	treeHelper.LoadTree()
	dataManager := data.TrackLoadedList(accounts[:], creds[:])
	accountHelper := account.Account{
		Tree: &treeHelper,
		DataManager: dataManager,
		LogFunc: Log,
	}
	tree := uiHolder.NewTree(&treeHelper)
	addAccForm := uiHolder.NewAddAccForm()
	// accDetails := ui.NewAccDetails(&accountHelper) // Sample (will be removed once the Storage is refined)

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

	uiHolder.Account = &accountHelper
	uiHolder.Tree = &treeHelper
	uiHolder.DataManager = dataManager
	uiHolder.Grid = grid
	uiHolder.UsedItem = addAccForm
	uiHolder.CurrentItem = addAccForm
	uiHolder.AddAccountForm = addAccForm

	treeHelper.CreateChildContent	= uiHolder.NewAccDetails
	treeHelper.SelectChild = uiHolder.ChangeBox
	// uiHolder.AccountDetails = accDetails

	// Vim keys
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			// Start to Store data to file here
			// Problem: Cannot type q when entering text in field
			// Sol: Add this when the user is in the form. Press ESC to unfocus
			// textView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
			// 		if event.Rune() == 'w' {
			// 				return tcell.NewEventKey(tcell.KeyRune, 'j', tcell.ModNone)
			// 		}
			// 		return event
			// })
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
