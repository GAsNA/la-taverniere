package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func no_live_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User
	
	roles := i.Member.Roles
	is_admin := false
	for i := 0; i < len(roles); i++ {
		if is_role_admin(roles[i]) {
			is_admin = true
			break
		}
	}

	// CAN'T USE THIS COMMAND
	if !is_admin {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Flags:		discordgo.MessageFlagsEphemeral,
					Content:	"You do not have the right to use this command.",
				},
			},)
		if err != nil { log.Fatal(err) }

		log_message(sess, "tried make a no live announcement, but <@" + author.ID + "> to not have the right.")

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
			err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData {
						Flags:		discordgo.MessageFlagsEphemeral,
						Content:	"The date does not have the good format. Use dd/mm/yyyy.",
					},
				},)
			if err != nil { log.Fatal(err) }

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

	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Flags:		discordgo.MessageFlagsEphemeral,
				Content:	"No live message made.",
			},
		},)
	if err != nil { log.Fatal(err) }

	log_message(sess, "added a no live message to <#" + no_live_chan_id + ">")
}
