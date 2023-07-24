package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/forPelevin/gomoji"
)

func check_reaction(reaction string, emoji_name *string, emoji_id *string) bool {
	find_all := gomoji.FindAll(reaction)
	if len(find_all) > 1 { return false }

	if len(find_all) == 1 {
		reaction_without_emoji := gomoji.RemoveEmojis(reaction)
		if len(reaction_without_emoji) > 0 { return false }
		
		if len(reaction_without_emoji) == 0 {
			*emoji_name = reaction
			return true
		}
	}

	if len(find_all) == 0 {
		if !strings.HasPrefix(reaction, "<:") { return false }
		reaction = strings.TrimLeft(reaction, "<:")

		if reaction[len(reaction) - 1] != '>' { return false }
		reaction = strings.TrimRight(reaction, ">")

		parts := strings.Split(reaction, ":")
		if len(parts) != 2 { return false }

		for i := 0; i < len(parts[1]); i++ {
			if parts[1][i] < '0' || parts[1][i] > '9' { return false }
		}

		*emoji_name = parts[0]
		*emoji_id = parts[1]
		return true
	}

	return false
}

func is_an_handler(link string, reaction string, role *discordgo.Role) bool {
	for i := 0; i < len(list_handler_reaction); i++ {
		if list_handler_reaction[i].link == link && list_handler_reaction[i].reaction == reaction &&
			list_handler_reaction[i].role.ID == role.ID {
			return true
		}
	}
	
	return false
}

func handler_reaction_for_role_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
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

	link_message := optionMap["link_message"].StringValue()
	reaction := optionMap["reaction"].StringValue()
	role := optionMap["role"].RoleValue(nil, guild_id)
	
	// VERIF LINK
	var message_guild_id string
	var message_channel_id string
	var message_id string
	if !get_discord_message_ids(link_message, &message_guild_id, &message_channel_id, &message_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "The link of the message is not at the good format.")
		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the link of the message is not at the good format.", author)

		return
	}
	if message_guild_id != guild_id {
		ephemeral_response_for_interaction(sess, i.Interaction, "The message linked is not from this guild.")
		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the message linked is not from this guild.", author)

		return
	}
	channel, err := sess.Channel(message_channel_id)
	if err != nil || channel.GuildID != guild_id {
		ephemeral_response_for_interaction(sess, i.Interaction, "The message linked is not from an existing channel in this guild.")
		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the message linked is not from an existing channel in this guild.", author)

		return
	}
	_, err = sess.ChannelMessage(message_channel_id, message_id)
	if err != nil {
		ephemeral_response_for_interaction(sess, i.Interaction, "The message linked does not exist.")
		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the message linked does not exist.", author)

		return
	}

	// VERIF REACTION
	emoji_name := ""
	emoji_id := ""
	if !check_reaction(reaction, &emoji_name, &emoji_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "The reaction is not at the good format.")
		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the reaction is not at the good format.", author)

		return
	}

	// VERIF IF ROLE IS @everyone
	role_id := role.ID
	if role_id == guild_id {
		ephemeral_response_for_interaction(sess, i.Interaction, "You can't choose the @everyone for the role")
		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the role chose was @everyone.", author)

		return
	}

	// VERIF IF HANDLER ALREADY EXISTS
	if is_an_handler(link_message, reaction, role) {
		ephemeral_response_for_interaction(sess, i.Interaction, "This handler was already made.")
		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the hanlder already exists.", author)

		return
	}

	// ADD IN HANDLER
	new_handler := handler_reaction{
		link: link_message,
		message_id: message_id,
		reaction: reaction,
		reaction_id: emoji_id,
		reaction_name: emoji_name,
		role: role,
		role_id: role_id,
		guild_id: guild_id,
	}
	list_handler_reaction = append(list_handler_reaction, new_handler)
	
	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	ephemeral_response_for_interaction(sess, i.Interaction, "Handler add to " + link_message + " with reaction " + reaction + " for role <@&" + role_id + ">")

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, "adds a handler to " + link_message + " with reaction " + reaction + " for role <@&" + role_id + ">.", author)
}
