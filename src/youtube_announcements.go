package main

import (
	"log"
	"fmt"
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

func youtube_announcements(sess *discordgo.Session) {
	var last_video *youtube.SearchResult
	
	for true {
		developerKey := get_env_var("YOUTUBE_API_KEY")

		youtube_client := &http.Client{
			Transport: &transport.APIKey{Key: developerKey},
		}

		service, err := youtube.New(youtube_client)
		if err != nil {
			log.Fatalf("Error creating new YouTube client: %v", err)
		}

		channelId := get_env_var("YOUTUBE_CHANNEL_ID")
	
		// MAKE THE API CALL TO YOUTUBE
		call := service.Search.List([]string{"snippet"}).
			MaxResults(1).
			ChannelId(channelId).
			Order("date")
		response, err := call.Do()
		if err != nil {
			log.Fatalf("Error doing the request: %v", err)
		}

		video := response.Items[0]
	
		fmt.Println("VIDEO:")
		fmt.Println("\tid: " + video.Id.VideoId + "\t\ttitle: " + video.Snippet.Title)

		if last_video == nil {
			last_video = video
		} else if video.Id.VideoId != last_video.Id.VideoId {
			last_video = video
			send_youtube_video_announcement(sess, last_video)
		}

		// sleep for 1min
		time.Sleep((60 * 1000) * time.Millisecond)
	}
}
