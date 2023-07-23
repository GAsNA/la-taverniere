package main

import (
	"time"
	"log"

	"github.com/bwmarrin/discordgo"
)

func log_message(sess *discordgo.Session, log_str string, user ...*discordgo.User) {
	message := sess.State.User.Username + " " + log_str

	if len(user) > 0 {
		message += "\nRequested by <@" + user[0].ID + ">"
	}

	// SEND LOG MESSAGE IN APPROPRIATE CHANNEL
	logs_chan_id := get_env_var("LOGS_CHAN_ID")
	embed := discordgo.MessageEmbed{
		Title:       "Log",
		Description: message,
		Timestamp: time.Now().Format(time.RFC3339),
		Color: get_color_by_name("Blue").code,
	}

	_, err := sess.ChannelMessageSendEmbed(logs_chan_id, &embed)
	if err != nil { log.Fatal(err) }
}
