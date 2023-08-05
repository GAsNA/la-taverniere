package main

import (
	"log"
	"os"
	"strings"
	"math"

	"github.com/bwmarrin/discordgo"
	"github.com/forPelevin/gomoji"
)

func get_env_var(key string) string {
	return os.Getenv(key)
}

func is_the_bot(idAuthorMessage string, idBot string) bool {
	if idAuthorMessage == idBot { return true }
	return false
}

func is_role_admin(id_role string) bool {
	admin_role_ids_env := get_env_var("ADMIN_ROLE_IDS")
	admin_role_ids := strings.Split(admin_role_ids_env, ",")

	for i := 0; i < len(admin_role_ids); i++ {
		if admin_role_ids[i] == id_role { return true }
	}

	return false
}

func is_good_format_date(date string) bool {
	if len(date) != 10 {
		return false
	}

	if date[2] != '/' || date[5] != '/' {
		return false
	}

	for i := 0; i < len(date); i++ {
		if i == 2 || i == 5 {
			continue
		}
		if date[i] < '0' || date[i] > '9' {
			return false
		}
	}

	return true
}

func get_discord_message_ids(link string, guild_id *string, channel_id *string, message_id *string) bool {
	discord_link := get_env_var("DISCORD_LINK")

	if !strings.HasPrefix(link, discord_link + "/channels/") { return false }
	link = strings.TrimLeft(link, discord_link + "/channels/")

	parts := strings.Split(link, "/")
	if len(parts) != 3 { return false }

	for i := 0; i < len(parts); i++ {
		for j := 0; j < len(parts[i]); j++ {
			if parts[i][j] < '0' || parts[i][j] > '9' { return false }
		}
	}

	*guild_id = parts[0]
	*channel_id = parts[1]
	*message_id = parts[2]
	
	return true
}

func ephemeral_response_for_interaction(sess *discordgo.Session, interaction *discordgo.Interaction, message string) {
	err := sess.InteractionRespond(interaction, &discordgo.InteractionResponse {
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData {
				Flags:		discordgo.MessageFlagsEphemeral,
				Content:	message,
			},
		},)
	if err != nil { log.Fatal(err) }
}

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

func is_a_registered_handler(link string, reaction string, role *discordgo.Role) bool {
	var handlers []handler_reaction_role
	err := db.NewSelect().Model(&handlers).
			Where("msg_link = ? AND reaction = ? AND role_id = ?", link, reaction, role.ID).
			Scan(ctx)
	if err != nil { log.Fatal(err) }

	if len(handlers) > 0 { return true }

	return false
}

func calcul_level_with_nb_messages(nb_msg int64) int64 {
		// level calcul with
			// (1 + racine(1 + 8 * 15 * nb_msg / 50)) / 2
	return int64((1.0 + math.Sqrt(1.0 + (8.0 * 15.0 * float64(nb_msg) / 50.0))) / 2.0)
}
