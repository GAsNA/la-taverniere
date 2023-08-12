package main

import (
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
  
	"github.com/bwmarrin/discordgo"
)

func main() {

	// INIT DATABASE
	run_database()

	// INIT BOT
	token := get_env_var("DISCORD_BOT_TOKEN")

	sess, err := discordgo.New("Bot " + token)
	if err != nil { log.Fatal(err) }

	// LIST SLASH COMMANDS
	list_slash_commands(sess)

	// ACTIONS FOR INTERACTIONS
	sess.AddHandler(func (sess *discordgo.Session, i *discordgo.InteractionCreate,) {
		switch i.Type {

			case discordgo.InteractionApplicationCommand:
				data := i.ApplicationCommandData()

				switch data.Name {
					case "help":
						help_command(sess, i)
					case "config":
						config_command(sess, i)
					case "blacklist":
						blacklist_command(sess, i)
					case "kick":
						kick_command(sess, i)
					case "no-live":
						no_live_command(sess, i)
					case "message":
						message_command(sess, i)
					case "handler-reaction-for-role":
						handler_reaction_for_role_command(sess, i)
					case "level":
						level_command(sess, i)
				}

			case discordgo.InteractionMessageComponent:
				data := i.MessageComponentData()

				switch data.CustomID {
					case "success-reset-level":
						reset_level(sess, i, true)
					case "fail-reset-level":
						reset_level(sess, i, false)
				}
		}
	})

	// HANDLERS
	sess.AddHandler(new_guild_joined)
	sess.AddHandler(new_message_posted)
	sess.AddHandler(handler_reaction_to_add_role)
	sess.AddHandler(handler_reaction_to_delete_role)

	// TURN ON
	err = sess.Open()
	if err != nil { log.Fatal(err) }

	fmt.Println("The bot is online!")

	// SET ACTIVITY
	set_activity(sess, 0, "Running the tavern")

	// CHECK FOR YOUTUBE ACTIVITY 
	//youtube_announcements(sess)

	// CHECK SIGNAL TO STOP
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
 
	err = sess.Close()
	if err != nil { log.Fatal(err) }
}
