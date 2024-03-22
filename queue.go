package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

func (v *VoiceInstance) PlayQueue(song Song) {
	// add song to queue
	v.QueueAdd(song)

	if v.speaking {
		// the bot is playing
		return
	}
	go func() {
		v.audioMutex.Lock()
		defer v.audioMutex.Unlock()
		for {

			if len(v.queue) == 0 {
				// dg.UpdateStatus(0, o.DiscordStatus)
				ChMessageSend(v.nowPlaying.ChannelID, "[**Music**] End of queue!")
				return
			}
			v.nowPlaying = v.QueueGetSong()
			go ChMessageSend(v.nowPlaying.ChannelID, "[**Music**] Playing, **`"+
				v.nowPlaying.Title+"`  -  `("+v.nowPlaying.Duration+")`  -  **<"+v.nowPlaying.User+">\n") //*`"+ v.nowPlaying.User +"`***")
			// If monoserver
			// if o.DiscordPlayStatus {
			// 	dg.UpdateStatus(0, v.nowPlaying.Title)
			// }
			v.stop = false
			v.skip = false
			v.speaking = true
			v.pause = false
			v.voice.Speaking(true)

			v.DCA(v.nowPlaying.VideoURL)

			// v.QueueRemoveFisrt()
			v.IncreaseCurrentQNum()
			if v.stop {
				v.QueueRemove()
			}
			v.stop = false
			v.skip = false
			v.speaking = false
			v.voice.Speaking(false)
		}
	}()
}

// DCA
func (v *VoiceInstance) DCA(url string) {
	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 128
	opts.Application = "lowdelay"
	// opts.FrameRate = 96000
	// opts.BufferedFrames = 200

	encodeSession, err := dca.EncodeFile(url, opts)
	if err != nil {
		log.Println("FATAL: Failed creating an encoding session: ", err)
	}
	v.encoder = encodeSession
	done := make(chan error)
	stream := dca.NewStream(encodeSession, v.voice, done)
	v.stream = stream
	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				log.Println("FATAL: An error occured", err)
			}
			// Clean up incase something happened and ffmpeg is still running
			encodeSession.Cleanup()
			return
		}
	}
}

// QueueGetSong
func (v *VoiceInstance) QueueGetSong() (song Song) {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	// if len(v.queue) != 0 {
	// 	return v.queue[0]
	// }

	fmt.Printf("v.currentQNum: %+v", v.currentQNum)
	fmt.Printf("len(v.queue): %+v", len(v.queue))

	if v.currentQNum+1 > len(v.queue) {
		v.currentQNum = 0
	}

	return v.queue[v.currentQNum]
	// return
}

// QueueAdd
func (v *VoiceInstance) QueueAdd(song Song) {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	v.queue = append(v.queue, song)

	jsonStr, err := json.Marshal(v.queue)
	if err != nil {
		log.Fatalln(err)
	}

	file, _ := os.OpenFile("./queue/"+v.guildID+".json", os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	_, err = file.Write(jsonStr)
	if err != nil {
		log.Fatalln(err)
	}
}

// QueueRemoveFirst
func (v *VoiceInstance) QueueRemoveFisrt() {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	if len(v.queue) != 0 {
		v.queue = v.queue[1:]
	}
}

// QueueRemoveIndex
func (v *VoiceInstance) QueueRemoveIndex(k int) {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	if len(v.queue) != 0 && k <= len(v.queue) {
		v.queue = append(v.queue[:k], v.queue[k+1:]...)
	}
}

// QueueRemoveUser
func (v *VoiceInstance) QueueRemoveUser(user string) {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	queue := v.queue
	v.queue = []Song{}
	if len(v.queue) != 0 {
		for _, q := range queue {
			if q.User != user {
				v.queue = append(v.queue, q)
			}
		}
	}
}

// QueueRemoveLast
func (v *VoiceInstance) QueueRemoveLast() {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	if len(v.queue) != 0 {
		v.queue = append(v.queue[:len(v.queue)-1], v.queue[len(v.queue):]...)
	}
}

// QueueClean
func (v *VoiceInstance) QueueClean() {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	// hold the actual song in the queue
	v.queue = v.queue[:1]
}

// QueueRemove
func (v *VoiceInstance) QueueRemove() {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	v.queue = []Song{}
}

// IncreaseCurrentQNum
func (v *VoiceInstance) IncreaseCurrentQNum() {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()
	v.currentQNum = v.currentQNum + 1
}

// LoadQueueFromFile
func (v *VoiceInstance) LoadQueueFromFile() {
	v.queueMutex.Lock()
	defer v.queueMutex.Unlock()

	file, _ := os.OpenFile("./queue/"+v.guildID+".json", os.O_APPEND|os.O_CREATE|os.O_RDONLY, 0644)
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}

	if len(byteValue) > 0 {
		songList := []Song{}
		err = json.Unmarshal(byteValue, &songList)
		if err != nil {
			log.Fatalln(err)
		}
		v.queue = songList
	}
}

