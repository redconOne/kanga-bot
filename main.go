package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const (
	PREFIX string = "!KangaBot"
	YELLOW string = "\033[33m"
	GREEN  string = "\033[32m"
	RED    string = "\033[31m"
	RESET  string = "\033[0m"
)

func main() {
	godotenv.Load()

	token := os.Getenv("BOT_TOKEN")
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}

	sess.AddHandler(messageHandler)
	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	u, err := sess.User("@me")
	if err != nil {
		log.Fatal(err)
	}

	BotName := u.Username

	err = sess.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer sess.Close()

	log.Printf("%s%s%s%s", GREEN, BotName, " is online!", RESET)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.GuildID == "" {
		log.Println("DM recieved")

		s.ChannelMessageSend(m.ChannelID, "Why hello there!")

		return
	}

	log.Println("Server message recieved")

	// TODO: ensure case insensitivity for args[0] aka (!kangabot, !kangaBot, !KaNgAbOt)
	args := strings.Split(m.Content, " ")
	if args[0] != PREFIX || len(args) < 2 {
		return
	}

	switch args[1] {
	case "hello":
		log.Println("Responding with test handshake")
		err := helloWorldHandler(s, m.ChannelID)
		if err != nil {
			log.Printf("%s%s%s", RED, "failed to respond with secret handshake", RESET)
			log.Fatal(err)
		}
	}
}

func helloWorldHandler(s *discordgo.Session, ChannelID string) error {
	_, err := s.ChannelMessageSend(ChannelID, "world")
	return err
}
