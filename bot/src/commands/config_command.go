package main

import (
	"github.com/bwmarrin/discordgo"
)

func config_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User

	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	if !is_admin(sess, i.Member, guild_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "You do not have the right to use this command.")
		log_message(sess, "tried to use the config command, but <@" + author.ID + "> to not have the right.")

		return
	}

	ephemeral_response_for_interaction(sess, i.Interaction, "I received!")
}
