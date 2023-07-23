package main

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/youtube/v3"
)

func send_youtube_video_announcement(sess *discordgo.Session, video *youtube.SearchResult) {
	message := ""
	channel_id := get_env_var("VIDEO_CHAN_ID")

	ping_role_ids_env := get_env_var("PING_YOUTUBE_VIDEO_ROLE_IDS")
	ping_role_ids := strings.Split(ping_role_ids_env, ",")

	for i := 0; i < len(ping_role_ids); i++ {
		message += "<@&" + ping_role_ids[i] + ">"
	}

	if len(ping_role_ids) > 0 {
		message += "\n"
	}

	switch video.Snippet.LiveBroadcastContent {
		case "upcoming":
			message += "Une vidéo se prépare sur la chaine de  " + video.Snippet.ChannelTitle + "...\n"
		case "live":
			message += "Une vidéo de " + video.Snippet.ChannelTitle + " est en live !\n"
		default:
			message += video.Snippet.ChannelTitle + " a posté une nouvelle vidéo. Enjoy!\n"
	}

	message += get_env_var("YOUTUBE_LINK") + "/watch?v=" + video.Id.VideoId
	_, err := sess.ChannelMessageSend(channel_id, message)
	if err != nil { log.Fatal(err) }

	log_message(sess, "made a youtube announcement in <#" + channel_id + ">.")
}

func call_api_youtube_video(service *youtube.Service, youtube_channel_id string, last_video **youtube.SearchResult, sess *discordgo.Session) {
	call := service.Search.List([]string{"snippet"}).
		MaxResults(1).
		ChannelId(youtube_channel_id).
		Order("date")
	response, err := call.Do()
	if err != nil { log.Fatal(err) }

	if len(response.Items) > 0 {
		video := response.Items[0]

		if *last_video == nil {
			*last_video = video
		} else if video.Id.VideoId != (*last_video).Id.VideoId {
			*last_video = video
			send_youtube_video_announcement(sess, *last_video)
		}
	}
}
