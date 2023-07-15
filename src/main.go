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
				blacklist_command(sess, i, data)
		}
	})

	sess.AddHandler(func (sess *discordgo.Session, e *discordgo.Event,) {
		fmt.Println("LOG TYPE EVENT : " + e.Type)
		// WHEN KICK GET GUILD_AUDIT_LOG_ENTRY_CREATE
		// WHEN BAN GET GUILD_AUDIT_LOG_ENTRY_CREATE, GUILD_BAN_ADD
		// WHEN REMOVE BAN GET GUILD_AUDIT_LOG_ENTRY_CREATE, GUILD_BAN_REMOVE

		fmt.Printf("TYPE: %T\n", e.Struct)

		// KICK HOW ??
		//if e.Type == discordgo.guildAuditLogEntryCreateEventType {
			//fmt.Println("IS GUILD_AUDIT_LOG_ENTRY_CREATE")
			//auditlog := e.Struct.(*discordgo.GuildAuditLog)
			/*if *(auditlog).ActionType == discordgo.AuditLogActionMemberKick {
				fmt.Println("IS TYPE KICK")
				//fmt.Println("TARGET ID: " + auditlog.TargetID)
				//fmt.Println("USER ID: " + auditlog.UserID)
				//fmt.Println("REASON: " + auditlog.Reason)
			}*/
		//}
	})
	
	sess.Identify.Intents = discordgo.IntentsAll

	// TURN ON
	error_open := sess.Open()
	if error_open != nil { log.Fatal(error_open) }

	fmt.Println("The bot is online!")

	// SET ACTIVITY
	set_activity(sess, 0, "Running the tavern")

	// CHECK SIGNAL TO STOP
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop
 
	error_open = sess.Close()
	if err != nil { log.Fatal(err) }
}
