package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func config_channels(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string) {
	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options[0].Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	action_id := optionMap["action"].IntValue()
	channel_id := optionMap["channel"].ChannelValue(nil).ID

	action := get_action_db_by_id(action_id)

	// PROVISIONAL
	if guild_id != get_env_var("GUILD_ID") && (action_id == get_action_db_by_name("Youtube Live Announcements").id || action_id == get_action_db_by_name("Youtube Video Announcements").id) {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "This option in this command is not open for now...")
		return
	}

	// VERIFICATION IF ENTER ALREADY EXISTS
	var channels_for_actions []channel_for_action
	err := db.NewSelect().Model(&channels_for_actions).
			Where("action_id = ? AND guild_id = ?", action_id, guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

		// UPDATE IF ALREADY EXIST
	if len(channels_for_actions) > 0 {
		if channels_for_actions[0].Channel_ID != channel_id {
			new_channel_for_action := channels_for_actions[0]
			new_channel_for_action.Channel_ID = channel_id
			_, err = db.NewUpdate().Model(&new_channel_for_action).
						Column("channel_id").
						Where("id = ?", new_channel_for_action.ID).
						Exec(ctx)
			if err != nil { log.Println(err); return }

			interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Channel <# " + channel_id + "> updated for action \"" + action.name + "\".")
			log_message(sess, guild_id, "changed to channel <#" + channel_id + "> for action \"" + action.name + "\"", author)
			log.Println("Channel id <#" + channel_id + " has been added for action \"" + action.name + "\" on guild id " + guild_id)
		} else {
			interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "This configuration already exists.")
		}
		return
	}

	// CREATE ENTER IN DB
	new_channel_for_action := &channel_for_action{
								Channel_ID: channel_id, Action_ID: action_id,
								Guild_ID: guild_id,
							}
	_, err = db.NewInsert().Model(new_channel_for_action).Ignore().Exec(ctx)
	if err != nil { log.Println(err); return }

	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Channel <# " + channel_id + "> added for action \"" + action.name + "\".")
	log_message(sess, guild_id, "added channel <#" + channel_id + "> for action \"" + action.name + "\"", author)

	log.Println("Channel id <#" + channel_id + " has been had for action \"" + action.name + "\" on guild id " + guild_id)
}

