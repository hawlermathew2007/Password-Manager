package ui

import (
	"fmt"
	"os"
	"tools/account"
	"tools/data"
	// "tools/security"
	"tools/tree"
	"tools/debug"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

type Page struct {
	Mode 			 int // 0: Normal; 1: Insert
	UIManager  *Item
	Logger 		 *debug.Logger
	ErrHandler *debug.Handler 
	App     	 *tview.Application
	Page 			 *tview.Pages
	LogSpace	 *tview.TextView
	ErrorModal *tview.Modal
	// Dashboard
	// Enter Master Pass
	// Weak Password table | Data Breaches Detection (coming soon)
	// Modals (use map)
}

func (page *Page) LogToBox(log string) {
	go func() {
		page.App.QueueUpdateDraw(func() {
			fmt.Fprintln(page.LogSpace, log)
		})
	}()
	// Remember to store the log to .log // Work with data module
	// .log has access.log for what the user done, errors.log for the program errors
}

func (page *Page) ErrorPopup(title, message string) {
	go func() {
		page.App.QueueUpdateDraw(func() {
			page.ErrorModal.SetTitle(fmt.Sprintf(" %s ", title))
			page.ErrorModal.SetText(message)
			page.Page.ShowPage("error-modal")
		})
	}()
}

func (page *Page) NewLogin() *tview.Flex {

	passwordInput := tview.NewInputField().
		SetLabel("").
		SetPlaceholder("Enter master password").
		SetMaskCharacter('*').
		SetFieldWidth(30).
		SetFieldBackgroundColor(tcell.ColorDarkSlateGray).
		SetFieldTextColor(tcell.ColorWhite).
		SetPlaceholderTextColor(tcell.ColorGray)

	titleText := tview.NewTextView().
		SetText("Master Password").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorWhite)

	passwordInput.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			password := passwordInput.GetText()
			if password == "mathew" {
				// Verify here
				// Prevent brute force too why not
				passwordInput.SetText("")

				page.Logger.Log(debug.MasterLogin, debug.LogContext{})

				vaultMaster, uiManager := page.NewVault()

				page.UIManager = uiManager
				page.Page.AddPage("vault", vaultMaster, true, false)
				page.Page.SwitchToPage("vault")

				_ = password
				page.Mode = 0

				return

			} else {
				page.ErrHandler.Handle(debug.GetUserError(debug.WrongMasterPwd))
				passwordInput.SetText("")  // clear the field
			}
		}
	})

	inputWrapper := tview.NewFlex().
    AddItem(nil, 0, 1, false).
    AddItem(passwordInput, 30, 0, true).
    AddItem(nil, 0, 1, false)

	inner := tview.NewFlex().SetDirection(tview.FlexRow).
    AddItem(titleText, 1, 0, false).
    AddItem(nil, 1, 0, false).
    AddItem(inputWrapper, 1, 0, true)

	outer := tview.NewFlex().SetDirection(tview.FlexRow).
    AddItem(nil, 0, 1, false).
    AddItem(
			tview.NewFlex().
				AddItem(nil, 0, 1, false).
				AddItem(inner, 0, 1, true).
				AddItem(nil, 0, 1, false),
			0, 1, true,
    ).
    AddItem(nil, 0, 1, false)

	outer.SetBackgroundColor(tcell.ColorBlack)
	return outer 
}

func (page *Page) NewVault() (*tview.Grid, *Item) {
	
	root := tview.NewTreeNode("Root")

	// Load ENV VAR, will be removed
	err := godotenv.Load()
	if err != nil {
		return nil, nil
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

	// Banner // Log out will Remove the vault page ofc // and spawn login page
	tcross :=  tview.NewTextView().
	   SetText("TCROSS\n\n ESC [aa] Add Account 	[ad] Add Domain 	[sp] Scan A Weak Password	 [spa] Scan Weak All Passwords	 [lv] View logs   [lg] Log out") // Should be change to Flex

	// Instruction
	text := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetText("Navigate with h, j, k, l\nSelect an account or a Category and press d to delete.")

	// Log Fields
	logText := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetText("Your logs and errors will appear here.\n").
		SetChangedFunc(func() {
			page.App.Draw()
		})

	log := tview.NewGrid().
		SetBorders(true).
		SetRows(0).
		SetColumns(0).
		AddItem(logText, 0, 0, 1, 1, 0, 0, false)
	
	page.LogSpace = logText

	// Boxes
	uiHolder := Item{}
	treeHelper := tree.Tree{
		Root: root,
		NodesList: make(map[string]*tview.TreeNode),
		ChildsList: make(map[*tview.TreeNode]uuid.UUID),
		ChildNodeList: domain,
		LogFunc: page.Logger.Log,
	}
	// dataHelper
	treeHelper.LoadTree()
	dataManager := data.TrackLoadedList(accounts[:], creds[:])
	accountHelper := account.Account{
		Tree: &treeHelper,
		DataManager: dataManager,
		LogFunc: page.Logger.Log,
	}
	tree := uiHolder.NewTree(&treeHelper)
	addAccForm := uiHolder.NewAddAccForm()
	// accDetails := ui.NewAccDetails(&accountHelper) // Sample (will be removed once the Storage is refined)

	// Grid Layout
	grid := tview.NewGrid().
		SetRows(4, 0, 5).
		SetColumns(35,0,0).
		AddItem(tree, 1, 0, 1, 1, 25, 20, false).
		AddItem(addAccForm, 1, 1, 1, 2, 25, 50, false).
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
	uiHolder.TreeBox = tree

	treeHelper.CreateChildContent	= uiHolder.NewAccDetails
	treeHelper.SelectChild = uiHolder.ChangeBox

	return grid, &uiHolder
}
