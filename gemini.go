package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
)

// GeminiMessageCreateHandler
func GeminiMessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if !strings.HasPrefix(m.Content, viper.GetString("Bot.Prefix")) {
		return
	}

	csMutex.Lock()

	content := strings.Replace(m.Content, viper.GetString("Bot.Prefix"), "", 1)
	resp, err := cs.SendMessage(
		context.Background(),
		genai.Text(content),
	)
	if err != nil {
		_, err = s.ChannelMessageSend(m.ChannelID, err.Error())
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			// Handle Part
			blob, _ := part.(genai.Blob)
			if blob.Data != nil {
				_, err = s.ChannelFileSend(
					m.ChannelID,
					"",
					bytes.NewBuffer(blob.Data),
				)
				if err != nil {
					log.Fatal(err)
				}
			}

			text, ok := part.(genai.Text)
			if ok {
				response := fmt.Sprintf("%s", text)
				arrResponse := ChunksString(response, 2000)
				for _, str := range arrResponse {
					_, err = s.ChannelMessageSend(m.ChannelID, str)
					if err != nil {
						log.Fatal(err)
					}
				}
			}
		}
	}

	csMutex.Unlock()

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