func config_admins(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string) {
	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options[0].Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	role_id := optionMap["role"].RoleValue(nil, guild_id).ID

	// VERIFICATION IF ENTER ALREADY EXISTS
	var roles_admins []role_admin
	err := db.NewSelect().Model(&roles_admins).
			Where("role_id = ? AND guild_id = ?", role_id, guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

	// IF NOT EXISTS INTEGERATE TO DB
	if len(roles_admins) == 0 {
		new_role_admin := &role_admin{Role_ID: role_id, Guild_ID: guild_id}
		_, err = db.NewInsert().Model(new_role_admin).Ignore().Exec(ctx)
		if err != nil { log.Println(err); return }

		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Role <@&" + role_id + "> is now concidered as admin.")
		log_message(sess, guild_id, "added role <@&" + role_id + "> as admin.", author)
		log.Println("Role id <@&" + role_id + "> has been added as admin on guild id " + guild_id)
		return
	}

	// IF EXISTS REMOVE FROM DB
	del_role_admin := roles_admins[0]
	_, err = db.NewDelete().Model(&del_role_admin).
				Where("id = ?", del_role_admin.ID).
				Exec(ctx)
	if err != nil { log.Println(err); return }

	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Role <@&" + role_id + "> is now removed from admin roles.")
	log_message(sess, guild_id, "removed role <@&" + role_id + "> from admin.", author)
	log.Println("Role id <@&" + role_id + "> has been removed from admin on guild id " + guild_id)
}

func config_youtube_roles_live(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string, role_id string) {
	// VERIFICATION IF ENTER ALREADY EXISTS
	var youtube_live_roles []youtube_live_role
	err := db.NewSelect().Model(&youtube_live_roles).
			Where("role_id = ? AND guild_id = ?", role_id, guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

	// IF NOT EXISTS INTEGERATE TO DB
	if len(youtube_live_roles) == 0 {
		new_youtube_live_role := &youtube_live_role{Role_ID: role_id, Guild_ID: guild_id}
		_, err = db.NewInsert().Model(new_youtube_live_role).Ignore().Exec(ctx)
		if err != nil { log.Println(err); return }

		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Role <@&" + role_id + "> will now be ping for each youtube live announcements.")
		log_message(sess, guild_id, "added role <@&" + role_id + "> as ping for youtube live announcements.", author)
		log.Println("Role id <@&" + role_id + "> has been added as ping for yt live announcements on guild id " + guild_id)
		return
	}

	// IF EXISTS REMOVE FROM DB
	del_youtube_live_role := youtube_live_roles[0]
	_, err = db.NewDelete().Model(&del_youtube_live_role).
				Where("id = ?", del_youtube_live_role.ID).
				Exec(ctx)
	if err != nil { log.Println(err); return }

	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Role <@&" + role_id + "> will no longer be ping for youtube live announcements.")
	log_message(sess, guild_id, "removed role <@&" + role_id + "> from list of role ping for youtube live announcements.", author)
	log.Println("Role id <@&" + role_id + "> has been removed as ping for yt live announcements on guild id " + guild_id)
}

func config_youtube_roles_video(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string, role_id string) {
	// VERIFICATION IF ENTER ALREADY EXISTS
	var youtube_video_roles []youtube_video_role
	err := db.NewSelect().Model(&youtube_video_roles).
			Where("role_id = ? AND guild_id = ?", role_id, guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

	// IF NOT EXISTS INTEGERATE TO DB
	if len(youtube_video_roles) == 0 {
		new_youtube_video_role := &youtube_video_role{Role_ID: role_id, Guild_ID: guild_id}
		_, err = db.NewInsert().Model(new_youtube_video_role).Ignore().Exec(ctx)
		if err != nil { log.Println(err); return }

		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Role <@&" + role_id + "> will now be ping for each youtube video announcements.")
		log_message(sess, guild_id, "added role <@&" + role_id + "> as ping for youtube video announcements.", author)
		log.Println("Role id <@&" + role_id + "> has been added as ping for yt video announcements on guild id " + guild_id)
		return
	}

	// IF EXISTS REMOVE FROM DB
	del_youtube_video_role := youtube_video_roles[0]
	_, err = db.NewDelete().Model(&del_youtube_video_role).
				Where("id = ?", del_youtube_video_role.ID).
				Exec(ctx)
	if err != nil { log.Println(err); return }

	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "Role <@&" + role_id + "> will no longer be ping for youtube video announcements.")
	log_message(sess, guild_id, "removed role <@&" + role_id + "> from list of role ping for youtube video announcements.", author)
	log.Println("Role id <@&" + role_id + "> has been removed as ping for yt video announcements on guild id " + guild_id)
}

func config_youtube_roles(sess *discordgo.Session, i *discordgo.InteractionCreate, author *discordgo.User, guild_id string) {
	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options[0].Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	youtube_announcements := optionMap["youtube-announcements"].StringValue()
	role_id := optionMap["role"].RoleValue(nil, guild_id).ID

	// WHICH ANNOUNCEMENTS
	switch youtube_announcements {
		case "live":
			config_youtube_roles_live(sess, i, author, guild_id, role_id)
		case "video":
			config_youtube_roles_video(sess, i, author, guild_id, role_id)
	}
}

func config_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User

	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	admin, err := is_admin(sess, i.Member, guild_id)
	if err != nil { log.Println(err); return }
	if !admin {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried to use the config command, but <@" + author.ID + "> to not have the right.")

		return
	}

	// WHICH SUBCOMMAND
	switch i.ApplicationCommandData().Options[0].Name {
		case "channels":
			config_channels(sess, i, author, guild_id)
		case "admins":
			config_admins(sess, i, author, guild_id)
		case "youtube-roles":
			// PROVISIONAL
			if guild_id != get_env_var("GUILD_ID") {
				interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "This command is not open for now...")
				return
			}
			config_youtube_roles(sess, i, author, guild_id)
	}
}
