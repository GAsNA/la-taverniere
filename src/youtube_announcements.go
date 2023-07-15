package main

import (
	"log"
	"net/http"
	"time"
	"strings"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/googleapi/transport"
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
			message += "A video is brewing on the channel of " + video.Snippet.ChannelTitle + "...\n"
		case "live":
			message += "A video of " + video.Snippet.ChannelTitle + " is live!\n"
		default:
			message += video.Snippet.ChannelTitle + " posted a new video. Enjoy!\n"
	}

	message += "https://www.youtube.com/watch?v=" + video.Id.VideoId
	_, err := sess.ChannelMessageSend(channel_id, message)
	if err != nil { log.Fatal(err) }

	log_message(sess, "made a youtube announcement in <#" + channel_id + ">.")
}

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
			message += "A live is brewing on the channel of " + live.Snippet.ChannelTitle + "...\n"
		case "live":
			message += live.Snippet.ChannelTitle + " is live! Come see!\n"
	}

	message += "https://www.youtube.com/watch?v=" + live.Id.VideoId
	_, err := sess.ChannelMessageSend(channel_id, message)
	if err != nil { log.Fatal(err) }

	log_message(sess, "made a live announcement in <#" + channel_id + ">.")
}

func get_service() *youtube.Service {
	developerKey := get_env_var("YOUTUBE_API_KEY")

	youtube_client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(youtube_client)
	if err != nil {
		log.Fatalf("Error creating new YouTube client: %v", err)
	}

	return service
}

func call_youtube_video(service *youtube.Service, youtube_channel_id string, last_video *youtube.SearchResult, sess *discordgo.Session) {
	call := service.Search.List([]string{"snippet"}).
		MaxResults(1).
		ChannelId(youtube_channel_id).
		Order("date")
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error doing the request: %v", err)
	}

	if len(response.Items) > 0 {
		video := response.Items[0]

		if last_video == nil {
			last_video = video
		} else if video.Id.VideoId != last_video.Id.VideoId {
			last_video = video
			send_youtube_video_announcement(sess, last_video)
		}
	}
}

func call_youtube_live(service *youtube.Service, youtube_channel_id string, last_live *youtube.SearchResult, sess *discordgo.Session) {
	call := service.Search.List([]string{"snippet"}).
		MaxResults(1).
		ChannelId(youtube_channel_id).
		Order("date").
		EventType("live").
		Type("video")
	response, err := call.Do()
	if err != nil {
		log.Fatalf("Error doing the request: %v", err)
	}

	if len(response.Items) > 0 {
		live := response.Items[0]

		if last_live == nil || live.Id.VideoId != last_live.Id.VideoId {
			last_live = live
			send_youtube_live_announcement(sess, last_live)
		}
	}
}

func youtube_announcements(sess *discordgo.Session) {
	var last_video *youtube.SearchResult
	var last_live *youtube.SearchResult

	service := get_service()
	youtube_channel_id := get_env_var("YOUTUBE_CHANNEL_ID")
	
	for true {
		// MAKE THE API CALL TO YOUTUBE FOR VIDEO
		call_youtube_video(service, youtube_channel_id, last_video, sess)

		// sleep for 30 seconds
		time.Sleep((30 * 1000) * time.Millisecond)

		// MAKE THE API CALL TO YOUTUBE FOR LIVE
		call_youtube_live(service, youtube_channel_id, last_live, sess)

		// sleep for 30 seconds
		time.Sleep((30 * 1000) * time.Millisecond)
	}
}
