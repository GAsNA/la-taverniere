package main

import (
	"github.com/bwmarrin/discordgo"
)

func handlerHello(sess *discordgo.Session) {
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID { return }

		if m.Content == "Hello" {
			s.ChannelMessageSend(m.ChannelID, "world!")
		}
	})
}

func handlerWorld(sess *discordgo.Session) {
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID { return }

		if m.Content == "world!" {
			s.ChannelMessageSend(m.ChannelID, "Hello")
		}
	})
}

