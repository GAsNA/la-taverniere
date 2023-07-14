package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func salope_command(sess *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Content:	"Kentaro, cette petite salope, il a voulu la commande je lui donne",
			},
		},)
	if err != nil { log.Fatal(err) }
}
