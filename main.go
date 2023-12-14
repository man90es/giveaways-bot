package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/man90es/giveaways-bot/raffle"
)

func getDiscordSession() *discordgo.Session {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	APIToken := os.Getenv("DISCORD_API_TOKEN")

	session, err := discordgo.New("Bot " + APIToken)
	if err != nil {
		log.Fatal(err)
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}

	return session
}

func main() {
	session := getDiscordSession()
	defer session.Close()

	raffle.RunRaffle(session)
}
