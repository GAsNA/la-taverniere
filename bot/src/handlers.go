package main

import (
	"log"
	"math"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

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
	msg_id := m.MessageReaction.MessageID
	reaction_id := m.MessageReaction.Emoji.ID
	reaction_name := m.MessageReaction.Emoji.Name

	user_id := m.MessageReaction.UserID

	var handlers []handler_reaction_role
	err := db.NewSelect().Model(&handlers).
			Where("msg_id = ? AND ((reaction_id IS NOT NULL AND reaction_id = ?) AND reaction_name = ?)", msg_id, reaction_id, reaction_name).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	if len(handlers) > 0 {
		this_handler := handlers[0]

		err = sess.GuildMemberRoleAdd(this_handler.Guild_ID, user_id, this_handler.Role_ID)
		if err != nil { log.Fatal(err) }

		log_message(sess, "add the role <@&" + this_handler.Role_ID + "> to <@" + user_id + ">")
	}
}

func handler_reaction_to_delete_role(sess *discordgo.Session, m *discordgo.MessageReactionRemove,) {
	msg_id := m.MessageReaction.MessageID
	reaction_id := m.MessageReaction.Emoji.ID
	reaction_name := m.MessageReaction.Emoji.Name

	user_id := m.MessageReaction.UserID

	var handlers []handler_reaction_role
	err := db.NewSelect().Model(&handlers).
			Where("msg_id = ? AND ((reaction_id IS NOT NULL AND reaction_id = ?) AND reaction_name = ?)", msg_id, reaction_id, reaction_name).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	if len(handlers) > 0 {
		this_handler := handlers[0]

		err = sess.GuildMemberRoleRemove(this_handler.Guild_ID, user_id, this_handler.Role_ID)
		if err != nil { log.Fatal(err) }

		log_message(sess, "removes the role <@&" + this_handler.Role_ID + "> to <@" + user_id + ">")
	}
}

func new_message_posted(sess *discordgo.Session, m *discordgo.MessageCreate) {
	guild_id := m.Message.GuildID
	user_id := m.Message.Author.ID

	if is_the_bot(user_id, sess.State.User.ID) { return }

	var users []nb_msg
	err := db.NewSelect().Model(&users).
			Where("user_id = ? AND guild_id = ?", user_id, guild_id).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	if len(users) == 0 {
		new_user := &nb_msg{User_ID: user_id, Guild_ID: guild_id}
		_, err = db.NewInsert().Model(new_user).Ignore().Exec(ctx)
		if err != nil { log.Println(err) }
		if err == nil { log.Println("User id " + user_id + " registered with guild id " + guild_id + " in nb_msg table!") }
	} else {
		user := users[0]
		user.Nb_Msg = user.Nb_Msg + 1
		_, err := db.NewUpdate().Model(&user).Column("nb_msg").Where("id = ?", user.ID).Exec(ctx)
		if err != nil { log.Println(err) }
		if err == nil { log.Println("Nb messages of user id " + user_id + " updated in nb_msg table!") }

		// levels calcul with
			// (1 + racine(1 + 8 * nb_msg / 50)) / 2
		levels := ((1 + math.Sqrt(float64(1 + 8 * user.Nb_Msg / 50))) / 2)
		levels_int := int(levels)
		log_message(sess, "gives xp to <@" + user_id + ">. They are levels" + strconv.Itoa(levels_int) + ".")
	}
}
