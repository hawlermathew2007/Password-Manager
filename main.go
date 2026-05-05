package main

import (
	// "encoding/json"
	"fmt"
	// "tools/account"
	// "tools/tree"
	// "tools/data"
	"tools/ui"
	"tools/debug"
	// "github.com/joho/godotenv"
	// "github.com/google/uuid"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// How to make data in Storage hard to dump the hash?
// Domain should be unique
// Should add Error handling
// Domain ID is kinda redundant
// Ya need to add errors handler for data that might not always be available

// Works have to be done first (holidays)
// - Be able to enter Master Password
// - The Vim keys should be implemented fully
// - View logs (Modal)
// - Complete the Account Details
// - The parent node (Domain) should list the accounts
// - Scan for Weak Password (Scan all / Single account) # Multi later (Pop-up table)
// - Store => Encryption process
// - Config file
// - Log activities
// - Validations: Domain unique, Error handling

// Roadmap:
// [x] Log out options
// [?] Fix the Vim issues for now & Make the Vim keys works like it shoule be
// [ ] Improve the workflow
// [?] Error Handling (=> Shows with Modal only)
// [x] Pop-up a Modal to view Logs
// [ ] Apply Password Policy while Entering Password
// [ ] List out the account in the Domain (Nodes)
// [ ] Research about Encryption and Secure Storage
//     [ ] Master Password to Encryption Key (AES-KDF)
//     [ ] Decrypt Vault (Use AES key)
//     [ ]
//     [ ] Typing Obfuscation (Protect Master Pwd when entering)
//     [ ] Clipbroad Clearing
//     [ ] Memory Clearing
//     [ ]
// [ ] Do the Secure part after research: Encryption & Secure Storage
// [ ] Logging files
// [ ] Config files (This is gonna be a nightmare ig)
// [ ] Secure Password Generator
// [ ] Autofill & Autosave (How does Bitdefender do that?)

// Maybe further Research
// [ ] MFA - IDK dude
// [ ] Special table that can help check the Data breaches/leaked and show what leaked and also other account that likely their

// Today Objectives:
// [*] Enter Master Pass
// [ ] Be able to navigate back to Add Account
// [ ] Fix Vim keys intercept issues (2nd)
// [ ] Pop the modal up first (store the pass - simulate an encryption key)
// [ ] Should add log out options (1st)
// [ ] Create a Modal for showing Errors
// [ ] List out account in Domain (Refer: Tree)

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

// Some Global list should be here (Domains)
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


func main() {

	// App Initial Setup
	app := tview.NewApplication()
	pages := tview.NewPages()

	pageManager := ui.Page{
		App: app,
		Page: pages,
		Mode: 1,
	}

	logger, err := debug.NewLogger(
		".log/access.log",
		".log/error.log",
		pageManager.LogToBox,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	pageManager.Logger = logger

	pageManager.ErrorModal = tview.NewModal().
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(_ int, _ string) {
			pageManager.Page.HidePage("error-modal")
		})

	pageManager.ErrHandler = debug.NewHandler(logger, pageManager.ErrorPopup)

	loginMaster := pageManager.NewLogin()
	pageManager.Page.
		AddPage("login", loginMaster, true, true).
		AddPage("error-modal", pageManager.ErrorModal, false, false)

	// Vim keys
	var keyBuffer string
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		var uiManager *ui.Item
    if pageManager.UIManager != nil {
        uiManager = pageManager.UIManager
    } else {
        uiManager = &ui.Item{}
    }

		switch event.Key() {
			case tcell.KeyEsc:
        pageManager.Mode = 0
				keyBuffer = ""
        app.SetFocus(nil) // unfocus everything
        return nil // swallow event
		}

		if pageManager.Mode == 0 {
			
			r := event.Rune()
			if r != 0 {
				keyBuffer += string(r)
				switch keyBuffer {
					case "q":
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
					case "lo":
						pageManager.Page.RemovePage("vault")
						pageManager.Page.SwitchToPage("login")
						pageManager.Mode = 1
						keyBuffer = ""
						return nil
					case "0":
						if uiManager.Tree != nil {
							app.SetFocus(uiManager.TreeBox)
							pageManager.Mode = 1
							keyBuffer = ""
						}
						return nil
					case "1":
						if uiManager.CurrentItem != nil {
							app.SetFocus(uiManager.CurrentItem)
							pageManager.Mode = 1
							keyBuffer = ""
						}
						return nil
				}
				if len(keyBuffer) > 2 {
						keyBuffer = string(keyBuffer[len(keyBuffer)-1])
				}
			}
		}
		return event
	})
	
	if err := app.SetRoot(pageManager.Page, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
