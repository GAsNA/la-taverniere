package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func kick_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
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
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Flags:		discordgo.MessageFlagsEphemeral,
					Content:	"You do not have the right to use this command.",
				},
			},)
		if err != nil { log.Fatal(err) }

		log_message(sess, "tried to add someone to the blacklist, but <@" + author.ID + "> to not have the right.")

		return
	}

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	user_to_kick_id := optionMap["user_to_kick"].UserValue(nil).ID
	user_to_kick := "<@" + user_to_kick_id + ">"
	reason := optionMap["reason"].StringValue()

	//CAN'T BAN IF USER TO BLACKLIST IS THE BOT
	if user_to_kick_id == sess.State.User.ID {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Flags:		discordgo.MessageFlagsEphemeral,
					Content:	"You can't kick the bot.",
				},
			},)
		if err != nil { log.Fatal(err) }

		log_message(sess, "tried to kick someone but can't kick themself.", author)

		return 
	}

	// BAN USER
	guild_id := get_env_var("DISCORD_GUILD_ID")
	err := sess.GuildMemberDeleteWithReason(guild_id, user_to_kick_id, reason)
	if err != nil { log.Fatal(err) }

	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Flags:		discordgo.MessageFlagsEphemeral,
				Content:	"User " + user_to_kick + " has been kick.",
			},
		},)
	if err != nil { log.Fatal(err) }

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, "kicked " + user_to_kick + ".", author)
}
