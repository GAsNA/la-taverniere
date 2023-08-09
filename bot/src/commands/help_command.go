package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func help_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User
	bot := sess.State.User

	title := "👋 Hi! Let me introduce myself..."
	description := "I am La Tavernivere and I am a multi-functionnal bot. I am running multiple taverns ; you better not play with me... 💪\nYou can see below all my supported commands and actions.\nEnjoy !"

	embed := discordgo.MessageEmbed {
		Title:			title,
		Description:	description,
		Color:			get_color_by_name("Bot Red").code,
		Footer:			&discordgo.MessageEmbedFooter {
			Text:		"Requested by " + author.Username,
			IconURL:	author.AvatarURL(""),
		},
		Thumbnail:		&discordgo.MessageEmbedThumbnail {
			URL:		bot.AvatarURL(""),
		},
		Author:	&discordgo.MessageEmbedAuthor {
			Name:		bot.Username,
			IconURL:	bot.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField {
			{
				Name:	"⚙️ /config",
				Value:	"Configurate me by giving me channels for each action I can do or by giving me roles that I can interpret as admins.",
				Inline: true,
			},
			{
				Name:	"☠️ /blacklist",
				Value:	"Ban a user and a message of blacklist to the guild.",
				Inline: true,
			},
			{
				Name:	"💥 /kick",
				Value:	"Kick a user.",
				Inline: true,
			},
			{
				Name:	"💬 /message",
				Value:	"Send a custom message to a channel through me.",
				Inline: true,
			},
			{
				Name:	"📡 /handler-reaction-for-role",
				Value:	"Add or delete an handler that adds a role to each person using the chosen reaction to the chosen message.",
				Inline: true,
			},
			{
				Name:	"🤖 Automatic actions",
				Value:	"- Levels for messages posted by each person.\n- Log messages for each of my actions.",
				Inline: false,
			},
		},
	}

	channel_id := i.ChannelID

	_, err := sess.ChannelMessageSendEmbed(channel_id, &embed)
	if err != nil { log.Fatal(err) }

	ephemeral_response_for_interaction(sess, i.Interaction, "Check out my presentation!")
}
