package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
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
)

func main() {
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
