package app

import (
	"fmt"
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

	// Fetch default digimon detail
	// This could be any digimon, here we use "Greymon" as an example
	// You can change this to any other digimon name or ID as needed
	digimonDetail, err := services.GetDigimonByName("Greymon")
	if err != nil {
		log.Println("Failed to fetch digimon detail:", err)
		return nil
	}
	a.digimon = digimonDetail

	// Setup the digimon block and list
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

	// Setup left block
	imageFlex := tview.NewImage()
	if image := services.GetImageByURL(a.digimon.Images[0].Href); image != nil {
		imageFlex.SetImage(image)
		leftBlock.AddItem(imageFlex, 18, 0, false)
	} else {
		log.Println("Failed to load digimon image")
	}
	
	digimonName := tview.NewTextView().
		SetText(fmt.Sprintf("Name: %s", a.digimon.Name))
	leftBlock.AddItem(digimonName, 0, 1, false)

	digimonReleaseDate := tview.NewTextView().
		SetText(fmt.Sprintf("Release Date: %s", a.digimon.ReleaseDate))
	leftBlock.AddItem(digimonReleaseDate, 0, 1, false)

	var levels string
	if len(a.digimon.Levels) > 0 {
		for i, t := range a.digimon.Levels {
			if i > 0 {
				levels += ", "
			}
			levels += t.Level
		}
	} else {
		levels = "Unknown"
	}
	digimonLevel := tview.NewTextView().
		SetText(fmt.Sprintf("Levels: %s", levels))
	leftBlock.AddItem(digimonLevel, 0, 1, false)

	var types string
	if len(a.digimon.Types) > 0 {
		for i, t := range a.digimon.Types {
			if i > 0 {
				types += ", "
			}
			types += t.Type
		}
	} else {
		types = "Unknown"
	}
	digimonTypes := tview.NewTextView().
		SetText(fmt.Sprintf("Types: %s", types))
	leftBlock.AddItem(digimonTypes, 0, 1, false)

	var attributes string
	if len(a.digimon.Attributes) > 0 {
		for i, t := range a.digimon.Attributes {
			if i > 0 {
				attributes += ", "
			}
			attributes += t.Attribute
		}
	} else {
		attributes = "Unknown"
	}
	digimonAttributes := tview.NewTextView().
		SetText(fmt.Sprintf("Attributes: %s", attributes))
	leftBlock.AddItem(digimonAttributes, 0, 1, false)

	fieldBlock := tview.NewFlex().SetDirection(tview.FlexColumn)
	for _, field := range a.digimon.Fields {
		fieldImage := tview.NewImage()
		if image := services.GetImageByURL(field.Image); image != nil {
			fieldImage.SetImage(image).SetFocusFunc(func() {
				log.Println("Field image clicked:", field.Image)
			})
		} else {
			log.Println("Failed to load field image:", field.Image)
		}
		fieldBlock.AddItem(fieldImage, 0, 1, false)
	}
	leftBlock.AddItem(fieldBlock, 3, 0, false)

	// Right block
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

	skillBlock := tview.NewFlex()
	skillBlock.SetDirection(tview.FlexRow)
	skillBlock.SetBorder(true)
	skillBlock.SetTitle("Skills").SetTitleAlign(tview.AlignLeft).SetTitleColor(tcell.ColorWhite)
	for _, skill := range a.digimon.Skills {
		skillText := fmt.Sprintf("%s: %s", skill.Skill, skill.Description)
		skillItem := tview.NewTextView().
			SetText(skillText).
			SetTextColor(tcell.ColorWhiteSmoke)
		skillItem.SetWrap(true)
		skillBlock.AddItem(skillItem, 0, 1, false)
	}
	rightBlock.AddItem(skillBlock, 0, 1, false)

	block.AddItem(leftBlock, 0, 1, false)
	block.AddItem(rightBlock, 0, 1, false)
}
