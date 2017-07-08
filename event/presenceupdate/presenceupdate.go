package presenceupdate

import (
	"fmt"
	"time"

	"sagiri-bot/globalstate"

	"github.com/bwmarrin/discordgo"
)

func PresenceUpdate(s *discordgo.Session, event *discordgo.PresenceUpdate) {
	checkOnline(s, event)
}

func checkOnline(s *discordgo.Session, event *discordgo.PresenceUpdate) {
	guildOnlineStates := globalstate.MemberStates[event.GuildID]
	channelID := "164000870298681345"
	id := event.User.ID
	if guildOnlineStates[id] == nil {
		fmt.Println("cannot find user id:", id)
		return
	}
	defer func() { guildOnlineStates[id].Status = event.Status }()
	if guildOnlineStates[id].Status == discordgo.StatusOffline {
		welcomeBack(s, channelID, guildOnlineStates[id].Name)
		return
	}

	if guildOnlineStates[id].Status == discordgo.StatusInvisible {
		welcomeBack(s, channelID, guildOnlineStates[id].Name)
		return
	}
}

func welcomeBack(s *discordgo.Session, channelID, name string) {
	//fmt.Println("おかえり " + name)
	msg, err := s.ChannelMessageSend(channelID, "おかえり "+name)
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		time.Sleep(10 * time.Second)
		s.ChannelMessageDelete(channelID, msg.ID)
	}()
}
