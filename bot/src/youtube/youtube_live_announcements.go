package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/youtube/v3"
)

func send_youtube_live_announcement(sess *discordgo.Session, live *youtube.SearchResult) {
	message := ""
	channel_id := get_env_var("LIVE_CHAN_ID")

	ping_role_ids_env := get_env_var("PING_YOUTUBE_LIVE_ROLE_IDS")
	ping_role_ids := strings.Split(ping_role_ids_env, ",")

	for i := 0; i < len(ping_role_ids); i++ {
		message += "<@&" + ping_role_ids[i] + ">"
	}

	if len(ping_role_ids) > 0 {
		message += "\n"
	}

	switch live.Snippet.LiveBroadcastContent {
		case "upcoming":
			message += "Un live se prépare sur la chaine de " + live.Snippet.ChannelTitle + "...\n"
		case "live":
			message += live.Snippet.ChannelTitle + " est en live ! Viens voir !\n"
	}

	message += get_env_var("YOUTUBE_LINK") + "/watch?v=" + live.Id.VideoId
	_, err := sess.ChannelMessageSend(channel_id, message)
	if err != nil { log.Fatal(err) }

	guild_id := get_env_var("GUILD_ID")
	log_message(sess, guild_id, "made a live announcement in <#" + channel_id + ">.")
}

func call_api_youtube_live(service *youtube.Service, youtube_channel_id string, last_live **youtube.SearchResult, sess *discordgo.Session) {
	call := service.Search.List([]string{"snippet"}).
		MaxResults(1).
		ChannelId(youtube_channel_id).
		Order("date").
		EventType("live").
		Type("video")
	response, err := call.Do()
	if err != nil { log.Fatal(err) }

	if len(response.Items) > 0 {
		live := response.Items[0]

		if *last_live == nil || live.Id.VideoId != (*last_live).Id.VideoId {
			*last_live = live
			send_youtube_live_announcement(sess, *last_live)
		}
	}
}
