package main

import (
	"context"
	"fmt"
	"log"

	"fyne.io/fyne/canvas"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	// Initialize Fyne application
	myApp := app.New()
	myWindow := myApp.NewWindow("CIEL")
	image := canvas.NewImageFromFile("F:/CEIL_AI/Raphael_Abstract.webp")
	image.FillMode = canvas.ImageFillOriginal

	// Chat history box
	chatHistory := container.NewVBox()
	scrollContainer := container.NewVScroll(chatHistory)
	scrollContainer.SetMinSize(fyne.NewSize(400, 300))

	// Initialize UI elements
	promptEntry := widget.NewEntry()
	promptEntry.SetPlaceHolder("Type your message...")

	// Function to handle the generationgo
	generateResponse := func() {
		// Get the user's prompt
		prompt := promptEntry.Text
		if prompt == "" {
			return
		}

		// Display user prompt in chat
		userLabel := widget.NewLabelWithStyle("You: "+prompt, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		chatHistory.Add(container.NewVBox(userLabel))
		promptEntry.SetText("") // Clear the input after sending

		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyCvFhU02TMOfF83R4_W9FP-V3OnagnHk-U"))
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-1.5-flash")
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			log.Fatal(err)
		}

		// Display AI response in chat
		if len(resp.Candidates) > 0 {
			aiResponse := fmt.Sprintf("CIEL: %v", resp.Candidates[0].Content.Parts[0])
			aiLabel := widget.NewLabelWithStyle(aiResponse, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
			chatHistory.Add(container.NewVBox(aiLabel))
			scrollContainer.ScrollToBottom()
		} else {
			aiLabel := widget.NewLabelWithStyle("AI: No content generated", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
			chatHistory.Add(container.NewVBox(aiLabel))
		}
	}

	// Button to generate response
	generateButton := widget.NewButtonWithIcon("", theme.ConfirmIcon(), generateResponse)
	generateButton.Importance = widget.HighImportance

	// Input layout with dynamic resizing
	inputContainer := container.New(layout.NewBorderLayout(nil, nil, nil, generateButton), promptEntry, generateButton)

	// Layout for the UI
	mainContainer := container.NewBorder(nil, inputContainer, nil, nil, scrollContainer)

	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(400, 400))
	myWindow.ShowAndRun()
}
