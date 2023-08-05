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
		data := i.ApplicationCommandData()

		switch data.Name {
			case "blacklist":
				blacklist_command(sess, i)
			case "kick":
				kick_command(sess, i)
			case "no-live":
				no_live_command(sess, i)
			case "who-are-this-people":
				people_command(sess, i)
			case "salope":
				salope_command(sess, i)
			case "message":
				message_command(sess, i)
			case "add-handler-reaction-for-role":
				add_handler_reaction_for_role_command(sess, i)
			case "delete-handler-reaction-for-role":
				delete_handler_reaction_for_role_command(sess, i)
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
