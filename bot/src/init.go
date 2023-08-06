package main

import (
	"log"
	"fmt"
	"database/sql"
	"context"

	_ "github.com/lib/pq"
	"github.com/bwmarrin/discordgo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
)

var (
	ctx 	= context.Background()
	db		*bun.DB
)

func run_database() {
	// INIT DB
	host	:= get_env_var("POSTGRES_HOST")
	user_pg	:= get_env_var("POSTGRES_USER")
	password:= get_env_var("POSTGRES_PASSWORD")
	dbname	:= get_env_var("POSTGRES_DB")

	psqlconn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user_pg, password, dbname)
	sqldb, err := sql.Open("postgres", psqlconn)
	if err != nil { log.Fatal(err) }

	if err = sqldb.Ping(); err != nil { log.Fatal(err) }	

	// GET BUN DB
	db = bun.NewDB(sqldb, sqlitedialect.New())

	log.Println("The database is connected")

	// CREATION TABLES
		// Guild
	_, err = db.NewCreateTable().Model((*guild)(nil)).IfNotExists().Exec(ctx)
	if err != nil { log.Fatal(err) }
		// Handler_Reaction_Role
	_, err = db.NewCreateTable().Model((*handler_reaction_role)(nil)).
				ForeignKey(`("guild_id") REFERENCES "guild" ("guild_id") ON DELETE CASCADE`).
				IfNotExists().Exec(ctx)
	if err != nil { log.Fatal(err) }
		// Level
	_, err = db.NewCreateTable().Model((*level)(nil)).
				ForeignKey(`("guild_id") REFERENCES "guild" ("guild_id") ON DELETE CASCADE`).
				IfNotExists().Exec(ctx)
	if err != nil { log.Fatal(err) }
		// Action
	_, err = db.NewCreateTable().Model((*action)(nil)).
				IfNotExists().Exec(ctx)
	if err != nil { log.Fatal(err) }
		// Channel_For_Action
	_, err = db.NewCreateTable().Model((*channel_for_action)(nil)).
				ForeignKey(`("action_id") REFERENCES "action" ("id") ON DELETE CASCADE`).
				ForeignKey(`("guild_id") REFERENCES "guild" ("guild_id") ON DELETE CASCADE`).
				IfNotExists().Exec(ctx)
	if err != nil { log.Fatal(err) }
		// Role_Admin
	_, err = db.NewCreateTable().Model((*role_admin)(nil)).
				ForeignKey(`("guild_id") REFERENCES "guild" ("guild_id") ON DELETE CASCADE`).
				IfNotExists().Exec(ctx)
	if err != nil { log.Fatal(err) }

	// ENTER ACTIONS IN DB
	for i := 0; i < len(actions_db); i++ {
		var actions []action
		err = db.NewSelect().Model(&actions).
				Where("name = ?", actions_db[i].name).
				Scan(ctx)
		if err != nil { log.Fatal(err) }

		if len(actions) == 0 {
			new_action := &action{Name: actions_db[i].name}
			_, err := db.NewInsert().Model(new_action).Ignore().Exec(ctx)
			if err != nil { log.Fatal(err) }

			err = db.NewSelect().Model(&actions).
					Where("name = ?", actions_db[i].name).
					Scan(ctx)
			if err != nil { log.Fatal(err) }
		}

		actions_db[i].id = actions[0].ID
	}
}

func list_slash_commands(sess *discordgo.Session) {
	app_id := get_env_var("DISCORD_APP_ID")

	actions := []*discordgo.ApplicationCommandOptionChoice{}
	for i := 0; i < len(actions_db); i++ {
		ac_action := discordgo.ApplicationCommandOptionChoice {
						Name: actions_db[i].name,
						Value: actions_db[i].id,
					}
		actions = append(actions, &ac_action)
	}

	colors := []*discordgo.ApplicationCommandOptionChoice{}
	all_colors := get_colors()
	for i := 0; i < len(all_colors); i++ {
		ac_color := discordgo.ApplicationCommandOptionChoice{
						Name: all_colors[i].name,
						Value: all_colors[i].code,
					}
		colors = append(colors, &ac_color)
	}
	
	_, err := sess.ApplicationCommandBulkOverwrite(app_id, "", []*discordgo.ApplicationCommand{
		{
			Name:			"config",
			Description:	"Configure the bot.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:			"config-channels",
					Description:	"Give a channel associate to an action.",
					Type:			discordgo.ApplicationCommandOptionSubCommand,
					Options:		[]*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionInteger,
							Name:        "action",
							Description: "Action that can be performed by the bot",
							Required:    true,
							Choices:	 actions,
						},
						{
							Type:        discordgo.ApplicationCommandOptionChannel,
							Name:        "channel",
							Description: "Channel associated to this action.",
							Required:    true,
						},
					},
				},
				{
					Name:			"config-admins",
					Description:	"Give a role to determine as administrator role for the bot.",
					Type:			discordgo.ApplicationCommandOptionSubCommand,
					Options:		[]*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionRole,
							Name:        "role",
							Description: "Administrator role",
							Required:    true,
						},
					},
				},
			},
		},
		{
			Name:			"blacklist",
			Description:	"Ban a user and send a message of blacklist to the serv",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
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
					Name:        "user",
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
						{ Name: "Ray", Value: "Ray", },
						{ Name: "Feitan", Value: "Feitan", },
						{ Name: "Ukyim", Value: "Ukyim", },
						{ Name: "Kentaro", Value: "Kentaro", },
						{ Name: "GAsNa", Value: "GAsNa", },
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
			Name:			"add-handler-reaction-for-role",
			Description:	"Add a handler to add a role to each person using the chosen reaction on the chosen message",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "link",
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
		{
			Name:			"delete-handler-reaction-for-role",
			Description:	"Delete a handler that adds a role to each person using the chosen reaction on the chosen message",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "link",
					Description: "Link of the message concerned by the handler",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "reaction",
					Description: "Reaction concerned by the handler",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role",
					Description: "Role concerned by the handler",
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
