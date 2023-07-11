package main

import (
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
  
	"github.com/bwmarrin/discordgo"
)
func list_slash_commands(sess *discordgo.Session) {
	appID := getEnvVar("DISCORD_APP_ID")
	guildID := getEnvVar("DISCORD_GUILD_ID")
	
	_, err := sess.ApplicationCommandBulkOverwrite(appID, guildID, []*discordgo.ApplicationCommand{
		{
			Name:        "hello-world",
			Description: "Showcase of a basic slash command",
		},
	})
	if err != nil { log.Fatal(err) }
}

func main() {

	// INIT BOT
	token := getEnvVar("DISCORD_BOT_TOKEN")

	sess, err := discordgo.New("Bot " + token)
	if err != nil { log.Fatal(err) }

	// LIST SLASH COMMANDS
	list_slash_commands(sess)

	// ADD HANDLER FOR SLASH COMMANDS
	sess.AddHandler(func (sess *discordgo.Session, i *discordgo.InteractionCreate,) {
		data := i.ApplicationCommandData()

		switch data.Name {
			case "hello-world":
				slash_command_hello_world(sess, i, data)
		}
	})

	// TURN ON
	error_open := sess.Open()
	if error_open != nil { log.Fatal(err) }

	fmt.Println("The bot is online!")

	// CHECK SIGNAL TO STOP
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
 
	error_open = sess.Close()
	if err != nil { log.Fatal(err) }
}
