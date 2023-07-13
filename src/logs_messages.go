package main

import (
	"time"
	"log"

	"github.com/bwmarrin/discordgo"
)

func log_message(sess *discordgo.Session, log_str string) {

	// SEND log MESSAGE IN APPROPRIATE CHANNEL
	logs_chan_id := get_env_var("LOGS_CHAN_ID")
	embed := discordgo.MessageEmbed{
		Title:       "Log",
		Description: sess.State.User.Username + " " + log_str,
		Timestamp: time.Now().Format(time.RFC3339),
		Color: 3447003,
		Footer: &discordgo.MessageEmbedFooter {
			Text: "Requested by " + sess.State.User.Username,
					IconURL: sess.State.User.AvatarURL(""),
		},
	}

	_, err := sess.ChannelMessageSendEmbed(logs_chan_id, &embed)
	if err != nil { log.Fatal(err) }
}
