package main

import (
	"log"
	"fmt"
	"database/sql"
	
	_ "github.com/lib/pq"
	"github.com/bwmarrin/discordgo"
)

var (
	host	= get_env_var("POSTGRES_HOST")
    user	= get_env_var("POSTGRES_USER")
    password= get_env_var("POSTGRES_PASSWORD")
    dbname	= get_env_var("POSTGRES_DB")
)

var db	*sql.DB

func run_database() {
	psqlconn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, dbname)
	db, err := sql.Open("postgres", psqlconn)
	if err != nil { log.Fatal(err) }

	defer db.Close()

	if err = db.Ping(); err != nil { log.Fatal(err) }

	log.Println("The database is connected")

	// CREATE TABLE
	createTest := `CREATE TABLE IF NOT EXISTS test_users
	(
		id SERIAL PRIMARY KEY NOT NULL,
		name TEXT,
		roll_number INT
	)`
	_, err = db.Exec(createTest)
	if err != nil { log.Fatal(err) }

	// INSERT IN TABLE
	insertDynStmt := `INSERT INTO test_users(name, roll_number)
	SELECT $1, $2
	WHERE
	NOT EXISTS (
		SELECT name FROM test_users WHERE name = $1
	);
	`
    _, err = db.Exec(insertDynStmt, "Jack", 21)
	if err != nil { log.Fatal(err) }

	// SELECT IN TABLE
	rows, err := db.Query(`SELECT * FROM test_users`)
	if err != nil { log.Fatal(err) }

	defer rows.Close()
	for rows.Next() {
		id := 0
		name := ""
		roll_number := 0

		err = rows.Scan(&id, &name, &roll_number)
		if err != nil { log.Fatal(err) }

		fmt.Println(id, name, roll_number)
	}
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
