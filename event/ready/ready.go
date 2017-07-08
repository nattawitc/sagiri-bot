package ready

import (
	"github.com/bwmarrin/discordgo"
)

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "!help")
}
