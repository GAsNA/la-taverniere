package main

import (
	"log"
	"fmt"
	"os"
	"os/signal"
	"syscall"
  
	"github.com/bwmarrin/discordgo"
)

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
