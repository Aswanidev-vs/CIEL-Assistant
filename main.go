package main

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	htgotts "github.com/hegedustibor/htgo-tts"
	handlers "github.com/hegedustibor/htgo-tts/handlers"
	voices "github.com/hegedustibor/htgo-tts/voices"
)

func main() {
	// Initialize Fyne application
	myApp := app.New()
	myWindow := myApp.NewWindow("CIEL")

	// Chat history box inside a scrollable container
	chatHistory := container.NewVBox()
	scrollContainer := container.NewVScroll(chatHistory)
	scrollContainer.SetMinSize(fyne.NewSize(800, 500))

	// Initialize UI elements
	promptEntry := widget.NewEntry()
	promptEntry.SetPlaceHolder("Type your message...")

	// Function to handle the generation and speech
	generateResponse := func() {
		// Get the user's prompt
		prompt := promptEntry.Text
		if prompt == "" {
			return
		}

		// Display user prompt in chat
		userLabel := widget.NewLabelWithStyle("You: "+prompt, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
		userLabel.Wrapping = fyne.TextWrapWord
		chatHistory.Add(container.NewVBox(userLabel))
		promptEntry.SetText("") // Clear the input after sending

		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyCvFhU02TMOfF83R4_W9FP-V3OnagnHk-U"))
		if err != nil {
			log.Fatalf("Failed to create GenAI client: %v", err)
		}
		defer client.Close()

		model := client.GenerativeModel("gemini-1.5-flash")
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			log.Fatalf("Failed to generate content: %v", err)
		}

		// Display AI response in chat
		if len(resp.Candidates) > 0 {
			displayResponse := fmt.Sprintf("CIEL: %v", resp.Candidates[0].Content.Parts[0])

			// Define the response without the prefix for TTS
			aiResponse := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

			// Clean up special symbols from the AI response
			cleanedResponse := cleanSymbols(aiResponse)

			// Display AI response in chat
			aiLabel := widget.NewLabelWithStyle(displayResponse, fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
			aiLabel.Wrapping = fyne.TextWrapWord
			chatHistory.Add(container.NewVBox(aiLabel))
			scrollContainer.ScrollToBottom()

			// Generate and save the AI response audio
			speech := htgotts.Speech{Folder: "audio", Language: voices.English, Handler: &handlers.Native{}}

			// Split the response into chunks and speak each chunk
			chunks := splitIntoChunks(cleanedResponse, 100) // Splitting into chunks of 100 characters
			for _, chunk := range chunks {
				err = speech.Speak(chunk)
				if err != nil {
					log.Printf("Error generating speech: %v", err)
					return
				}
			}

			// Adjust the playback speed of the audio file
			audioFilePath := "audio/voice.mp3"                                 // Replace with the actual generated audio file path
			adjustedAudioFilePath := "audio/voice_adjusted.mp3"                // Output file with adjusted speed
			err = adjustAudioSpeed(audioFilePath, adjustedAudioFilePath, 1.25) // Adjust speed to 1.25x
			if err != nil {
				log.Printf("Error adjusting audio speed: %v", err)
				return
			}

			// Play the adjusted audio file
			err = playAudio(adjustedAudioFilePath)
			if err != nil {
				log.Printf("Error playing adjusted audio: %v", err)
				return
			}

		} else {
			aiLabel := widget.NewLabelWithStyle("CIEL: No content generated", fyne.TextAlignLeading, fyne.TextStyle{Italic: true})
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

	// Ensure the window size is responsive
	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

// cleanSymbols removes special symbols from the text except when explicitly asked
func cleanSymbols(text string) string {
	// Regular expression to match special symbols
	regex := regexp.MustCompile(`[^\w\s]`)
	return regex.ReplaceAllString(text, "")
}

// splitIntoChunks splits a string into chunks of a specified size
func splitIntoChunks(s string, chunkSize int) []string {
	var chunks []string
	var currentChunk strings.Builder

	words := strings.Fields(s)

	for _, word := range words {
		// Check if adding the next word would exceed the chunk size
		if currentChunk.Len()+len(word)+1 > chunkSize {
			chunks = append(chunks, currentChunk.String())
			currentChunk.Reset()
		}

		if currentChunk.Len() > 0 {
			currentChunk.WriteString(" ")
		}
		currentChunk.WriteString(word)
	}

	// Add any remaining text as the last chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// adjustAudioSpeed uses ffmpeg to adjust the speed of the audio file
func adjustAudioSpeed(inputPath, outputPath string, speedFactor float64) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-filter:a", fmt.Sprintf("atempo=%v", speedFactor), outputPath)
	return cmd.Run()
}

// playAudio plays the audio file (customize this function based on your needs)
func playAudio(filePath string) error {
	// Use your preferred method to play the audio file
	// Here, you could use a command-line tool like mpg123, afplay, etc.
	cmd := exec.Command("mpg123", filePath)
	return cmd.Run()
}

func addPauses(text string) string {
	var result strings.Builder

	for i, ch := range text {
		switch ch {
		case '.':
			result.WriteString(".")
			// Add a long pause after periods
			if i+1 < len(text) && text[i+1] != ' ' && text[i+1] != '\n' {
				result.WriteString("\n\n")
			}
		case ',':
			result.WriteString(",")
			// Add a short pause after commas if not followed by a space
			if i+1 < len(text) && text[i+1] != ' ' {
				result.WriteString(" ")
			}
		case '!':
			result.WriteString("!")
			// Add a medium pause with a newline for emphasis
			if i+1 < len(text) && text[i+1] != '\n' {
				result.WriteString("\n")
			}
		case '?':
			result.WriteString("?")
			// Add a medium pause with a newline for emphasis
			if i+1 < len(text) && text[i+1] != '\n' {
				result.WriteString("\n")
			}
		case ':', ';':
			result.WriteString(string(ch))
			// Add a medium pause with a newline after colons and semicolons
			if i+1 < len(text) && text[i+1] != '\n' {
				result.WriteString("\n")
			}
		case '\n':
			// Add a long pause for new paragraphs
			result.WriteString("\n\n")
		case ' ', '\t':
			// Directly append whitespace without modifying it
			result.WriteString(string(ch))
		default:
			// Append other characters directly
			result.WriteRune(ch)
		}
	}

	return result.String()
}
