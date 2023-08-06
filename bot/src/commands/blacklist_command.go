package main

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func blacklist_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User
	
	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	if !is_admin(sess, i.Member, guild_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "You do not have the right to use this command.")
		log_message(sess, "tried to add someone to the blacklist, but <@" + author.ID + "> to not have the right.")

		return
	}

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	user_to_blacklist_id := optionMap["user"].UserValue(nil).ID
	user_to_blacklist := "<@" + user_to_blacklist_id + ">"
	reason := optionMap["reason"].StringValue()

	//CAN'T BAN IF USER TO BLACKLIST IS THE BOT
	if user_to_blacklist_id == sess.State.User.ID {
		ephemeral_response_for_interaction(sess, i.Interaction, "You can't ban and add to the blacklist the bot.")
		log_message(sess, "can't ban and add themself to the blacklist.", author)

		return 
	}

	// BAN USER
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

	_, err = sess.ChannelMessageSendEmbed(blacklist_chan_id, &embed)
	if err != nil { log.Fatal(err) }

	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	ephemeral_response_for_interaction(sess, i.Interaction, "User " + user_to_blacklist + " added to blacklist.")

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, "banned " + user_to_blacklist + " and added them to the blacklist.", author)
}
