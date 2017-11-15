package whyhurry

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"sagiri-bot/event/messagecreate"

	"github.com/bwmarrin/discordgo"
)

var (
	stop        = make(map[string]*chan int)
	buffer      = make([][]byte, 0)
	delMsg      = make(map[string][]string)
	delMsgMutex = &sync.Mutex{}
	stopMutex   = &sync.Mutex{}
)

func init() {
	loadSound()
	messagecreate.AddCommand("whyhurry", "play whyhurry", HottoDogu)
}

func HottoDogu(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Find the channel that the message came from.
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Could not find guild.
		return
	}

	// Look for the message sender in that guild's current voice states.
	for _, vs := range g.VoiceStates {
		if vs.UserID == m.Author.ID {
			err = playSound(s, g.ID, vs.ChannelID)
			if err != nil {
				fmt.Println("Error playing sound:", err)
			}
			return
		}
	}
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound() error {
	file, err := os.Open("resource/whyhurry.dca")
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)
	key := guildID + ":" + channelID
	stopCh := make(chan int)
	stopMutex.Lock()
	stop[key] = &stopCh
	stopMutex.Unlock()
	// Send the buffer data.
	func() {
		for _, buff := range buffer {
			select {
			case <-stopCh:
				return
			default:
				vc.OpusSend <- buff
			}
		}
	}()
	close(stopCh)
	delete(stop, key)

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	return nil
}
