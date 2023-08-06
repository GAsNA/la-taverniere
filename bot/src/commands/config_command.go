package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func config_channels(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string) {
	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options[0].Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	action_id := optionMap["action"].IntValue()
	channel_id := optionMap["channel"].ChannelValue(nil).ID

	// VERIFICATION IF ENTER ALREADY EXISTS
	var channels_for_actions []channel_for_action
	err := db.NewSelect().Model(&channels_for_actions).
			Where("action_id = ? AND guild_id = ?", action_id, guild_id).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

		// UPDATE IF ALREADY EXIST
	if len(channels_for_actions) > 0 {
		if channels_for_actions[0].Channel_ID != channel_id {
			new_channel_for_action := channels_for_actions[0]
			new_channel_for_action.Channel_ID = channel_id
			_, err = db.NewUpdate().Model(&new_channel_for_action).
						Column("channel_id").
						Where("id = ?", new_channel_for_action.ID).
						Exec(ctx)
			if err != nil { log.Fatal(err) }

			ephemeral_response_for_interaction(sess, i.Interaction, "Channel updated for this action.")
			log_message(sess, guild_id, "updated a channel for an action", author)
		} else {
			ephemeral_response_for_interaction(sess, i.Interaction, "This configuration already exists.")
			log_message(sess, guild_id, "tried to add a configuration for the channel <#" + channel_id + ">, but the configuration already exists.", author)
		}
		return
	}

	// CREATE ENTER IN DB
	new_channel_for_action := &channel_for_action{
								Channel_ID: channel_id, Action_ID: action_id,
								Guild_ID: guild_id,
							}
	_, err = db.NewInsert().Model(new_channel_for_action).Ignore().Exec(ctx)
	if err != nil { log.Fatal(err) }

	ephemeral_response_for_interaction(sess, i.Interaction, "Channel added for this action.")
	log_message(sess, guild_id, "added a channel for an action", author)
}

func config_admins(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string) {
	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options[0].Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	role_id := optionMap["role"].RoleValue(nil, guild_id).ID

	// VERIFICATION IF ENTER ALREADY EXISTS
	var roles_admins []role_admin
	err := db.NewSelect().Model(&roles_admins).
			Where("role_id = ? AND guild_id = ?", role_id, guild_id).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	// IF NOT EXISTS INTEGERATE TO DB
	if len(roles_admins) == 0 {
		new_role_admin := &role_admin{Role_ID: role_id, Guild_ID: guild_id}
		_, err = db.NewInsert().Model(new_role_admin).Ignore().Exec(ctx)
		if err != nil { log.Fatal(err) }

		ephemeral_response_for_interaction(sess, i.Interaction, "Role <@&" + role_id + "> is now concidered as admin.")
		log_message(sess, guild_id, "added role <@&" + role_id + "> as admin.", author)
		return
	}

	// IF EXISTS REMOVE FROM DB
	del_role_admin := roles_admins[0]
	_, err = db.NewDelete().Model(&del_role_admin).
				Where("id = ?", del_role_admin.ID).
				Exec(ctx)
	if err != nil { log.Fatal(err) }

	ephemeral_response_for_interaction(sess, i.Interaction, "Role <@&" + role_id + "> is now removed from admin roles.")
	log_message(sess, guild_id, "removes role <@&" + role_id + "> from admin.", author)
}

func config_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User

	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	if !is_admin(sess, i.Member, guild_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried to use the config command, but <@" + author.ID + "> to not have the right.")

		return
	}

	// WHICH SUBCOMMAND
	switch i.ApplicationCommandData().Options[0].Name {
		case "config-channels":
			config_channels(sess, i, author, guild_id)
		case "config-admins":
			config_admins(sess, i, author, guild_id)
	}
}
