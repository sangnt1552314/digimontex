package app

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sangnt1552314/digimontex/internal/models"
	"github.com/sangnt1552314/digimontex/internal/services"
)

type App struct {
	*tview.Application
	digimon      *models.DigimonDetail
	digimonBlock *tview.Flex
}

func NewApp() *App {
	app := &App{
		Application:  tview.NewApplication(),
		digimon:      &models.DigimonDetail{},
		digimonBlock: tview.NewFlex(),
	}

	app.EnableMouse(true)

	app.setupBindings()

	root := tview.NewFlex()
	app.setupLayout(root)

	app.Application.SetRoot(root, true)

	return app
}

func (a *App) Run() error {
	return a.Application.Run()
}

func (a *App) setupBindings() {
	a.Application.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			a.Stop()
			return nil
		}
		return event
	})
}

func (a *App) setupLayout(root *tview.Flex) {
	root.SetDirection(tview.FlexRow).SetBorder(false)

	menu := a.setupMainMenu()

	mainContent := a.setupMainContent()

	root.AddItem(mainContent, 0, 1, false)
	root.AddItem(menu, 3, 0, false)
}

func (a *App) setupMainMenu() tview.Primitive {
	menuFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	menuFlex.SetBorder(true).SetBorderColor(tcell.ColorDarkCyan)
	menuFlex.SetTitle("Options").SetTitleAlign(tview.AlignLeft).SetTitleColor(tcell.ColorWhite)

	exitButton := tview.NewButton("⏻ Exit")
	exitButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorRed))
	exitButton.SetSelectedFunc(func() {
		a.Application.Stop()
	})

	menuFlex.AddItem(exitButton, 9, 0, false)

	return menuFlex
}

func (a *App) setupMainContent() tview.Primitive {
	mainContent := tview.NewFlex().SetDirection(tview.FlexColumn)
	mainContent.SetBorder(false).SetTitle("DigimonTex").SetTitleAlign(tview.AlignCenter)

	digimonListBlock := tview.NewFlex()

	digimonDetail, err := services.GetDigimonByName("Greymon")
	if err != nil {
		log.Println("Failed to fetch digimon detail:", err)
		return nil
	}
	a.digimon = digimonDetail

	a.setupDigimonBlock(a.digimonBlock)
	a.setupListDigimonBlock(digimonListBlock)

	mainContent.AddItem(digimonListBlock, 0, 2, false)
	mainContent.AddItem(a.digimonBlock, 0, 8, false)

	return mainContent
}

func (a *App) setupListDigimonBlock(block *tview.Flex) {
	page := 0
	digimonList := tview.NewList()

	block.SetDirection(tview.FlexRow)
	block.SetBorder(true).SetBorderColor(tcell.ColorDarkCyan)

	navigationFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	leftButton := tview.NewButton("<<")
	leftButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack))
	rightButton := tview.NewButton(">>")
	rightButton.SetStyle(tcell.StyleDefault.Foreground(tcell.ColorGreen).Background(tcell.ColorBlack))
	leftButton.SetSelectedFunc(func() {
		if page > 0 {
			page--
			a.buildDigimonList(digimonList, models.DigimonSearchQueryParams{
				PageSize: 10,
				Page:     page,
			})
		}
	})
	rightButton.SetSelectedFunc(func() {
		page++
		a.buildDigimonList(digimonList, models.DigimonSearchQueryParams{
			PageSize: 10,
			Page:     page,
		})
	})
	navigationFlex.AddItem(leftButton, 0, 1, false)
	navigationFlex.AddItem(rightButton, 0, 1, false)

	a.buildDigimonList(digimonList, models.DigimonSearchQueryParams{
		PageSize: 10,
		Page:     page,
	})

	block.AddItem(digimonList, 0, 1, false)
	block.AddItem(navigationFlex, 1, 0, false)
}

func (a *App) buildDigimonList(list *tview.List, params models.DigimonSearchQueryParams) {
	list.SetBorder(false)
	list.Clear()

	digimons, err := services.GetDigimonList(params)

	list.SetMainTextColor(tcell.ColorOrange)
	list.SetSelectedTextColor(tcell.ColorBlack)
	list.SetSelectedBackgroundColor(tcell.ColorWhite)

	if err != nil {
		log.Println("Failed to fetch digimon list:", err)
		list.AddItem("Failed to fetch digimon list", "", 0, nil)
	}

	for _, digimon := range digimons {
		list.AddItem(digimon.Name, "", 0, func() {
			digimonDetail, err := services.GetDigimonByID(digimon.ID)
			if err != nil {
				log.Println("Failed to fetch digimon detail:", err)
				return
			}
			a.digimon = digimonDetail
			a.setupDigimonBlock(a.digimonBlock)
		})
	}
}

func (a *App) setupDigimonBlock(block *tview.Flex) {
	if a.digimon == nil {
		log.Println("No digimon data available to display")
		return
	}
	block.Clear()

	block.SetDirection(tview.FlexColumn)
	block.SetBorder(true).SetBorderColor(tcell.ColorDarkCyan)

	leftBlock := tview.NewFlex().SetDirection(tview.FlexRow)
	rightBlock := tview.NewFlex().SetDirection(tview.FlexRow)

	imageFlex := tview.NewImage()
	if image := services.GetImageByURL(a.digimon.Images[0].Href); image != nil {
		imageFlex.SetImage(image)
		leftBlock.AddItem(imageFlex, 0, 1, false)
	} else {
		log.Println("Failed to load digimon image")
	}

	var description string
	for _, descriptionItem := range a.digimon.Descriptions {
		if descriptionItem.Language == "en_us" {
			description = descriptionItem.Description
			break
		}
	}
	descriptionBlock := tview.NewFlex()
	descriptionBlock.SetBorder(true)
	descriptionBlock.SetTitle("Description").SetTitleAlign(tview.AlignLeft).SetTitleColor(tcell.ColorWhite)
	descriptionBlock.AddItem(tview.NewTextView().
		SetText(description).
		SetTextColor(tcell.ColorWhiteSmoke), 0, 1, false)

	rightBlock.AddItem(descriptionBlock, 0, 1, false)

	block.AddItem(leftBlock, 0, 1, false)
	block.AddItem(rightBlock, 0, 1, false)
}
