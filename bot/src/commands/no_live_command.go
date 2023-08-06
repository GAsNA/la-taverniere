package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func no_live_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User

	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND
	if !is_admin(sess, i.Member, guild_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried make a no live announcement, but <@" + author.ID + "> to not have the right.")

		return
	}

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	no_live_chan_id := get_env_var("NO_LIVE_CHAN_ID")
	ping_role_ids_env := get_env_var("PING_YOUTUBE_LIVE_ROLE_IDS")
	ping_role_ids := strings.Split(ping_role_ids_env, ",")

	message := ""

	for i := 0; i < len(ping_role_ids); i++ {
		message += "<@&" + ping_role_ids[i] + ">"
	}

	if len(ping_role_ids) > 0 {
		message += "\n"
	}

	if len(optionMap) > 0 {
		date := optionMap["date"].StringValue()

		if !is_good_format_date(date) {
			ephemeral_response_for_interaction(sess, i.Interaction, "The date does not have the good format. Use dd/mm/yyyy.")

			return
		}

		message += "Pas de live youtube prévu jusqu'au " + date + ". Désolé !"

		_, err := sess.ChannelMessageSend(no_live_chan_id, message)
		if err != nil { log.Fatal(err) }
	} else {
		message += "Pas de live youtube aujourd'hui. Désolé !"

		_, err := sess.ChannelMessageSend(no_live_chan_id, message)
		if err != nil { log.Fatal(err) }
	}

	ephemeral_response_for_interaction(sess, i.Interaction, "No live message made.")

	log_message(sess, guild_id, "added a no live message to <#" + no_live_chan_id + ">.", author)
}
