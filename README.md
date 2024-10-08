CIEL: A Conversational AI Chatbot
CIEL is a conversational AI chatbot built using the Fyne framework and the Google Generative AI API. It allows users to interact with a chatbot that can generate human-like responses to their input.

Features
Conversational interface: Users can type messages to the chatbot and receive responses in real-time.
Generative AI: The chatbot uses the Google Generative AI API to generate human-like responses to user input.
Dynamic resizing: The chatbot's input field and response area resize dynamically to accommodate user input and responses.
Requirements
Go 1.17 or later
Fyne framework (install with go get fyne.io/fyne/v2)
Google Generative AI API key (obtain from the Google Cloud Console)
Installation
Clone the repository: git clone https://github.com/your-username/ciel.git
Install dependencies: go get
Set the GOOGLE_APPLICATION_CREDENTIALS environment variable to the path of your Google Generative AI API key file.
Run the application: go run main.go
Usage
Launch the application and type a message in the input field.
Press the "Send" button or press Enter to send the message to the chatbot.
The chatbot will respond with a generated message.
Continue interacting with the chatbot by typing messages and receiving responses.
Notes
This code uses a hardcoded API key for demonstration purposes only. In a production environment, you should use a secure method to store and retrieve your API key.
The chatbot's responses are generated based on the user's input and may not always be accurate or relevant.
This code is for educational purposes only and should not be used in production without proper testing and validation.
