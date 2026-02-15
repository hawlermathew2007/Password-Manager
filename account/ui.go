package account

import (
	"github.com/rivo/tview"
  "github.com/gdamore/tcell/v2"
)

func (acc *Account) NewAddAccForm() *tview.Form {
	form := tview.NewForm()
	form.
		AddInputField("Domain", "", 30, nil, nil).
		AddInputField("Account", "", 30, nil, nil).
		AddPasswordField("Password", "", 30, '*', nil).
		AddButton("Add", func() {
			domain := form.GetFormItemByLabel("Domain").(*tview.InputField).GetText()
			account := form.GetFormItemByLabel("Account").(*tview.InputField).GetText()
			// password := form.GetFormItemByLabel("Password").(*tview.InputField).GetText()
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

func (acc *Account) NewAccDetails() *tview.Frame {

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




