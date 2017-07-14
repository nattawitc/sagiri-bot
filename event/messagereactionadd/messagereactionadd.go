package messagereactionadd

import (
	"fmt"
	"sagiri-bot/cmdnav"
	"sagiri-bot/logger"

	"github.com/bwmarrin/discordgo"
)

var ()

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func MessageReactionAdd(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.UserID == s.State.User.ID {
		return
	}
	page := cmdnav.GetNextPage(m.MessageID, m.Emoji.Name, m.UserID)
	if page == nil {
		fmt.Println("page nil")
		return
	}

	go func() {
		msg, err := cmdnav.WritePage(s, m.ChannelID, page)
		if err != nil {
			logger.PrintError(err)
			return
		}
		cmdnav.AddPage(msg.ID, page)
	}()
	go func() {
		s.ChannelMessageDelete(m.ChannelID, m.MessageID)
	}()
	//b, _ := json.MarshalIndent(m, "", "  ")
	//fmt.Println(string(b))
	//fmt.Printf("%x\n", m.Emoji.Name)

	//// Ignore all messages created by the bot itself
	//// This isn't required in this specific example but it's a good practice.
	//if m.UserID == s.State.User.ID {
	//	return
	//}

}
