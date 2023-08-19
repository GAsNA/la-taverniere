package main

import (
	"time"
	"log"

	"github.com/bwmarrin/discordgo"
)

func log_message(sess *discordgo.Session, guild_id string, log_str string, user ...*discordgo.User) {
	// VERIF IF LOG CHANNEL IS SET
	action := get_action_db_by_name("Logs")
	var channels_for_actions []channel_for_action
	err := db.NewSelect().Model(&channels_for_actions).
			Where("action_id = ? AND guild_id = ?", action.id, guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

	if len(channels_for_actions) == 0 { return }

	// LOG MESSAGE
	logs_chan_id := channels_for_actions[0].Channel_ID
	
	message := "I " + log_str

	if len(user) > 0 {
		message += "\nRequested by <@" + user[0].ID + ">"
	}

	// SEND LOG MESSAGE IN APPROPRIATE CHANNEL
	embed := discordgo.MessageEmbed{
		Title:       "Log",
		Description: message,
		Timestamp: time.Now().Format(time.RFC3339),
		Color: get_color_by_name("Blue").code,
	}

	_, err = sess.ChannelMessageSendEmbed(logs_chan_id, &embed)
	if err != nil { log.Println(err) }
}

func interaction_respond(sess *discordgo.Session, interaction *discordgo.Interaction, type_res discordgo.InteractionResponseType, is_ephemeral bool, message string, opt ...interface{}) {
	data := &discordgo.InteractionResponseData {
		Content: message,
	}

	if is_ephemeral { data.Flags = discordgo.MessageFlagsEphemeral }

	data.Components = []discordgo.MessageComponent {}
	data.Files = []*discordgo.File {}

	for i := 0; i < len(opt); i++ {
		if _, ok := opt[i].(discordgo.MessageComponent); ok {
			data.Components = append(data.Components, opt[i].(discordgo.MessageComponent))
		} else if _, ok := opt[i].(*discordgo.File); ok {
			data.Files = append(data.Files, opt[i].(*discordgo.File))
		} else {
			log.Print("Type not known: ")
			log.Println(opt[i])
		}
	}

	err := sess.InteractionRespond(interaction, &discordgo.InteractionResponse {
			Type: type_res,
			Data: data,
		},)
	if err != nil { log.Println(err) }
}

func levels_message(sess *discordgo.Session, levels_chan_id string, user *level, level int) error {
	message := "<@" + user.User_ID + "> reached lvl." + strconv.Itoa(level) + "!"
	embed := discordgo.MessageEmbed{
		Title:       "New level!",
		Description: message,
		Timestamp: time.Now().Format(time.RFC3339),
		Color: get_color_by_name("Aqua").code,
	}

	_, err := sess.ChannelMessageSendEmbed(levels_chan_id, &embed)
	if err != nil { return err }

	return nil
}
