package voicestateupdate

import (
	"sagiri-bot/globalstate"
	"sagiri-bot/logger"

	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	guildOnlineStates := globalstate.MemberVoiceStates[event.GuildID]
	id := event.UserID
	if guildOnlineStates[id] == nil {
		logger.PrintError("cannot find user id:", id)
		return
	}
	memberState := guildOnlineStates[id]
	defer func() { memberState.State = event }()
	if memberState.State.ChannelID != event.ChannelID {
		moveChan(s, event, memberState)
	}

}

func moveChan(s *discordgo.Session, event *discordgo.VoiceStateUpdate, memberState *globalstate.MemberVoiceState) {
	ch, err := s.Channel(event.ChannelID)
	if err != nil {
		logger.PrintError(err)
		return
	}
	if memberState.State.ChannelID == "" {
		sendMessage(s, memberState.Name+" has joined \""+ch.Name+"\" voice channel")
	} else {
		if event.ChannelID != "" {
			sendMessage(s, memberState.Name+" has moved to \""+ch.Name+"\" voice channel")
		} else {
			sendMessage(s, memberState.Name+" has left voice channel")
		}
	}
}

func sendMessage(s *discordgo.Session, msg string) {
	channelID := "333291751924039681"
	_, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		logger.PrintError(err)
	}
}