// GetQueue
func GetQueue(v *VoiceInstance, m *discordgo.MessageCreate) {
	log.Println("INFO:", m.Author.Username, "send 'queue list'")
	message := "[**Music**] My songs are:\n\nNow Playing: **`" + v.nowPlaying.Title + "`  -  `(" +
		v.nowPlaying.Duration + ")`  -  " + v.nowPlaying.User + "**\n"

	queue := v.queue
	if len(queue) != 0 {
		var duration TimeDuration
		for i, q := range queue {
			message = message + "\n**`[" + strconv.Itoa(i+1) + "]`  -  `" + q.Title + "`  -  `(" + q.Duration + ")`  -  " + q.User + "**"
			d := strings.Split(q.Duration, ":")

			switch len(d) {
			case 2:
				// mm:ss
				ss, _ := strconv.Atoi(d[1])
				duration.Second = duration.Second + ss
				mm, _ := strconv.Atoi(d[0])
				duration.Minute = duration.Minute + mm
			case 3:
				// hh:mm:ss
				ss, _ := strconv.Atoi(d[2])
				duration.Second = duration.Second + ss
				mm, _ := strconv.Atoi(d[1])
				duration.Minute = duration.Minute + mm
				hh, _ := strconv.Atoi(d[0])
				duration.Hour = duration.Hour + hh
			case 4:
				// dd:hh:mm:ss
				ss, _ := strconv.Atoi(d[3])
				duration.Second = duration.Second + ss
				mm, _ := strconv.Atoi(d[2])
				duration.Minute = duration.Minute + mm
				hh, _ := strconv.Atoi(d[1])
				duration.Hour = duration.Hour + hh
				dd, _ := strconv.Atoi(d[0])
				duration.Day = duration.Day + dd
			}
		}
		t := AddTimeDuration(duration)
		message = message + "\n\nThe total duration: **`" +
			strconv.Itoa(t.Day) + "d` `" +
			strconv.Itoa(t.Hour) + "h` `" +
			strconv.Itoa(t.Minute) + "m` `" +
			strconv.Itoa(t.Second) + "s`**"
	}
	ChMessageSend(m.ChannelID, message)
	return
}

// Stop stop the audio
func (v *VoiceInstance) Stop() {
	v.stop = true
	if v.encoder != nil {
		v.encoder.Cleanup()
	}
}

func (v *VoiceInstance) Skip() bool {
	if v.speaking {
		if v.pause {
			return true
		} else {
			if v.encoder != nil {
				v.encoder.Cleanup()
			}
		}
	}
	return false
}

// Pause pause the audio
func (v *VoiceInstance) Pause() {
	v.pause = true
	if v.stream != nil {
		v.stream.SetPaused(true)
	}
}

// Resume resume the audio
func (v *VoiceInstance) Resume() {
	v.pause = false
	if v.stream != nil {
		v.stream.SetPaused(false)
	}
}

func (v *VoiceInstance) JumpTo(i int) bool {
	if v.speaking {
		if v.pause {
			return true
		}
		if v.encoder != nil {
			v.encoder.Cleanup()
			v.currentQNum = i - 2
			// v.stream.SetPaused(false)

			// v.stop = false
			// v.skip = false
			// v.speaking = true
			// v.pause = false
			// v.voice.Speaking(true)

			// v.DCA(v.nowPlaying.VideoURL)
		}
	}
	return false
}
