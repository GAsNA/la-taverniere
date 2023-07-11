package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func blacklist(sess *discordgo.Session) {
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if isTheBot(m.Author.ID, s.State.User.ID) { return }

		args := strings.Split(m.Content, " ")
		
		//ephemeral=True for only visible by you

		if beginWithPrefix(args[0]) && trimPrefix(args[0]) == "blacklist" /*len(args) == 3*/ {
			s.ChannelMessageSend(m.ChannelID, "world!")
		}
	})
}
