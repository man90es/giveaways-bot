package raffle

import (
	"log"
	"os"
	"slices"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/man90es/giveaways-bot/db"
	"github.com/man90es/giveaways-bot/utils"
)

func RunRaffle(session *discordgo.Session, eventID string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	channelID := os.Getenv("GIVEAWAY_CHANNEL_ID")
	guildID := os.Getenv("GIVEAWAY_GUILD_ID")
	participantRoleID := os.Getenv("GIVEAWAY_ROLE_ID")

	availablePrizes, err := db.GetAvailablePrizes()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	prize, err := utils.RandomChoice(availablePrizes)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	(&db.Prize{ID: prize.ID}).AssignEvent(eventID)

	members, err := session.GuildMembers(guildID, "", 1e3)
	if err != nil {
		log.Fatal(err)
		return
	}

	pastWinnerIDs, err := db.GetPastWinnerIDs()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// Select participants
	participants := []*discordgo.Member{}
	for _, member := range members {
		// User already won something
		if slices.Contains(pastWinnerIDs, member.User.ID) {
			continue
		}

		// User doesn't have a required role
		if !slices.Contains(member.Roles, participantRoleID) {
			continue
		}

		db.NewParticipantFromDiscordUser(member.User).Upsert()
		participants = append(participants, member)
	}

	winner, err := utils.RandomChoice(participants)
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	(&db.Prize{ID: prize.ID}).AssignWinner(winner.User.ID)
	session.ChannelMessageSend(channelID, winner.Mention()+" won "+prize.Name)
}
