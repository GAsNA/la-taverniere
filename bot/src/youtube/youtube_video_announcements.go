package main

import (
	"log"
	"net/http"
	"io"
	"encoding/xml"

	"github.com/bwmarrin/discordgo"
)

func send_youtube_video_announcement(sess *discordgo.Session, video Video, channel_id string, guild_id string) {
	message := ""
	
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
	guild_id := get_env_var("GUILD_ID")
	
	// CHECK IF CHANNEL IS SET
	var channels_for_actions []channel_for_action
	err := db.NewSelect().Model(&channels_for_actions).
			Where("action_id = ? AND guild_id = ?", get_action_db_by_name("Youtube Video Announcements").id, guild_id).
			Scan(ctx)
	if err != nil { log.Println(err); return }

	if len(channels_for_actions) == 0 { return }
	channel_id := channels_for_actions[0].Channel_ID

	// SEARCH IF NEW VIDEO UPLOADED
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
		} else if video.Link.Href != (*last_video).Link.Href {
			*last_video = video
			send_youtube_video_announcement(sess, *last_video, channel_id, guild_id)
		}
	}
}
