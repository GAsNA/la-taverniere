package main

import (
	"log"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func send_youtube_announcement(sess *discordgo.Session, video *youtube.SearchResult, type_video string) {
	message := ""
	channel_id := ""

	switch type_video {
		case "upcoming":
			message = "Something is brewing on the channel of " + video.Snippet.ChannelTitle + ".\n"
			channel_id = get_env_var("UPCOMING_CHAN_ID")
		case "live":
			message = video.Snippet.ChannelTitle + " is live !\n"
			channel_id = get_env_var("LIVE_CHAN_ID")
		default:
			message = video.Snippet.ChannelTitle + " posted a new video.\n"
			channel_id = get_env_var("VIDEO_CHAN_ID")
	}

	message += "https://www.youtube.com/watch?v=" + video.Id.VideoId
	_, err := sess.ChannelMessageSend(channel_id, message)
	if err != nil { log.Fatal(err) }

	log_message(sess, "made a youtube announcement in <#" + channel_id + ">.")
}

func youtube_announcements(sess *discordgo.Session) {
	var last_video *youtube.SearchResult
	var live *youtube.SearchResult
	var upcoming *youtube.SearchResult
	
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
		fmt.Println("id: " + video.Id.VideoId + "\ttitle: " + video.Snippet.Title)

		switch video.Snippet.LiveBroadcastContent {
			case "upcoming":
				if upcoming == nil || upcoming.Id.VideoId != video.Id.VideoId {
					upcoming = video
					send_youtube_announcement(sess, upcoming, "upcoming")
				}
			case "live":
				if live == nil || live.Id.VideoId != live.Id.VideoId {
					live = video
					send_youtube_announcement(sess, live, "live")
				}
			default:
				if last_video == nil {
					last_video = video
				} else if last_video.Id.VideoId != last_video.Id.VideoId && (live == nil || last_video.Id.VideoId != live.Id.VideoId) {
					last_video = video
					send_youtube_announcement(sess, last_video, "video")
				}
		}

		// sleep for 1min
		time.Sleep((60 * 1000) * time.Millisecond)
	}
}
