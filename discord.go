package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func TestDiscord() {
	fmt.Println(viper.GetString("Bot.Token"))
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + viper.GetString("Bot.Token"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Just like the ping pong example, we only care about receiving message
	// events in this example.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// DiscordConnect make a new connection to Discord
func DiscordConnect() (err error) {
	dg, err = discordgo.New("Bot " + viper.GetString("Bot.Token"))
	if err != nil {
		log.Println("FATAL: error creating Discord session,", err)
		return
	}

	log.Println("INFO: Bot is Opening")

	// dg.AddHandler(MessageCreateHandler)
	// dg.AddHandler(GuildCreateHandler)
	// dg.AddHandler(GuildDeleteHandler)
	// dg.AddHandler(ConnectHandler)

	dg.AddHandler(GeminiMessageCreateHandler)

	// Open Websocket
	err = dg.Open()
	if err != nil {
		log.Println("FATAL: Error Open():", err)
		return
	}

	_, err = dg.User("@me")
	if err != nil {
		// Login unsuccessful
		log.Println("FATAL:", err)
		return
	}

	// Login successful
	log.Println("INFO: Bot is now running. Press CTRL-C to exit.")
	initRoutine()
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// In this example, we only care about messages that are "ping".
	if m.Content != "ping" {
		return
	}

	// We create the private channel with the user who sent the message.
	channel, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		// If an error occurred, we failed to create the channel.
		//
		// Some common causes are:
		// 1. We don't share a server with the user (not possible here).
		// 2. We opened enough DM channels quickly enough for Discord to
		//    label us as abusing the endpoint, blocking us from opening
		//    new ones.
		fmt.Println("error creating channel:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Something went wrong while sending the DM!",
		)
		return
	}
	// Then we send the message through the channel we created.
	_, err = s.ChannelMessageSend(channel.ID, "Pong!")
	if err != nil {
		// If an error occurred, we failed to send the message.
		//
		// It may occur either when we do not share a server with the
		// user (highly unlikely as we just received a message) or
		// the user disabled DM in their settings (more likely).
		fmt.Println("error sending DM message:", err)
		s.ChannelMessageSend(
			m.ChannelID,
			"Failed to send you a DM. "+
				"Did you disable DM in your privacy settings?",
		)
	}
}

// MessageCreateHandler
func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// _, err = s.ChannelMessageSend(channel.ID, "Pong!")
	ChMessageSend(m.ChannelID, "Pong!")

	return

	// if !strings.HasPrefix(m.Content, viper.GetString("Bot.Prefix")) {
	// 	return
	// }

	// /*
	//   // Method with memory (volatile)
	//   guildID := SearchGuild(m.ChannelID)
	//   v := voiceInstances[guildID]
	//   owner, _:= s.Guild(guildID)
	//   content := strings.Replace(m.Content, o.DiscordPrefix, "", 1)
	//   command := strings.Fields(content)
	//   if len(command) == 0 {
	//     return
	//   }
	//   if owner.OwnerID == m.Author.ID {
	//     if strings.HasPrefix(command[0], "ignore") {
	//       ignore[m.ChannelID] = true
	//       ChMessageSend(m.ChannelID, "[**Music**] `Ignoring` comands in this channel!")
	//     }
	//     if strings.HasPrefix(command[0], "unignore") {
	//       if ignore[m.ChannelID] == true {
	//         delete(ignore, m.ChannelID)
	//         ChMessageSend(m.ChannelID, "[**Music**] `Unignoring` comands in this channel!")
	//       }
	//     }
	//   }
	//   if ignore[m.ChannelID] == true {
	//     return
	//   }
	// */
	// // Method with database (persistent)

	// guildID := SearchGuild(m.ChannelID)
	// v := voiceInstances[guildID]
	// // owner, _ := s.Guild(guildID)

	// content := strings.Replace(m.Content, viper.GetString("Bot.Prefix"), "", 1)
	// command := strings.Fields(content)
	// if len(command) == 0 {
	// 	return
	// }
	// fmt.Printf("command: %+v\n", command)

	// // if owner.OwnerID == m.Author.ID {
	// // 	if strings.HasPrefix(command[0], "ignore") {
	// // 		err := PutDB(m.ChannelID, "true")
	// // 		if err == nil {
	// // 			ChMessageSend(m.ChannelID, "[**Music**] `Ignoring` comands in this channel!")
	// // 		} else {
	// // 			log.Println("FATAL: Error writing in DB,", err)
	// // 		}
	// // 	}
	// // 	if strings.HasPrefix(command[0], "unignore") {
	// // 		err := PutDB(m.ChannelID, "false")
	// // 		if err == nil {
	// // 			ChMessageSend(m.ChannelID, "[**Music**] `Unignoring` comands in this channel!")
	// // 		} else {
	// // 			log.Println("FATAL: Error writing in DB,", err)
	// // 		}
	// // 	}
	// // }
	// // if GetDB(m.ChannelID) == "true" {
	// // 	return
	// // }

	// switch command[0] {
	// case "join":
	// 	JoinReporter(v, m, s)
	// case "play", "p", "เล่น", "ล":
	// 	PlayMusic(v, m, s)
	// case "queue", "q", "คิว", "ค":
	// 	GetQueue(v, m)
	// case "jump", "j":
	// 	JumpTo(v, m, s)
	// default:
	// 	return
	// }
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
