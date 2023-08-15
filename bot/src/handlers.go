package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func guild_joined(sess *discordgo.Session, gc *discordgo.GuildCreate) {
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

func guild_left(sess *discordgo.Session, gd *discordgo.GuildDelete) {
	guild_id := gd.Guild.ID

	var guilds []guild
	err := db.NewSelect().Model(&guilds).Where("guild_id = ?", guild_id).Scan(ctx)
	if err != nil { log.Println(err) }

	if len(guilds) > 0 {
		del_level := &level{}
		_, err = db.NewDelete().Model(del_level).
					Where("guild_id = ?", guild_id).
					Exec(ctx)
		if err != nil { log.Fatal(err) }

		del_handler_reaction_role := &handler_reaction_role{}
		_, err = db.NewDelete().Model(del_handler_reaction_role).
					Where("guild_id = ?", guild_id).
					Exec(ctx)
		if err != nil { log.Fatal(err) }

		del_channel_for_action := &channel_for_action{}
		_, err = db.NewDelete().Model(del_channel_for_action).
					Where("guild_id = ?", guild_id).
					Exec(ctx)
		if err != nil { log.Fatal(err) }

		del_role_admin := &role_admin{}
		_, err = db.NewDelete().Model(del_role_admin).
					Where("guild_id = ?", guild_id).
					Exec(ctx)
		if err != nil { log.Fatal(err) }

		del_youtube_live_role := &youtube_live_role{}
		_, err = db.NewDelete().Model(del_youtube_live_role).
					Where("guild_id = ?", guild_id).
					Exec(ctx)
		if err != nil { log.Fatal(err) }

		del_youtube_video_role := &youtube_video_role{}
		_, err = db.NewDelete().Model(del_youtube_video_role).
					Where("guild_id = ?", guild_id).
					Exec(ctx)
		if err != nil { log.Fatal(err) }

		del_guild := guilds[0]
		_, err = db.NewDelete().Model(&del_guild).
					Where("guild_id = ?", del_guild.Guild_ID).
					Exec(ctx)
		if err != nil { log.Fatal(err) }
		if err == nil { log.Println("Guild id " + guild_id + " deleted.") }
	}
}

func message_posted(sess *discordgo.Session, m *discordgo.MessageCreate) {
	guild_id := m.Message.GuildID
	author := m.Message.Author
	user_id := author.ID

	if author.Bot { return }

	if m.Message.Type == discordgo.MessageTypeGuildMemberJoin { return }

	var channels_for_actions []channel_for_action
	err := db.NewSelect().Model(&channels_for_actions).
			Where("action_id = ? AND guild_id = ?", get_action_db_by_name("Levels").id, guild_id).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	if len(channels_for_actions) == 0 { return }

	channel_id := channels_for_actions[0].Channel_ID

	var users []level
	err = db.NewSelect().Model(&users).
			Where("user_id = ? AND guild_id = ?", user_id, guild_id).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	if len(users) == 0 {
		level_calculated := calcul_level_with_nb_messages(1)

		new_user := &level{User_ID: user_id, Guild_ID: guild_id, Level: level_calculated}
		_, err = db.NewInsert().Model(new_user).Ignore().Exec(ctx)
		if err != nil {
			log.Println(err)
		} else { log.Println("User id " + user_id + " registered with guild id " + guild_id + " in level table!") }

		if int(level_calculated) > 0 { levels_message(sess, channel_id, new_user, int(level_calculated)) }
	} else {
		user := users[0]
		user.Nb_Msg += 1
		
		level_calculated := calcul_level_with_nb_messages(user.Nb_Msg)

		if level_calculated > user.Level {
			if int(level_calculated) > int(user.Level) { levels_message(sess, channel_id, &user, int(level_calculated)) }
			user.Level = level_calculated
		}
		
		_, err := db.NewUpdate().Model(&user).Column("nb_msg", "level").Where("id = ?", user.ID).Exec(ctx)
		if err != nil { log.Println(err) }
		if err == nil { log.Println("Nb messages of user id " + user_id + " with guild id " + guild_id + " updated in level table!") }
	}
}

func handler_reaction_to_add_role(sess *discordgo.Session, m *discordgo.MessageReactionAdd,) {	
	msg_id := m.MessageReaction.MessageID
	reaction_id := m.MessageReaction.Emoji.ID
	reaction_name := m.MessageReaction.Emoji.Name

	user_id := m.MessageReaction.UserID

	guild_id := m.MessageReaction.GuildID

	var handlers []handler_reaction_role
	err := db.NewSelect().Model(&handlers).
			Where("msg_id = ? AND ((reaction_id IS NOT NULL AND reaction_id = ?) AND reaction_name = ?)", msg_id, reaction_id, reaction_name).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	if len(handlers) > 0 {
		this_handler := handlers[0]

		err = sess.GuildMemberRoleAdd(this_handler.Guild_ID, user_id, this_handler.Role_ID)
		if err != nil {
			log_message(sess, guild_id, "is probably too low in the guild and can't give the role <@&" + this_handler.Role_ID + ">. Try to give her a higher place.\nIf the problem persists, please try to contact her owner.")
			log.Println(err)
		}

		log_message(sess, guild_id, "adds the role <@&" + this_handler.Role_ID + "> to <@" + user_id + ">")
	}
}

func handler_reaction_to_delete_role(sess *discordgo.Session, m *discordgo.MessageReactionRemove,) {
	msg_id := m.MessageReaction.MessageID
	reaction_id := m.MessageReaction.Emoji.ID
	reaction_name := m.MessageReaction.Emoji.Name

	user_id := m.MessageReaction.UserID

	guild_id := m.MessageReaction.GuildID

	var handlers []handler_reaction_role
	err := db.NewSelect().Model(&handlers).
			Where("msg_id = ? AND ((reaction_id IS NOT NULL AND reaction_id = ?) AND reaction_name = ?)", msg_id, reaction_id, reaction_name).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	if len(handlers) > 0 {
		this_handler := handlers[0]

		err = sess.GuildMemberRoleRemove(this_handler.Guild_ID, user_id, this_handler.Role_ID)
		if err != nil {
			log_message(sess, guild_id, "is probably too low in the guild and can't remove the role <@&" + this_handler.Role_ID + ">. Try to give her a higher place.\nIf the problem persists, please try to contact her owner.")
			log.Println(err)
		}

		log_message(sess, guild_id, "removes the role <@&" + this_handler.Role_ID + "> to <@" + user_id + ">")
	}
}
