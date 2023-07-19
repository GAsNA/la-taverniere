package main

import (
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
  
	"github.com/bwmarrin/discordgo"
)

const (
	// COLORS INT
	BLACK = 0
	AQUA = 1752220
	DARK_AQUA = 1146986
	GREEN = 3066993
	DARK_GREEN = 2067276
	BLUE = 3447003
	DARK_BLUE = 2123412
	PURPLE = 10181046
	DARK_PURPLE = 7419530
	PINK = 15277667
	DARK_PINK = 11342935
	GOLD = 15844367
	DARK_GOLD = 15844367
	ORANGE = 15105570
	DARK_ORANGE = 11027200
	RED = 15158332
	DARK_RED = 10038562
	GREY = 9807270
	DARK_GREY = 9936031
	DARKER_GREY = 8359053
	LIGHT_GREY = 12370112
	NAVY = 3426654
	DARK_NAVY = 2899536
	YELLOW = 2899536
	WHITE = 16777215
)

func list_slash_commands(sess *discordgo.Session) {
	appID := get_env_var("DISCORD_APP_ID")
	guildID := get_env_var("DISCORD_GUILD_ID")
	
	_, err := sess.ApplicationCommandBulkOverwrite(appID, guildID, []*discordgo.ApplicationCommand{
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
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name: "Black",
							Value: BLACK,
						},
						{
							Name: "Aqua",
							Value: AQUA,
						},
						{
							Name: "Dark aqua",
							Value: DARK_AQUA,
						},
						{
							Name: "Green",
							Value: GREEN,
						},
						{
							Name: "Dark green",
							Value: DARK_GREEN,
						},
						{
							Name: "Blue",
							Value: BLUE,
						},
						{
							Name: "Dark blue",
							Value: DARK_BLUE,
						},
						{
							Name: "Purple",
							Value: PURPLE,
						},
						{
							Name: "Dark purple",
							Value: DARK_PURPLE,
						},
						{
							Name: "Pink",
							Value: DARK_PINK,
						},
						{
							Name: "Gold",
							Value: GOLD,
						},
						{
							Name: "Dark gold",
							Value: DARK_GOLD,
						},
						{
							Name: "Orange",
							Value: ORANGE,
						},
						{
							Name: "Dark orange",
							Value: DARK_ORANGE,
						},
						{
							Name: "Red",
							Value: RED,
						},
						{
							Name: "Dark red",
							Value: DARK_RED,
						},
						{
							Name: "Grey",
							Value: GREY,
						},
						{
							Name: "Dark grey",
							Value: DARK_GREY,
						},
						{
							Name: "Darker grey",
							Value: DARKER_GREY,
						},
						{
							Name: "Light grey",
							Value: LIGHT_GREY,
						},
						{
							Name: "Navy",
							Value: NAVY,
						},
						{
							Name: "Dark navy",
							Value: DARK_NAVY,
						},
						{
							Name: "Yellow",
							Value: YELLOW,
						},
						{
							Name: "White",
							Value: WHITE,
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionAttachment,
					Name:        "thumbnail",
					Description: "Thumbnail of your embed. If embed is false, the image is just an attachement",
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
			case "no-live":
				no_live_command(sess, i)
			case "who-are-this-people":
				people_command(sess, i)
			case "salope":
				salope_command(sess, i)
			case "message":
				message_command(sess, i)
		}
	})

	// TURN ON
	error_open := sess.Open()
	if error_open != nil { log.Fatal(error_open) }

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
 
	error_open = sess.Close()
	if err != nil { log.Fatal(err) }
}
