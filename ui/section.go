package ui

import (
	"tools/account"
	"tools/data"
	"tools/debug"

	// "tools/security"
	"tools/tree"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
)

type Item struct {
	Tree 						*tree.Tree
	Account 				*account.Account
	DataManager			*data.Manager
	Grid 						*tview.Grid
	UsedItem 				tview.Primitive
	CurrentItem 		tview.Primitive
	TreeBox 				*tview.TreeView
	AddAccountForm 	*tview.Form
	AccountDetails 	*tview.Frame
	// ListDomains
	// ListAccs
}

// Account
func (item *Item) NewAddAccForm() *tview.Form {
	form := tview.NewForm()
	form.
		AddInputField("Domain", "", 30, nil, nil).
		AddInputField("Account", "", 30, nil, nil).
		AddPasswordField("Password", "", 30, '*', nil).
		AddTextArea("Notes", "", 45, 6, 30, nil).
		AddButton("Add", func() {
			domain := form.GetFormItemByLabel("Domain").(*tview.InputField).GetText()
			account := form.GetFormItemByLabel("Account").(*tview.InputField).GetText()
			password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			notes := form.GetFormItemByLabel("Notes").(*tview.TextArea).GetText()
			
			item.Account.AddAccount(account, password, domain, notes)
		}).
		AddButton("Scan", func() {
			item.Account.LogFunc(debug.PwdScanned, debug.LogContext{
				Username: "Test",
				Domain: "Test again",
				ScanResult: "Again",
			})
		})

	form.
		SetBorder(true).
		SetTitle(" Add Account [1] ").
		SetTitleAlign(tview.AlignLeft)
	
	item.AddAccountForm = form
	return form
}

// args: domain, username, notes, ID
// Work with data module
func (item *Item) NewAccDetails(accountID uuid.UUID) *tview.Frame {

	siteText := item.DataManager.LoadedAccountsList[accountID].Domain
	noteText := item.DataManager.LoadedAccountsList[accountID].Notes
	userText := item.DataManager.LoadedAccountsList[accountID].Username
	passText := item.DataManager.LoadedCredsList[accountID]

	CreateField := func(labelText string, fieldText string) tview.Primitive {

		label := tview.NewTextView().
			SetText(labelText).
			SetTextColor(tcell.ColorGray)

		field := tview.NewInputField().
			SetText(fieldText).
			SetAcceptanceFunc(func(text string, ch rune) bool {
				return false
			})		

		switch labelText {
			case "Password":

				visible := false

				field.
					SetMaskCharacter('*').
					SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
						if event.Key() == tcell.KeyEnter {

							visible = !visible

							if visible {
								// pwd := sec.DecryptPass(accountID)
								field.SetMaskCharacter(0) // unmask
								// field.SetText(pwd)
							} else {
								field.SetMaskCharacter('*') // mask again
								// field.SetText(labelText)
							}

							item.Account.LogFunc(debug.PwdShown, debug.LogContext{
								Username: userText,
								Domain: siteText,
							})

							return nil
						}
						if event.Rune() == 'c' {
							// Copy to Clipboard // may need to check for Clipboard software?
							item.Account.LogFunc(debug.PwdCopied, debug.LogContext{
								Username: userText,
								Domain: siteText,
							})
							return nil
						}
						return event
					})
		}

		container := tview.NewFlex().
			SetDirection(tview.FlexRow).
			AddItem(label, 1, 0, false).
			AddItem(field, 1, 0, true)

		return container
	}

	// Start your day here
	if item.DataManager == nil {
		panic("The item.DataManager is nil.")
	}

	siteField := CreateField("Site", siteText)
	noteField := CreateField("Notes", noteText)
	userField := CreateField("Username", userText)
	passField := CreateField("Password", passText)
	// "HWHHASHEDHAHASKJDF" will be used to obfuscate the pwd length

	editBtn := tview.NewButton("Edit").
		SetSelectedFunc(func() {})

	deleteBtn := tview.NewButton("Delete").
		SetSelectedFunc(func() {})

	flexBtn := tview.NewFlex().
		AddItem(editBtn, 10, 0, false).
		AddItem(nil, 2, 0, false).
		AddItem(deleteBtn, 10, 0, false)

	userWrapper := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(userField, 40, 0, false)

	grid := tview.NewGrid().
		SetRows(2, 2, 2, 1).
		SetColumns(40, 40).
		SetGap(1, 2)
		// SetRowSpacing(1).
		// SetColumnSpacing(4)

	grid.
		AddItem(siteField, 0, 0, 1, 2, 0, 0, false).
		AddItem(noteField, 1, 0, 1, 1, 0, 0, false).
		AddItem(passField, 1, 1, 1, 1, 0, 0, true).
		AddItem(userWrapper, 2, 0, 1, 2, 0, 0, false).
		AddItem(flexBtn, 3, 0, 1, 1, 0, 0, false)

	frame := tview.NewFrame(grid).
		SetBorders(0, 0, 0, 0, 0, 0)

	frame.
		SetBorder(true).
		SetTitle(" Account Details [1] ").
		SetTitleAlign(tview.AlignLeft)
	
	item.AccountDetails = frame

	return frame
}

// Tree
func (item *Item) NewTree(tree *tree.Tree) *tview.TreeView {
	_tree := tview.NewTreeView().
		SetRoot(tree.Root).
		SetCurrentNode(tree.Root)
	
	_tree.
		SetBorder(true).
		SetTitle(" Accounts Overview [0] ").
		SetTitleAlign(tview.AlignLeft)

	return _tree
}

// Management Functions
func (item *Item) ChangeBox(options string) {
	switch options {
		case "addAcc":
			// Change to Add Account
			item.Grid.RemoveItem(item.UsedItem)
			item.Grid.AddItem(item.AddAccountForm, 1, 1, 1, 2, 25, 50, true)
			item.UsedItem = item.CurrentItem
			item.CurrentItem = item.AddAccountForm

		case "accDet":
			// Change to Account Details in DomainName/CategoryName
			item.Grid.RemoveItem(item.UsedItem)
			item.Grid.AddItem(item.AccountDetails, 1, 1, 1, 2, 25, 50, true)
			item.UsedItem = item.CurrentItem
			item.CurrentItem = item.AccountDetails

			// Change to List Domain in Root (1st)
			// Change to List Account in Domain 
			// Change to Table of Leaked and Weak Password
	}
}
