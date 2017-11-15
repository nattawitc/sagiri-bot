package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "sagiri-bot/config"
	"sagiri-bot/event/guildcreate"
	"sagiri-bot/event/messagecreate"
	_ "sagiri-bot/event/messagecreate/anime"
	_ "sagiri-bot/event/messagecreate/hottodogu"
	_ "sagiri-bot/event/messagecreate/smash"
	_ "sagiri-bot/event/messagecreate/whyhurry"
	"sagiri-bot/event/messagereactionadd"
	"sagiri-bot/event/presenceupdate"
	"sagiri-bot/event/ready"
	"sagiri-bot/event/voicestateupdate"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func addHandler(dg *discordgo.Session) {
	dg.AddHandler(ready.Ready)
	dg.AddHandler(messagecreate.MessageCreate)
	dg.AddHandler(messagereactionadd.MessageReactionAdd)
	dg.AddHandler(guildcreate.GuildCreate)
	dg.AddHandler(presenceupdate.PresenceUpdate)
	dg.AddHandler(voicestateupdate.VoiceStateUpdate)
}

func main() {
	token := viper.GetString("token")
	if token == "" {
		fmt.Println("No token provided.")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register callbacks for the events.
	addHandler(dg)

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
