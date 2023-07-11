package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func blacklist_command(sess *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	author := i.Member.User

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	user_to_blacklist := "<@" + optionMap["user_to_blacklist"].UserValue(nil).ID + ">"
	reason := optionMap["reason"].StringValue()

	// SEND BLACKLIST MESSAGE IN APPROPRIATE CHANNEL
	blacklist_chan_id := get_env_var("BLACKLIST_CHAN_ID")
	embed := discordgo.MessageEmbed{
		Title:       "Blacklisted user",
		Description: "This user has been blacklisted",
		/*Color: 16711680,*/
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

	_, err := sess.ChannelMessageSendEmbed(blacklist_chan_id, &embed)
	if err != nil { log.Fatal(err) }

	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	err = sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Flags:		discordgo.MessageFlagsEphemeral,
				Content:	"User " + user_to_blacklist + " added to blacklist",
			},
		},)
	if err != nil { log.Fatal(err) }
}
