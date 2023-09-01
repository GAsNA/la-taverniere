package main

import (
	"net/http"
	"time"
	"encoding/xml"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type Feed struct {
	XMLName		xml.Name	`xml:"feed"`
	Videos		[]Video		`xml:"entry"`
}

type Video struct {
	XMLName		xml.Name	`xml:"entry"`
	Author		Author		`xml:"author"`
	Link		Link		`xml:"link"`
}

type Author	struct {
	XMLName		xml.Name	`xml:"author"`
	Name		string		`xml:"name"`		
}

type Link struct {
	XMLName		xml.Name	`xml:"link"`
	Href		string		`xml:"href,attr"`
}

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
//	var last_video *youtube.SearchResult
//	var last_live *youtube.SearchResult

//	service, err := get_service()
//	if err != nil { log.Println(err); return }

	var last_video Video
	
	youtube_channel_id := get_env_var("YOUTUBE_CHANNEL_ID")
	
	for true {
		// MAKE THE API CALL TO YOUTUBE FOR VIDEO
		call_api_youtube_video(youtube_channel_id, &last_video, sess)

		// sleep for 30 seconds
		time.Sleep((30 * 1000) * time.Millisecond)

	/*	// MAKE THE API CALL TO YOUTUBE FOR LIVE
		call_api_youtube_live(service, youtube_channel_id, &last_live, sess)

		// sleep for 30 seconds
		time.Sleep((30 * 1000) * time.Millisecond)*/
	}
}
