package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// DiscordConnect make a new connection to Discord
func DiscordConnect() (err error) {
	dg, err = discordgo.New("Bot " + viper.GetString("Bot.Token"))
	if err != nil {
		log.Println("FATA: error creating Discord session,", err)
		return
	}
	log.Println("INFO: Bot is Opening")
	dg.AddHandler(MessageCreateHandler)
	dg.AddHandler(GuildCreateHandler)
	dg.AddHandler(GuildDeleteHandler)
	dg.AddHandler(ConnectHandler)
	// Open Websocket
	err = dg.Open()
	if err != nil {
		log.Println("FATA: Error Open():", err)
		return
	}
	_, err = dg.User("@me")
	if err != nil {
		// Login unsuccessful
		log.Println("FATA:", err)
		return
	} // Login successful
	// log.Println("INFO: Bot user test")
	// log.Println("INFO: Bot is now running. Press CTRL-C to exit.")
	// purgeRoutine()
	initRoutine()
	// dg.UpdateStatus(0, o.DiscordStatus)
	return nil
}

func initRoutine() {
	songSignal = make(chan PkgSong)
	go GlobalPlay(songSignal)
}

func GlobalPlay(songSig chan PkgSong) {
	for {
		select {
		case song := <-songSig:
			go song.v.PlayQueue(song.data)
		}
	}
}

// MessageCreateHandler
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !strings.HasPrefix(m.Content, viper.GetString("Bot.Prefix")) {
		return
	}

	/*
	  // Method with memory (volatile)
	  guildID := SearchGuild(m.ChannelID)
	  v := voiceInstances[guildID]
	  owner, _:= s.Guild(guildID)
	  content := strings.Replace(m.Content, o.DiscordPrefix, "", 1)
	  command := strings.Fields(content)
	  if len(command) == 0 {
	    return
	  }
	  if owner.OwnerID == m.Author.ID {
	    if strings.HasPrefix(command[0], "ignore") {
	      ignore[m.ChannelID] = true
	      ChMessageSend(m.ChannelID, "[**Music**] `Ignoring` comands in this channel!")
	    }
	    if strings.HasPrefix(command[0], "unignore") {
	      if ignore[m.ChannelID] == true {
	        delete(ignore, m.ChannelID)
	        ChMessageSend(m.ChannelID, "[**Music**] `Unignoring` comands in this channel!")
	      }
	    }
	  }
	  if ignore[m.ChannelID] == true {
	    return
	  }
	*/
	// Method with database (persistent)

	guildID := SearchGuild(m.ChannelID)
	v := voiceInstances[guildID]
	// owner, _ := s.Guild(guildID)

	content := strings.Replace(m.Content, viper.GetString("Bot.Prefix"), "", 1)
	command := strings.Fields(content)
	if len(command) == 0 {
		return
	}
	fmt.Printf("command: %+v\n", command)

	// if owner.OwnerID == m.Author.ID {
	// 	if strings.HasPrefix(command[0], "ignore") {
	// 		err := PutDB(m.ChannelID, "true")
	// 		if err == nil {
	// 			ChMessageSend(m.ChannelID, "[**Music**] `Ignoring` comands in this channel!")
	// 		} else {
	// 			log.Println("FATA: Error writing in DB,", err)
	// 		}
	// 	}
	// 	if strings.HasPrefix(command[0], "unignore") {
	// 		err := PutDB(m.ChannelID, "false")
	// 		if err == nil {
	// 			ChMessageSend(m.ChannelID, "[**Music**] `Unignoring` comands in this channel!")
	// 		} else {
	// 			log.Println("FATA: Error writing in DB,", err)
	// 		}
	// 	}
	// }
	// if GetDB(m.ChannelID) == "true" {
	// 	return
	// }

	switch command[0] {
	case "join":
		JoinReporter(v, m, s)
	case "play", "p", "เล่น", "ล":
		PlayMusic(v, m, s)
	case "queue", "q", "คิว", "ค":
		GetQueue(v, m)
	case "jump", "j":
		JumpTo(v, m, s)
	default:
		return
	}
}

// SearchVoiceChannel search the voice channel id into from guild.
func SearchVoiceChannel(user string) (voiceChannelID string) {
	for _, g := range dg.State.Guilds {
		for _, v := range g.VoiceStates {
			if v.UserID == user {
				return v.ChannelID
			}
		}
	}
	return ""
}

// ConnectHandler
func ConnectHandler(s *discordgo.Session, connect *discordgo.Connect) {
	log.Println("INFO: Connected!!")
	// s.UpdateStatus(0, o.DiscordStatus)
}

// GuildCreateHandler
func GuildCreateHandler(s *discordgo.Session, guild *discordgo.GuildCreate) {
	log.Println("INFO: Guild Create:", guild.ID)
}

// GuildDeleteHandler
func GuildDeleteHandler(s *discordgo.Session, guild *discordgo.GuildDelete) {
	log.Println("INFO: Guild Delete:", guild.ID)
	// v := voiceInstances[guild.ID]
	// if v != nil {
	// 	v.Stop()
	// 	time.Sleep(200 * time.Millisecond)
	// 	mutex.Lock()
	// 	delete(voiceInstances, guild.ID)
	// 	mutex.Unlock()
	// }
}

// AddTimeDuration calculate the total time duration
func AddTimeDuration(t TimeDuration) (total TimeDuration) {
	total.Second = t.Second % 60
	t.Minute = t.Minute + t.Second/60
	total.Minute = t.Minute % 60
	t.Hour = t.Hour + t.Minute/60
	total.Hour = t.Hour % 24
	total.Day = t.Day + t.Hour/24
	return
}

// ChMessageSendEmbed
func ChMessageSendEmbed(textChannelID, title, description string) {
	embed := discordgo.MessageEmbed{}
	embed.Title = title
	embed.Description = description
	embed.Color = 0xb20000
	for i := 0; i < 10; i++ {
		_, err := dg.ChannelMessageSendEmbed(textChannelID, &embed)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		// msgToPurgeQueue(msg)
		break
	}
}

// ChMessageSendHold send a message
func ChMessageSendHold(textChannelID, message string) {
	for i := 0; i < 10; i++ {
		_, err := dg.ChannelMessageSend(textChannelID, message)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
}

// ChMessageSend send a message and auto-remove it in a time
func ChMessageSend(textChannelID, message string) {
	for i := 0; i < 10; i++ {
		_, err := dg.ChannelMessageSend(textChannelID, message)
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}
		// msgToPurgeQueue(msg)
		break
	}
}

func SearchGuild(textChannelID string) (guildID string) {
	channel, _ := dg.Channel(textChannelID)
	guildID = channel.GuildID
	return
}
