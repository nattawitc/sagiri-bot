package guildcreate

import (
	"sagiri-bot/globalstate"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func GuildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}
	memberState := make(map[string]*globalstate.MemberState)
	// Store Presence status
	for _, presence := range event.Presences {
		memberState[presence.User.ID] = &globalstate.MemberState{
			Status: presence.Status,
		}
	}
	// Store Member data
	for _, member := range event.Members {
		if memberState[member.User.ID] == nil {
			memberState[member.User.ID] = &globalstate.MemberState{
				Status: discordgo.StatusOffline,
			}
		}
		memberState[member.User.ID].Name = member.User.Username
	}
	globalstate.MemberStates[event.ID] = memberState

	// Store Voice State
	for _, voice := range event.VoiceStates {
		memberState[voice.UserID].VoiceState = &discordgo.VoiceStateUpdate{
			voice,
		}
	}
}
