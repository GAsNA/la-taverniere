package main

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func blacklist_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
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

	user_to_blacklist_id := optionMap["user_to_blacklist"].UserValue(nil).ID
	user_to_blacklist := "<@" + user_to_blacklist_id + ">"
	reason := optionMap["reason"].StringValue()

	//CAN'T BAN IF USER TO BLACKLIST IS THE BOT
	if user_to_blacklist_id == sess.State.User.ID {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Flags:		discordgo.MessageFlagsEphemeral,
					Content:	"You can't ban and add to the blacklist the bot.",
				},
			},)
		if err != nil { log.Fatal(err) }

		log_message(sess, "can't ban and add themself to the blacklist. Requested by <@" + author.ID + ">.")

		return 
	}

	// BAN USER
	guild_id := get_env_var("DISCORD_GUILD_ID")
	err := sess.GuildBanCreateWithReason(guild_id, user_to_blacklist_id, reason, 0)
	if err != nil { log.Fatal(err) }

	// SEND BLACKLIST MESSAGE IN APPROPRIATE CHANNEL
	blacklist_chan_id := get_env_var("BLACKLIST_CHAN_ID")
	embed := discordgo.MessageEmbed{
		Title:       "Blacklisted user",
		Description: "This user has been blacklisted",
		Timestamp: time.Now().Format(time.RFC3339),
		Color: get_color_by_name("Black").code,
		Fields: []*discordgo.MessageEmbedField {
			{
				Name:  "User",
				Value: user_to_blacklist,
			},
			{
				Name:  "Reason",
				Value: reason,
			},
		},
		Footer: &discordgo.MessageEmbedFooter {
			Text: "Requested by " + author.Username,
			IconURL: author.AvatarURL(""),
		},
	}

	_, err_msg := sess.ChannelMessageSendEmbed(blacklist_chan_id, &embed)
	if err_msg != nil { log.Fatal(err_msg) }

	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Flags:		discordgo.MessageFlagsEphemeral,
				Content:	"User " + user_to_blacklist + " added to blacklist",
			},
		},)
	if err != nil { log.Fatal(err) }

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, "banned someone and added them to the blacklist.")
}
