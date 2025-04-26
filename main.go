package main

import (
	"fmt"
	"netlos/models"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + "")
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(models.CreateMessageTicket)
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionModalSubmit:
			models.InteractionModal(i.ModalSubmitData(), s, i)
		case discordgo.InteractionMessageComponent:
			if h, ok := models.ComponentsHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
			models.ButtonInteraction(s, i)
		}
	})

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
