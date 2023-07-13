package main

import (
	"flag"
	"log"
	"fmt"
	"net/http"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

func youtube_announcements() {
	flag.Parse()

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

	videos := response.Items
	
	fmt.Println("VIDEOS:")
	for _, item := range videos {
		fmt.Println("id: " + item.Id.VideoId + "\ttitle: " + item.Snippet.Title)
	}
}
//time.Sleep((60 * 1000) * time.Millisecond)
