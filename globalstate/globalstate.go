package globalstate

import (
	"github.com/bwmarrin/discordgo"
)

var (
	MemberStates      map[string]map[string]*MemberState
	MemberVoiceStates map[string]map[string]*MemberVoiceState
)

type MemberState struct {
	Name   string
	Status discordgo.Status
}

type MemberVoiceState struct {
	Name  string
	State *discordgo.VoiceStateUpdate
}

func init() {
	MemberStates = make(map[string]map[string]*MemberState)
	MemberVoiceStates = make(map[string]map[string]*MemberVoiceState)
}
