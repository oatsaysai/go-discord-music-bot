package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// JoinReporter
func JoinReporter(v *VoiceInstance, m *discordgo.MessageCreate, s *discordgo.Session) *VoiceInstance {
	log.Println("INFO:", m.Author.Username, "send 'join'")
	voiceChannelID := SearchVoiceChannel(m.Author.ID)
	if voiceChannelID == "" {
		log.Println("ERROR: Voice channel id not found.")
		ChMessageSend(m.ChannelID, "[**Music**] <@"+m.Author.ID+"> You need to join a voice channel!")
		return nil
	}
	if v != nil {
		log.Println("INFO: Voice Instance already created.")
	} else {
		guildID := SearchGuild(m.ChannelID)
		// create new voice instance
		mutex.Lock()
		v = new(VoiceInstance)
		voiceInstances[guildID] = v
		v.guildID = guildID
		v.session = s
		mutex.Unlock()
		//v.InitVoice()
	}
	var err error
	v.voice, err = dg.ChannelVoiceJoin(v.guildID, voiceChannelID, false, false)
	if err != nil {
		// v.Stop()
		log.Println("ERROR: Error to join in a voice channel: ", err)
		return nil
	}
	v.voice.Speaking(false)

	// Load queue list
	v.LoadQueueFromFile()

	log.Println("INFO: New Voice Instance created")

	ChMessageSend(m.ChannelID, "[**Music**] I've joined a voice channel!")
	return v
}

func PlayMusic(v *VoiceInstance, m *discordgo.MessageCreate, s *discordgo.Session) {
	log.Println("INFO:", m.Author.Username, "send 'play'")
	if v == nil {
		// log.Println("INFO: The bot is not joined in a voice channel")
		// ChMessageSend(m.ChannelID, "[**Music**] I need join in a voice channel!")
		// return
		v = JoinReporter(v, m, s)
		// time.Sleep(500 * time.Millisecond)
	}
	if len(strings.Fields(m.Content)) < 2 {
		ChMessageSend(m.ChannelID, "[**Music**] You need specify a name or URL.")
		return
	}
	// if the user is not a voice channel not accept the command
	voiceChannelID := SearchVoiceChannel(m.Author.ID)
	if v.voice.ChannelID != voiceChannelID {
		ChMessageSend(m.ChannelID, "[**Music**] <@"+m.Author.ID+"> You need to join in my voice channel for send play!")
		return
	}
	// send play my_song_youtube
	command := strings.SplitAfter(m.Content, strings.Fields(m.Content)[0])
	query := strings.TrimSpace(command[1])
	song, err := YoutubeFind(query, v, m)
	if err != nil || song.data.ID == "" {
		log.Println("ERROR: Youtube search: ", err)
		ChMessageSend(m.ChannelID, "[**Music**] I can't found this song!")
		return
	}
	//***`"+ song.data.User +"`***
	ChMessageSend(m.ChannelID, "[**Music**] **`"+song.data.User+"`** has added , **`"+
		song.data.Title+"`** to the queue. **`("+song.data.Duration+")` `["+strconv.Itoa(len(v.queue)+1)+"]`**")
	go func() {
		songSignal <- song
	}()
}

func JumpTo(v *VoiceInstance, m *discordgo.MessageCreate, s *discordgo.Session) {
	log.Println("INFO:", m.Author.Username, "send 'jump'")
	if v == nil {
		log.Println("INFO: The bot is not joined in a voice channel")
		ChMessageSend(m.ChannelID, "[**Music**] I need join in a voice channel!")
		return
	}
	if len(strings.Fields(m.Content)) < 2 {
		ChMessageSend(m.ChannelID, "[**Music**] You need specify a name or URL.")
		return
	}
	// if the user is not a voice channel not accept the command
	voiceChannelID := SearchVoiceChannel(m.Author.ID)
	if v.voice.ChannelID != voiceChannelID {
		ChMessageSend(m.ChannelID, "[**Music**] <@"+m.Author.ID+"> You need to join in my voice channel for send play!")
		return
	}
	// send play my_song_youtube
	command := strings.Split(m.Content, strings.Fields(m.Content)[0])

	i, err := strconv.Atoi(strings.TrimSpace(command[1]))
	if err != nil {
		log.Println("ERROR: ", err)
		// ChMessageSend(m.ChannelID, "[**Music**] I can't found this song!")
		return
	}
	v.JumpTo(i)
	// strconv.Itoa(command[1])

	// query := strings.TrimSpace(command[1])
	// song, err := YoutubeFind(query, v, m)
	// if err != nil || song.data.ID == "" {
	// 	log.Println("ERROR: Youtube search: ", err)
	// 	ChMessageSend(m.ChannelID, "[**Music**] I can't found this song!")
	// 	return
	// }
	//***`"+ song.data.User +"`***
	// ChMessageSend(m.ChannelID, "[**Music**] **`"+song.data.User+"`** has added , **`"+
	// 	song.data.Title+"`** to the queue. **`("+song.data.Duration+")` `["+strconv.Itoa(len(v.queue)+1)+"]`**")
	// go func() {
	// 	songSignal <- song
	// }()
}
