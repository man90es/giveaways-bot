package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/man90es/giveaways-bot/db"
	"github.com/man90es/giveaways-bot/raffle"
)

func getDiscordSession() *discordgo.Session {
	session, err := discordgo.New("Bot " + db.Config[db.ConfigKeyDiscordAPIToken])
	if err != nil {
		log.Fatal("Error occured while trying to create a Discord session: ", err.Error())
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	session.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = session.Open()
	if err != nil {
		log.Fatal("Error occured while trying to open a Discord session: ", err.Error())
	}

	return session
}

func main() {
	err := db.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err.Error())
		return
	}

	session := getDiscordSession()
	defer session.Close()

	events, err := session.GuildScheduledEvents(db.Config[db.ConfigKeyGiveawayGuildlID], false)
	if err != nil {
		log.Println("Error occured while trying to retrieve Discord events: ", err.Error())
		return
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
			log.Println("Error occured while trying to retrieve an event from the DB: ", err.Error())
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
