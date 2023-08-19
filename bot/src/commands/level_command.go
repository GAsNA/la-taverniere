package main

import (
	"log"
	"strings"
	"os"

	"github.com/bwmarrin/discordgo"
)

func reset_level(sess *discordgo.Session, i *discordgo.InteractionCreate, reset bool) {
	guild_id := i.Interaction.GuildID
	author := i.Member.User
	user_id := author.ID

	content := i.Interaction.Message.Content
	if strings.Contains(content, "<@") && strings.Contains(content, ">") {
		parts := strings.Split(content, "<@")
		parts = strings.Split(parts[1], ">")
		user_id = parts[0]
	}

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	admin, err := is_admin(sess, i.Member, guild_id)
	if err != nil { log.Println(err); return }
	if author.ID != user_id && !admin {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried to use the config command, but <@" + author.ID + "> to not have the right.")

		return
	}

	message := "Action canceled..."
	if reset {
		var levels []level
		err = db.NewSelect().Model(&levels).
				Where("user_id = ? AND guild_id = ?", user_id, guild_id).
				Scan(ctx)
		if err != nil { log.Println(err); return }

		if len(levels) == 0 {
			if author.ID == user_id {
				message = "You do not"
			} else {
				message = "User <@" + user_id + "> doesn't"
			}
			message += " have a level yet."
		} else {
			del_level := levels[0]
			_, err = db.NewDelete().Model(&del_level).
						Where("id = ?", del_level.ID).
						Exec(ctx)
			if err != nil { log.Println(err); return }

			message = "The level has been reset!"
			log_message(sess, guild_id, "reset the level of <@" + user_id + ">.", author)
			log.Println("Level if user id " + user_id + " has been reset on guild id " + guild_id)
		}
	}

	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseUpdateMessage, true, message)
}

func ask_reset_level(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, user *discordgo.User) {
	message := "Are you sure you want to reset"
	if user.ID == author.ID {
		message += " your level?"
	} else {
		message += " the level of <@" + user.ID + ">?"
	}



	component := discordgo.ActionsRow {
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: discordgo.ComponentEmoji{
								Name: "✅",
							},
							Label: "Yes, I'm sure.",
							Style: discordgo.SuccessButton,
							CustomID: "success-reset-level",
						},
						discordgo.Button{
							Emoji: discordgo.ComponentEmoji{
								Name: "❌",
							},
							Label: "No, I'm not.",
							Style: discordgo.SecondaryButton,
							CustomID: "fail-reset-level",
						},
					},
				}

	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, message, component)
}

func display_level(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string, user *discordgo.User) {
	var levels []level
	err := db.NewSelect().Model(&levels).
			Where("user_id = ? AND guild_id = ?", user.ID, guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

	if len(levels) == 0 {
		message := ""
		if user.ID == author.ID {
			message = "You don't"
		} else {
			message = "This person doesn't"
		}
		message += " have a level yet."
		
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, message)

		return
	}
	
	this_level := levels[0]

	username := user.Username
	guild, err := sess.Guild(guild_id)
	guild_name := guild.Name
	name_file := "level-" + username + "-" + guild_name + ".png"
	link_image_level, err := get_image_level(name_file, username, user.AvatarURL("80"), guild_name, this_level.Level)
	if err != nil { log.Println(err); return }

	reader_file, err := os.Open(link_image_level)
	if err != nil { log.Println(err); return }

	file :=	&discordgo.File{
				Name:			name_file,
				ContentType:	"image/png",
				Reader:			reader_file,
			}

	message := ""

	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, false, message, file)
	
	// DELETE LOCAL FILE
	err = reader_file.Close()
	if err != nil { log.Println(err); return }
	err = os.Remove(link_image_level)
	if err != nil { log.Println(err); return }
}

func level_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User

	guild_id := i.Interaction.GuildID

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	//user, reset
	reset := false
	if _, ok := optionMap["reset"]; ok {
		reset = optionMap["reset"].BoolValue()
	}
	user := author
	if _, ok := optionMap["user"]; ok {
		user = optionMap["user"].UserValue(sess)
	}

	if reset {
		ask_reset_level(sess, i, author, user)
	} else {
		display_level(sess, i, author, guild_id, user)
	}
}
