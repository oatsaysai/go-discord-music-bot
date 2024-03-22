package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

// GeminiMessageCreateHandler
func GeminiMessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if !strings.HasPrefix(m.Content, viper.GetString("Bot.Prefix")) {
		return
	}

	content := strings.Replace(m.Content, viper.GetString("Bot.Prefix"), "", 1)

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(viper.GetString("Gimini.APIKey")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	resp, err := model.GenerateContent(ctx, genai.Text(content))
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, err.Error())
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	for _, candidate := range resp.Candidates {
		fmt.Println(candidate)
		response := fmt.Sprintf("%s", candidate.Content.Parts[0])
		arrResponse := ChunksString(response, 2000)

		for _, str := range arrResponse {
			_, err = s.ChannelMessageSend(m.ChannelID, str)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return
}

func ChunksString(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}
