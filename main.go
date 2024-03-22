package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
)

var (
	dg             *discordgo.Session
	voiceInstances = map[string]*VoiceInstance{}
	mutex          sync.Mutex
	songSignal     chan PkgSong
	// purgeTime      int64
	// purgeQueue     []PurgeMessage
	// radioSignal    chan PkgRadio
	// ytdl           youtube.Youtube

	model *genai.GenerativeModel
	cs    *genai.ChatSession
)

func main() {

	// Create Gemini chat session
	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(viper.GetString("Gimini.APIKey")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// For text-only input, use the gemini-pro model
	model := client.GenerativeModel("gemini-pro")
	model.SafetySettings = []*genai.SafetySetting{
		{
			Category:  genai.HarmCategoryHateSpeech,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategorySexuallyExplicit,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryDangerousContent,
			Threshold: genai.HarmBlockNone,
		},
		{
			Category:  genai.HarmCategoryHarassment,
			Threshold: genai.HarmBlockNone,
		},
	}

	// Initialize the chat
	cs = model.StartChat()

	DiscordConnect()
	<-make(chan struct{})
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	viper.SetDefault("Log.Level", "debug")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("unable to read config: %v\n", err)
		os.Exit(1)
	}
}
