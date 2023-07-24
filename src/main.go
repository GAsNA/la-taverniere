package main

import (
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
  
	"github.com/bwmarrin/discordgo"
)

type color struct {
	name	string
	code	int
}

type handler_reaction struct {
    link			string
	message_id		string
    reaction		string
	reaction_id		string
	reaction_name	string
    role			*discordgo.Role
	role_id			string
	guild_id		string
}

var list_handler_reaction []handler_reaction = []handler_reaction{}

// In function because go doesn't allow const blobal array
func get_colors() []color {
	return []color{
		{ name: "Black", code: 0, },
		{ name: "Aqua", code: 1752220, },
		{ name: "Dark aqua", code: 1146986, },
		{ name: "Green", code: 3066993, },
		{ name: "Dark green", code: 2067276, },
		{ name: "Blue", code: 3447003, },
		{ name: "Dark blue", code: 2123412, },
		{ name: "Purple", code: 10181046, },
		{ name: "Dark purple", code: 7419530, },
		{ name: "Pink", code: 15277667, },
		{ name: "Dark pink", code: 11342935, },
		{ name: "Gold", code: 15844367, },
		{ name: "Dark gold", code: 15844367, },
		{ name: "Orange", code: 15105570, },
		{ name: "Dark orange", code: 11027200, },
		{ name: "Red", code: 15158332, },
		{ name: "Dark red", code: 10038562, },
		{ name: "Grey", code: 9807270, },
		{ name: "Dark grey", code: 9936031, },
		{ name: "Darker grey", code: 8359053, },
		{ name: "Light grey", code: 12370112, },
		{ name: "Navy", code: 3426654, },
		{ name: "Dark navy", code: 2899536, },
		{ name: "Yellow", code: 2899536, },
		{ name: "White", code: 16777215, },
	}
}

func get_color_by_name(name string) color {
	colors := get_colors()
	for i := 0; i < len(colors); i++ {
		if colors[i].name == name { return colors[i] }
	}
	return colors[0]
}

func list_slash_commands(sess *discordgo.Session) {
	app_id := get_env_var("DISCORD_APP_ID")

	colors := []*discordgo.ApplicationCommandOptionChoice{}
	all_colors := get_colors()
	for i := 0; i < len(all_colors); i++ {
		ac_color := discordgo.ApplicationCommandOptionChoice{ Name: all_colors[i].name, Value: all_colors[i].code, }
		colors = append(colors, &ac_color)
	}
	
	_, err := sess.ApplicationCommandBulkOverwrite(app_id, "", []*discordgo.ApplicationCommand{
		{
			Name:			"blacklist",
			Description:	"Ban a user and send a message of blacklist to the serv",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user_to_blacklist",
					Description: "User you want to ban",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "Reason of the ban",
					Required:    true,
				},
			},
		},
		{
			Name:			"kick",
			Description:	"Kick a user with a reason",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user_to_kick",
					Description: "User you want to Kick",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reason",
					Description: "Reason of the kick",
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
					Description: "People you want to get a description",
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
		{
			Name:			"message",
			Description:	"Send a message to a choose channel",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "channel",
					Description: "Channel where you want to send this message",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "message",
					Description: "Message you want to send",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "title",
					Description: "Title of your message",
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "embed",
					Description: "If you want to send your message to an embed (false by default)",
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "color",
					Description: "Color of your embed (ignored if embed is false)",
					Choices:	 colors,
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "thumbnail",
					Description: "Thumbnail of your embed (ignored if embed is false)",
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "attachment",
					Description: "Attachment",
				},
			},
		},
		{
			Name:			"handler_reaction_for_role",
			Description:	"Add a handler to add a role to each person using the chosen reaction on the chosen message",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "link_message",
					Description: "Link of the message on which you want to add the handler",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reaction",
					Description: "Reaction you want to handler",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "Role you want to add",
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
			case "handler_reaction_for_role":
				handler_reaction_for_role_command(sess, i)
		}
	})

	// HANDLER FOR REACTION ADDED
	sess.AddHandler(func (sess *discordgo.Session, m *discordgo.MessageReactionAdd,) {	
		for i := 0; i < len(list_handler_reaction); i++ { 
			this_handler := list_handler_reaction[i]

			if m.MessageReaction.MessageID != this_handler.message_id { continue }

			if (this_handler.reaction_id != "" && m.MessageReaction.Emoji.ID != this_handler.reaction_id) ||
				(m.MessageReaction.Emoji.Name != this_handler.reaction_name) { continue }

			err := sess.GuildMemberRoleAdd(this_handler.guild_id, m.MessageReaction.UserID, this_handler.role_id)
			if err != nil { log.Fatal(err) }

			log_message(sess, "add the role <@&" + this_handler.role_id + "> to <@" + m.MessageReaction.UserID + ">")
			break
		}
	})

	// HANDLER FOR REACTION DELETED
	sess.AddHandler(func (sess *discordgo.Session, m *discordgo.MessageReactionRemove,) {
		for i := 0; i < len(list_handler_reaction); i++ { 
			this_handler := list_handler_reaction[i]

			if m.MessageReaction.MessageID != this_handler.message_id { continue }

			if (this_handler.reaction_id != "" && m.MessageReaction.Emoji.ID != this_handler.reaction_id) ||
				(m.MessageReaction.Emoji.Name != this_handler.reaction_name) { continue }

			err := sess.GuildMemberRoleRemove(this_handler.guild_id, m.MessageReaction.UserID, this_handler.role_id)
			if err != nil { log.Fatal(err) }

			log_message(sess, "removes the role <@&" + this_handler.role_id + "> to <@" + m.MessageReaction.UserID + ">")
			break
		}
	})

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
