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
	appID := get_env_var("DISCORD_APP_ID")
	guildID := get_env_var("DISCORD_GUILD_ID")
	
	_, err := sess.ApplicationCommandBulkOverwrite(appID, guildID, []*discordgo.ApplicationCommand{
		{
			Name:			"blacklist",
			Description:	"Add a blacklist message in the appropriate channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user_to_blacklist",
					Description: "User option",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "String option",
					Required:    true,
				},
			},
		},
		{
			Name:			"no-live",
			Description:	"Make an annoucement for no live today or until a date given in parameter",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "date",
					Description: "Date to give until there will be no live",
				},
			},
		},
    {
			Name:        "who-are-this-people",
			Description: "Want to know something about this people?",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type: discordgo.ApplicationCommandOptionString,
					Name:        "people",
					Description: "People option",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Ray",
							Value: "Ray",
						},
						{
							Name: "Feitan",
							Value: "Feitan",
						},
						{
							Name: "Ukyim",
							Value: "Ukyim",
						},
						{
							Name: "Kentaro",
							Value: "Kentaro",
						},
						{
							Name: "GAsNa",
							Value: "GAsNa",
						},
					},
				},
			},
		},
		{
			Name:			"salope",
			Description:	"Suprise!",
		},
	})
	if err != nil { log.Fatal(err) }
}

func set_activity(sess *discordgo.Session, idle int, name string) {
	err := sess.UpdateGameStatus(idle, name)
	if err != nil { log.Fatal(err) }
}

func main() {

	// INIT BOT
	token := get_env_var("DISCORD_BOT_TOKEN")

	sess, err := discordgo.New("Bot " + token)
	if err != nil { log.Fatal(err) }

	// LIST SLASH COMMANDS
	list_slash_commands(sess)

	// ADD HANDLER FOR SLASH COMMANDS
	sess.AddHandler(func (sess *discordgo.Session, i *discordgo.InteractionCreate,) {
		data := i.ApplicationCommandData()

		switch data.Name {
			case "blacklist":
				blacklist_command(sess, i)
			case "no-live":
				no_live_command(sess, i)
			case "who-are-this-people":
				people_command(sess, i)
			case "salope":
				salope_command(sess, i)
		}
	})

	// TURN ON
	error_open := sess.Open()
	if error_open != nil { log.Fatal(error_open) }

	fmt.Println("The bot is online!")

	// SET ACTIVITY
	set_activity(sess, 0, "Running the tavern")

	// CHECK FOR YOUTUBE ACTIVITY 
	youtube_announcements(sess)

	// CHECK SIGNAL TO STOP
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
 
	error_open = sess.Close()
	if err != nil { log.Fatal(err) }
}
