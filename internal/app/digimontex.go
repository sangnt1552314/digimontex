package app

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"sync"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/sangnt1552314/digimontex/internal/models"
	"github.com/sangnt1552314/digimontex/internal/services"
	"github.com/sangnt1552314/digimontex/internal/services/cache"
)

type App struct {
	*tview.Application
	digimon      *models.DigimonDetail
	digimonBlock *tview.Flex
	cache        *cache.DigimonCache
	loadingMutex sync.RWMutex
	isLoading    bool
}

func NewApp() *App {
	app := &App{
		Application:  tview.NewApplication(),
		digimon:      &models.DigimonDetail{},
		digimonBlock: tview.NewFlex(),
		cache:        cache.NewDigimonCache(10),
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

	exitButton := tview.NewButton("Exit")
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
	a.setupListDigimonBlock(digimonListBlock)

	go func() {
		// Fetch default digimon detail
		// This could be any digimon, here we use "Greymon" as an example
		// You can change this to any other digimon name or ID as needed
		digimonDetail, err := services.GetDigimonByName("Greymon")
		if err != nil {
			log.Println("Failed to fetch digimon detail:", err)
			return
		}

		a.QueueUpdateDraw(func() {
			a.digimon = digimonDetail
			if digimonDetail.ID > 0 {
				a.cache.Put(digimonDetail.ID, digimonDetail)
			}
			a.setupDigimonBlock(a.digimonBlock)
		})
	}()

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

	// Show loading state
	list.AddItem("Loading...", "", 0, nil)

	// Use goroutine for API call
	go func() {
		digimons, err := services.GetDigimonList(params)

		a.QueueUpdateDraw(func() {
			list.Clear()
			list.SetMainTextColor(tcell.ColorOrange)
			list.SetSelectedTextColor(tcell.ColorBlack)
			list.SetSelectedBackgroundColor(tcell.ColorWhite)

			if err != nil {
				log.Println("Failed to fetch digimon list:", err)
				list.AddItem("Failed to fetch digimon list", "", 0, nil)
			}

			for _, digimon := range digimons {
				currentDigimon := digimon
				list.AddItem(currentDigimon.Name, "", 0, func() {
					a.loadDigimonDetail(currentDigimon.ID)
				})
			}
		})
	}()
}

func (a *App) setupDigimonBlock(block *tview.Flex) {
	if a.digimon == nil {
		log.Println("No digimon data available to display")
		return
	}

	block.Clear()

	a.loadingMutex.RLock()
	isLoading := a.isLoading
	a.loadingMutex.RUnlock()

	if isLoading {
		block.SetDirection(tview.FlexColumn)
		block.SetBorder(true).SetBorderColor(tcell.ColorDarkCyan)

		loadingText := tview.NewTextView().
			SetText("Loading Digimon details...").
			SetTextAlign(tview.AlignCenter).
			SetTextColor(tcell.ColorYellow)

		block.AddItem(loadingText, 0, 1, false)
		return
	}

	block.SetDirection(tview.FlexColumn)
	block.SetBorder(true).SetBorderColor(tcell.ColorDarkCyan)

	leftBlock := tview.NewFlex().SetDirection(tview.FlexRow)
	rightBlock := tview.NewFlex().SetDirection(tview.FlexRow)

	// Setup left block
	imagesFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	imageFlex := tview.NewImage()
	if image := services.GetImageByURL(a.digimon.Images[0].Href); image != nil {
		imageFlex.SetImage(image).SetAlign(0, 0)
		imagesFlex.AddItem(imageFlex, 0, 8, false)
	} else {
		a.loadFallbackImage(imageFlex, imagesFlex)
	}

	fieldBlock := tview.NewFlex().SetDirection(tview.FlexRow)
	for _, field := range a.digimon.Fields {
		fieldImage := tview.NewImage()
		if image := services.GetImageByURL(field.Image); image != nil {
			fieldImage.SetImage(image)
		} else {
			log.Println("Failed to load field image:", field.Image)
		}
		fieldBlock.AddItem(fieldImage, 0, 1, false)
	}
	imagesFlex.AddItem(fieldBlock, 0, 1, false)
	leftBlock.AddItem(imagesFlex, 0, 8, false)

	digimonName := tview.NewTextView().
		SetText(fmt.Sprintf("Name: %s", a.digimon.Name)).
		SetTextColor(tcell.ColorGold)
	leftBlock.AddItem(digimonName, 1, 0, false)

	digimonReleaseDate := tview.NewTextView().
		SetText(fmt.Sprintf("Release Date: %s", a.digimon.ReleaseDate)).
		SetTextColor(tcell.ColorSilver)
	leftBlock.AddItem(digimonReleaseDate, 1, 0, false)

	digimonLevel := tview.NewTextView().
		SetText(fmt.Sprintf("Levels: %s", a.getDigimonLevels())).
		SetTextColor(tcell.ColorGreen)
	leftBlock.AddItem(digimonLevel, 1, 0, false)

	digimonTypes := tview.NewTextView().
		SetText(fmt.Sprintf("Types: %s", a.getDigimonTypes())).
		SetTextColor(tcell.ColorPurple)
	leftBlock.AddItem(digimonTypes, 1, 0, false)

	digimonAttributes := tview.NewTextView().
		SetText(fmt.Sprintf("Attributes: %s", a.getDigimonAttributes())).
		SetTextColor(tcell.ColorLightCyan)
	leftBlock.AddItem(digimonAttributes, 1, 0, false)

	// Right block
	descriptionBlock := tview.NewFlex()
	descriptionBlock.SetBorder(true).SetBorderColor(tcell.ColorBlue)
	descriptionBlock.SetTitle("Description").SetTitleAlign(tview.AlignLeft).SetTitleColor(tcell.ColorOrange)
	descriptionBlock.AddItem(tview.NewTextView().
		SetText(a.getDigimonDescription()).
		SetTextColor(tcell.ColorLightCyan), 0, 1, false)

	rightBlock.AddItem(descriptionBlock, 0, 1, false)

	skillBlock := tview.NewFlex()
	skillBlock.SetDirection(tview.FlexRow)
	skillBlock.SetBorder(true).SetBorderColor(tcell.ColorRed)
	skillBlock.SetTitle("Skills").SetTitleAlign(tview.AlignLeft).SetTitleColor(tcell.ColorOrange)

	skillsTextView := tview.NewTextView().
		SetText(a.getDigimonSkills()).SetWrap(true)
	skillsTextView.SetTextColor(tcell.ColorYellow)
	skillBlock.AddItem(skillsTextView, 0, 1, false)

	rightBlock.AddItem(skillBlock, 0, 1, false)

	block.AddItem(leftBlock, 0, 1, false)
	block.AddItem(rightBlock, 0, 1, false)
}

func (a *App) loadDigimonDetail(digimonID int) {
	// Check if already loading
	a.loadingMutex.RLock()
	if a.isLoading {
		a.loadingMutex.RUnlock()
		return
	}
	a.loadingMutex.RUnlock()

	// Set loading state
	a.loadingMutex.Lock()
	a.isLoading = true
	a.loadingMutex.Unlock()

	// Check cache first
	a.setupLoadingState()

	// Use goroutine for API call
	go func() {
		digimonDetail, err := services.GetDigimonByID(digimonID)

		// Update UI on main thread
		a.QueueUpdateDraw(func() {
			a.loadingMutex.Lock()
			a.isLoading = false
			a.loadingMutex.Unlock()

			if err != nil {
				log.Println("Failed to fetch digimon detail:", err)
				return
			}

			// Cache the result
			a.cache.Put(digimonID, digimonDetail)

			// Update UI
			a.digimon = digimonDetail
			a.setupDigimonBlock(a.digimonBlock)
		})
	}()
}

func (a *App) setupLoadingState() {
	a.digimonBlock.Clear()
	a.digimonBlock.SetDirection(tview.FlexColumn)
	a.digimonBlock.SetBorder(true).SetBorderColor(tcell.ColorDarkCyan)

	loadingText := tview.NewTextView().
		SetText("Loading Digimon details...").
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorYellow)

	a.digimonBlock.AddItem(loadingText, 0, 1, false)
}

