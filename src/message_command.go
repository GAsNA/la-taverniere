package main

import (
	"log"
	"time"
	"strings"
	"net/http"
	"io"

	"github.com/bwmarrin/discordgo"
)

func get_io_reader(URL string) io.Reader {
	response, err := http.Get(URL)
	if err != nil { log.Fatal(err) }

	if response.StatusCode != 200 { log.Fatal("Received non 200 response code")	}
	
	return response.Body
}

func message_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User
	
	roles := i.Member.Roles
	is_admin := false
	for i := 0; i < len(roles); i++ {
		if is_role_admin(roles[i]) {
			is_admin = true
			break
		}
	}

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	if !is_admin {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Flags:		discordgo.MessageFlagsEphemeral,
					Content:	"You do not have the right to use this command.",
				},
			},)
		if err != nil { log.Fatal(err) }

		log_message(sess, "tried to send a message via the bot, but <@" + author.ID + "> to not have the right.")

		return
	}

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

		// Options required
	channel_id := optionMap["channel"].ChannelValue(nil).ID
	channel := "<#" + channel_id + ">"
	message := optionMap["message"].StringValue()
	message = strings.ReplaceAll(message, "\\n", "\n")

		// Option not required
	is_embed := false
	if _, ok := optionMap["embed"]; ok {
		is_embed = optionMap["embed"].BoolValue()
	}
	title := ""
	if _, ok := optionMap["title"]; ok {
		title = optionMap["title"].StringValue()
	}
	color := 0
	if _, ok := optionMap["color"]; ok {
		color = int(optionMap["color"].IntValue())
	}
	url_thumbnail := ""
	var attachment *discordgo.MessageAttachment = nil
	if _, ok := optionMap["thumbnail"]; ok {
		attachment_id := optionMap["thumbnail"].Value.(string)
		attachment = i.ApplicationCommandData().Resolved.Attachments[attachment_id]
		url_thumbnail = attachment.URL
	}

	// SEND EMBED IF EMBED IS TRUE ELSE MESSAGE
	if is_embed {
		embed := discordgo.MessageEmbed {
			Title:			title,
			Description:	message,
			Timestamp:		time.Now().Format(time.RFC3339),
			Color:			color,
			Thumbnail:		&discordgo.MessageEmbedThumbnail {
				URL:	url_thumbnail,
			},
		}

		_, err_msg := sess.ChannelMessageSendEmbed(channel_id, &embed)
		if err_msg != nil { log.Fatal(err_msg) }
	} else {
		message = "### " + title + "\n" + message

		files := []*discordgo.File {}
		if url_thumbnail != "" {
			files = []*discordgo.File {
				{
					Name:			attachment.Filename,
					ContentType:	attachment.ContentType,
					Reader:			get_io_reader(url_thumbnail),
				},
			}
		}

		messageSend := discordgo.MessageSend {
			Content:		message,
			Embeds:			[]*discordgo.MessageEmbed {},
			TTS:			true,
			Components:		[]discordgo.MessageComponent {},
			Files:			files,
		}

		_, err_msg := sess.ChannelMessageSendComplex(channel_id, &messageSend)
		if err_msg != nil { log.Fatal(err_msg) }
	}

	// SUCCESSFULLY SENT
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Flags:		discordgo.MessageFlagsEphemeral,
				Content:	"Message successfully send to " + channel + ".",
			},
		},)
	if err != nil { log.Fatal(err) }

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, "send a message to " + channel + ".")
}
