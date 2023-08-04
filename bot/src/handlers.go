package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type handler_reaction struct {
    link			string
	message_id		string
    reaction		string
	reaction_id		string
	reaction_name	string
    role			*discordgo.Role
	role_id			string
	guild_id		string
}

var list_handler_reaction []handler_reaction = []handler_reaction{}

func new_guild_joined(sess *discordgo.Session, gc *discordgo.GuildCreate) {
	guild_id := gc.Guild.ID

	var guilds []guild
	err := db.NewSelect().Model(&guilds).Where("guild_id = ?", guild_id).Scan(ctx)
	if err != nil { log.Println(err) }

	if len(guilds) == 0 {
		new_guild := &guild{Guild_ID: guild_id}
		_, err = db.NewInsert().Model(new_guild).Ignore().Exec(ctx)
		if err != nil { log.Println(err) }
		if err == nil { log.Println("Guild id " + guild_id + " registered!") }
	}
}

func handler_reaction_to_add_role(sess *discordgo.Session, m *discordgo.MessageReactionAdd,) {	
	for i := 0; i < len(list_handler_reaction); i++ { 
		this_handler := list_handler_reaction[i]

		if m.MessageReaction.MessageID != this_handler.message_id { continue }

		if (this_handler.reaction_id != "" && m.MessageReaction.Emoji.ID != this_handler.reaction_id) ||
			(m.MessageReaction.Emoji.Name != this_handler.reaction_name) { continue }

		err := sess.GuildMemberRoleAdd(this_handler.guild_id, m.MessageReaction.UserID, this_handler.role_id)
		if err != nil { log.Fatal(err) }

		log_message(sess, "add the role <@&" + this_handler.role_id + "> to <@" + m.MessageReaction.UserID + ">")
		break
	}
}

func handler_reaction_to_delete_role(sess *discordgo.Session, m *discordgo.MessageReactionRemove,) {
	for i := 0; i < len(list_handler_reaction); i++ { 
		this_handler := list_handler_reaction[i]

		if m.MessageReaction.MessageID != this_handler.message_id { continue }

		if (this_handler.reaction_id != "" && m.MessageReaction.Emoji.ID != this_handler.reaction_id) ||
			(m.MessageReaction.Emoji.Name != this_handler.reaction_name) { continue }

		err := sess.GuildMemberRoleRemove(this_handler.guild_id, m.MessageReaction.UserID, this_handler.role_id)
		if err != nil { log.Fatal(err) }

		log_message(sess, "removes the role <@&" + this_handler.role_id + "> to <@" + m.MessageReaction.UserID + ">")
		break
	}
}