func (a *App) getDigimonDescription() string {
	var description string
	for _, descriptionItem := range a.digimon.Descriptions {
		if descriptionItem.Language == "en_us" {
			description = descriptionItem.Description
			break
		} else {
			description = "No description available in English"
		}
	}
	if description == "" {
		description = "No description available"
	}
	return description
}

func (a *App) getDigimonSkills() string {
	var skillsText string
	if len(a.digimon.Skills) == 0 {
		skillsText = "No skills available"
	} else {
		skillsText = ""
	}
	for _, skill := range a.digimon.Skills {
		if skill.Skill == "" {
			continue
		}
		var skillText string
		if skill.Description == "" {
			skillText = fmt.Sprintf("%s", skill.Skill)
		} else {
			skillText = fmt.Sprintf("%s: %s", skill.Skill, skill.Description)
		}
		skillsText += skillText + "\n"
	}
	return skillsText
}

func (a *App) getDigimonLevels() string {
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
	return levels
}

func (a *App) getDigimonTypes() string {
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
	return types
}

func (a *App) getDigimonAttributes() string {
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
	return attributes
}

func (a *App) loadFallbackImage(imageFlex *tview.Image, imagesFlex *tview.Flex) {
	noImageFile, err := os.Open("assets/no-image.png")
	if err != nil {
		log.Println("Failed to open no-image.png:", err)
	} else {
		defer noImageFile.Close()
		noImage, err := png.Decode(noImageFile)
		if err != nil {
			log.Println("Failed to decode no-image.png:", err)
		} else {
			imageFlex.SetImage(noImage).SetAlign(0, 0)
		}
	}
	imagesFlex.AddItem(imageFlex, 0, 8, false)
}
