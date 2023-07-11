package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func blacklist_command(sess *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
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
	message := "Blacklist:\n- user: " + user_to_blacklist + "\n- reason: " + reason
	sess.ChannelMessageSend(blacklist_chan_id, message)

	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Flags:		discordgo.MessageFlagsEphemeral,
				Content:	"User " + user_to_blacklist + " added to blacklist",
			},
		},)
	if err != nil { log.Fatal(err) }
}
