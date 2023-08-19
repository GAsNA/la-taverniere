package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/youtube/v3"
)

func send_youtube_video_announcement(sess *discordgo.Session, video *youtube.SearchResult) {
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
	_, err = sess.ChannelMessageSend(channel_id, message)
	if err != nil { log.Println(err); return }

	log_message(sess, guild_id, "made a youtube announcement in <#" + channel_id + ">.")
	log.Println("A yt video announcement has been made in channel <#" + channel_id + "> on guild id " + guild_id)
}

func call_api_youtube_video(service *youtube.Service, youtube_channel_id string, last_video **youtube.SearchResult, sess *discordgo.Session) {
	call := service.Search.List([]string{"snippet"}).
		MaxResults(1).
		ChannelId(youtube_channel_id).
		Order("date")
	response, err := call.Do()
	if err != nil { log.Println(err); return
	 }

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
