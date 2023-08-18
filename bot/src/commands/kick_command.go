package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func kick_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User
	
	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	if !is_admin(sess, i.Member, guild_id) {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried to kick someone, but <@" + author.ID + "> to not have the right.")

		return
	}

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	user_to_kick_id := optionMap["user"].UserValue(nil).ID
	user_to_kick := "<@" + user_to_kick_id + ">"
	reason := optionMap["reason"].StringValue()

	//CAN'T BAN IF USER TO BLACKLIST IS THE BOT
	if user_to_kick_id == sess.State.User.ID {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "You can't kick the bot.")
		return 
	}

	// BAN USER
	err := sess.GuildMemberDeleteWithReason(guild_id, user_to_kick_id, reason)
	if err != nil { log.Fatal(err) }

	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "User " + user_to_kick + " has been kick.")

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, guild_id, "kicked " + user_to_kick + ".", author)
}
