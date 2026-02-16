package ui

import (
	"fmt"
	"tools/tree"
	"tools/account"
	"github.com/rivo/tview"
  "github.com/gdamore/tcell/v2"
)

type Item struct {
	Grid 						*tview.Grid
	UsedItem 				tview.Primitive
	CurrentItem 		tview.Primitive
	AddAccount 			*tview.Form
	AccountDetails 	*tview.Frame
	// ListDomains
	// ListAccs
	// SmallPassTable
}

type Page struct {

}

// Account
func NewAddAccForm(acc *account.Account) *tview.Form {
	form := tview.NewForm()
	form.
		AddInputField("Domain", "", 30, nil, nil).
		AddInputField("Account", "", 30, nil, nil).
		AddPasswordField("Password", "", 30, '*', nil).
		AddTextArea("Notes", "", 45, 6, 30, nil).
		AddButton("Add", func() {
			domain := form.GetFormItemByLabel("Domain").(*tview.InputField).GetText()
			account := form.GetFormItemByLabel("Account").(*tview.InputField).GetText()
			// password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
			// notes := form.GetFormItemByLabel("Notes").(*tview.InputField).GetText()
			// Should Deal with Password ofc
			
			acc.AddAccount(domain, account)
		}).
		AddButton("Scan", func() {
			acc.LogFunc("Scanning your Pass...")
		})

	form.
		SetBorder(true).
		SetTitle(" Add Account [1] ").
		SetTitleAlign(tview.AlignLeft)

	return form
}

// args: domain, username, notes, ID
func NewAccDetails(acc *account.Account) *tview.Frame {

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
				field.
					SetMaskCharacter('*').
					SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
						if event.Key() == tcell.KeyEnter {
							acc.LogFunc("Show Password")
							return nil
						}
						if event.Rune() == 'c' {
							acc.LogFunc("Copied Password")
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

	siteField := CreateField("Site", "adds.com")
	noteField := CreateField("Notes", "Just a test password")
	userField := CreateField("Username", "Administrator")
	passField := CreateField("Password", "THJDOWHASHEDSKKD")

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

	_tree.SetSelectedFunc(func(node *tview.TreeNode) {
    tree.LogFunc(fmt.Sprintf("Selected: %s", node.GetText()))
		item.ChangeBox("accDet")
		// Create NewAccDetails here (using ui module)
		// Should be able to provide the ID
		// Note: Shoudl check if the nodeName is in Domain or is Domain carefully
	})
	return _tree
}

// Management Functions
func (item *Item) ChangeBox(options string) {
	switch options {
		case "addAcc":
			// Change to Add Account
			item.Grid.RemoveItem(item.UsedItem)
			item.Grid.AddItem(item.AddAccount, 1, 1, 1, 2, 25, 50, true)
			item.UsedItem = item.AddAccount

		case "accDet":
			// Change to Account Details in DomainName/CategoryName
			item.Grid.RemoveItem(item.UsedItem)
			item.Grid.AddItem(item.AccountDetails, 1, 1, 1, 2, 25, 50, true)
			item.UsedItem = item.AccountDetails

			// Change to List Domain in Root (1st)
			// Change to List Account in Domain 
			// Change to Table of Leaked and Weak Password
	}
}

func (page *Page) ChangePage() {
	// Page Layout for providing Modal, Dashboard, Table (Scan Pass)
}
