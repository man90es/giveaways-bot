package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/man90es/giveaways-bot/db"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	guildID := os.Getenv("GIVEAWAY_GUILD_ID")

	session := getDiscordSession()
	defer session.Close()

	events, err := session.GuildScheduledEvents(guildID, false)
	if err != nil {
		log.Fatal(err)
	}

	for _, event := range events {
		db.NewEventFromDiscordEvent(event).Upsert()
	}

	session.AddHandler(func(_ *discordgo.Session, newEvent *discordgo.GuildScheduledEventCreate) {
		db.NewEventFromDiscordEvent(newEvent.GuildScheduledEvent).Upsert()
	})

	session.AddHandler(func(_ *discordgo.Session, updatedEvent *discordgo.GuildScheduledEventUpdate) {
		event, err := db.GetEventByID(updatedEvent.ID)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		isTrigger := event.TriggersRaffle
		wasScheduled := event.Status == discordgo.GuildScheduledEventStatusScheduled
		isActive := updatedEvent.Status == discordgo.GuildScheduledEventStatusActive

		if isTrigger && wasScheduled && isActive {
			raffle.RunRaffle(session, event.ID)
		}

		db.NewEventFromDiscordEvent(updatedEvent.GuildScheduledEvent).Upsert()
	})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Shutting down")
}
