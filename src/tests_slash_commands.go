package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func slash_command_hello_world(sess *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Content: "Hello world!",
			},
		},)
	if err != nil { log.Fatal(err) }
	//sess.ChannelMessageSend("1127992339227496599", "Hello world!")
}
