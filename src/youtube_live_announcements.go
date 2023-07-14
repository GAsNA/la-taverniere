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

func youtube_live_announcements(sess *discordgo.Session) {
	var last_live *youtube.SearchResult
	
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
			Order("date").
			EventType("live").
			Type("video")
		response, err := call.Do()
		if err != nil {
			log.Fatalf("Error doing the request: %v", err)
		}

		if len(response.Items) > 0 {
			live := response.Items[0]
	
			fmt.Println("LIVE:")
			fmt.Println("\tid: " + live.Id.VideoId + "\t\ttitle: " + live.Snippet.Title)

			if live == nil || live.Id.VideoId != last_live.Id.VideoId {
				last_live = live
				send_youtube_live_announcement(sess, last_live)
			}
		} else {
			fmt.Println("NO LIVE DETECTED")
		}

		// sleep for 1min
		time.Sleep((60 * 1000) * time.Millisecond)
	}
}
