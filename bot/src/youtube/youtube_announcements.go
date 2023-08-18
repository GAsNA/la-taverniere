package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func get_service() (*youtube.Service, error) {
	developerKey := get_env_var("YOUTUBE_API_KEY")

	youtube_client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(youtube_client)
	if err != nil { return nil, err }

	return service, nil
}

func youtube_announcements(sess *discordgo.Session) {
	var last_video *youtube.SearchResult
	var last_live *youtube.SearchResult

	service, err := get_service()
	if err != nil { log.Println(err); return }
	
	youtube_channel_id := get_env_var("YOUTUBE_CHANNEL_ID")
	
	for true {
		// MAKE THE API CALL TO YOUTUBE FOR VIDEO
		call_api_youtube_video(service, youtube_channel_id, &last_video, sess)

		// sleep for 30 seconds
		time.Sleep((30 * 1000) * time.Millisecond)

		// MAKE THE API CALL TO YOUTUBE FOR LIVE
		call_api_youtube_live(service, youtube_channel_id, &last_live, sess)

		// sleep for 30 seconds
		time.Sleep((30 * 1000) * time.Millisecond)
	}
}
