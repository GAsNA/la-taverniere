package main

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func reset_level(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string, user *discordgo.User) {
	// GET OPTIONS AND MAP
	/*options := i.ApplicationCommandData().Options[0].Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	*/

	/*var channels_for_actions []channel_for_action
	err := db.NewSelect().Model(&channels_for_actions).
			Where("action_id = ? AND guild_id = ?", action_id, guild_id).
			Scan(ctx)
	if err != nil { log.Fatal(err) }*/

	/*new_channel_for_action := channels_for_actions[0]
	new_channel_for_action.Channel_ID = channel_id
	_, err = db.NewUpdate().Model(&new_channel_for_action).
				Column("channel_id").
				Where("id = ?", new_channel_for_action.ID).
				Exec(ctx)
	if err != nil { log.Fatal(err) }*/

	/*new_channel_for_action := &channel_for_action{
								Channel_ID: channel_id, Action_ID: action_id,
								Guild_ID: guild_id,
							}
	_, err = db.NewInsert().Model(new_channel_for_action).Ignore().Exec(ctx)
	if err != nil { log.Fatal(err) }*/

	ephemeral_response_for_interaction(sess, i.Interaction, "Reset Level received!")
	//log_message(sess, guild_id, "added a channel for an action", author)
}

func display_level(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string, user *discordgo.User) {
	if user == nil {
		user = author
	}

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

	// CAN'T USE THIS COMMAND IF NOT ADMIN
/*	if !is_admin(sess, i.Member, guild_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried to use the config command, but <@" + author.ID + "> to not have the right.")

		return
	}
*/

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
	var user *discordgo.User
	if _, ok := optionMap["user"]; ok {
		user = optionMap["user"].UserValue(nil)
	}

	if reset {
		reset_level(sess, i, author, guild_id, user)
	} else {
		display_level(sess, i, author, guild_id, user)
	}
}
