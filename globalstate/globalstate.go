package globalstate

import (
	"github.com/bwmarrin/discordgo"
)

var (
	MemberStates map[string]map[string]*MemberState
)

type MemberState struct {
	Name   string
	Status discordgo.Status
}

func init() {
	MemberStates = make(map[string]map[string]*MemberState)
}
