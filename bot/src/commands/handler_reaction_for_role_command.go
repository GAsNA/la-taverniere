package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func handler_reaction_for_role_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User
	
	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	admin, err := is_admin(sess, i.Member, guild_id)
	if err != nil { log.Println(err); return }
	if !admin {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried to add a handler to a message to add a role with reaction, but <@" + author.ID + "> to not have the right.")

		return
	}

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	link_message := optionMap["link"].StringValue()
	reaction := optionMap["reaction"].StringValue()
	role := optionMap["role"].RoleValue(nil, guild_id)
	
	// VERIF LINK
	var message_guild_id string
	var message_channel_id string
	var message_id string
	if !get_discord_message_ids(link_message, &message_guild_id, &message_channel_id, &message_id) {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "The link of the message is not at the good format.")
		return
	}
	if message_guild_id != guild_id {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "The message linked is not from this guild.")
		return
	}
	channel, err := sess.Channel(message_channel_id)
	if err != nil || channel.GuildID != guild_id {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "The message linked is not from an existing channel in this guild.")
		return
	}
	_, err = sess.ChannelMessage(message_channel_id, message_id)
	if err != nil {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "The message linked does not exist.")
		return
	}

	// VERIF REACTION
	emoji_name := ""
	emoji_id := *new(string)
	if !check_reaction(reaction, &emoji_name, &emoji_id) {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "The reaction is not at the good format.")
		return
	}

	// VERIF IF ROLE IS @everyone
	role_id := role.ID
	if role_id == guild_id {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "You can't choose the @everyone for the role")
		return
	}

	// VERIF IF HANDLER ALREADY EXISTS
	registered_handler, err := is_a_registered_handler(link_message, reaction, role)
	if err != nil { log.Println(err); return }
	if registered_handler {
		// DELETE HANDLER IN DB
		del_handler := &handler_reaction_role {
			Msg_Link: link_message, Msg_ID: message_id,
			Reaction: reaction, Reaction_ID: emoji_id, Reaction_Name: emoji_name,
			Role_ID: role_id,
			Guild_ID: guild_id,
		}
		
		_, err = db.NewDelete().
					Model(del_handler).
					Where("msg_link = ? AND reaction = ? AND role_id = ?", link_message, reaction, role_id).
					Exec(ctx)
		if err != nil { log.Println(err); return }
	
		// RESPOND TO USER WITH EPHEMERAL MESSAGE
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Handler deleted to " + link_message + " with reaction " + reaction + " for role <@&" + role_id + ">")

		// ADD LOG IN LOGS CHANNEL
		log_message(sess, guild_id, "deleted a handler on " + link_message + " with reaction " + reaction + " for role <@&" + role_id + ">.", author)

		log.Println("Handler has been removed with reaction name " + emoji_name + " on link message " + link_message + " for role id " + role_id + " on guild id " + guild_id)
		return
	}
	
	// ADD HANDLER IN DB
	new_handler := &handler_reaction_role{
		Msg_Link: link_message, Msg_ID: message_id,
		Reaction: reaction, Reaction_ID: emoji_id, Reaction_Name: emoji_name,
		Role_ID: role_id,
		Guild_ID: guild_id,
	}
	_, err = db.NewInsert().Model(new_handler).Ignore().Exec(ctx)
	if err != nil { log.Println(err); return }

	log.Println("New Handler inserted in table!")
	
	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Handler added to " + link_message + " with reaction " + reaction + " for role <@&" + role_id + ">")

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, guild_id, "added a handler to " + link_message + " with reaction " + reaction + " for role <@&" + role_id + ">.", author)
	log.Println("Handler has been added with reaction name " + emoji_name + " on link message " + link_message + " for role id " + role_id + " on guild id " + guild_id)
}
