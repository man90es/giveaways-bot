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

func RunRaffle(session *discordgo.Session) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	channelID := os.Getenv("GIVEAWAY_CHANNEL_ID")
	guildID := os.Getenv("GIVEAWAY_GUILD_ID")
	participantRoleID := os.Getenv("GIVEAWAY_ROLE_ID")

	members, _ := session.GuildMembers(guildID, "", 1e3)
	if err != nil {
		log.Fatal(err)
	}

	// Select participants
	pastWinnerIDs := db.GetPastWinnerIDs()
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

		(&db.Participant{ID: member.User.ID, Username: member.User.Username}).Create()
		participants = append(participants, member)
	}

	winner, err := utils.RandomChoice(participants)
	if nil != err {
		log.Fatal(err.Error())
		return
	}

	availablePrizes := db.GetAvailablePrizes()
	prize, err := utils.RandomChoice(availablePrizes)
	if nil != err {
		log.Fatal(err.Error())
		return
	}

	(&db.Prize{ID: prize.ID}).AssignWinner(winner.User.ID)
	session.ChannelMessageSend(channelID, winner.Mention()+" won "+prize.Name)
}
