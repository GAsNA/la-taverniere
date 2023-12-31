package main

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func blacklist_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User
	
	guild_id := i.Interaction.GuildID

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	admin, err := is_admin(sess, i.Member, guild_id)
	if err != nil { log.Println(err); return }
	if !admin {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "You do not have the right to use this command.")
		log_message(sess, guild_id, "tried to add someone to the blacklist, but <@" + author.ID + "> to not have the right.")

		return
	}

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	user_to_blacklist_id := optionMap["user"].UserValue(nil).ID
	user_to_blacklist := "<@" + user_to_blacklist_id + ">"
	reason := optionMap["reason"].StringValue()

	//CAN'T BAN IF USER TO BLACKLIST IS THE BOT
	if user_to_blacklist_id == sess.State.User.ID {
		interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "You can't ban and add to the blacklist the bot.")
		return 
	}

	// BAN USER
	err = sess.GuildBanCreateWithReason(guild_id, user_to_blacklist_id, reason, 0)
	if err != nil { log.Println(err); return }

	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	interaction_respond(sess, i.Interaction, discordgo.InteractionResponseChannelMessageWithSource, true, "User " + user_to_blacklist + " has been ban.")

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, guild_id, "banned " + user_to_blacklist + ".", author)

	log.Println("User id " + user_to_blacklist_id + " has been ban of guild id " + guild_id)

	// SEND BLACKLIST MESSAGE IN APPROPRIATE CHANNEL
	var channels_for_actions []channel_for_action
	err = db.NewSelect().Model(&channels_for_actions).
			Where("action_id = ? AND guild_id = ?", get_action_db_by_name("Blacklist Logs").id, guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

	if len(channels_for_actions) == 0 { return }

	blacklist_chan_id := channels_for_actions[0].Channel_ID

	embed := discordgo.MessageEmbed{
		Title:       "Blacklisted user",
		Description: "This user has been blacklisted",
		Timestamp: time.Now().Format(time.RFC3339),
		Color: get_color_by_name("Black").code,
		Fields: []*discordgo.MessageEmbedField {
			{
				Name:  "User",
				Value: user_to_blacklist,
			},
			{
				Name:  "Reason",
				Value: reason,
			},
		},
		Footer: &discordgo.MessageEmbedFooter {
			Text: "Requested by " + author.Username,
			IconURL: author.AvatarURL(""),
		},
	}

	_, err = sess.ChannelMessageSendEmbed(blacklist_chan_id, &embed)
	if err != nil { log.Println(err) }
}
