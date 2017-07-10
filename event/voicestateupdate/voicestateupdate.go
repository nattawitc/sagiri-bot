package voicestateupdate

import (
	"sagiri-bot/globalstate"
	"sagiri-bot/logger"

	"github.com/bwmarrin/discordgo"
)

func VoiceStateUpdate(s *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	guildOnlineStates := globalstate.MemberStates[event.GuildID]
	id := event.UserID
	if guildOnlineStates[id] == nil {
		logger.PrintError("cannot find user id:", id)
		return
	}
	memberState := guildOnlineStates[id]
	defer func() { memberState.VoiceState = event }()
	if memberState.VoiceState == nil {
		newState(s, event, memberState)
		return
	}
	if memberState.VoiceState.ChannelID != event.ChannelID {
		moveChan(s, event, memberState)
	}
}

func newState(s *discordgo.Session, event *discordgo.VoiceStateUpdate, memberState *globalstate.MemberState) {
	ch, err := s.Channel(event.ChannelID)
	if err != nil {
		logger.PrintError(err)
		return
	}

	sendMessage(s, memberState.Name+" has joined \""+ch.Name+"\" voice channel")
}

func moveChan(s *discordgo.Session, event *discordgo.VoiceStateUpdate, memberState *globalstate.MemberState) {
	if event.ChannelID == "" {
		sendMessage(s, memberState.Name+" has left voice channel")
		return
	}
	ch, err := s.Channel(event.ChannelID)
	if err != nil {
		logger.PrintError(err)
		return
	}
	if memberState.VoiceState.ChannelID == "" {
		sendMessage(s, memberState.Name+" has joined \""+ch.Name+"\" voice channel")
	} else {
		sendMessage(s, memberState.Name+" has moved to \""+ch.Name+"\" voice channel")
	}
}

func sendMessage(s *discordgo.Session, msg string) {
	channelID := "333291751924039681"
	_, err := s.ChannelMessageSend(channelID, msg)
	if err != nil {
		logger.PrintError(err)
	}
}
