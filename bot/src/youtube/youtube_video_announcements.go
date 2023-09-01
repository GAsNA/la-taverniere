package main

import (
	"log"
	"net/http"
	"io"
	"encoding/xml"

	"github.com/bwmarrin/discordgo"
)

func send_youtube_video_announcement(sess *discordgo.Session, video Video) {
	message := ""
	channel_id := get_env_var("VIDEO_CHAN_ID")
	guild_id := get_env_var("GUILD_ID")

	var youtube_video_roles []youtube_video_role
	err := db.NewSelect().Model(&youtube_video_roles).
			Where("guild_id = ?", guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

	var ping_role_ids []string
	for i := 0; i < len(youtube_video_roles); i++ {
		ping_role_ids = append(ping_role_ids, youtube_video_roles[i].Role_ID)
	}

	for i := 0; i < len(ping_role_ids); i++ { message += "<@&" + ping_role_ids[i] + ">" }

	if len(ping_role_ids) > 0 { message += "\n" }

	message += video.Author.Name + " a posté une nouvelle vidéo. Enjoy!\n"
	
	message += video.Link.Href
	_, err = sess.ChannelMessageSend(channel_id, message)
	if err != nil { log.Println(err); return }

	log_message(sess, guild_id, "made a youtube announcement in <#" + channel_id + ">.")
	log.Println("A yt video announcement has been made in channel " + channel_id + " on guild id " + guild_id)
}

func call_api_youtube_video(youtube_channel_id string, last_video *Video, sess *discordgo.Session) {
	response, err := http.Get(get_env_var("YOUTUBE_LINK") + "/feeds/videos.xml?channel_id=" + youtube_channel_id)
	if err != nil { log.Println(err); return }

	content, err := io.ReadAll(response.Body)
	response.Body.Close()
	if err != nil { log.Println(err); return }

	var feed Feed

	xml.Unmarshal(content, &feed)

	if len(feed.Videos) > 0 {
		video := feed.Videos[0]

		if *last_video == (Video{}) {
			*last_video = video
			send_youtube_video_announcement(sess, *last_video)
		} else if video.Link.Href != (*last_video).Link.Href {
			*last_video = video
			send_youtube_video_announcement(sess, *last_video)
		}
	}
}
