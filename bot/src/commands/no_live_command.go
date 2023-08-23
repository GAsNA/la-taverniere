package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func no_live_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User

	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND
	admin, err := is_admin(sess, i.Member, guild_id)
	if err != nil { log.Println(err); return }
	if !admin {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried make a no live announcement, but <@" + author.ID + "> to not have the right.")

		return
	}

	// PROVISIONAL
	if guild_id != get_env_var("GUILD_ID") {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "This command is not open for now...")
		return
	}
	
	// VERIFICATION IF CHANNEL IS CONFIGURATE
	var channels_for_actions []channel_for_action
	err = db.NewSelect().Model(&channels_for_actions).
			Where("action_id = ? AND guild_id = ?", get_action_db_by_name("Youtube Live Announcements").id, guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

	if len(channels_for_actions) == 0 {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "This command needs to be configurate with ``/config``. Choose the action ``Youtube Live Announcements``.")
		return
	}

	no_live_chan_id := channels_for_actions[0].Channel_ID

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var youtube_live_roles []youtube_live_role
	err = db.NewSelect().Model(&youtube_live_roles).
			Where("guild_id = ?", guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }


	var ping_role_ids []string
	for i := 0; i < len(youtube_live_roles); i++ {
		ping_role_ids = append(ping_role_ids, youtube_live_roles[i].Role_ID)
	}

	message := ""

	for i := 0; i < len(ping_role_ids); i++ {
		message += "<@&" + ping_role_ids[i] + ">"
	}

	if len(ping_role_ids) > 0 {
		message += "\n"
	}

	if len(optionMap) > 0 {
		date := optionMap["date"].StringValue()

		if !is_good_format_date(date) {
			interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "The date does not have the good format. Use dd/mm/yyyy.")

			return
		}

		message += "Pas de live youtube prévu jusqu'au " + date + ". Désolé !"

		_, err = sess.ChannelMessageSend(no_live_chan_id, message)
		if err != nil { log.Println(err); return }
	} else {
		message += "Pas de live youtube aujourd'hui. Désolé !"

		_, err = sess.ChannelMessageSend(no_live_chan_id, message)
		if err != nil { log.Println(err); return }
	}

	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "No live message made.")

	log_message(sess, guild_id, "added a no live message to <#" + no_live_chan_id + ">.", author)
}
