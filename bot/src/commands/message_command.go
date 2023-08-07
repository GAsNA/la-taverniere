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
	
	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	if !is_admin(sess, i.Member, guild_id) {
		ephemeral_response_for_interaction(sess, i.Interaction, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried to send a message via the bot, but <@" + author.ID + "> to not have the right.")

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
	message = strings.ReplaceAll(message, "\\t", "\t")

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
	if _, ok := optionMap["thumbnail"]; ok {
		thumbnail_id := optionMap["thumbnail"].Value.(string)
		thumbnail := i.ApplicationCommandData().Resolved.Attachments[thumbnail_id]
		url_thumbnail = thumbnail.URL
	}
	var attachment *discordgo.MessageAttachment = nil
	if _, ok := optionMap["attachment"]; ok {
		attachment_id := optionMap["attachment"].Value.(string)
		attachment = i.ApplicationCommandData().Resolved.Attachments[attachment_id]
	}

	message_to_send := ""
	embeds := []*discordgo.MessageEmbed {}

	if is_embed {
		embeds = []*discordgo.MessageEmbed {
			{
				Title:			title,
				Description:	message,
				Timestamp:		time.Now().Format(time.RFC3339),
				Color:			color,
				Thumbnail:		&discordgo.MessageEmbedThumbnail {
					URL:	url_thumbnail,
				},
			},
		}
	} else {
		message_to_send = "### " + title + "\n" + message
	}

	files := []*discordgo.File {}
	if attachment != nil {
		files = []*discordgo.File {
			{
				Name:			attachment.Filename,
				ContentType:	attachment.ContentType,
				Reader:			get_io_reader(attachment.URL),
			},
		}
	}

	data_to_send := discordgo.MessageSend {
		Content:		message_to_send,
		Embeds:			embeds,
		TTS:			true,
		Components:		[]discordgo.MessageComponent {},
		Files:			files,
	}

	msg_send, err := sess.ChannelMessageSendComplex(channel_id, &data_to_send)
	if err != nil { log.Fatal(err) }

	// RESPONSE MESSAGE FOR SUCCESSFULLY SENT
	ephemeral_response_for_interaction(sess, i.Interaction, "Message successfully send to " + channel + ".")

	link_msg_send := get_env_var("DISCORD_LINK") + "/channels/" + i.Interaction.GuildID + "/" + msg_send.ChannelID + "/" + msg_send.ID
	// ADD LOG IN LOGS CHANNEL
	log_message(sess, guild_id, "send this message " + link_msg_send + " to " + channel + ".", author)
}
