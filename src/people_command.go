package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func people_command(sess *discordgo.Session, i *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) {
	
	options := i.ApplicationCommandData().Options
	optionsMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionsMap[opt.Name] = opt
	}

	people_name := optionsMap["people"].StringValue()

	message := people_name + " : "

	switch people_name {
		case "Ray":
			message += "c'est le dark sasuke du serv"
		case "Feitan":
			message += "c'est la petite chienne de GAsNa"
		case "Ukyim":
			message += "ca va elle est gentille, vous pouvez lui faire confiance"
		case "Kentaro":
			message += "c'est le casse-couille du serv"
		case "GAsNa":
			message += "c'est moi qui ait code le bot, tu penses bien que je vais rien dire de mal sur moi"
	}

	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Content:	message,
			},
		},)
	if err != nil { log.Fatal(err) }
}
