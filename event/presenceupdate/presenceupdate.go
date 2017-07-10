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
	memberState := guildOnlineStates[id]
	defer func() { guildOnlineStates[id].Status = event.Status }()
	if online(memberState.Status) != online(event.Presence.Status) {
		if online(event.Presence.Status) {
			_, err := s.ChannelMessageSend(channelID, memberState.Name+" has become online")
			if err != nil {
				logger.PrintError(err)
			}

		} else {
			_, err := s.ChannelMessageSend(channelID, memberState.Name+" has become offline")
			if err != nil {
				logger.PrintError(err)
			}
		}
	}
}

func online(status discordgo.Status) bool {
	switch status {
	case discordgo.StatusOnline:
		return true
	case discordgo.StatusIdle:
		return true
	case discordgo.StatusDoNotDisturb:
		return true
	case discordgo.StatusInvisible:
		return false
	case discordgo.StatusOffline:
		return false
	default:
		return false
	}
}
