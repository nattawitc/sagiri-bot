package messagecreate

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	commands []command
)

type command struct {
	command     string
	description string
	handler     func(*discordgo.Session, *discordgo.MessageCreate)
}

func AddCommand(cmd, desc string, handler func(*discordgo.Session, *discordgo.MessageCreate)) {
	if commands == nil {
		commands = make([]command, 0, 0)
	}
	commands = append(commands, command{
		command:     "!" + cmd,
		description: desc,
		handler:     handler,
	})
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!help") {
		help(s, m)
	}

	for _, cmd := range commands {
		if strings.HasPrefix(m.Content, cmd.command) {
			cmd.handler(s, m)
		}
	}
}

func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	msg := "```autoit\n"
	msg += "!help - list all command\n"
	length := len(msg)
	for _, cmd := range commands {
		addMsg := cmd.command + " - " + cmd.description + "\n"
		if length+len(addMsg) > 1500 {
			s.ChannelMessageSend(m.ChannelID, msg)
			msg := "```autoit\n"
			length = len(msg)
		}
		msg += addMsg
	}
	msg += "```"
	s.ChannelMessageSend(m.ChannelID, msg)
}
