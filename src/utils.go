package main

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/bwmarrin/discordgo"
)

func get_env_var(key string) string {
	err := godotenv.Load(".env")
	if err != nil { log.Fatal(err) }

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
