package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"log"
	"regexp"
)

// GeminiMessageCreateHandler
func GeminiMessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID { // Ignore the bot's own messages
		return
	}

	csMutex.Lock()

	for _, user := range m.Mentions {
		if user.ID == s.State.User.ID { // Check if bot was mentioned

			pattern := regexp.MustCompile(`<@\d+>`)
			msg := pattern.ReplaceAllString(m.Content, "")
			msg = msg[1:]
			fmt.Printf("%s sent message: %s\n", m.Author.ID, msg)

			resp, err := cs.SendMessage(
				context.Background(),
				genai.Text(msg),
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
						messageReference := &discordgo.MessageReference{
							MessageID: m.Message.ID,
							ChannelID: m.ChannelID,
						}
						response := fmt.Sprintf("<@%s> %s", m.Author.ID, text)
						arrResponse := ChunksString(response, 2000)
						for _, str := range arrResponse {
							replyMessage := &discordgo.MessageSend{
								Content:   str,
								Reference: messageReference,
							}
							_, err = s.ChannelMessageSendComplex(m.ChannelID, replyMessage)
							if err != nil {
								log.Fatal(err)
							}
						}
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
