package main

import (
	"github.com/bwmarrin/discordgo"
)

func config_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	ephemeral_response_for_interaction(sess, i.Interaction, "I received!");
}
