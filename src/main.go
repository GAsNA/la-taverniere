package main

import (
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/bwmarrin/discordgo"
)

func getEnvVar(key string) string {
	err := godotenv.Load(".env")
	if err != nil { log.Fatal(err) }

	return os.Getenv(key)
}

func handlerHello(sess *discordgo.Session) {
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID { return }

		if m.Content == "Hello" {
			s.ChannelMessageSend(m.ChannelID, "world!")
		}
	})
}

func handlerWorld(sess *discordgo.Session) {
	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID { return }

		if m.Content == "world!" {
			s.ChannelMessageSend(m.ChannelID, "Hello")
		}
	})
}

func main() {
	discordToken := getEnvVar("DISCORD_TOKEN")

	sess, err := discordgo.New("Bot " + discordToken)
	if err != nil { log.Fatal(err) }

	handlerHello(sess)
	handlerWorld(sess)

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil { log.Fatal(err) }

	defer sess.Close()

	fmt.Println("The bot is online!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
