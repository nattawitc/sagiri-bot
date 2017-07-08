package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "sagiri-bot/config"
	"sagiri-bot/event/guildcreate"
	"sagiri-bot/event/messagecreate"
	_ "sagiri-bot/event/messagecreate/slam"
	"sagiri-bot/event/presenceupdate"
	"sagiri-bot/event/ready"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func main() {
	token := viper.GetString("token")
	if token == "" {
		fmt.Println("No token provided. Please run: airhorn -t <bot token>")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready.Ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messagecreate.MessageCreate)

	// Register guildCreate as a callback for the guildCreate events.
	dg.AddHandler(guildcreate.GuildCreate)

	// Register presenceUpdate as a callback for the presenceUpdate events.
	dg.AddHandler(presenceupdate.PresenceUpdate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Sagiri is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
