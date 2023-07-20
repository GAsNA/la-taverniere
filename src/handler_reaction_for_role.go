package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/forPelevin/gomoji"
)

func check_reaction(reaction string, emoji_name *string, emoji_id *string) bool {
	find_all := gomoji.FindAll(reaction)
	if len(find_all) > 1 { return false }

	if len(find_all) == 1 {
		reaction_without_emoji := gomoji.RemoveEmojis(reaction)
		if len(reaction_without_emoji) > 0 { return false }
		
		if len(reaction_without_emoji) == 0 {
			*emoji_name = reaction
			return true
		}
	}

	if len(find_all) == 0 {
		if !strings.HasPrefix(reaction, "<:") { return false }
		reaction = strings.TrimLeft(reaction, "<:")

		if reaction[len(reaction) - 1] != '>' { return false }
		reaction = strings.TrimRight(reaction, ">")

		parts := strings.Split(reaction, ":")
		if len(parts) != 2 { return false }

		for i := 0; i < len(parts[1]); i++ {
			if parts[1][i] < '0' || parts[1][i] > '9' { return false }
		}

		*emoji_name = parts[0]
		*emoji_id = parts[1]
		return true
	}

	return false
}

func handler_reaction_for_role_command(sess *discordgo.Session, i *discordgo.InteractionCreate) {
	author := i.Member.User
	
	roles := i.Member.Roles
	is_admin := false
	for i := 0; i < len(roles); i++ {
		if is_role_admin(roles[i]) {
			is_admin = true
			break
		}
	}

	// CAN'T USE THIS COMMAND IF NOT ADMIN
	if !is_admin {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Flags:		discordgo.MessageFlagsEphemeral,
					Content:	"You do not have the right to use this command.",
				},
			},)
		if err != nil { log.Fatal(err) }

		log_message(sess, "tried to add someone to the blacklist, but <@" + author.ID + "> to not have the right.")

		return
	}

	// GET OPTIONS AND MAP
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	guild_id := get_env_var("DISCORD_GUILD_ID")

	link_message := optionMap["link_message"].StringValue()
	reaction := optionMap["reaction"].StringValue()
	role := optionMap["role"].RoleValue(nil, guild_id)
	
	// VERIF LINK
	message_id := get_message_id(link_message, guild_id)
	if message_id == "" {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Flags:		discordgo.MessageFlagsEphemeral,
					Content:	"The link of the message is not at the good format or the message is not in this guild.",
				},
			},)
		if err != nil { log.Fatal(err) }

		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the link of the message is not at the good format or the message is not in this guild.")

		return
	}

	// VERIF REACTION
	emoji_name := ""
	emoji_id := ""
	if !check_reaction(reaction, &emoji_name, &emoji_id) {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Flags:		discordgo.MessageFlagsEphemeral,
					Content:	"The reaction is not at the good format.",
				},
			},)
		if err != nil { log.Fatal(err) }

		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the reaction is not at the good format.")

		return
	}

	// VERIF IF ROLE IS @everyone
	role_id := role.ID
	if role_id == guild_id {
		err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Flags:		discordgo.MessageFlagsEphemeral,
					Content:	"You can't choose the @everyone for the role",
				},
			},)
		if err != nil { log.Fatal(err) }

		log_message(sess, "tried to add a handler to a message to add a role with reaction, but the role chose was @everyone.")

		return
	}

	// TODO VERIF LES CONDITIONS IDENTIQUES
	// OR TODO SUPP HANDLER

	sess.AddHandler(func (sess *discordgo.Session, m *discordgo.MessageReactionAdd,) {
		if m.MessageReaction.MessageID != message_id { return }

		if (emoji_id != "" && m.MessageReaction.Emoji.ID != emoji_id) ||
			(m.MessageReaction.Emoji.Name != emoji_name) { return }

		err_add_role := sess.GuildMemberRoleAdd(guild_id, m.MessageReaction.UserID, role_id)
		if err_add_role != nil { log.Fatal(err_add_role) }

		log_message(sess, "add the role <@&" + role_id + "> to <@" + m.MessageReaction.UserID + ">")
	})

	// TODO HANDLER FOR SUPP REACTION

	// RESPOND TO USER WITH EPHEMERAL MESSAGE
	err := sess.InteractionRespond(i.Interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Flags:		discordgo.MessageFlagsEphemeral,
				Content:	"Handler add to " + link_message + " with reaction " + reaction + " for role <@&" + role_id + ">",
			},
		},)
	if err != nil { log.Fatal(err) }

	// ADD LOG IN LOGS CHANNEL
	log_message(sess, "add a handler for add a role with reaction")
}
