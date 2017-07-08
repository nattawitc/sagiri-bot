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
	guildState := make(map[string]*globalstate.MemberState)
	// Store Presence status
	for _, presence := range event.Presences {
		guildState[presence.User.ID] = &globalstate.MemberState{
			Status: presence.Status,
		}
		//fmt.Println(presence.User.ID)
	}
	// Store Member data
	for _, member := range event.Members {
		if guildState[member.User.ID] == nil {
			guildState[member.User.ID] = &globalstate.MemberState{
				Status: discordgo.StatusOffline,
			}
		}
		guildState[member.User.ID].Name = member.User.Username
	}
	//data, _ := json.Marshal(guildState)
	//fmt.Println(string(data))
	globalstate.MemberStates[event.ID] = guildState
}
