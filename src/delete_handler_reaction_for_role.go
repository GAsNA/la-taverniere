package main

import (
	"github.com/bwmarrin/discordgo"
)

func delete_handler_reaction_for_role_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User
	
	roles := i.Member.Roles
	is_admin := false
	for i := 0; i < len(roles); i++ {
		if is_role_admin(roles[i]) {
			is_admin = true
			break
		}
	}

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	if !is_admin {
		ephemeral_response_for_interaction(sess, i.Interaction, "You do not have the right to use this command.")
		log_message(sess, "tried to add a handler to a message to add a role with reaction, but <@" + author.ID + "> to not have the right.")

		return
	}

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	guild_id := i.Interaction.GuildID

	link_message := optionMap["link"].StringValue()
	reaction := optionMap["reaction"].StringValue()
	role := optionMap["role"].RoleValue(nil, guild_id)
	
	// VERIF LINK
	var message_guild_id string
	var message_channel_id string
	var message_id string
	if !get_discord_message_ids(link_message, &message_guild_id, &message_channel_id, &message_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "The link of the message is not at the good format.")
		log_message(sess, "tried to delete a handler on a message that adds a role with reaction, but the link of the message is not at the good format.", author)

		return
	}
	if message_guild_id != guild_id {
		ephemeral_response_for_interaction(sess, i.Interaction, "The message linked is not from this guild.")
		log_message(sess, "tried to delete a handler on a message that adds a role with reaction, but the message linked is not from this guild.", author)

		return
	}
	channel, err := sess.Channel(message_channel_id)
	if err != nil || channel.GuildID != guild_id {
		ephemeral_response_for_interaction(sess, i.Interaction, "The message linked is not from an existing channel in this guild.")
		log_message(sess, "tried to delete a handler on a message that adds a role with reaction, but the message linked is not from an existing channel in this guild.", author)

		return
	}
	_, err = sess.ChannelMessage(message_channel_id, message_id)
	if err != nil {
		ephemeral_response_for_interaction(sess, i.Interaction, "The message linked does not exist.")
		log_message(sess, "tried to delete a handler on a message that adds a role with reaction, but the message linked does not exist.", author)

		return
	}

	// VERIF REACTION
	emoji_name := ""
	emoji_id := ""
	if !check_reaction(reaction, &emoji_name, &emoji_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "The reaction is not at the good format.")
		log_message(sess, "tried to delete a handler on a message that adds a role with reaction, but the reaction is not at the good format.", author)

		return
	}

	// VERIF IF ROLE IS @everyone
	role_id := role.ID
	if role_id == guild_id {
		ephemeral_response_for_interaction(sess, i.Interaction, "You can't choose the @everyone for the role")
		log_message(sess, "tried to delete a handler on a message that adds a role with reaction, but the role chose was @everyone.", author)

		return
	}

	// VERIF IF HANDLER ALREADY EXISTS
	if !is_an_handler(link_message, reaction, role) {
		ephemeral_response_for_interaction(sess, i.Interaction, "This handler does not exists.")
		log_message(sess, "tried to delete a handler on a message that adds a role with reaction, but the hanlder does not exist.", author)

		return
	}

	// DELETE THE HANDLER
	for i := 0; i < len(list_handler_reaction); i++ {
		this_handler := list_handler_reaction[i]
		if link_message == this_handler.link && reaction == this_handler.reaction && role_id == this_handler.role.ID {
			copy(list_handler_reaction[i:], list_handler_reaction[i+1:])
			list_handler_reaction[len(list_handler_reaction)-1] = handler_reaction{}
			list_handler_reaction = list_handler_reaction[:len(list_handler_reaction)-1]
			break
		}
	}
	
	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	ephemeral_response_for_interaction(sess, i.Interaction, "Handler deleted to " + link_message + " with reaction " + reaction + " for role <@&" + role_id + ">")

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, "deletes a handler on " + link_message + " with reaction " + reaction + " for role <@&" + role_id + ">.", author)
}
