package presenceupdate

import (
	"sagiri-bot/globalstate"
	"sagiri-bot/logger"

	"github.com/bwmarrin/discordgo"
)

func PresenceUpdate(s *discordgo.Session, event *discordgo.PresenceUpdate) {
	checkOnline(s, event)
}

func checkOnline(s *discordgo.Session, event *discordgo.PresenceUpdate) {
	guildOnlineStates := globalstate.MemberStates[event.GuildID]
	channelID := "333291751924039681"
	id := event.User.ID
	if guildOnlineStates[id] == nil {
		logger.PrintError("cannot find user id:", id)
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
	_, err := s.ChannelMessageSend(channelID, name+" has become online")
	if err != nil {
		logger.PrintError(err)
	}
}
