package main

import (
	"log"
	"math/rand"
	"os"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	channelID := os.Getenv("GIVEAWAY_CHANNEL_ID")
	APIToken := os.Getenv("DISCORD_API_TOKEN")
	guildID := os.Getenv("GIVEAWAY_GUILD_ID")
	participantRoleID := os.Getenv("GIVEAWAY_ROLE_ID")

	session, err := discordgo.New("Bot " + APIToken)
	if err != nil {
		log.Fatal(err)
	}

	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	members, _ := session.GuildMembers(guildID, "", 1e3)
	if err != nil {
		log.Fatal(err)
	}

	// Select members with specified role
	participants := []*discordgo.Member{}
	for _, member := range members {
		if slices.Contains(member.Roles, participantRoleID) {
			participants = append(participants, member)
		}
	}

	// Select and mention a random participant
	n := rand.Intn(len(participants))
	session.ChannelMessageSend(channelID, participants[n].User.Mention())
}
