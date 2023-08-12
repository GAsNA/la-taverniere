package main

import (
	"log"
	"strconv"
	"strings"

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
	if author.ID != user_id && !is_admin(sess, i.Member, guild_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried to use the config command, but <@" + author.ID + "> to not have the right.")

		return
	}

	message := "Action canceled..."
	if reset {
		var levels []level
		err := db.NewSelect().Model(&levels).
				Where("user_id = ? AND guild_id = ?", user_id, guild_id).
				Scan(ctx)
		if err != nil { log.Fatal(err) }

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
			if err != nil { log.Fatal(err) }

			message = "The level has been reset!"
			log_message(sess, guild_id, "reset the level of <@" + user_id + ">.", author)
		}
	}

	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: message,
					Flags: discordgo.MessageFlagsEphemeral,
				},
			})
	if err != nil { log.Fatal(err) }
}

func ask_reset_level(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, user *discordgo.User) {
	message := "Are you sure you want to reset"
	if user.ID == author.ID {
		message += " your level?"
	} else {
		message += " the level of <@" + user.ID + ">?"
	}

	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message,
					Flags: discordgo.MessageFlagsEphemeral,
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
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
						},
					},
				},
			})
	if err != nil { log.Fatal(err) }
}

func display_level(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string, user *discordgo.User) {
	var levels []level
	err := db.NewSelect().Model(&levels).
			Where("user_id = ? AND guild_id = ?", user.ID, guild_id).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	message := ""
	if len(levels) == 0 {
		if user.ID == author.ID {
			message = "You don't"
		} else {
			message = "This person doesn't"
		}
		message += " have a level yet."
	} else {
		this_level := levels[0]
		if user.ID == author.ID {
			message = "You are"
		} else {
			message = "This person is"
		}
		message += " lvl." + strconv.Itoa(int(this_level.Level)) + "."
	}
	
	ephemeral_response_for_interaction(sess, i.Interaction, message)
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
		user = optionMap["user"].UserValue(nil)
	}

	if reset {
		ask_reset_level(sess, i, author, user)
	} else {
		display_level(sess, i, author, guild_id, user)
	}
}
